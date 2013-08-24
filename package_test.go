package yfinance

import (
	"testing"
)

func TestOneSymbol(t* testing.T){
	yf := Interface{}
	results, err := yf.GetPrices([]string{"MSFT"}, Date(2009, 1, 1), Date(2009, 12, 31))
	if (err != nil){
		t.Error("failed to get prices", err)
	}
	if len(results.Symbols) != 1 || results.Symbols[0] != "MSFT"{
		t.Error("invalid number of symbols found: ", results.Symbols)
	}
	t.Log("results got: ", results)
	first, err := results.PriceAt("MSFT", results.From)
	if err == nil{
		t.Error("found nonexisting price", first)
	}
	prices, found := results.Prices("MSFT")
	if !found || len(prices) != 252{
		t.Error("can't get all prices", prices, err)
	}
	first, err = results.PriceAt("MSFT", Date(2009, 1, 2))
	if (first.Date.Year() != 2009 || first.Date.Month() != 1 || first.Date.Day() != 2){
		t.Error("invalid date on the first record: ", first.Date)
	}
	if (first.AdjustedClose != 1803){
		t.Error("invalid adjusted close on the first record: ", first.AdjustedClose)
	}
	last, err := results.PriceAt("MSFT", results.To)
	if err != nil{
		t.Error("didn't found existing price", last)
	}
	if (last.Date.Year() != 2009 || last.Date.Month() != 12 || last.Date.Day() != 31){
		t.Error("invalid date on the last record: ", last.Date)
	}
	if (last.AdjustedClose != 2766){
		t.Error("invalid adjusted close on the last record: ", last.AdjustedClose)
	}
}
func TestTwoSymbols(t* testing.T){
	yf := Interface{}
	results, err := yf.GetPrices([]string{"MSFT", "GOOG"}, Date(2009, 1, 1), Date(2009, 12, 31))
	if (err != nil){
		t.Error("failed to get prices", err)
	}
	if len(results.Symbols) != 2{
		t.Error("invalid number of symbols found: ", results.Symbols)
	}
	t.Log("results got: ", results)
	ms1, _ := results.PriceAt("MSFT", Date(2009, 1, 2))
	go1, _ := results.PriceAt("GOOG", Date(2009, 1, 2))
	if (ms1.AdjustedClose != 1803 || go1.AdjustedClose != 32132){
		t.Error("invalid adjusted close on the first record: ", ms1, go1)
	}
	ms2, _ := results.PriceAt("MSFT", results.To)
	go2, _ := results.PriceAt("GOOG", results.To)
	if (ms2.AdjustedClose != 2766 || go2.AdjustedClose != 61998){
		t.Error("invalid adjusted close on the last record: ", ms2, go2)
	}
}
func TestNonExistingSymbol(t* testing.T){
	yf := Interface{}
	results, err := yf.GetPrices([]string{"blahblahcallmeifthisexists"}, Date(2009, 1, 1), Date(2009, 12, 31))
	if (err == nil){
		t.Error("succeeded in getting prices, got results: ", len(results.Symbols))
	}
}
func TestBadDate(t* testing.T){
	yf := Interface{}
	results, err := yf.GetPrices([]string{"MSFT"}, Date(2009, 1, 1), Date(2007, 12, 31))
	if (err == nil){
		t.Error("succeeded in getting prices, got results: ", len((*results).Symbols))
	}
}

func TestReadBadSymbol(t* testing.T){
	yf := Interface{}
	results, _:= yf.GetPrices([]string{"MSFT"}, Date(2009, 1, 1), Date(2007, 12, 31))
	ok := false
	defer func(){
		recover()
		ok = true
	}()
	results.PriceAt("haventAskedForThisSymbol", Date(2009, 1, 1))
	if (!ok){
		t.Error("didn't pacnicked on nonexisting symbol")
	}
}