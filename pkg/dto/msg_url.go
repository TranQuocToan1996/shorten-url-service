package dto

import (
	"encoding/json"
)

type URLMessage struct {
	URL         string `json:"url"`
	CallbackURL string `json:"callback_url,omitempty"`
}

func (m *URLMessage) Bytes() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *URLMessage) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}
