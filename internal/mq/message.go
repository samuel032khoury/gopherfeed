package mq

import (
	"encoding/json"
	"fmt"
)

type EmailMessage struct {
	TemplateFile string          `json:"template_file"`
	Username     string          `json:"username"`
	Email        string          `json:"email"`
	Data         json.RawMessage `json:"data"`
}

func NewEmailMessage(templateFile, username, email string, data any) (*EmailMessage, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}
	return &EmailMessage{
		TemplateFile: templateFile,
		Username:     username,
		Email:        email,
		Data:         jsonData,
	}, nil
}

func (m *EmailMessage) ToBytes() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email message: %w", err)
	}
	return bytes, nil
}

func FromBytes(data []byte) (*EmailMessage, error) {
	var msg EmailMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal email message: %w", err)
	}
	return &msg, nil
}
