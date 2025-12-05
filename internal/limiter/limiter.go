package limiter

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var period int = 10
var limit int = 5
var db *sql.DB

func StartLimiter() {
	err := LoadEnvVariables()

	if err != nil {
		log.Printf("Failed to load enviroment variables: %w", err)
	}

	err = InitDB()

	if err != nil {
		log.Printf("Error connecting to DB: %w", err)
	}

	err = CreateTables()

	if err != nil {
		log.Printf("Error creating tables: %w", err)
	}
}

func LoadEnvVariables() error {
	var err error

	limit, err = strconv.Atoi(os.Getenv("RATE_LIMIT"))

	if err != nil {
		return err
	}

	period, err = strconv.Atoi(os.Getenv("RATE_PERIOD"))

	return err

}

func CreateTables() error {
	var err error

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS rate_limits(ip TEXT NOT NULL, ts TIMESTAMPTZ NOT NULL DEFAULT NOW());")

	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_rate_limits ON rate_limits(ip, ts);")

	return err
}

func InitDB() error {
	dbUser, dbPass, dbHost, dbPort, dbName := os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	log.Printf("Connecting with: %s", dataSource)

	var err error
	db, err = sql.Open("postgres", dataSource)

	if err != nil {
		return err
	}

	return db.Ping()
}

func CloseDB() error {
	err := db.Close()

	if err != nil {
		return err
	}

	return nil
}

func GetPeriod() int {
	return period
}

func CleanLogs(key string) error {
	_, err := db.Exec("DELETE FROM rate_limits WHERE ip = $1 AND ts < NOW() - ($2 || ' seconds')::interval", key, period)

	return err
}

func AddLog(key string) error {
	_, err := db.Exec("INSERT INTO rate_limits(ip) VALUES ($1)", key)

	return err
}

func CanAccess(ip string) (bool, error) {
	CleanLogs(ip)

	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM rate_limits WHERE ip=$1", ip).Scan(&count)

	if err != nil {
		return false, err
	}

	return count < limit, nil
}

func GetOldestLog(ip string) (time.Time, error) {
	if err := CleanLogs(ip); err != nil {
		return time.Time{}, err
	}

	rows, err := db.Query("SELECT ts FROM rate_limits WHERE ip=$1", ip)

	if err != nil {
		return time.Time{}, err
	}

	var oldest time.Time

	for rows.Next() {
		var ts time.Time

		err := rows.Scan(&ts)

		if err != nil {
			return time.Time{}, err
		}

		if oldest.IsZero() || ts.Before(oldest) {
			oldest = ts
		}
	}

	return oldest, nil
}
