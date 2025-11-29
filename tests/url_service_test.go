package tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"shorten/model"
	"shorten/pkg/config"
	"shorten/pkg/dto"
	"shorten/service"
)

// --- Mocks ---
type mockProducer struct {
	callPayload []byte
	err         error
}

func (m *mockProducer) Publish(_ context.Context, _ string, payload []byte) error {
	m.callPayload = payload
	return m.err
}
func (m *mockProducer) Close() error { return nil }

type mockRepo struct {
	byCode     *model.ShortenURL
	byLongURL  *model.ShortenURL
	errByCode  error
	errByLong  error
	saveErr    error
	saveCalled bool
}

func (r *mockRepo) GetByCode(code string) (*model.ShortenURL, error) {
	return r.byCode, r.errByCode
}

func (r *mockRepo) GetByLongURL(longURL string) (*model.ShortenURL, error) {
	return r.byLongURL, r.errByLong
}

func (r *mockRepo) Save(su *model.ShortenURL) error {
	r.saveCalled = true
	return r.saveErr
}

type mockCache struct {
	getData []byte
	getErr  error
	setArgs []interface{}
}

func (c *mockCache) Get(_ context.Context, key string) ([]byte, error) {
	return c.getData, c.getErr
}

func (c *mockCache) Set(_ context.Context, key string, value []byte, expiration time.Duration) error {
	c.setArgs = []interface{}{key, value, expiration}
	return nil
}
func (c *mockCache) Delete(_ context.Context, _ string) error { return nil }

type mockEncoder struct {
	toCode string
	err    error
}

func (e *mockEncoder) Encode(_ string) (string, error) { return e.toCode, e.err }

func makeConfig() config.Config {
	return config.Config{
		QUEUE_NAME:    "queue",
		REDIRECT_HOST: "https://short/",
	}
}

func TestSubmitURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		valid   bool
		errOut  error
		prodErr error
	}{
		{"valid", "https://abc.com/def", true, nil, nil},
		{"invalid_url", "ftp://abc.com", false, fmt.Errorf("invalid URL: ftp://abc.com"), nil},
		{"prod_fail", "https://any.com/1", true, errors.New("fail"), errors.New("fail")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &mockProducer{err: tt.prodErr}
			us := service.NewURLService(makeConfig(), &mockRepo{}, p, &mockCache{}, service.NewBase62Encoder(makeConfig()))
			err := us.SubmitURL(context.Background(), tt.url)
			if (err != nil) != (tt.errOut != nil) {
				t.Fatalf("err mismatch, got %v, want %v", err, tt.errOut)
			}
			if tt.errOut != nil && err.Error() != tt.errOut.Error() {
				t.Fatalf("err value: got %v, want %v", err, tt.errOut)
			}
			if tt.valid && p.callPayload == nil {
				t.Errorf("expected payload published")
			}
		})
	}
}

func TestGetDecode(t *testing.T) {
	result := &model.ShortenURL{LongURL: "foo", Code: "bar"}
	data, _ := json.Marshal(result)
	cfg := makeConfig()

	tests := []struct {
		name        string
		cache       *mockCache
		repo        *mockRepo
		inputURL    string
		errExpected bool
	}{
		{"cache_hit_and_unmarshal_ok", &mockCache{getData: data}, &mockRepo{}, cfg.REDIRECT_HOST + "bar", false},
		{"cache_hit_unmarshal_fail_db_ok", &mockCache{getData: []byte("notjson"), getErr: nil}, &mockRepo{byCode: result}, cfg.REDIRECT_HOST + "bar", false},
		{"cache_miss_db_ok", &mockCache{getErr: errors.New("not found")}, &mockRepo{byCode: result}, cfg.REDIRECT_HOST + "bar", false},
		{"cache_miss_db_fail", &mockCache{getErr: errors.New("not found")}, &mockRepo{byCode: nil, errByCode: errors.New("xx")}, cfg.REDIRECT_HOST + "bar", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := service.NewURLService(cfg, tt.repo, &mockProducer{}, tt.cache, service.NewBase62Encoder(makeConfig()))
			res, err := us.GetDecode(context.Background(), tt.inputURL)
			if tt.errExpected && err == nil {
				t.Fatalf("expected error, got none")
			}
			if !tt.errExpected && err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}
			if !tt.errExpected && res == nil {
				t.Fatalf("expected result, got nil")
			}
		})
	}
}

func TestHandleShortenURL(t *testing.T) {
	goodMsg := (&dto.URLMessage{URL: "https://a.com"}).Bytes()
	badMsg, _ := json.Marshal(map[string]any{"bad": "val"})

	tests := []struct {
		name    string
		payload []byte
		repo    *mockRepo
		encoder *mockEncoder
		errOut  bool
	}{
		{"valid_new", goodMsg, &mockRepo{}, &mockEncoder{toCode: "abc"}, false},
		{"invalid_url", (&dto.URLMessage{URL: "ftp://zz"}).Bytes(), &mockRepo{}, &mockEncoder{}, true},
		{"unmarshal_fail", badMsg, &mockRepo{}, &mockEncoder{}, true},
		{"url_already_exists", goodMsg, &mockRepo{byLongURL: &model.ShortenURL{LongURL: "https://a.com"}}, &mockEncoder{}, false},
		{"encode_fail", goodMsg, &mockRepo{}, &mockEncoder{err: errors.New("encfail")}, true},
		{"save_fail", goodMsg, &mockRepo{saveErr: errors.New("serr")}, &mockEncoder{toCode: "xxx"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewURLService(makeConfig(), tt.repo, &mockProducer{}, &mockCache{}, tt.encoder)
			err := svc.HandleShortenURL(context.Background(), "q", tt.payload)
			if (err != nil) != tt.errOut {
				t.Fatalf("got err %v want error? %v", err, tt.errOut)
			}
		})
	}
}
