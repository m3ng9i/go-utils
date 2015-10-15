package http

import "net/http"
import "net/url"
import "strings"
import "fmt"


// RedirectToHTTPS redirect a request to a corresponding https url
func RedirectToHTTPS(port uint) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        host := fmt.Sprintf("%s:%d", strings.Split(r.Host, ":")[0], port)
        base := url.URL {
            Scheme: "https",
            Host: host,
        }
        newURL := base.ResolveReference(r.URL).String()
        http.Redirect(w, r, newURL, http.StatusTemporaryRedirect)
    }
}
