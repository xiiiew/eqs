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