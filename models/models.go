package models

import (
	"encoding/json"
	"encoding/xml"
	"log"
)

type Model struct {
	Path       string
	DataFormat interface{}
}

func (m *Model) FromXML(xmlStr *[]byte) (interface{}, error) {
	v := m.DataFormat
	err := xml.Unmarshal(*xmlStr, &v)
	if err != nil {
		log.Println("Error unmarshalling from XML", err)
		return nil, err
	}
	return v, nil
}

func (m *Model) ToJSON(v interface{}) ([]byte, error) {
	result, err := json.Marshal(v)
	if err != nil {
		log.Println("Error marshalling to JSON", err)
		return []byte{}, err
	}
	return result, nil
}
