package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	successfulWrites = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "postgres_rows_written_total",
			Help: "Number of successful writes",
		},
		[]string{"dbname"},
	)
	failureWrites = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "postgres_write_errors_total",
			Help: "Number of failed writes",
		},
		[]string{"dbname"},
	)
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func init() {
	prometheus.MustRegister(successfulWrites)
	prometheus.MustRegister(failureWrites)
}

func main() {
	dbConfig := initConfig()

	db := connectToDatabase(dbConfig)
	defer db.Close()

	bootstrapDb(db)

	startMetricsServer()

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for range ticker.C {
		if rand.Intn(100) < 10 {
			simulateDatabaseError(dbConfig.Name)
		} else {
			addNewRow(db)
		}
	}
}

func initConfig() *DBConfig {
	return &DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "devuser"),
		Password: getEnv("DB_PASSWORD", "devpass"),
		Name:     getEnv("DB_NAME", "devdb"),
	}
}

func connectToDatabase(config *DBConfig) *sql.DB {
	connString := getConnString(config)
	db, err := sql.Open("pgx", connString)

	if err != nil {
		log.Fatal("Failed to connect to the database with error: ", err)
	}

	log.Println("Waiting for database...")

	for i := range 30 {
		if err := db.Ping(); err == nil {
			log.Println("Database connection established!")
			return db
		}

		log.Printf("Database is not yet ready, retrying (%d/30)", i)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Database connection failed after 30 attempts. Check: host=%s port=%s dbname=%s user=%s",
		config.Host, config.Port, config.Name, config.User)

	return nil
}

func bootstrapDb(db *sql.DB) error {
	newTableSql := `
		CREATE TABLE IF NOT EXISTS devdb (
			id SERIAL PRIMARY KEY,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			data TEXT DEFAULT ''
		)
	`

	if _, err := db.Exec(newTableSql); err != nil {
		log.Fatal("Failed to create table with error: ", err)
	}

	log.Println("Database is ready!")
	return nil
}

func addNewRow(db *sql.DB) {
	insertRowSql := `
		INSERT INTO devdb (data) VALUES ($1)
	`

	log.Println("Adding row")
	randomData := fmt.Sprintf("%d data", time.Now().UnixNano())

	_, err := db.Exec(insertRowSql, randomData)

	if err != nil {
		log.Println("Failed to add new row:", err)
		failureWrites.WithLabelValues("devdb").Inc()
	} else {
		log.Println("New row added.")
		successfulWrites.WithLabelValues("devdb").Inc()
	}
}

func simulateDatabaseError(dbname string) {
	errors := []string{
		"database connection failed",
		"database timeout",
		"database unavailable",
		"database error",
	}

	error := errors[rand.Intn(len(errors))]
	log.Println("Simulating database error:", error)
	failureWrites.WithLabelValues(dbname).Inc()
}

func startMetricsServer() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Metrics server is ready on port :8080.")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}

func getConnString(config *DBConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)
}

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)

	if value == "" {
		value = defaultValue
	}

	return value
}
