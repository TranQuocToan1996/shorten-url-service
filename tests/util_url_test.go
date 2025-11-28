package tests

import (
	"testing"

	utilsurl "shorten/pkg/utils/url_utils"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name   string
		rawURL string
		want   bool
	}{
		{
			name:   "valid_https_with_host",
			rawURL: "https://example.com/path",
			want:   true,
		},
		{
			name:   "without_scheme",
			rawURL: "example.org/resource",
			want:   false,
		},
		{
			name:   "invalid_empty_host",
			rawURL: "http://",
			want:   false,
		},
		{
			name:   "invalid_unparsable",
			rawURL: "://missing-scheme.com",
			want:   false,
		},
		{
			name:   "invalid_empty_string",
			rawURL: "",
			want:   false,
		},
		{
			name:   "valid_https_with_host_and_query_and_fragment",
			rawURL: "http://example.com/path?a=1#x",
			want:   true,
		},
		{
			name:   "localhost_valid",
			rawURL: "http://192.168.1.1/page",
			want:   true,
		},
		{
			name:   "invalid_ftp_scheme",
			rawURL: "ftp://example.com",
			want:   false,
		},
		{
			name:   "invalid",
			rawURL: "abc",
			want:   false,
		},
		{
			name:   "valid_ip4",
			rawURL: "https://8.8.8.8",
			want:   true,
		},
		{
			name:   "valid_ip6",
			rawURL: "https://[2001:4860:4860::8888]",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilsurl.IsValidURL(tt.rawURL)
			if got != tt.want {
				t.Fatalf("IsValidURLLoose(%q) = %v, want %v", tt.rawURL, got, tt.want)
			}
		})
	}
}
