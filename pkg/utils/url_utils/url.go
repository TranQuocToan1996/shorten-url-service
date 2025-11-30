package url_utils

import (
	"fmt"
	"net"
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
	host := u.Hostname()
	ips, err := net.LookupIP(host)
	if err != nil {
		return false
	}

	for _, ip := range ips {
		if isPrivateIP(ip) {
			return false
		}
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

func ValidateWebhookCallback(callbackURL string) error {
	u, err := url.Parse(callbackURL)
	if err != nil {
		return err
	}

	if u.Scheme != "https" {
		return fmt.Errorf("callback must use https")
	}

	host := u.Hostname()
	ips, err := net.LookupIP(host)
	if err != nil {
		return err
	}

	for _, ip := range ips {
		if isPrivateIP(ip) {
			return fmt.Errorf("callback resolves to private IP: %s", ip)
		}
	}

	return nil
}

var privateCIDRs []*net.IPNet

func init() {
	blocks := []string{
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
		"0.0.0.0/8",
		"100.64.0.0/10",
		"224.0.0.0/4",
		"240.0.0.0/4",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
		"::/128",
		"ff00::/8",
	}

	for _, cidr := range blocks {
		_, block, _ := net.ParseCIDR(cidr)
		privateCIDRs = append(privateCIDRs, block)
	}
}

func isPrivateIP(ip net.IP) bool {
	for _, block := range privateCIDRs {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
