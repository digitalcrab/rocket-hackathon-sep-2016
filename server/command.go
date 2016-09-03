package main

type (
	WsCommand struct {
		Car       string `json:"car"`
		Speed     byte   `json:"speed"`
		Direction string `json:"direction"`
	}
)
