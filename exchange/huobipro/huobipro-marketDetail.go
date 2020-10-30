/*
获取市场概要
api文档: https://huobiapi.github.io/docs/spot/v1/cn/#7c47ef3411
*/
package huobipro

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
type HuobiproMarketDetail struct {
	Ch   string
	Ts   int
	Tick struct {
		Amount float64
		Open   float64
		Close  float64
		High   float64
		Id     int
		Count  int
		Low    float64
		Vol    float64
	}
}

func (h *HuobiproWsConn) StartMarketDetail() {
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
func (h *HuobiproWsConn) subscribeMarketDetail() bool {
	symbol := h.symbol.ToLowerWithSep("")
	sub := fmt.Sprintf("market.%s.detail", symbol)
	message, _ := json.Marshal(
		map[string]string{
			"sub": sub,
			"id":  h.id,
		},
	)
	if !h.writeMessage(message) {
		return false
	}

	return true
}

/*
获取市场概要数据
*/
func (h *HuobiproWsConn) readMarketDetail() {
	defer func() {
		go h.StartMarketDetail()
	}()

	for {
		unzipData, err := h.readMessage()
		if err != nil {
			h.ErrChan <- err
		}

		dataStr := string(unzipData)
		if strings.Contains(dataStr, "subbed") { // 订阅成功消息
			continue
		} else if strings.Contains(dataStr, "ping") { // ping
			h.writeMessage(unzipData)
		} else if strings.Contains(dataStr, "pong") { // pong
			continue
		} else {
			rt := HuobiproMarketDetail{
				Ch: "",
				Ts: 0,
				Tick: struct {
					Amount float64
					Open   float64
					Close  float64
					High   float64
					Id     int
					Count  int
					Low    float64
					Vol    float64
				}{Amount: 0, Open: 0, Close: 0, High: 0, Id: 0, Count: 0, Low: 0, Vol: 0},
			}
			err := json.Unmarshal(unzipData, &rt)
			if err != nil {
				h.ErrChan <- err
			}else {
				h.OutChan <- rt
			}
		}
	}
}
