// database/database.go - 데이터베이스 연결 및 유틸리티
// GORM 기반 데이터베이스 유틸리티
package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 데이터베이스 설정
type Config struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	Charset         string
	ParseTime       bool
	Loc             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	LogLevel        logger.LogLevel
}

// DefaultConfig 기본 데이터베이스 설정
func DefaultConfig() *Config {
	return &Config{
		Host:            "localhost",
		Port:            3306,
		Username:        "root",
		Password:        "",
		Database:        "test",
		Charset:         "utf8mb4",
		ParseTime:       true,
		Loc:             "Local",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		LogLevel:        logger.Info,
	}
}

// DSN 데이터베이스 연결 문자열 생성
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset, c.ParseTime, c.Loc)
}

// Connect 데이터베이스 연결
func Connect(config *Config) (*gorm.DB, error) {
	if config == nil {
		config = DefaultConfig()
	}

	dsn := config.DSN()
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 연결 풀 설정
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return db, nil
}

// ConnectWithFallback 데이터베이스 연결 (실패해도 서버 구동 가능)
func ConnectWithFallback(config *Config) *gorm.DB {
	db, err := Connect(config)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Server will continue without database connection")
		return nil
	}

	log.Println("Successfully connected to database")
	return db
}

// HealthCheck 데이터베이스 연결 상태 확인
func HealthCheck(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// IsConnected 데이터베이스 연결 여부 확인
func IsConnected(db *gorm.DB) bool {
	return HealthCheck(db) == nil
}

// Paginate GORM 페이징 스코프
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 20
		}
		if pageSize > 100 {
			pageSize = 100
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// Transaction 트랜잭션 실행
func Transaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	return db.Transaction(fn)
}

// SafeTransaction 안전한 트랜잭션 실행 (데이터베이스 연결이 없어도 에러 없이 처리)
func SafeTransaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	if db == nil {
		log.Println("Warning: Database connection is nil, skipping transaction")
		return nil
	}

	return Transaction(db, fn)
}