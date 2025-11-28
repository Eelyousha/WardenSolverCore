package main

import (
	"encoding/json"
	"fmt"
)

// BlacklistEntry представляет запись в чёрном списке
type BlacklistEntry struct {
	SocialNetwork string `json:"social_network"`
	UserID        string `json:"user_id"`
	Nickname      string `json:"nickname"`
}

// GenerateBlacklistEntry создаёт новую запись для чёрного списка
func GenerateBlacklistEntry(socialNetwork, userID, nickname string) BlacklistEntry {
	return BlacklistEntry{
		SocialNetwork: socialNetwork,
		UserID:        userID,
		Nickname:      nickname,
	}
}

// ToJSON преобразует запись в JSON
func (e BlacklistEntry) ToJSON() string {
	data, _ := json.Marshal(e)
	return string(data)
}

// Validate проверяет валидность записи
func (e BlacklistEntry) Validate() error {
	if e.SocialNetwork == "" {
		return fmt.Errorf("social_network cannot be empty")
	}
	if e.UserID == "" {
		return fmt.Errorf("user_id cannot be empty")
	}
	if e.Nickname == "" {
		return fmt.Errorf("nickname cannot be empty")
	}
	return nil
}
