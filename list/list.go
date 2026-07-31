package list

import "sync"

// перенесём в Entity
type StockScrapeInfo struct {
	StockTag    string
	AnomalySize uint64
}

var lock sync.Mutex
var stocksInfo = []StockScrapeInfo{
	{"LKOH", 80000000},
	{"SBER", 80000000},
	{"GAZP", 80000000},
	{"ROSN", 70000000},
	{"NVTK", 80000000},
	{"TCSG", 80000000},
}

// шиза на счёт гонки данных scrapper-ом и командой бота
func GetStocksInfoList() []StockScrapeInfo {
	lock.Lock()
	defer lock.Unlock()
	answer := make([]StockScrapeInfo, len(stocksInfo))
	copy(answer, stocksInfo)
	return answer
}
