## Setup

### Running the Server
 * go run rpcServer.go

### Running the Client

# Buy Stocks

```
go run rpcClient.go Request "stockSymbolAndPercentage":"GOOG:50%,YHOO:50%" "budget":1000

```
# Get Stocks Info
```
go run rpcClient.go Request "tradeId":1
```