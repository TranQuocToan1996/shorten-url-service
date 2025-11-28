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
