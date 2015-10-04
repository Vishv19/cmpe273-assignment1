package main

import (
    "net/http"
    "github.com/gorilla/rpc"
    "github.com/gorilla/rpc/json"
    "strconv"
    "yahoofinance"
    "math"
    "strings"
)

var tradeStorage [100]TradeStorage
var tradecount int = 1

type Args struct {
    StockSymbolAndPercentage string
    Budget int
}

type Reply struct {
    TradeId int
    Stocks string
    UnvestedAmount float64
}

type TradeArgs struct {
    TradeId int
}

type TradeReply struct {
    Stocks string
    CurrentMarketValue float64
    UnvestedAmount float64
}

type StockInfo struct {
    companyName string
    percentage int
}

type StockStorage struct {
    companyName string
    numberOfStock int
    stockPrice float64
}

type TradeStorage struct {
    stockList []StockStorage
    tradeId int
    uninvestedAmount float64
}

type StockService struct {}

func splitStockData(stockInfo string) []StockInfo{
    var stockData []string = strings.Split(stockInfo, ",")
    lengthOfStock := len(stockData)
    splitData := make([]StockInfo, lengthOfStock, lengthOfStock)

    for i := 0; i < lengthOfStock; i++ {
        var leftData []string = strings.Split(stockData[i], ":")
        leftData[1] = strings.TrimSuffix(leftData[1], "%")
        splitData[i].companyName = leftData[0]
        splitData[i].percentage, _ = strconv.Atoi(leftData[1])
    }
    return splitData
}

func buyStock(stockInfo []StockInfo, priceList []float64, budget int) Reply{
    totalItems := len(stockInfo)
    tradeStorage[tradecount].stockList = make([]StockStorage, totalItems, totalItems)
    var responseData string
    var reply Reply

    for i := 0; i < totalItems; i++ {
        var expectedInvestment float64 = float64((budget * stockInfo[i].percentage)/100)
        var numberOfStock int
        var uninvestedMoney float64
        if expectedInvestment > priceList[i] {
            numberOfStock = int(expectedInvestment/priceList[i])
            uninvestedMoney = math.Mod(expectedInvestment, priceList[i])
        } else {
            numberOfStock = 0
            uninvestedMoney = expectedInvestment
        }
        responseData += "'" + stockInfo[i].companyName + ":" + strconv.Itoa(numberOfStock) + ":" + "$" + strconv.FormatFloat(priceList[i], 'f', -1, 64) + "'"
        if i!=(totalItems-1) {
            responseData  += ","
        }
        tradeStorage[tradecount].stockList[i].companyName = stockInfo[i].companyName
        tradeStorage[tradecount].stockList[i].numberOfStock = numberOfStock
        tradeStorage[tradecount].stockList[i].stockPrice = priceList[i]
        tradeStorage[tradecount].uninvestedAmount += uninvestedMoney 
    }
    tradeStorage[tradecount].tradeId = tradecount
    reply.TradeId = tradecount
    reply.Stocks = responseData
    reply.UnvestedAmount = tradeStorage[tradecount].uninvestedAmount
    tradecount++
    return reply
}

func getStock(tradeId int) TradeReply{
    var nameOfCompanies string
    var totalStocks int = len(tradeStorage[tradeId].stockList)
    var currentMarketValue float64
    var responseData string
    var reply TradeReply

    for i := 0; i < (totalStocks-1); i++ {
        nameOfCompanies += tradeStorage[tradeId].stockList[i].companyName + ","
    }
    nameOfCompanies += tradeStorage[tradeId].stockList[totalStocks-1].companyName
    var currentPriceList []float64 = yahoofinance.ReturnStockPrice(nameOfCompanies)

    for i := 0; i < totalStocks; i++ {
        companyName := tradeStorage[tradeId].stockList[i].companyName
        numberOfStock := tradeStorage[tradeId].stockList[i].numberOfStock
        buyingPrice := tradeStorage[tradeId].stockList[i].stockPrice
        currentMarketValue += float64(numberOfStock) * currentPriceList[i]
        responseData += "'" + companyName + ":" + strconv.Itoa(numberOfStock) + ":"

        if(buyingPrice < currentPriceList[i]) {
            responseData += "+" + "$" + strconv.FormatFloat(currentPriceList[i], 'f', -1, 64) + "'"
        } else if(buyingPrice > currentPriceList[i]) {
            responseData += "-" + "$" + strconv.FormatFloat(currentPriceList[i], 'f', -1, 64) + "'"
        } else {
            responseData += "$" + strconv.FormatFloat(currentPriceList[i], 'f', -1, 64) + "'"
        }
        if i!=(totalStocks-1) {
            responseData  += ","
        }
    }
    reply.Stocks = responseData
    reply.CurrentMarketValue = currentMarketValue
    reply.UnvestedAmount = tradeStorage[tradeId].uninvestedAmount

    return reply
}

func (s *StockService) ReturnStock(r *http.Request, args *Args, reply *Reply) error {
    var receivedData string = args.StockSymbolAndPercentage
    var budget int = args.Budget
    var stockInfo []StockInfo = splitStockData(receivedData)

    var nameOfCompanies string
    for i := 0; i < len(stockInfo)-1; i++ {
        nameOfCompanies += stockInfo[i].companyName + "," 
    }
    nameOfCompanies += stockInfo[len(stockInfo)- 1].companyName

    var priceList []float64 = yahoofinance.ReturnStockPrice(nameOfCompanies)
    var replyFromStockMarket Reply = buyStock(stockInfo, priceList, budget)
    reply.TradeId = replyFromStockMarket.TradeId
    reply.Stocks = replyFromStockMarket.Stocks
    reply.UnvestedAmount = replyFromStockMarket.UnvestedAmount
    return nil
}

func (s *StockService) ReturnPortfolio(r *http.Request, args *TradeArgs, reply *TradeReply) error {
    tradeId := args.TradeId

    var replyFromSystem TradeReply = getStock(tradeId)
    reply.Stocks = replyFromSystem.Stocks
    reply.CurrentMarketValue = replyFromSystem.CurrentMarketValue
    reply.UnvestedAmount = replyFromSystem.UnvestedAmount
    return nil
}

func main() {
    s := rpc.NewServer()
    s.RegisterCodec(json.NewCodec(), "application/json")
    s.RegisterService(new(StockService), "")
    http.Handle("/stock", s)
	http.ListenAndServe(":8080", nil)
}