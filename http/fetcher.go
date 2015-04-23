package http

import "net"
import "net/http"
import "io/ioutil"
import "time"
import "golang.org/x/net/proxy"


type ProxyConfig struct {
    Addr string
    Username string
    Password string
}


func Socks5Client(conf ProxyConfig) (client *http.Client, err error) {

    var proxyAuth *proxy.Auth
    if conf.Username != "" {
        proxyAuth = new(proxy.Auth)
        proxyAuth.User = conf.Username
        proxyAuth.Password = conf.Password
    }

    dialer, err := proxy.SOCKS5(
        "tcp",
        conf.Addr,
        proxyAuth,
        &net.Dialer {
            Timeout: 30 * time.Second,
            KeepAlive: 30 * time.Second,
        })
    if err != nil {
        return
    }

    var transport http.RoundTripper = &http.Transport {
        Proxy: http.ProxyFromEnvironment,
        Dial: dialer.Dial,
        TLSHandshakeTimeout: 10 * time.Second,
    }

    client = &http.Client { Transport: transport }
    return
}


type Fetcher struct {
    Client *http.Client
    Headers map[string]string
}


func (this *Fetcher) FetchAll(url string) (b []byte, err error) {

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return
    }

    if this.Headers != nil {
        for key, value := range this.Headers {
            req.Header.Set(key, value)
        }
    }

    var client *http.Client
    if this.Client != nil {
        client = this.Client
    } else {
        client = http.DefaultClient
    }

    resp, err := client.Do(req)
    defer resp.Body.Close()
    if err != nil {
        return
    }

    b, err = ioutil.ReadAll(resp.Body)
    return
}


func NewFetcher(client *http.Client, headers map[string]string) *Fetcher {
    return &Fetcher { Client: client, Headers: headers }
}
