package apiDriver

import (
	"scrapper-bot/entity"
)

type ApiDriver struct {
	token string
}

func NewApiDriver(token string) *ApiDriver {
	return &ApiDriver{token: token}
}

func (d *ApiDriver) GetStockInfo(tag string) *entity.StockInfo {
	panic("implement me")
}
