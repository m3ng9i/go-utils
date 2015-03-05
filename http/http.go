package http

import "net/http"
import "strings"


// Get client IP
func GetIP(r *http.Request) string {
    ip := r.Header.Get("X-Real-IP")
    if ip == "" {
        ip = r.Header.Get("X-Forwarded-For")
        if ip == "" {
            ip = r.RemoteAddr
        }
    }
    return strings.Split(ip, ":")[0]
}
