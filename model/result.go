package model

type Result struct {
	IsSuccess bool        `json:"isSuccess"`
	Data      interface{} `json:"data"`
}
