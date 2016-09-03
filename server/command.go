package main

type (
	WsCommand struct {
		Car        string   `json:"car"`
		Speed      byte     `json:"speed"`
		Directions []string `json:"directions"`
	}
)
