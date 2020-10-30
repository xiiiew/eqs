package eqsModels

import (
	"fmt"
	"strings"
)

/*
交易对

如交易对BTC/USDT:
Base: BTC
Quote: USDT
Sep: /
*/
type EqsSymbol struct {
	Base  string
	Quote string
	Sep   string
}

/*
交易对转换成字符串
*/
func (s *EqsSymbol) ToString() string {
	return fmt.Sprintf("%s%s%s", s.Base, s.Sep, s.Quote)
}

/*
交易对转换成字符串,替换sep
*/
func (s *EqsSymbol) ToStringWithSep(sep string) string {
	return fmt.Sprintf("%s%s%s", s.Base, sep, s.Quote)
}

/*
交易对转换成大写字符串
*/
func (s *EqsSymbol) ToUpper() string {
	base := strings.ToUpper(s.Base)
	quote := strings.ToUpper(s.Quote)
	return fmt.Sprintf("%s%s%s", base, s.Sep, quote)
}

/*
交易对转换成大写字符串,替换sep
*/
func (s *EqsSymbol) ToUpperWithSep(sep string) string {
	base := strings.ToUpper(s.Base)
	quote := strings.ToUpper(s.Quote)
	return fmt.Sprintf("%s%s%s", base, sep, quote)
}

/*
交易对转换成小写字符串
*/
func (s *EqsSymbol) ToLower() string {
	base := strings.ToLower(s.Base)
	quote := strings.ToLower(s.Quote)
	return fmt.Sprintf("%s%s%s", base, s.Sep, quote)
}

/*
交易对转换成小写字符串,替换sep
*/
func (s *EqsSymbol) ToLowerWithSep(sep string) string {
	base := strings.ToLower(s.Base)
	quote := strings.ToLower(s.Quote)
	return fmt.Sprintf("%s%s%s", base, sep, quote)
}
