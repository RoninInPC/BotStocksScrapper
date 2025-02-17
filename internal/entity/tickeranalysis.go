package entity

import (
	"cmp"
	"slices"
)

type StocksInfo []StockInfo

func (s StocksInfo) Len() int {
	return len(s)
}

func (s StocksInfo) SumVolume() float64 {
	var sum = 0.0
	for _, v := range s {
		sum += v.Volume
	}
	return sum
}

func (s StocksInfo) AvgVolume() float64 {
	return s.SumVolume() / float64(s.Len())
}

func (s StocksInfo) MedianVolume() float64 {
	slices.SortFunc(s, func(a, b StockInfo) int {
		return cmp.Compare(a.Volume, b.Volume)
	})
	return s[s.Len()/2].Volume
}

func (s StocksInfo) MaxVolume() float64 {
	return slices.MaxFunc(s, func(a, b StockInfo) int {
		return cmp.Compare(a.Volume, b.Volume)
	}).Volume
}

func (s StocksInfo) MinVolume() float64 {
	return slices.MinFunc(s, func(a, b StockInfo) int {
		return cmp.Compare(a.Volume, b.Volume)
	}).Volume
}

func (s StocksInfo) AvgPrice() float64 {
	var sum = 0.0
	for _, v := range s {
		sum += v.Stock.Price
	}
	return sum / float64(s.Len())
}

func (s StocksInfo) MedianPrice() float64 {
	slices.SortFunc(s, func(a, b StockInfo) int {
		return cmp.Compare(a.Stock.Price, b.Stock.Price)
	})
	return s[s.Len()/2].Stock.Price
}

func (s StocksInfo) SumLots() int64 {
	answer := int64(0)
	for _, v := range s {
		answer += v.LotsCount
	}
	return answer
}

func (s StocksInfo) AvgLots() float64 {
	return float64(s.SumLots()) / float64(s.Len())
}

func (s StocksInfo) MedianLots() int64 {
	slices.SortFunc(s, func(a, b StockInfo) int {
		return cmp.Compare(a.LotsCount, b.LotsCount)
	})
	return s[s.Len()/2].LotsCount
}

func (s StocksInfo) MaxLots() int64 {
	return slices.MaxFunc(s, func(a, b StockInfo) int {
		return cmp.Compare(a.LotsCount, b.LotsCount)
	}).LotsCount
}

func (s StocksInfo) MinLots() int64 {
	return slices.MinFunc(s, func(a, b StockInfo) int {
		return cmp.Compare(a.LotsCount, b.LotsCount)
	}).LotsCount
}

func (s StocksInfo) ToString(prefix, postfix string) string {
	answer := ""
	for _, v := range s {
		answer += prefix + v.StringShort() + postfix
	}
	return answer
}

type StocksByMoveType map[StockMoveType]StocksInfo

func (s StocksByMoveType) Add(info StockInfo) StocksByMoveType {
	_, ok := s[info.StockMove]
	if ok {
		s[info.StockMove] = append(s[info.StockMove], info)
	} else {
		s[info.StockMove] = StocksInfo{info}
	}
	return s
}

func (s StocksByMoveType) Copy() StocksByMoveType {
	ans := make(StocksByMoveType)
	for key, value := range s {
		ans[key] = value
	}
	return ans
}

func (s StocksByMoveType) ToSlice() StocksInfo {
	answer := make(StocksInfo, 0)
	for _, v := range s {
		answer = append(answer, v...)
	}
	return answer
}

func (s StocksByMoveType) ToString() string {
	answer := ""
	for t, v := range s {
		answer += t.String() + ":\n" + v.ToString("\t", "\n")
	}
	return answer
}

type TickerAnalysis map[string]StocksByMoveType

func (t TickerAnalysis) Add(info StockInfo) TickerAnalysis {
	_, ok := t[info.Stock.Ticker]
	if ok {
		t[info.Stock.Ticker] = t[info.Stock.Ticker].Add(info)
	} else {
		t[info.Stock.Ticker] = make(StocksByMoveType)
		t[info.Stock.Ticker] = t[info.Stock.Ticker].Add(info)
	}
	return t
}
func (t TickerAnalysis) GetByTicker(ticker string) StocksByMoveType {
	return t[ticker]
}
