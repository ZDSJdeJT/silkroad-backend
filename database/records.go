package database

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

const (
	UploadsDir   = "./data/uploads/"
	ChunksDir    = UploadsDir + "chunks/"
	DataDir      = UploadsDir + "data/"
	TextFilename = "data.txt"
)

func DeleteExpiredRecords(isFile bool) {
	db, err := OpenDBConnection()
	if err != nil {
		log.Printf("Error opening database connection: %s.", err)
		return
	}
	var sql string
	if isFile {
		sql = "DELETE FROM records WHERE expire_at < ? AND filename != \"\" RETURNING id"
	} else {
		sql = "DELETE FROM records WHERE expire_at < ? AND filename = \"\" RETURNING id"
	}
	rows, err := db.Raw(sql, time.Now()).Rows()
	if err != nil {
		log.Printf("Error deleteing records: %s.", err)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
			return
		}
	}()

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			log.Printf("Something wrong: %s.", err)
			return
		}
		err := os.RemoveAll(DataDir + id.String())
		if err != nil {
			fmt.Println("Failed to delete directory:", err)
			return
		}
	}
}

func DeleteExpiredTextRecords() {
	DeleteExpiredRecords(false)
	log.Printf("Successfully deleted expired text records")
}

func DeleteExpiredFileRecords() {
	DeleteExpiredRecords(true)
	log.Printf("Successfully deleted expired file records")
}
