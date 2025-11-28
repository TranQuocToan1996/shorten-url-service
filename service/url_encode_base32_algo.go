package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"math/big"
	"strconv"

	"shorten/pkg/config"
)

type URLEncoder interface {
	Encode(originalURL string) (string, error)
}

// Refactor to use interface and config
type base62Encoder struct {
	secretKey      string
	shortURLLength int
}

const defaultShortURLLength = 8

func NewBase62Encoder(cfg config.Config) URLEncoder {
	shortURLLength, err := strconv.Atoi(cfg.SHORT_URL_LENGTH)
	if err != nil {
		shortURLLength = defaultShortURLLength
	}
	return &base62Encoder{
		secretKey:      cfg.SECRET_KEY,
		shortURLLength: shortURLLength,
	}
}

// Base62 characters
const base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// ShortenURL generates deterministic, secure short code
func (e *base62Encoder) Encode(originalURL string) (string, error) {
	// Create HMAC hash
	h := hmac.New(sha256.New, []byte(e.secretKey))
	h.Write([]byte(originalURL))
	hash := h.Sum(nil)

	// Encode first bytes of hash to Base62
	shortCode := encodeBase62(hash[:6], e.shortURLLength) // 6 bytes → 8 chars
	return shortCode, nil
}

// Encode a byte slice into Base62 string
func encodeBase62(data []byte, length int) string {
	num := new(big.Int).SetBytes(data)
	result := make([]byte, 0, length)
	base := big.NewInt(62)
	for len(result) < length {
		mod := new(big.Int)
		num.DivMod(num, base, mod)
		result = append(result, base62[mod.Int64()])
	}
	// Reverse to get correct order
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}
