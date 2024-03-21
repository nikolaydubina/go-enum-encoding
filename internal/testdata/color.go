package main

import (
	"encoding/json"
	"fmt"
)

type Color struct{ c uint }

//go:generate go-enum-encoding -type=Color
var (
	Undefined      = Color{}            // json:"-"
	Red            = Color{1}           // json:"red"
	Green          = Color{2}           // json:"green"
	Blue           = Color{3}           // json:"blue"
	Purple, Orange = Color{4}, Color{5} // json:"blue"
)

type V struct {
	Color Color `json:"color"`
}

func main() {
	var v V
	s := `{"color": "red"}`
	json.Unmarshal([]byte(s), &v)
	fmt.Println(v)
}
