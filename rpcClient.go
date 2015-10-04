package main

import (
    "fmt"
    "flag"
    "encoding/json"
    "strings"
    "bytes"
    "strconv"
    "net/http"
    "io/ioutil"
)

var jsonId int = 1

type Params []struct {
    StockSymbolAndPercentage string `json:"stockSymbolAndPercentage"`
    Budget int `json:"budget"`
}

type BuyStock struct {
    Method string `json:"method"`
    Params []struct {
        StockSymbolAndPercentage string `json:"stockSymbolAndPercentage"`
        Budget int `json:"budget"`
    } `json:"params"`
    ID int `json:"id"`
}

type BuyStockResponse struct {
    Result struct {
        Stocks string  `json:"Stocks"`
        TradeID int     `json:"TradeId"`
        UnvestedAmount float64 `json:"UnvestedAmount"`
    } `json:"result"`
    Error  interface{} `json:"error"`
    ID int `json:"id"`
}

type GetPortfolio struct {
    Method string `json:"method"`
    Params []struct {
        TradeID int `json:"tradeId"`
    } `json:"params"`
    ID int `json:"id"`
}

type GetPortfolioResponse struct {
    Result struct {
        CurrentMarketValue float64 `json:"CurrentMarketValue"`
        Stocks string  `json:"Stocks"`
        UnvestedAmount float64 `json:"UnvestedAmount"`
    } `json:"result"`
    Error  interface{} `json:"error"`
    ID int `json:"id"`
}

func main() {
    flag.Parse()
    numberOfArgs := len(flag.Args())

    if(numberOfArgs == 3) { //if number of commandline arguments are three, then we buystocks

        budget := strings.Split(flag.Arg(2), ":")
        stockArgs := strings.Split(flag.Arg(1), ",")
        var stockInput string   //spit received argument and aggregate stock input
        var totalPercent int    //to check if enter values sum up to 100

        firstStockArg := strings.Split(stockArgs[0], ":")
        stockInput += firstStockArg[1] + ":" + firstStockArg[2] + ","
        stockPercent, _ := strconv.Atoi(strings.TrimSuffix(firstStockArg[2], "%"))
        totalPercent += stockPercent

        for i := 1; i < len(stockArgs); i++ {
            stockInput += stockArgs[i]
            percent := strings.Split(stockArgs[i], ":")
            stockPercentValue, _ := strconv.Atoi(strings.TrimSuffix(percent[1], "%"))
            totalPercent += stockPercentValue
            if i!=len(stockArgs)-1 {
                stockInput += ","
            }
        }
        if(totalPercent!=100) {
            fmt.Println("The total percent is not 100")
        } else {
            var buyStock BuyStock
            buyStock.ID = jsonId
            buyStock.Params = make(Params, 1, 1)
            buyStock.Params[0].Budget, _ = strconv.Atoi(budget[1])
            buyStock.Method = "StockService.ReturnStock"    //buystocks
            buyStock.Params[0].StockSymbolAndPercentage = stockInput
            jsondata, err := json.Marshal(buyStock)
            if err != nil {
                fmt.Println("error:", err)
            }

            url := "http://localhost:8080/stock"

            req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsondata))
            req.Header.Set("Content-Type", "application/json")

            client := &http.Client{}
            resp, err := client.Do(req)
            if err != nil {
                panic(err)
            }
            defer resp.Body.Close()

            body, _ := ioutil.ReadAll(resp.Body)
            var data BuyStockResponse

            err = json.Unmarshal(body, &data)

            if err != nil {
                panic(err)
            }
            fmt.Println("'tradeId':", data.Result.TradeID)
            fmt.Println("'stocks':", data.Result.Stocks)
            fmt.Println("'unvestedAmount':", data.Result.UnvestedAmount)
        }

        jsonId++
    } else if(numberOfArgs == 2){   //if number of commandline arguments are two, then we return stock data
        type Params []struct {
            TradeID int `json:"tradeId"`
        }
        tradeId := strings.Split(flag.Arg(1), ":")
        var portfolio GetPortfolio
        portfolio.ID = jsonId
        portfolio.Params = make(Params, 1, 1)
        portfolio.Params[0].TradeID, _ = strconv.Atoi(tradeId[1])
        portfolio.Method = "StockService.ReturnPortfolio"

        jsondata, err := json.Marshal(portfolio)
        if err != nil {
            fmt.Println("error:", err)
        }

        url := "http://localhost:8080/stock"

        req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsondata))
        req.Header.Set("Content-Type", "application/json")

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()

        body, _ := ioutil.ReadAll(resp.Body)
        var data GetPortfolioResponse

        err = json.Unmarshal(body, &data) // here!

        if err != nil {
            panic(err)
        }
        fmt.Println("'stocks':", data.Result.Stocks)
        fmt.Println("'currentMarketValue':", data.Result.CurrentMarketValue)
        fmt.Println("'unvestedAmount':", data.Result.UnvestedAmount)
        jsonId++
    } else {
        fmt.Println("Number of Arguments should either be 2 or 3")
    }

}