package models

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

func SkillQueue() Model {
	return Model{
		Path:       "char/SkillQueue.xml.aspx",
		DataFormat: new(SkillQueueFormat),
	}
}
