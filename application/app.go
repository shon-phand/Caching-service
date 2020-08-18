package application

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shon-phand/CryptoServer/controllers"
	"github.com/shon-phand/CryptoServer/datasources/mysql/currency_db"
	"github.com/shon-phand/CryptoServer/logger"
	"github.com/shon-phand/CryptoServer/utils/errors"
)

var (
	r = gin.New()
)

func init() {
	//fmt.Println("updating database, it will take 18 seconds to sync all data")
	//go services.UpdateDatabase()
	// if err != nil {
	// 	panic("error in synching data")
	// }

	// fmt.Println("update complete")
	err := currency_db.Client.Ping()
	if err != nil {
		logger.Error(errors.StatusInternalServerError("error in pinging database"), err)
		panic(err)
	}
	log.Printf("connected to DB")

}

func StartApplication() {
	//gin.SetMode(gin.ReleaseMode)
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/currency/", controllers.GetAllCurrency())
	r.GET("/currency/:symbol", controllers.GetCurrency())
	r.GET("/internal/currency/UpdateDb", controllers.UpdateDB())
	//go services.UpdateDatabase()
	gin.SetMode(gin.ReleaseMode)
	server1 := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 20 * time.Second,
		ErrorLog:     nil,
	}

	if err := server1.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
