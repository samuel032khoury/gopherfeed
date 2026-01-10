package email

import (
	"encoding/json"
	"fmt"
)

type Email struct {
	To           string          `json:"to"`
	TemplatePath string          `json:"template_path"`
	Data         json.RawMessage `json:"data"`
}

func New(to string, templatePath string, data any) (*Email, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}
	return &Email{
		To:           to,
		TemplatePath: templatePath,
		Data:         jsonData,
	}, nil
}

func (m *Email) ToBytes() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email message: %w", err)
	}
	return bytes, nil
}

func FromBytes(data []byte) (*Email, error) {
	var msg Email
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal email message: %w", err)
	}
	return &msg, nil
}

type Sender interface {
	Send(to string, templatePath string, data any) error
}
