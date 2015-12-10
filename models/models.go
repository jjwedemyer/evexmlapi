package models

import (
	"encoding/json"
	"log"
)

type Model interface {
	Path() string
	// ToStruct([]byte) (interface{}, error)
	CacheDuration() int64
}

func JSON(m Model) ([]byte, error) {
	result, err := json.Marshal(m)
	if err != nil {
		log.Printf("Error marshalling %q to JSON: %q", m.Path(), err)
		return []byte{}, err
	}
	return result, nil
}
