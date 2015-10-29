package models

import (
	"encoding/xml"
	"log"
)

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

func NewSkillQueue() SkillQueue {
	return SkillQueue{
		path: "char/SkillQueue.xml.aspx",
	}
}

type SkillQueue struct {
	path string
}

func (m SkillQueue) GetCachedUntil(v interface{}) string {
	return v.(*SkillQueueFormat).CachedUntil
}

func (m SkillQueue) Path() string {
	return m.path
}

func (m SkillQueue) FromXML(xmlStr *[]byte) (interface{}, error) {
	v := new(SkillQueueFormat)
	err := xml.Unmarshal(*xmlStr, &v)
	if err != nil {
		log.Println("Error unmarshalling from XML", err)
		return nil, err
	}
	return v, nil
}
