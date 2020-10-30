package huobipro

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/binary"
	"github.com/xiiiew/eqs/eqsWebsocket"
	"github.com/xiiiew/eqs/eqsModels"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/url"
	"time"
)

const DefaultScheme = "wss"
const DefaultHost = "api.huobi.pro"
const DefaultPath = "/ws"

type HuobiproWsConn struct {
	scheme  string
	host    string
	path    string
	symbol  eqsModels.EqsSymbol
	id      string
	OutChan chan interface{}
	ErrChan chan error
	wsConn  *eqsWebsocket.Connection
}

/*
使用默认地址, 大陆网络慎用
*/
func NewDefauldHuobiproWsConn(symbol eqsModels.EqsSymbol, id string) *HuobiproWsConn {
	outChan := make(chan interface{}, 100)
	errChan := make(chan error, 100)
	return &HuobiproWsConn{
		scheme:  DefaultScheme,
		host:    DefaultHost,
		path:    DefaultPath,
		symbol:  symbol,
		id:      id,
		OutChan: outChan,
		ErrChan: errChan,
	}
}

/*
替换默认地址
api文档: https://huobiapi.github.io/docs/spot/v1/cn/#7c47ef3411
*/
func NewHuobiproWsConnWithHost(scheme string, host string, symbol eqsModels.EqsSymbol, id string) *HuobiproWsConn {
	outChan := make(chan interface{}, 100)
	errChan := make(chan error, 100)
	return &HuobiproWsConn{
		scheme:  scheme,
		host:    host,
		path:    DefaultPath,
		symbol:  symbol,
		id:      id,
		OutChan: outChan,
		ErrChan: errChan,
	}
}

/*
创建websocket连接
*/
func (h *HuobiproWsConn) createConnection() bool {
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
func (h *HuobiproWsConn) writeMessage(bytes []byte) bool {
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
func (h *HuobiproWsConn) readMessage() ([]byte, error) {
	data, err := h.wsConn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return h.unzip(data)
}

/*
解压
*/
func (h *HuobiproWsConn) unzip(bytesData []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.LittleEndian, bytesData)
	if err != nil {
		return nil, err
	}

	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	unzipData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return unzipData, nil
}

/*
ping
*/
func (h *HuobiproWsConn) ping() {
	for {
		msg := []byte(fmt.Sprintf("{\"ping\" : %d}", time.Now().Unix()))
		if ! h.writeMessage(msg) {
			return
		}
		time.Sleep(5)
	}
}
