package types

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const (
	DbConnString string = "user=postgres password=postgres dbname=SQO sslmode=disable"
	DbDialect    string = "postgres"
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

type HttpRequest struct {
	Path       string
	UserAgent  string
	BaseURL    string
	Params     map[string]string
	DataFormat interface{}
}

type Cache interface {
	Read(key Key) (string, error)
	Write(key Key, data []byte)
}

type Key struct {
	CharacterId string
	KeyId       string
	Path        string
}

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	db, err := sql.Open(DbDialect, DbConnString)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Read(key Key) (string, error) {
	return "", nil
}

func (db *DB) Write(key Key, data []byte) error {
	return nil
}
