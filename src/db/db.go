package db

import (
	"fmt"

	"app/models"

	"gitlab.online-fx.com/go-packages/gormdb"
	"gitlab.online-fx.com/go-packages/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	config *models.Config
)

func Init(configData *models.Config) (*gorm.DB, error) {
	config = configData

	//db, err := gormdb.AddPostgres(models.ServiceDB, config.Database)
	dataSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		config.Database.Address, config.Database.User, config.Database.Password, config.Database.DB, config.Database.Port, config.Timezone)
	logger.Infof("DB Config: Host=%s, Database=%s, TimeZone=%s", config.Database.Address, config.Database.DB, config.Timezone)
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{})

	gormdb.AddClient(models.ServiceDB, db)
	logger.Infof("Postgres initialized (db: %s)", config.Database.DB)

	if err != nil {
		return nil, err
	}

	if configData.LogLevel == models.LevelDebug {
		db = db.Debug()
	}

	return db, nil
}

func RunMigrations() {
	gormdb.ApplyMigrationsForClient(models.ServiceDB, models.GetModels()...)
}
