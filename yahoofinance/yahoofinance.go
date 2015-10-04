package yahoofinance

import (
    "net/http"
    "encoding/json"
    "strconv"
    "io/ioutil"
)

type FinanceData struct {
    List struct {
        Meta struct {
            Count int `json:"count"`
            Start int `json:"start"`
            Type  string `json:"type"`
        } `json:"meta"`
        Resources []struct {
            Resource struct {
                Classname string `json:"classname"`
                Fields struct {
                    Name string `json:"name"`
                    Price string `json:"price"`
                    Symbol string `json:"symbol"`
                    Ts string `json:"ts"`
                    Type string `json:"type"`
                    Utctime string `json:"utctime"`
                    Volume string `json:"volume"`
                } `json:"fields"`
            } `json:"resource"`
        } `json:"resources"`
    } `json:"list"`
}

func ReturnStockPrice(companyName string) []float64{
    url := "http://finance.yahoo.com/webservice/v1/symbols/" + companyName + "/quote?format=json"
    response, err := http.Get(url)
    defer response.Body.Close()

    if err != nil {
        panic(err)
    }

     // read json http response
    jsonData, err := ioutil.ReadAll(response.Body)

    if err != nil {
        panic(err)
    }

    var financeData FinanceData

    err = json.Unmarshal(jsonData, &financeData) //unmarshall data extracted from yahoo finance api

    if err != nil {
        panic(err)
    }

    var length int = len(financeData.List.Resources)
    priceList := make([]float64, length, length)

    for i:= 0; i < length; i++ {
        priceList[i], _ = strconv.ParseFloat(financeData.List.Resources[i].Resource.Fields.Price, 64)
    }

    return priceList //return the price list related to costs
}