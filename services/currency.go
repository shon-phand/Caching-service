package services

import (
	"github.com/shon-phand/CryptoServer/domain"
	"github.com/shon-phand/CryptoServer/utils/errors"
)

type currencyService struct{}

var (
	CurrencyService currencyServiceInterface = &currencyService{}
)

type currencyServiceInterface interface {
	GetCurrency(string) (*domain.Currency, *errors.RestErr)
	GetAllCurrency() (*domain.Currencies, *errors.RestErr)
}

func (cs *currencyService) GetCurrency(curr string) (*domain.Currency, *errors.RestErr) {

	obj := &domain.Currency{}
	// fetching from cache
	res, err := obj.Get(curr)

	if err != nil {

		// not found in cache
		res, err = obj.GetCurrencyFromDB(curr)
		if err != nil {
			return nil, err
		}
		//go res.SaveToCache()
	}

	return res, nil

}

func (cs *currencyService) GetAllCurrency() (*domain.Currencies, *errors.RestErr) {

	obj := &domain.Currencies{}
	//res, err := obj.GetAll()
	//if err != nil {
	res, err := obj.GetAllFromDB()
	if err != nil {
		return nil, err
	}
	//}

	return res, nil

}
