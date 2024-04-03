package repository
import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"mywon/students_reports/constants"
	"os"
	"strconv"
	"time"
)

type SQLConnDetails struct {
	Pool     *sql.DB
	PgSchema string
}

var Pool *sql.DB

func InitDbPool() {
	var port, user, password, dbname, sslmode string
	host := os.Getenv(constants.POSTGRES_HOST)
	port = os.Getenv(constants.POSTGRES_PORT)
	user = os.Getenv(constants.POSTGRES_USER)
	password = os.Getenv(constants.POSTGRES_PASSWORD)
	dbname = os.Getenv(constants.POSTGRES_DB_NAME)
	sslmode = os.Getenv(constants.POSTGRES_SSLMODE)

	databaseURL := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=" + sslmode

	maxOpenConnection, err := strconv.Atoi(os.Getenv(constants.POSTGRES_MAX_CONN))
	if err != nil {
		log.Println(err)
		maxOpenConnection = 5
	}
	maxIdleTime, err := strconv.Atoi(os.Getenv(constants.POSTGRES_MAX_IDLE_TIME))
	if err != nil {
		log.Println(err)
		maxIdleTime = 5
	}
	maxConnectionLifetime, err := strconv.Atoi(os.Getenv(constants.POSTGRES_MAX_LIFETIME))
	if err != nil {
		log.Println(err)
		maxConnectionLifetime = 2
	}

	config, err := sql.Open(constants.DRIVER_NAME, databaseURL)
	fmt.Println("$$$$")
	if err != nil {
		fmt.Println("$$$$123")
		log.Print(err)
		log.Print(constants.CONFIG_ERR)
	}
	config.SetMaxIdleConns(maxOpenConnection)
	config.SetConnMaxLifetime(time.Duration(maxConnectionLifetime) * time.Minute)
	config.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Minute)

	err = config.Ping()
	if err != nil {
		//log.Print(err)
		log.Print(constants.POSTGRES_NOT_CONNECTED)
	} else {
		log.Println(constants.POSTGRES_CONNECTED)
	}

	Pool = config

}

func GetPool() *sql.DB {
	if Pool == nil {
		InitDbPool()
	}

	return Pool
}