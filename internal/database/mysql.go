package database

import (
	"fmt"
	"time"

	"github.com/hacker4257/pet_charity/internal/config"
	"github.com/hacker4257/pet_charity/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL() error {
	cfg := config.Global.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	var logLevel logger.LogLevel
	if config.Global.Server.Mode == "debug" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("connect mysql failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB failed: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	DB = db
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Organization{},
		&model.Pet{},
		&model.PetImage{},
		&model.Adoption{},
		&model.Donation{},
		&model.Rescue{},
		&model.RescueFollow{},
		&model.RescueImage{},
		&model.RescueClaim{},
		&model.SmsCode{},
		&model.Message{},
		&model.UserActivityLogs{},
		&model.Notification{},
	)
}

func Close() {
	//关闭MySQL
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
	//关闭redis
	if RDB != nil {
		RDB.Close()
	}
}
