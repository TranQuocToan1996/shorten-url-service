package dto

import (
	"encoding/json"
	"fmt"
)

type URLMessage struct {
	URL string `json:"url"`
}

func (m *URLMessage) Bytes() []byte {
	return fmt.Appendf(nil, "{\"url\": %s}", m.URL)
}

func (m *URLMessage) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}
