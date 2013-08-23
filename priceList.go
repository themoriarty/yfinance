package yfinance

import (
	"fmt"
	"time"
	"sort"
)

type PriceList struct{
	Symbols []string
	From time.Time
	To time.Time

	prices map[string][]Price
}

func (p *PriceList) PriceAt(symbol string, at time.Time) (Price, error) {
	// if we have prices on 1st and 5th, and asked for a price on 2nd, return price at the end of 1st
	prices, found := p.prices[symbol]
	if !found{
		panic(fmt.Sprint("can't find symbol", symbol, "in", p.Symbols))
	}
	idx := sort.Search(len(prices), func(i int) bool{
		return prices[i].Date.Equal(at) || prices[i].Date.After(at)
	})
	if idx == len(prices){
		return Price{}, fmt.Errorf("can't find price for", symbol, ", at time", at)
	}
	if prices[idx].Date.After(at){
		idx--
		if idx < 0{
			return Price{}, fmt.Errorf("can't find price for", symbol, ", at time", at)
		}
	}
	return prices[idx], nil
}

func priceList(from time.Time, to time.Time, pricesBySymbol map[string][]Price) (ret *PriceList){
	ret = &PriceList{nil, from, to, make(map[string][]Price)}
	for symbol, prices := range(pricesBySymbol){
		ret.Symbols = append(ret.Symbols, symbol)
		ret.prices[symbol] = prices
		sort.Sort(sortProxy{ret.prices[symbol]})
	}
	return
}

type sortProxy struct{
	Data []Price
}
func (sp sortProxy) Len() int{
	return len(sp.Data)
}
func (sp sortProxy) Less(i, j int) bool{
	return sp.Data[i].Date.Before(sp.Data[j].Date)
}
func (sp sortProxy) Swap(i, j int){
	tmp := sp.Data[i]
	sp.Data[i] = sp.Data[j]
	sp.Data[j] = tmp
}

