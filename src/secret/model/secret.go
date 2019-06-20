package model

import (
	"encoding/json"
	"time"
)

// Secret TODO: add comment
type Secret struct {
	Hash           string    `json:"hash" xml:"hash"`
	SecretText     string    `json:"secretText" xml:"secretText"`
	CreatedAt      time.Time `json:"createdAt" xml:"createdAt"`
	ExpiresAt      time.Time `json:"expiresAt" xml:"expiresAt"`
	RemainingViews int       `json:"remainingViews" xml:"remainingViews"`
}

// MarshalBinary -
func (s Secret) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary -
func (s *Secret) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return nil
}
