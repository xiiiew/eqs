# eqs

eqs是一个包括火币(Huobipro)，币安(Binance),，Ok(Okex)等在内的各头部交易所的实时行情爬虫，使用websocket对接，实时获取最新行情

---

## 安装

`go get github.com/xiiiew/eqs`

## 火币

#### 获取行情数据

```go
package main

import (
	"fmt"
	"github.com/xiiiew/eqs/eqsModels"
	"github.com/xiiiew/eqs/exchange/huobipro"
)

func main() {
	symbol := eqsModels.EqsSymbol{
		Base:  "BTC",
		Quote: "USDT",
		Sep:   "/",
	}

	// 使用默认host, 需要科学上网
	//huobiproConn := huobipro.NewDefauldHuobiproWsConn(symbol, "001")

	huobiproConn := huobipro.NewHuobiproWsConnWithHost("wss", "api.huobi.me", symbol, "001")

	go huobiproConn.StartMarketDetail()

	// 获取行情数据
	go func() {
		for {
			outData := <-huobiproConn.OutChan
			data := outData.(huobipro.HuobiproMarketDetail)
			fmt.Printf("%+v\n", data)
		}
	}()

	// 获取报错信息
	go func() {
		for {
			fmt.Println(<-huobiproConn.ErrChan)
		}
	}()

	select {}
}
```

## Okex

#### 获取市场成交

```go
package main

import (
	"fmt"
	"github.com/xiiiew/eqs/eqsModels"
	"github.com/xiiiew/eqs/exchange/okex"
)

func main() {
	symbol := eqsModels.EqsSymbol{
		Base:  "BTC",
		Quote: "USDT",
		Sep:   "/",
	}

	okexConn := Okex.NewOkexWsConnWithHost("wss","real.okex.com:8443", symbol)
	go okexConn.StartSpotTrade()

	// 获取行情数据
	go func() {
		for {
			outData := <-okexConn.OutChan
			data := outData.(Okex.OkexSpotTrade)
			fmt.Printf("%+v\n", data)
		}
	}()

	// 获取报错信息
	go func() {
		for {
			fmt.Println(<-okexConn.ErrChan)
		}
	}()

	select {}
}
```

## 币安

#### 获取市场逐笔成交

```go
package main

import (
	"fmt"
	"github.com/xiiiew/eqs/eqsModels"
	"github.com/xiiiew/eqs/exchange/binance"
)

func main() {
	symbol := eqsModels.EqsSymbol{
		Base:  "BTC",
		Quote: "USDT",
		Sep:   "/",
	}

	binanceConn := binance.NewDefauldBinanceWsConn(symbol)
	go binanceConn.StartTrade()

	// 获取行情数据
	go func() {
		for {
			outData := <-binanceConn.OutChan
			data := outData.(binance.BinanceTrade)
			fmt.Printf("%+v\n", data)
		}
	}()

	// 获取报错信息
	go func() {
		for {
			fmt.Println(<-binanceConn.ErrChan)
		}
	}()

	select {}
}
```
