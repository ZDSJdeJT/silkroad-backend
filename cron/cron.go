package cron

import (
	"github.com/robfig/cron/v3"
	"log"
	"silkroad-backend/database"
	"silkroad-backend/utils"
)

func Start() {
	c := cron.New()
	_, err := c.AddFunc("0 1 * * *", database.DeleteExpiredTextRecords)
	_, err = c.AddFunc("0 2 * * *", database.DeleteExpiredFileRecords)
	_, err = c.AddFunc("0 3 * * *", utils.DeleteOldChunks)
	if err != nil {
		log.Printf("Error adding cron function: %s.", err)
		return
	}
	c.Start()
}
