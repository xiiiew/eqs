/*
获取市场逐笔交易
api文档: https://binance-docs.github.io/apidocs/spot/cn/#2b149598d9
*/
package binance

import (
	"encoding/json"
	"fmt"
	"strings"
)

/*
"e": "trade",     // 事件类型
"E": 123456789,   // 事件时间
"s": "BNBBTC",    // 交易对
"t": 12345,       // 交易ID
"p": "0.001",     // 成交价格
"q": "100",       // 成交笔数
"b": 88,          // 买方的订单ID
"a": 50,          // 卖方的订单ID
"T": 123456785,   // 成交时间
"m": true,        // 买方是否是做市方。如true，则此次成交是一个主动卖出单，否则是一个主动买入单。
"M": true         // 请忽略该字段
*/
type BinanceTrade struct {
	EL string `json:"e"`
	E  int    `json:"E"`
	SL string `json:"s"`
	TL int    `json:"t"`
	PL string `json:"p"`
	QL string `json:"q"`
	BL int    `json:"b"`
	AL int    `json:"a"`
	T  int    `json:"T"`
	ML bool   `json:"m"`
	M  bool   `json:"M"`
}

func (h *BinanceWsConn) StartTrade() {
	// 创建ws
	for {
		if h.createConnection() {
			break
		}
	}
	defer h.wsConn.Close()

	// 订阅
	for {
		if h.subscribeTrade() {
			break
		}
	}

	// ping
	go h.ping()

	go h.readTrade()

	select {}
}

/*
订阅频道
*/
func (h *BinanceWsConn) subscribeTrade() bool {
	symbol := h.symbol.ToLowerWithSep("")
	streamName := fmt.Sprintf("%s@trade", symbol)
	message, _ := json.Marshal(map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{streamName},
		"id":     1,
	})
	if !h.writeMessage(message) {
		return false
	}

	return true
}

/*
获取市场逐笔成交
*/
func (h *BinanceWsConn) readTrade() {
	defer func() {
		go h.StartTrade()
	}()

	for {
		data, err := h.readMessage()
		if err != nil {
			h.ErrChan <- err
		}

		dataStr := string(data)
		if strings.Contains(dataStr, "result") { // 订阅成功消息
			continue
		} else {
			rt := BinanceTrade{}
			err := json.Unmarshal(data, &rt)
			if err != nil {
				h.ErrChan <- err
			} else {
				h.OutChan <- rt
			}
		}
	}
}
