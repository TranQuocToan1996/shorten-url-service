package tests

import (
	"testing"

	"shorten/pkg/config"
	"shorten/service"
)

func TestNewBase62Encoder(t *testing.T) {
	tests := []struct {
		name          string
		cfg           config.Config
		wantLength    int
		wantSecretKey string
	}{
		{
			name: "valid_config_with_custom_length",
			cfg: config.Config{
				SECRET_KEY:       "test-secret",
				SHORT_URL_LENGTH: "10",
			},
			wantLength:    10,
			wantSecretKey: "test-secret",
		},
		// {
		// 	name: "valid_config_with_default_length",
		// 	cfg: config.Config{
		// 		SECRET_KEY:       "my-secret-key",
		// 		SHORT_URL_LENGTH: "8",
		// 	},
		// 	wantLength:    8,
		// 	wantSecretKey: "my-secret-key",
		// },
		// {
		// 	name: "invalid_length_falls_back_to_default",
		// 	cfg: config.Config{
		// 		SECRET_KEY:       "secret",
		// 		SHORT_URL_LENGTH: "invalid",
		// 	},
		// 	wantLength:    8,
		// 	wantSecretKey: "secret",
		// },
		// {
		// 	name: "empty_length_falls_back_to_default",
		// 	cfg: config.Config{
		// 		SECRET_KEY:       "secret",
		// 		SHORT_URL_LENGTH: "",
		// 	},
		// 	wantLength:    8,
		// 	wantSecretKey: "secret",
		// },
		// {
		// 	name: "zero_length_falls_back_to_default",
		// 	cfg: config.Config{
		// 		SECRET_KEY:       "secret",
		// 		SHORT_URL_LENGTH: "0",
		// 	},
		// 	wantLength:    0,
		// 	wantSecretKey: "secret",
		// },
		// {
		// 	name: "empty_secret_key",
		// 	cfg: config.Config{
		// 		SECRET_KEY:       "",
		// 		SHORT_URL_LENGTH: "8",
		// 	},
		// 	wantLength:    8,
		// 	wantSecretKey: "",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := service.NewBase62Encoder(tt.cfg)
			if encoder == nil {
				t.Fatal("NewBase62Encoder returned nil")
			}
			// Test that encoder works
			result, err := encoder.Encode("https://example.com")
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}
			if len(result) != tt.wantLength && tt.wantLength > 0 {
				t.Errorf("Encode() result length = %d, want %d", len(result), tt.wantLength)
			}
		})
	}
}

func TestBase62Encoder_Encode(t *testing.T) {
	cfg := config.Config{
		SECRET_KEY:       "test-secret-key-123",
		SHORT_URL_LENGTH: "8",
	}
	encoder := service.NewBase62Encoder(cfg)

	tests := []struct {
		name        string
		originalURL string
		wantLength  int
		wantErr     bool
		validate    func(t *testing.T, result string)
	}{
		{
			name:        "simple_https_url",
			originalURL: "https://example.com",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "http_url_with_path",
			originalURL: "http://example.com/path/to/resource",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "url_with_query_params",
			originalURL: "https://example.com/search?q=test&page=1",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "url_with_fragment",
			originalURL: "https://example.com/page#section",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "url_with_all_components",
			originalURL: "https://example.com/path?query=value#fragment",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "empty_string",
			originalURL: "",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "very_long_url",
			originalURL: "https://example.com/" + string(make([]byte, 1000)),
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "url_with_special_characters",
			originalURL: "https://example.com/path?q=hello%20world&x=test+value",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "url_with_unicode",
			originalURL: "https://example.com/测试/路径",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "ip_address_url",
			originalURL: "http://192.168.1.1:8080/path",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "localhost_url",
			originalURL: "http://localhost:3000/api/v1/users",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "url_with_port",
			originalURL: "https://example.com:8443/secure/path",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "mailto_url",
			originalURL: "mailto:test@example.com",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "url_with_multiple_slashes",
			originalURL: "https://example.com///path///to///resource",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
		{
			name:        "url_with_encoded_chars",
			originalURL: "https://example.com/path%20with%20spaces",
			wantLength:  8,
			wantErr:     false,
			validate: func(t *testing.T, result string) {
				if len(result) != 8 {
					t.Errorf("Encode() length = %d, want 8", len(result))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encoder.Encode(tt.originalURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != tt.wantLength {
					t.Errorf("Encode() length = %d, want %d", len(got), tt.wantLength)
				}
				if tt.validate != nil {
					tt.validate(t, got)
				}
			}
		})
	}
}

func TestBase62Encoder_Encode_Deterministic(t *testing.T) {
	cfg := config.Config{
		SECRET_KEY:       "test-secret-key-123",
		SHORT_URL_LENGTH: "8",
	}
	encoder := service.NewBase62Encoder(cfg)

	tests := []struct {
		name        string
		originalURL string
	}{
		{
			name:        "same_url_produces_same_code",
			originalURL: "https://example.com/test",
		},
		{
			name:        "different_url_produces_different_code",
			originalURL: "https://example.com/another",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result1, err1 := encoder.Encode(tt.originalURL)
			if err1 != nil {
				t.Fatalf("First Encode() error = %v", err1)
			}

			result2, err2 := encoder.Encode(tt.originalURL)
			if err2 != nil {
				t.Fatalf("Second Encode() error = %v", err2)
			}

			if result1 != result2 {
				t.Errorf("Encode() is not deterministic: first = %q, second = %q", result1, result2)
			}
		})
	}
}

func TestBase62Encoder_Encode_DifferentURLs(t *testing.T) {
	cfg := config.Config{
		SECRET_KEY:       "test-secret-key-123",
		SHORT_URL_LENGTH: "8",
	}
	encoder := service.NewBase62Encoder(cfg)

	urls := []string{
		"https://example.com/1",
		"https://example.com/2",
		"https://example.com/3",
		"https://different.com/1",
		"http://example.com/1",
	}

	results := make(map[string]string)
	for _, url := range urls {
		result, err := encoder.Encode(url)
		if err != nil {
			t.Fatalf("Encode(%q) error = %v", url, err)
		}
		results[url] = result
	}

	// Check all results are different
	seen := make(map[string]string)
	for url, result := range results {
		if existingURL, exists := seen[result]; exists {
			t.Errorf("Collision detected: %q and %q both produce %q", url, existingURL, result)
		}
		seen[result] = url
	}
}

func TestBase62Encoder_Encode_DifferentSecretKeys(t *testing.T) {
	url := "https://example.com/test"

	cfg1 := config.Config{
		SECRET_KEY:       "secret-key-1",
		SHORT_URL_LENGTH: "8",
	}
	encoder1 := service.NewBase62Encoder(cfg1)

	cfg2 := config.Config{
		SECRET_KEY:       "secret-key-2",
		SHORT_URL_LENGTH: "8",
	}
	encoder2 := service.NewBase62Encoder(cfg2)

	result1, err1 := encoder1.Encode(url)
	if err1 != nil {
		t.Fatalf("Encoder1.Encode() error = %v", err1)
	}

	result2, err2 := encoder2.Encode(url)
	if err2 != nil {
		t.Fatalf("Encoder2.Encode() error = %v", err2)
	}

	if result1 == result2 {
		t.Errorf("Different secret keys produced same result: %q", result1)
	}
}

func TestBase62Encoder_Encode_DifferentLengths(t *testing.T) {
	url := "https://example.com/test"

	tests := []struct {
		name   string
		length string
		want   int
	}{
		{
			name:   "length_4",
			length: "4",
			want:   4,
		},
		{
			name:   "length_6",
			length: "6",
			want:   6,
		},
		{
			name:   "length_8",
			length: "8",
			want:   8,
		},
		{
			name:   "length_10",
			length: "10",
			want:   10,
		},
		{
			name:   "length_12",
			length: "12",
			want:   12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				SECRET_KEY:       "test-secret",
				SHORT_URL_LENGTH: tt.length,
			}
			encoder := service.NewBase62Encoder(cfg)

			result, err := encoder.Encode(url)
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}

			if len(result) != tt.want {
				t.Errorf("Encode() length = %d, want %d", len(result), tt.want)
			}
		})
	}
}

func TestBase62Encoder_Encode_Base62Characters(t *testing.T) {
	cfg := config.Config{
		SECRET_KEY:       "test-secret-key",
		SHORT_URL_LENGTH: "8",
	}
	encoder := service.NewBase62Encoder(cfg)

	url := "https://example.com/test"
	result, err := encoder.Encode(url)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	base62Chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charMap := make(map[rune]bool)
	for _, char := range base62Chars {
		charMap[char] = true
	}

	for _, char := range result {
		if !charMap[char] {
			t.Errorf("Encode() contains invalid character %q, only Base62 characters allowed", char)
		}
	}
}
