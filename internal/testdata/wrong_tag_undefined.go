package main

type BadUndefined1 struct{ c uint }

//go:generate go-enum-encoding -type=BadUndefined1
var (
	UndefinedBadUndefined1 = BadUndefined1{}  // json:"something"
	ABadUndefined1         = BadUndefined1{1} // json:"red"
	BBadUndefined1         = BadUndefined1{2} // json:"blue"
)
