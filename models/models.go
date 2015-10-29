package models

import (
	"encoding/json"
	"log"
)

type Model interface {
	Path() string
	GetCachedUntil(interface{}) string
	FromXML(*[]byte) (interface{}, error)
}

func ToJSON(v Model) ([]byte, error) {
	result, err := json.Marshal(v)
	if err != nil {
		log.Println("Error marshalling to JSON", err)
		return []byte{}, err
	}
	return result, nil
}
