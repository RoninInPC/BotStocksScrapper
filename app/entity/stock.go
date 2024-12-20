package entity

type StockAdd struct {
	StockName string
	Type      string // SALE или BUY
	NumPrice  int64
}
