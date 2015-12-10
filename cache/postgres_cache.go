package cache

import (
	"github.com/jovon/eve-xmlapi-go/db"
	_ "github.com/lib/pq"
)

type PostgresCache struct {
	db *data.DB
	table string
}

func (pc PostgresCache) New(db *data.DB, table string) *PostgresCache {
	return &PostgresCache{db: db, table: table}
}

// MergeCache runs a postgres function that tries to insert the record.
// On success it then deletes records with the same keyid, characterid, and apipath with a cachedUntil date
// 	less than the provided cachedUntil date.
// If the record exists it ignores the unique validation error and does nothing.
func (pc PostgresCache) Write(key string, value []byte, duration int64) int64 {
	var lastID int64
	// queryStr := `INSERT INTO ` + pc.table + ` 
	// err := pc.db.QueryRow(`Select insert_delete_cache($1::integer, $2::varchar(25), $3::varchar(50), $4::jsonb, $5::timestamp)`,
	// 	hr.Param("keyID"), hr.Param("characterID"), model.Path(), string(r), hr.CachedUntil()).Scan(&lastId)

	// if err != nil {
	// 	log.Fatal("Fatal merge cache: ", err)
	// }

	return lastID
}

func (pc PostgresCache) Read(key string) []byte {
	// var lastId int64 = 0
	// err := db.QueryRow(`Select id from cache Where keyid = $1 and characterid = $2 and apipath = $3 and cachedUntil = $4`,
	// 	hr.Param("keyID"), hr.Param("characterID"), model.Path(), hr.CachedUntil()).Scan(&lastId)
	return []byte{}
}