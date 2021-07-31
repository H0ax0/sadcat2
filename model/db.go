package model

import (
	"fmt"
	"time"

	"github.com/H0ax0/sadcat2/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Sql.Host,
		config.Sql.User,
		config.Sql.Password,
		config.Sql.DB,
		config.Sql.Port,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now()
		},
	})

	if err != nil {
		zap.S().Errorw("failed to open db", "error", err)
	}

	DB.AutoMigrate(&User{})
	Preload()
}
