package simple

import (
	"BotStocksScrapper/internal/entity"
	"sync"
)

type SimpleAnalysis struct {
	content entity.TickerAnalysis
	sync.Mutex
}

func Init() SimpleAnalysis {
	return SimpleAnalysis{content: make(entity.TickerAnalysis)}
}

func (s SimpleAnalysis) RemoveAll() bool {
	for t, _ := range s.content {
		s.RemoveTicker(t)
	}
	return true
}

func (s SimpleAnalysis) GetAnomaly(stock entity.Stock) (entity.StocksByMoveType, bool) {
	answer := s.GetInfoByTicker(stock.Ticker)
	if answer.ToSlice().SumVolume() >= stock.AnomalySize {
		return answer, true
	}
	return nil, false
}

func (s SimpleAnalysis) Add(info entity.StockInfo) bool {
	s.Lock()
	defer s.Unlock()
	s.content = s.content.Add(info)
	return true
}

func (s SimpleAnalysis) GetInfoByTicker(ticker string) entity.StocksByMoveType {
	s.Lock()
	defer s.Unlock()
	return s.content.GetByTicker(ticker)
}

func (s SimpleAnalysis) RemoveTicker(ticker string) bool {
	s.Lock()
	defer s.Unlock()
	clear(s.content[ticker])
	s.content[ticker] = make(entity.StocksByMoveType)
	return true
}
