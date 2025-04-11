package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func createDBFiles(projectName, db string) {
	// config.go
	configPath := filepath.Join(projectName, "internal", "config", "config.go")
	configCode := fmt.Sprintf(`package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func LoadConfig() *Config {
	return &Config{
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "%s"),
		DBName:     getEnv("DB_NAME", "mydb"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func (c *Config) GetDSN(dbType string) string {
	switch dbType {
	case "postgres":
		return fmt.Sprintf("postgres://%%s:%%s@%%s:%%s/%%s?sslmode=disable", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
	case "mysql":
		return fmt.Sprintf("%%s:%%s@tcp(%%s:%%s)/%%s", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
	case "mongo":
		return fmt.Sprintf("mongodb://%%s:%%s@%%s:%%s", c.DBUser, c.DBPassword, c.DBHost, c.DBPort)
	default:
		return ""
	}
}
`, map[string]string{"postgres": "5432", "mysql": "3306", "mongo": "27017"}[db])

	os.WriteFile(configPath, []byte(configCode), 0644)

	// db.go
	dbPath := filepath.Join(projectName, "pkg", "db", "db.go")
	var dbCode string

	if db == "mongo" {
		dbCode = fmt.Sprintf(`package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"%s/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB(cfg *config.Config, dbType string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.GetDSN(dbType)))
	if err != nil {
		log.Fatalf("Failed to create Mongo client: %%v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %%v", err)
	}

	fmt.Println("Connected to MongoDB successfully.")
	return client
}
`, projectName)
	} else {
		driver := map[string]string{
			"postgres": "_ \"github.com/lib/pq\"",
			"mysql":    "_ \"github.com/go-sql-driver/mysql\"",
		}[db]

		dbCode = fmt.Sprintf(`package db

import (
	"database/sql"
	"fmt"
	"log"

	"%s/internal/config"
	%s
)

func InitDB(cfg *config.Config, dbType string) *sql.DB {
	dsn := cfg.GetDSN(dbType)
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to %%s database: %%v", dbType, err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping %%s database: %%v", dbType, err)
	}

	fmt.Println("Database connection established successfully.")
	return db
}
`, projectName, driver)
		dbCode = strings.ReplaceAll(dbCode, "{{.ProjectName}}", projectName)
	}

	os.WriteFile(dbPath, []byte(dbCode), 0644)
}
