package types


type SkillQueueRowXML struct {		
		Position    string  `xml:"queuePosition,attr" json:"queuePosition"`
		TypeID		string 	`xml:"typeID,attr" json:"typeID"`
		Level 		string 	`xml:"level,attr" json:"level"`
		StartSP 	string 	`xml:"startSP,attr" json:"startSP"`
		EndSP 		string 	`xml:"endSP,attr" json:"endSP"`
		StartTime	string	`xml:"startTime,attr" json:"startTime"`
		EndTime		string	`xml:"endTime,attr" json:"endTime"`
	}
	
		
type SkillQueueXML struct {		
		CurrentTime string				`xml:"currentTime" json:"currentTime"`
		Rows 		[]SkillQueueRowXML 	`xml:"result>rowset>row" json:"rows"`
		CachedUntil string				`xml:"cachedUntil"	json:"cachedUntil"`
	}
