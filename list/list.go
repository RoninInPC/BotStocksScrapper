package list

import "sync"

type StockScrapeInfo struct {
	StockTag    string
	AnomalySize uint64
	Figi        string
}

var lock sync.Mutex
var stocksInfo = []StockScrapeInfo{
	{"LKOH", 80000000, "BBG004731032"},
	{"SBER", 80000000, "BBG004730N88"},
	{"GAZP", 80000000, "BBG004730RP0"},
	{"ROSN", 70000000, "BBG004731354"},
	{"NVTK", 80000000, "BBG00475KKY8"},
	{"TCSG", 80000000, "TCS00A107UL4"},
}

// шиза на счёт гонки данных scrapper-ом и командой бота
func GetStocksInfoList() []StockScrapeInfo {
	lock.Lock()
	defer lock.Unlock()
	answer := make([]StockScrapeInfo, len(stocksInfo))
	copy(answer, stocksInfo)
	return answer
}
