package database

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"silkroad-backend/models"
	"time"
)

func DeleteExpiredTextRecords() {
	db, err := OpenDBConnection()
	if err != nil {
		log.Printf("Error opening database connection: %s.", err)
		return
	}
	result := db.Where("expire_at < ? AND is_file = false", time.Now()).Delete(models.Record{})
	if result.Error != nil {
		log.Printf("Error deleteing text records: %s.", err)
	} else {
		log.Printf("%d outdated text records have been deleted", result.RowsAffected)
	}
}

func DeleteExpiredFileRecords() {
	db, err := OpenDBConnection()
	if err != nil {
		log.Printf("Error opening database connection: %s.", err)
		return
	}
	sql := "DELETE FROM records WHERE expire_at < ? AND is_file = true RETURNING id"
	rows, err := db.Raw(sql, time.Now()).Rows()
	if err != nil {
		log.Printf("Error deleteing file records: %s.", err)
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
		err := os.RemoveAll("./data/files/" + id.String())
		if err != nil {
			fmt.Println("Failed to delete directory:", err)
			return
		}
	}
	log.Printf("Successfully deleted outdated file records")
}
