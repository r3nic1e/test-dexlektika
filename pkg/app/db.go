package app

import (
	"log"

	"github.com/r3nic1e/test-dexlektika/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (app *App) ConnectDB(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Printf("Failed to connect to DB: %v", err)
		return err
	}

	app.db = db.Debug()
	err = app.migrate()
	if err != nil {
		log.Printf("Failed to migrate schema: %v", err)
		return err
	}
	return nil
}

func (app *App) migrate() error {
	return app.db.AutoMigrate(
		models.DenylistedIP{},
	)
}
