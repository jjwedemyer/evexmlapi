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

// func (m *Model) XMLToJSON(xmlStr *[]byte) ([]byte, string, error) {
// 	v := m.DataFormat
// 	err := xml.Unmarshal(*xmlStr, &v)
// 	if err != nil {
// 		fmt.Println("Error unmarshalling from XML", err)
// 		return []byte{}, "", err
// 	}

// 	cachedUntil := v.(*SkillQueueFormat).CachedUntil
// 	result, err := json.Marshal(v)
// 	if err != nil {
// 		fmt.Println("Error marshalling to JSON", err)
// 		return []byte{}, "", err
// 	}
// 	return result, cachedUntil, nil
// }

func (m *Model) GetCachedUntil(v interface{}) string {
	return v.(*SkillQueueFormat).CachedUntil
}
