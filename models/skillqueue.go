package models

import (
	"encoding/xml"
	"log"
)

type SkillQueue struct {
	path          string
	cacheDuration int64
}

func (m SkillQueue) New() SkillQueue {
	return SkillQueue{
		path:          "char/SkillQueue.xml.aspx",
		cacheDuration: 3600,
	}
}

type SkillQueueRowFormat struct {
	Position  string `xml:"queuePosition,attr" json:"queuePosition"`
	TypeID    string `xml:"typeID,attr" json:"typeID"`
	Level     string `xml:"level,attr" json:"level"`
	StartSP   string `xml:"startSP,attr" json:"startSP"`
	EndSP     string `xml:"endSP,attr" json:"endSP"`
	StartTime string `xml:"startTime,attr" json:"startTime"`
	EndTime   string `xml:"endTime,attr" json:"endTime"`
}

type SkillQueueFormat struct {
	CurrentTime string                `xml:"currentTime" json:"currentTime"`
	Rows        []SkillQueueRowFormat `xml:"result>rowset>row" json:"rows"`
	CachedUntil string                `xml:"cachedUntil" json:"cachedUntil"`
}


func (m SkillQueue) Path() string {
	return m.path
}

func (m SkillQueue) CacheDuration() int64 {
	return m.cacheDuration
}

func (m SkillQueue) ToStruct(xmlStr []byte) (*SkillQueueFormat, error) {
	if len(xmlStr) != 0 {
		v := new(SkillQueueFormat)
		err := xml.Unmarshal(xmlStr, v)
		if err != nil {		
			log.Fatalf("Error(%q) unmarshalling %q from XML: %q", err, m.Path(), xmlStr)		
			return nil, err
		}
		return v, nil
	}
	return nil, nil
}
