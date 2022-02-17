package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBConfig contains gorm dialector & config
type DBConfig struct {
	GormDialect gorm.Dialector
	GormConfig  gorm.Config
}

// GetSqliteConfig get the sqlite config
func GetSqliteConfig(dbName string) *DBConfig {
	return &DBConfig{
		GormDialect: sqlite.Open(dbName),
		GormConfig:  gorm.Config{},
	}
}
