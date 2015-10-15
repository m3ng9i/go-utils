package tls

// This file is a modification of go's source file: go/src/crypto/tls/generate_cert.go

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)


// EcdsaCurve is is used to generate a key
type EcdsaCurve string


const (
    DefCurve EcdsaCurve = ""        // DefCurve generates an RSA keypair of given bit size
    P224     EcdsaCurve = "P224"    // P224 Curve
    P256     EcdsaCurve = "P256"    // P256 Curve
    P384     EcdsaCurve = "P384"    // P384 Curve
    P521     EcdsaCurve = "P521"    // P521 Curve
)


type CertOption struct {
    PublicKey       string          // Path of public key.
    PrivateKey      string          // Path of private key.
    Host            string          // Comma-separated hostnames and IPs to generate a certificate for.
    ValidFrom       *time.Time      // Creation time.
    ValidFor        time.Duration   // Duration that certificate is valid for.
    IsCA            bool            // Whether this cert should be its own Certificate Authority.
    RsaBits         int             // Size of RSA key to generate. Ignored if EcdsaCurve is set.
    EcdsaCurve      EcdsaCurve      // ECDSA curve to use to generate a key. Valid values are P224, P256, P384, P521.
    Organization    string          // Organization name
}


func DefaultCertOption() *CertOption {
    now := time.Now()
    return &CertOption {
        PublicKey:      "cert.pem",
        PrivateKey:     "key.pem",
        Host:           "localhost,127.0.0.1",
        ValidFrom:      &now,
        ValidFor:       365*24*time.Hour*10, // 10 years
        IsCA:           false,
        RsaBits:        2048,
        EcdsaCurve:     DefCurve,
        Organization:   "Acme Co",
    }
}


func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}


func pemBlockForKey(priv interface{}) (*pem.Block, error) {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}, nil
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
            err = fmt.Errorf("Unable to marshal ECDSA private key: %v", err)
            return nil, err
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
	default:
		return nil, nil
	}
}


// MakeCert generate a self-signed X.509 certificate for a TLS server.
func MakeCert(option *CertOption) error {

    if len(option.Host) == 0 {
        return fmt.Errorf("Option Host could not be empty.")
    }

	var priv interface{}
    var err error
	switch option.EcdsaCurve {
        case DefCurve:
            priv, err = rsa.GenerateKey(rand.Reader, option.RsaBits)
        case P224:
            priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
        case P256:
            priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
        case P384:
            priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
        case P521:
            priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
        default:
            err = fmt.Errorf("Unrecognized elliptic curve: %q", option.EcdsaCurve)
	}
	if err != nil {
		return fmt.Errorf("Failed to generate private key: %s.", err)
	}

	var notBefore time.Time

    if option.ValidFrom == nil {
		notBefore = time.Now()
    } else {
        notBefore = *option.ValidFrom
    }

	notAfter := notBefore.Add(option.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)

    serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
        return fmt.Errorf("Failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{option.Organization},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(option.Host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if option.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

    derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
        return fmt.Errorf("Failed to create certificate: %s", err)
	}

    // Create public key
	certOut, err := os.Create(option.PublicKey)
	if err != nil {
        if e, ok := err.(*os.PathError); ok {
            err = e.Err
        }
        return fmt.Errorf("Failed to open %s for writing: %s", option.PublicKey, err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
    // written public key

    keyOut, err := os.OpenFile(option.PrivateKey, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
        if e, ok := err.(*os.PathError); ok {
            err = e.Err
        }
        return fmt.Errorf("Failed to open %s for writing: %s", option.PrivateKey, err)
	}
	defer keyOut.Close()
    pemBlock, err := pemBlockForKey(priv)
    if err != nil {
        return err
    }
	pem.Encode(keyOut, pemBlock)
    // written private key

    return nil
}

