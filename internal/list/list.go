package list

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type StockScrapeInfo struct {
	StockTag          string
	AnomalySizeVolume float64
	AnomalySizeSolo   float64
	Figi              string
}

const (
	filename = "/etc/project/save.json"
	mln      = 1000000
)

var lock sync.Mutex
var stocksInfo = []StockScrapeInfo{
	{"AKRN", 60 * mln, 50 * mln, ""},
	{"AFKS", 250 * mln, 100 * mln, ""},
	{"AFLT", 300 * mln, 100 * mln, ""},
	{"ALRS", 180 * mln, 100 * mln, ""},
	{"AMEZ", 20 * mln, 20 * mln, ""},
	{"APTK", 10 * mln, 20 * mln, ""},
	{"AQUA", 80 * mln, 60 * mln, ""},
	{"ASTR", 100 * mln, 60 * mln, ""},
	{"BANE", 60 * mln, 80 * mln, ""},
	{"BANEP", 100 * mln, 80 * mln, ""},
	{"BELU", 100 * mln, 50 * mln, ""},
	{"BSPB", 140 * mln, 100 * mln, ""},
	{"CBOM", 110 * mln, 60 * mln, ""},
	{"CHMF", 220 * mln, 100 * mln, ""},
	{"CNRU", 60 * mln, 50 * mln, ""},
	{"DOMRF", 170 * mln, 100 * mln, ""},
	{"ELMT", 100 * mln, 60 * mln, ""},
	{"EUTR", 100 * mln, 60 * mln, ""},
	{"ETLN", 100 * mln, 60 * mln, ""},
	{"ENPG", 80 * mln, 50 * mln, ""},
	{"FEES", 100 * mln, 60 * mln, ""},
	{"FESH", 100 * mln, 80 * mln, ""},
	{"FIXR", 60 * mln, 50 * mln, ""},
	{"FLOT", 200 * mln, 100 * mln, ""},
	{"GAZP", 850 * mln, 120 * mln, ""},
	{"GCHE", 50 * mln, 50 * mln, ""},
	{"GMKN", 350 * mln, 110 * mln, ""},
	{"GTRK", 70 * mln, 50 * mln, ""},
	{"GEMC", 50 * mln, 50 * mln, ""},
	{"HEAD", 100 * mln, 70 * mln, ""},
	{"HNFG", 50 * mln, 50 * mln, ""},
	{"HYDR", 100 * mln, 60 * mln, ""},
	{"IRAO", 100 * mln, 60 * mln, ""},
	{"IRKT", 100 * mln, 50 * mln, ""},
	{"LEAS", 70 * mln, 50 * mln, ""},
	{"LENT", 80 * mln, 60 * mln, ""},
	{"LKOH", 800 * mln, 120 * mln, ""},
	{"LSNG", 60 * mln, 50 * mln, ""},
	{"LSNGP", 60 * mln, 50 * mln, ""},
	{"LNZL", 15 * mln, 30 * mln, ""},
	{"LSRG", 80 * mln, 60 * mln, ""},
	{"MAGN", 250 * mln, 100 * mln, ""},
	{"MBNK", 80 * mln, 60 * mln, ""},
	{"MDMG", 80 * mln, 60 * mln, ""},
	{"MGNT", 450 * mln, 110 * mln, ""},
	{"MOEX", 300 * mln, 100 * mln, ""},
	{"MRKC", 50 * mln, 50 * mln, ""},
	{"MRKP", 50 * mln, 50 * mln, ""},
	{"MRKK", 10 * mln, 15 * mln, ""},
	{"MRKU", 50 * mln, 50 * mln, ""},
	{"MRKV", 50 * mln, 50 * mln, ""},
	{"MSNG", 60 * mln, 50 * mln, ""},
	{"MSRS", 10 * mln, 15 * mln, ""},
	{"MTLR", 400 * mln, 80 * mln, ""},
	{"MTLRP", 400 * mln, 80 * mln, ""},
	{"MTSS", 250 * mln, 100 * mln, ""},
	{"MVID", 60 * mln, 50 * mln, ""},
	{"NKHP", 50 * mln, 50 * mln, ""},
	{"NLMK", 250 * mln, 100 * mln, ""},
	{"NMTP", 70 * mln, 60 * mln, ""},
	{"NVTK", 350 * mln, 110 * mln, ""},
	{"OGKB", 70 * mln, 50 * mln, ""},
	{"OZON", 250 * mln, 100 * mln, ""},
	{"OZPH", 80 * mln, 60 * mln, ""},
	{"PHOR", 220 * mln, 100 * mln, ""},
	{"PIKK", 170 * mln, 90 * mln, ""},
	{"PLZL", 280 * mln, 100 * mln, ""},
	{"POSI", 100 * mln, 90 * mln, ""},
	{"PRMD", 80 * mln, 60 * mln, ""},
	{"RASP", 150 * mln, 80 * mln, ""},
	{"RAGR", 100 * mln, 80 * mln, ""},
	{"RENI", 80 * mln, 60 * mln, ""},
	{"RNFT", 220 * mln, 90 * mln, ""},
	{"ROSN", 500 * mln, 120 * mln, ""},
	{"RTKM", 150 * mln, 80 * mln, ""},
	{"RTKMP", 100 * mln, 70 * mln, ""},
	{"RUAL", 200 * mln, 90 * mln, ""},
	{"SBER", 900 * mln, 120 * mln, ""},
	{"SELG", 100 * mln, 70 * mln, ""},
	{"SFIN", 200 * mln, 80 * mln, ""},
	{"SGZH", 170 * mln, 80 * mln, ""},
	{"SIBN", 350 * mln, 100 * mln, ""},
	{"SMLT", 250 * mln, 90 * mln, ""},
	{"SNGS", 200 * mln, 90 * mln, ""},
	{"SNSGP", 300 * mln, 100 * mln, ""},
	{"SOFL", 100 * mln, 60 * mln, ""},
	{"SPBE", 150 * mln, 90 * mln, ""},
	{"SVAV", 60 * mln, 50 * mln, ""},
	{"SVCB", 200 * mln, 100 * mln, ""},
	{"T", 900 * mln, 120 * mln, ""},
	{"TATN", 230 * mln, 100 * mln, ""},
	{"TATNP", 200 * mln, 100 * mln, ""},
	{"TGKA", 60 * mln, 50 * mln, ""},
	{"TRMK", 150 * mln, 90 * mln, ""},
	{"TRNFP", 250 * mln, 100 * mln, ""},
	{"UGLD", 160 * mln, 80 * mln, ""},
	{"UNAC", 110 * mln, 50 * mln, ""},
	{"UPRO", 110 * mln, 60 * mln, ""},
	{"UWGN", 90 * mln, 50 * mln, ""},
	{"VKCO", 170 * mln, 90 * mln, ""},
	{"VSMO", 110 * mln, 60 * mln, ""},
	{"VTBR", 400 * mln, 110 * mln, ""},
	{"WUSH", 100 * mln, 50 * mln, ""},
	{"X5", 300 * mln, 100 * mln, ""},
	{"YDEX", 410 * mln, 100 * mln, ""},
	{"ZAYM", 60 * mln, 50 * mln, ""},
}

// шиза на счёт гонки данных scrapper-ом и командой бота
func GetStocksInfoList() []StockScrapeInfo {
	lock.Lock()
	defer lock.Unlock()
	answer := make([]StockScrapeInfo, len(stocksInfo))
	copy(answer, stocksInfo)
	return answer
}

func InitFile() error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return SaveInFile()
	}
	return LoadFromFile()
}

func SaveInFile() error {
	lock.Lock()
	defer lock.Unlock()
	byteValue, err := json.Marshal(stocksInfo)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, byteValue, 0644)
	return err
}

func LoadFromFile() error {
	lock.Lock()
	defer lock.Unlock()
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&stocksInfo)
	return err
}

func Find(info *StockScrapeInfo) int {
	lock.Lock()
	defer lock.Unlock()
	for j := 0; j < len(stocksInfo); j++ {
		if info.StockTag == stocksInfo[j].StockTag {
			return j
		}
	}
	return -1
}

func Set(info *StockScrapeInfo) {
	lock.Lock()
	defer lock.Unlock()
	i := Find(info)
	if i == -1 {
		stocksInfo = append(stocksInfo, *info)
		return
	}
	stocksInfo[i] = *info
}

func Remove(info *StockScrapeInfo) {
	lock.Lock()
	defer lock.Unlock()
	j := Find(info)
	if j == -1 {
		return
	}
	stocksInfo = append(stocksInfo[:j], stocksInfo[j+1:]...)
}
