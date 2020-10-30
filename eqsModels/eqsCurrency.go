package eqsModels

import "strings"

type EqsCurrency struct {
	Name string
}

func (obj *EqsCurrency) ToUpper() string {
	return strings.ToUpper(obj.Name)
}
