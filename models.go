package evexmlapi

type model struct {
	format interface{}
}

// skillQueueRowFormat is a data format for each row
type skillQueueRowFormat struct {
	Position  string `xml:"queuePosition,attr" json:"queuePosition"`
	TypeID    string `xml:"typeID,attr" json:"typeID"`
	Level     string `xml:"level,attr" json:"level"`
	StartSP   string `xml:"startSP,attr" json:"startSP"`
	EndSP     string `xml:"endSP,attr" json:"endSP"`
	StartTime string `xml:"startTime,attr" json:"startTime"`
	EndTime   string `xml:"endTime,attr" json:"endTime"`
}

// SkillQueueFormat is the top level format for the data
type SkillQueueFormat struct {
	CurrentTime string                `xml:"currentTime" json:"currentTime"`
	Rows        []skillQueueRowFormat `xml:"result>rowset>row" json:"rows"`
	CachedUntil string                `xml:"cachedUntil" json:"cachedUntil"`
}

// ServerStatusFormat is the top level format for the data
type ServerStatusFormat struct {
	CurrentTime   string `xml:"currentTime" json:"currentTime"`
	ServerOpen    bool   `xml:"result>serverOpen" json:"serverOpen"`
	OnlinePlayers int    `xml:"result>onlinePlayers" json:"onlinePlayers"`
	CachedUntil   string `xml:"cachedUntil" json:"cachedUntil"`
}

// APIKeyInfoFormat is the top level format for the data
type APIKeyInfoFormat struct {
	Keys        []KeysRowFormat `xml:"result>key" json:"keys"`
	CachedUntil string          `xml:"cachedUntil" json:"cachedUntil"`
	CurrentTime string          `xml:"currentTime" json:"currentTime"`
}

type KeysRowFormat struct {
	Characters []CharactersRowFormat `xml:"rowset>row" json:"characters"`
	AccessMask string                `xml:"accessMask,attr" json:"accessMask"`
	Type       string                `xml:"type,attr" json:"type"`
	Expires    string                `xml:"expires,attr" json:"expires"`
}
type CharactersRowFormat struct {
	CharacterID     string `xml:"characterID,attr" json:"characterID"`
	CharacterName   string `xml:"characterName,attr" json:"characterName"`
	CorporationID   string `xml:"corporationID,attr" json:"corporationID"`
	CorporationName string `xml:"corporationName,attr" json:"corporationName"`
	AllianceID      string `xml:"allianceID,attr" json:"allianceID"`
	AllianceName    string `xml:"allianceName,attr" json:"allianceName"`
	FactionID       string `xml:"faction,attr" json:"factionID"`
	FactionName     string `xml:"factionName,attr" json:"factionName"`
}
