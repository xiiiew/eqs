package binance

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	"github.com/xiiiew/eqs/eqsModels"
	"github.com/xiiiew/eqs/eqsWebsocket"
	"net/url"
	"time"
)

const DefaultScheme = "wss"
const DefaultHost = "stream.binance.com:9443"
const DefaultPath = "/ws"

type BinanceWsConn struct {
	scheme  string
	host    string
	path    string
	symbol  eqsModels.EqsSymbol
	OutChan chan interface{}
	ErrChan chan error
	wsConn  *eqsWebsocket.Connection
}

/*
使用默认地址, 大陆网络慎用
*/
func NewDefauldBinanceWsConn(symbol eqsModels.EqsSymbol, id string) *BinanceWsConn {
	outChan := make(chan interface{}, 100)
	errChan := make(chan error, 100)
	return &BinanceWsConn{
		scheme:  DefaultScheme,
		host:    DefaultHost,
		path:    DefaultPath,
		symbol:  symbol,
		OutChan: outChan,
		ErrChan: errChan,
	}
}

/*
替换默认地址
api文档: https://huobiapi.github.io/docs/spot/v1/cn/#7c47ef3411
*/
func NewBinanceWsConnWithHost(scheme string, host string, symbol eqsModels.EqsSymbol) *BinanceWsConn {
	outChan := make(chan interface{}, 100)
	errChan := make(chan error, 100)
	return &BinanceWsConn{
		scheme:  scheme,
		host:    host,
		path:    DefaultPath,
		symbol:  symbol,
		OutChan: outChan,
		ErrChan: errChan,
	}
}

/*
创建websocket连接
*/
func (h *BinanceWsConn) createConnection() bool {
	u := url.URL{
		Scheme: h.scheme,
		Host:   h.host,
		Path:   h.path,
	}
	dialer := &websocket.Dialer{
		TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}, // 禁用https证书验证
	}
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		h.ErrChan <- err
		return false
	}

	w, err := eqsWebsocket.InitConnection(conn)

	if err != nil {
		h.ErrChan <- err
		return false
	}

	h.wsConn = w
	return true
}

/*
发送数据
*/
func (h *BinanceWsConn) writeMessage(bytes []byte) bool {
	err := h.wsConn.WriteMessage(bytes)
	if err != nil {
		h.ErrChan <- err
		return false
	}
	return true
}

/*
读取数据
*/
func (h *BinanceWsConn) readMessage() ([]byte, error) {
	data, err := h.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return data, err
}

/*
ping
*/
func (h *BinanceWsConn) ping() {
	for {
		msg := []byte("{\"method\": \"GET_PROPERTY\", \"params\":[\"combined\"],\"id\":2}")
		if ! h.writeMessage(msg) {
			return
		}
		time.Sleep(5 * time.Second)
	}
}
