package validate

import "regexp"

// If ip is a valid IPv4 address
func IsValidIPv4(ip string) bool {
    if regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`).MatchString(ip) == false {
        return false
    }

    return true
}
