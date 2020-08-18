package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/shon-phand/CryptoServer/datasources/mysql/currency_db"
	"github.com/shon-phand/CryptoServer/datasources/redis"
	"github.com/shon-phand/CryptoServer/logger"
	"github.com/shon-phand/CryptoServer/utils/errors"
)

const (
	//queryInsertUser   = "INSERT INTO users (first_name,last_name,email,password,status,date_created) VALUES( ?,?,?,?,?,? )"
	queryGetUser    = "SELECT sym,id,ask,bid,last,open,low,high,feeCurrency FROM currency where sym=? "
	queryGetUserAll = "SELECT sym,id,ask,bid,last,open,low,high,feeCurrency FROM currency"
	//queryUpdateUser   = "UPDATE  currency SET first_name=?,last_name=?,email=? WHERE id=?;"
	//queryDeleteUser   = "DELETE from users WHERE id=?;"
	//queryFindByStatus = "SELECT id,first_name,last_name,email,status,date_created FROM users WHERE status=?;"
	expiration = 5
)

func (cr *Currency) Get(cur string) (*Currency, *errors.RestErr) {

	res, err := redis.Client.HGetAll("sym:" + cur).Result()
	if err != nil {
		return nil, errors.StatusInternalServerError("database error")
	}

	_, ok := res["id"]
	if !ok {
		return nil, errors.StatusNotFoundError("currency not found")
	}

	Response := Currency{}
	Response.ID = res["id"]
	Response.Ask = res["ask"]
	Response.Bid = res["bid"]
	Response.Last = res["last"]
	Response.Open = res["open"]
	Response.Low = res["low"]
	Response.High = res["high"]
	Response.FeeCurrency = res["feeCurrency"]
	SetExp(cur)
	return &Response, nil
}

func (cr *Currencies) GetAll() (*Currencies, *errors.RestErr) {

	hashes, err := redis.Client.Keys("sym:*").Result()

	if err != nil {
		return nil, errors.StatusInternalServerError("database error")
	}

	if len(hashes) == 0 {
		return nil, errors.StatusNotFoundError("no currencies found in the database, please sync all currencies first")
	}

	result := Currencies{}
	Response := Currency{}
	for _, v := range hashes {

		res, err := redis.Client.HGetAll(v).Result()
		if err != nil {
			return nil, errors.StatusInternalServerError("database error")
		}
		Response.ID = res["id"]
		Response.Ask = res["ask"]
		Response.Bid = res["bid"]
		Response.Last = res["last"]
		Response.Open = res["open"]
		Response.Low = res["low"]
		Response.High = res["high"]
		Response.FeeCurrency = res["feeCurrency"]
		result = append(result, Response)

	}

	return &result, nil
}

//////////////////// mysql DAO /////////////////////////////////

func (cr *Currency) GetCurrencyFromDB(cur string) (*Currency, *errors.RestErr) {
	stmt, err := currency_db.Client.Prepare(queryGetUser)
	if err != nil {
		return nil, errors.StatusInternalServerError("database error")
	}
	defer stmt.Close()
	result := stmt.QueryRow(cur)
	Response := Currency{}
	if err := result.Scan(&Response.Sym, &Response.ID, &Response.Ask, &Response.Bid, &Response.Last, &Response.Open, &Response.Low, &Response.High, &Response.FeeCurrency); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			logger.Info(errors.StatusNotFoundError("sym "+cur+" not found"), err)
			return nil, errors.StatusNotFoundError("sym " + cur + " not found")
		}
		logger.Error(errors.StatusInternalServerError("error in fetching data"), err)
		return nil, errors.StatusInternalServerError("database error")
	}
	Response.SaveToCache()
	return &Response, nil
}

func (cr Currencies) GetAllFromDB() (*Currencies, *errors.RestErr) {
	stmt, err := currency_db.Client.Prepare(queryGetUserAll)
	if err != nil {
		return nil, errors.StatusInternalServerError("database error")
	}
	defer stmt.Close()
	results := Currencies{}
	Response := Currency{}
	result, err := stmt.Query()
	if err != nil {
		return nil, errors.StatusInternalServerError("error in fetching all data")
	}
	for result.Next() {

		if err := result.Scan(&Response.Sym, &Response.ID, &Response.Ask, &Response.Bid, &Response.Last, &Response.Open, &Response.Low, &Response.High, &Response.FeeCurrency); err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				//logger.Info(errors.StatusNotFoundError("sym "+cur+" not found"), err)
				return nil, errors.StatusNotFoundError("no symbols: " + " not found")
			}
			logger.Error(errors.StatusInternalServerError("error in fetching all data"), err)
			return nil, errors.StatusInternalServerError("error in fetching all data")
		}

		results = append(results, Response)

	}
	results.SaveAllToCache()
	return &results, nil

}

////////////// cache dao /////////////////

func (crs *Currencies) SaveAllToCache() *errors.RestErr {
	// //results, err := cr.GetAllFromDB()
	// if err != nil {
	// 	fmt.Println("error in fetching from DB")
	// }
	for _, cr := range *crs {
		cr.SaveToCache()
	}

	return nil
}

func (cr *Currency) SaveToCache() *errors.RestErr {
	_, err := redis.Client.HSet("sym:"+cr.Sym, "id", cr.ID, "ask", cr.Ask, "bid", cr.Bid, "last", cr.Last, "open", cr.Open, "low", cr.Low, "high", cr.High, "feeCurrency", cr.FeeCurrency).Result()
	if err != nil {
		fmt.Println("error in saving in cache")
	}
	SetExp(cr.Sym)
	return nil
}

func SetExp(sym string) {

	err := redis.Client.Expire("sym:"+sym, expiration*time.Second).Err()
	if err != nil {
		fmt.Println("error in setting exp in cache")
	}

}
