package services

import (
	"github.com/shon-phand/CryptoServer/datasources/mysql/currency_db"
	"github.com/shon-phand/CryptoServer/logger"
	"github.com/shon-phand/CryptoServer/sync"
	"github.com/shon-phand/CryptoServer/utils/errors"
)

const (
	queryInsertUser = "INSERT INTO currency (sym,id,ask,bid,last,open,low,high,feeCurrency) VALUES( ?,?,?,?,?,?,?,?,? )"
)

func Save(curr sync.Currencies) {

	for _, cr := range curr {
		// _, err := redis.Client.HSet("sym:"+cr.Symbol, "id", cr.ID, "ask", cr.Ask, "bid", cr.Bid, "last", cr.Last, "open", cr.Open, "low", cr.Low, "high", cr.High, "feeCurrency", cr.FeeCurrency).Result()
		// if err != nil {
		// 	fmt.Println("error in saving in redis")
		// }

		stmt, err := currency_db.Client.Prepare(queryInsertUser)
		if err != nil {
			logger.Error(errors.StatusInternalServerError("error in preapre stmt"), err)
			//return errors.StatusInternalServerError(" database error")
		}
		defer stmt.Close()
		_, err = stmt.Exec(&cr.Symbol, &cr.ID, &cr.Ask, &cr.Bid, &cr.Last, &cr.Open, &cr.Low, &cr.High, &cr.FeeCurrency)
		//fmt.Println("Currnecy inserted: ", cr.Symbol)
		if err != nil {
			logger.Error(errors.StatusInternalServerError("error while saving currency"), err)
		}
	}

}
