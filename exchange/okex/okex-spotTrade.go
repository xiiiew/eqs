/*
获取市场交易
api文档: https://www.okex.me/docs/zh/#spot_ws-trade
*/
package Okex

import (
	"encoding/json"
	"fmt"
	"strings"
)

/*
字段		数据类型	描述
id		integer	unix时间，同时作为消息ID
amount	float	24小时成交量
count	integer	24小时成交笔数
open	float	24小时开盘价
close	float	最新价
low		float	24小时最低价
high	float	24小时最高价
vol		float	24小时成交额
*/
type OkexSpotTrade struct {
	Table string
	Data  []struct {
		InstrumentId string
		Price        string
		Side         string
		Size         string
		Timestamp    string
		TradeId      string
	}
}

func (h *OkexWsConn) StartMarketDetail() {
	// 创建ws
	for {
		if h.createConnection() {
			break
		}
	}
	defer h.wsConn.Close()

	// 订阅
	for {
		if h.subscribeMarketDetail() {
			break
		}
	}

	// ping
	go h.ping()

	go h.readMarketDetail()

	select {}
}

/*
订阅频道
*/
func (h *OkexWsConn) subscribeMarketDetail() bool {
	symbol := h.symbol.ToUpperWithSep("-")
	args := []string{
		fmt.Sprintf("spot/trade:%s", symbol),
	}
	message, _ := json.Marshal(map[string]interface{}{
		"op":   "subscribe",
		"args": args,
	})
	if !h.writeMessage(message) {
		return false
	}

	return true
}

/*
获取市场成交
*/
func (h *OkexWsConn) readMarketDetail() {
	defer func() {
		go h.StartMarketDetail()
	}()

	for {
		unzipData, err := h.readMessage()
		if err != nil {
			h.ErrChan <- err
		}

		dataStr := string(unzipData)
		if strings.Contains(dataStr, "event") { // 订阅成功消息
			continue
		} else {
			rt := OkexSpotTrade{
				Table: "",
				Data:  nil,
			}
			err := json.Unmarshal(unzipData, &rt)
			if err != nil {
				h.ErrChan <- err
			} else {
				h.OutChan <- rt
			}
		}
	}
}
