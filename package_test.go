package yfinance

import (
	"testing"
	"sort"
)

type SortProxy struct{
	Data []Price
}
func (sp SortProxy) Len() int{
	return len(sp.Data)
}
func (sp SortProxy) Less(i, j int) bool{
	return sp.Data[i].Date.Before(sp.Data[j].Date)
}
func (sp SortProxy) Swap(i, j int){
	tmp := sp.Data[i]
	sp.Data[i] = sp.Data[j]
	sp.Data[j] = tmp
}

func TestOneSymbol(t* testing.T){
	yf := Interface{}
	results, err := yf.GetPrices([]string{"MSFT"}, Date(2009, 1, 1), Date(2009, 12, 31))
	if (err != nil){
		t.Error("failed to get prices", err)
	}
	t.Log("results got: ", len(results))
	sort.Sort(SortProxy{results})
	if len(results) != 252{
		t.Error("invalid number of results: ", len(results))
	}
	if (results[0].Date.Year() != 2009 || results[0].Date.Month() != 1 || results[0].Date.Day() != 2){
		t.Error("invalid date on the first record: ", results[0].Date)
	}
	if (results[0].AdjustedClose != 1803){
		t.Error("invalid adjusted close on the first record: ", results[0].AdjustedClose)
	}
	li := len(results) - 1
	if (results[li].Date.Year() != 2009 || results[li].Date.Month() != 12 || results[li].Date.Day() != 31){
		t.Error("invalid date on the last record: ", results[li].Date)
	}
	if (results[li].AdjustedClose != 2766){
		t.Error("invalid adjusted close on the last record: ", results[li].AdjustedClose)
	}
}