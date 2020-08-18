package currency_db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shon-phand/CryptoServer/logger"
	"github.com/shon-phand/CryptoServer/utils/errors"
)

var (
	Client *sql.DB
)

func init() {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		"root",
		"mysqlpass",
		"mysql:3306",
		"currency_db",
	)
	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		logger.Error(errors.StatusInternalServerError("error in starting database"), err)
		panic(err)
	}
	err = Client.Ping()
	if err != nil {
		logger.Error(errors.StatusInternalServerError("error in pinging database"), err)
		panic(err)
	}

	log.Println("Database successully connected")
}
