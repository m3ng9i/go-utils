package dns

import "fmt"
import "errors"
import "time"

import mdns "github.com/miekg/dns"

/* Get a domain's IPs from a specific name server.

Parameters:
    domain      the domain you want to query
    nameserver  name server's IP address
    port        53 in general
    net         tcp or udp
    timeout     in seconds, can be omitted

Here's an example：
    r, e := ARecords("www.example.com", "8.8.8.8", 53, "tcp")
    if e != nil {
        fmt.Println(e)
    } else {
        fmt.Println(r)
    }
*/
func ARecords(domain, nameserver string, port uint16, net string, timeout ...uint8) ([]string, error) {
    var result []string

    if net != "tcp" && net != "udp" {
        return result, errors.New("The parameter 'net' should only be 'tcp' or 'udp'.")
    }

    msg := new(mdns.Msg)
    msg.SetQuestion(mdns.Fqdn(domain), mdns.TypeA)

    var client *mdns.Client
    if len(timeout) > 0 {
        tm := time.Duration(timeout[0]) * time.Second
        client = &mdns.Client { Net: net, DialTimeout: tm, ReadTimeout: tm, WriteTimeout: tm }
    } else {
        client = &mdns.Client { Net: net }
    }

    r, _, err := client.Exchange(msg, fmt.Sprintf("%s:%d", nameserver, port))
    if err != nil {
        return result, err
    }

    for _, i := range r.Answer {
        if t, ok := i.(*mdns.A); ok {
            result = append(result, t.A.String())
        }
    }

    return result, nil
}

