package url_utils

import (
	"net/url"
	"strings"
)

func IsValidURL(rawURL string) bool {
	if len(rawURL) > 2048 {
		return false
	}
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return false
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}
	if u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func IsBase62(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, ch := range s {
		switch {
		case ch >= '0' && ch <= '9':
		case ch >= 'A' && ch <= 'Z':
		case ch >= 'a' && ch <= 'z':
		default:
			return false
		}
	}

	return true
}
