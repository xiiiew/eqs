package Okex

import (
	"bytes"
	"compress/flate"
	"crypto/tls"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/xiiiew/eqs/eqsModels"
	"github.com/xiiiew/eqs/eqsWebsocket"
	"io/ioutil"
	"net/url"
	"time"
)

const DefaultScheme = "wss"
const DefaultHost = "real.okex.com:8443"
const DefaultPath = "/ws/v3"

type OkexWsConn struct {
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
func NewDefauldOkexWsConn(symbol eqsModels.EqsSymbol) *OkexWsConn {
	outChan := make(chan interface{}, 100)
	errChan := make(chan error, 100)
	return &OkexWsConn{
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
*/
func NewOkexWsConnWithHost(scheme string, host string, symbol eqsModels.EqsSymbol) *OkexWsConn {
	outChan := make(chan interface{}, 100)
	errChan := make(chan error, 100)
	return &OkexWsConn{
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
func (h *OkexWsConn) createConnection() bool {
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
func (h *OkexWsConn) writeMessage(bytes []byte) bool {
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
func (h *OkexWsConn) readMessage() ([]byte, error) {
	data, err := h.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return h.unzip(data)
}

/*
解压
*/
func (h *OkexWsConn) unzip(bytesData []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(bytesData))
	unzipData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return unzipData, nil
}

/*
ping
*/
func (h *OkexWsConn) ping() {
	for {
		msg := []byte(fmt.Sprintf("{\"ping\" : %d}", time.Now().Unix()))
		if ! h.writeMessage(msg) {
			return
		}
		time.Sleep(5 * time.Second)
	}
}
