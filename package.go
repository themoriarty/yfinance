package yfinance
import (
	"fmt"
	"net/http"
	"time"
	//"log"
	"strconv"
	"encoding/csv"
)

type Interface struct{
}

type Error struct{
	Msg string	
}
func (e Error) Error() string{
	return e.Msg
}

type Price struct{
	Date time.Time
	AdjustedClose int
}

func Date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func parseCsv(reader *csv.Reader, out chan Price) error{
	header, err := reader.Read()
	if (err != nil){
		return Error{fmt.Sprintf("Can't read header: %s", err)}
	}
	idx := make(map[string]int)
	for i, s := range(header){
		idx[s] = i
	}
	adjCloseIdx, ok := idx["Adj Close"];
	if (!ok){
		return fmt.Errorf("Can't find 'Adj Close' column in this: %s", idx)
	}
	dateIdx, ok := idx["Date"];
	if (!ok){
		return Error{"Can't find 'Date' column"}
	}
	allData, err := reader.ReadAll();
	if (err != nil){
		return fmt.Errorf("Can't read all data: %s", err)
	}
	for _, row := range(allData){
		if date, err := time.Parse("2006-01-02", row[dateIdx]); err == nil{
			if adjPrice, err := strconv.ParseFloat(row[adjCloseIdx], 64); err == nil{
				out <- Price{date, int(adjPrice * 100)}
			} else{
				return fmt.Errorf("Can't parse price from %s: %s", row[adjCloseIdx], err)
			}
		} else{
			return fmt.Errorf("Can't parse data from %s: %s", row[dateIdx], err)
		}
	}
	return nil
}

func (this Interface) GetPrices(symbols []string, from time.Time, to time.Time) ([]Price, error){
	outChan := make(chan Price)
	doneChan := make(chan error, len(symbols))
	activeWorkers := 0
	for _, s := range(symbols){
		go func(symbol string){
			url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&d=%d&e=%d&f=%d&g=d&a=%d&b=%d&c=%d&ignore=.csv",
			symbol, to.Month() - 1, to.Day(), to.Year(), from.Month() - 1, from.Day(), from.Year())
			res, err := http.Get(url)
			//log.Printf("asked for %s. err: %s, http code: %d\n", url, err, res.StatusCode)
			if res != nil{
				defer res.Body.Close()
			}
			if err == nil && res.StatusCode >= 200 && res.StatusCode < 300{				
				doneChan <- parseCsv(csv.NewReader(res.Body), outChan)
				return
			}
			doneChan <- fmt.Errorf("can't fetch data: %s (http status %d)", err, res.StatusCode)
		}(s)
		activeWorkers++
	}
	ret := make([]Price, 0)
	var retError error
	for activeWorkers > 0 {
		select{
		case err := <- doneChan:
			activeWorkers--
			if err != nil{
				retError = err
			}
		case r:= <- outChan:
			ret = append(ret, r)
		}
	}
	return ret, retError
}

