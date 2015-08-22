package http

import "net"
import "net/http"
import "strings"


// Get client IP.
func GetIP(r *http.Request) string {
    host, _, _ := net.SplitHostPort(r.RemoteAddr)
    return host
}


/*
Get a value of key in query string.

You should make sure to call (*http.Request).ParseForm() first, then to call this function.

Space before and after the value will be striped.
If the key appears more than once, the last value will be get.

parameters:
    r       *http.Request
    key     Key in query string
    defval  Default value of key, if not found the key, return the default value. If not provide defval, empty string will be used.
*/
func QueryValue(r *http.Request, key string, defval ...string) string {
    q := r.Form[key]
    length := len(q)
    if length == 0 {
        if len(defval) > 0 {
            return defval[0]
        } else {
            return ""
        }
    }
    return strings.TrimSpace(q[length -1])
}
