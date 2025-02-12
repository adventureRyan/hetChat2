package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 用于表示 JSON 响应的数据格式
type H struct {
	Code  int
	Msg   string
	Data  interface{}
	Rows  interface{}
	Total interface{}
}

// 通用响应函数 Resp
func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret)
}

// 列表响应函数 RespList
func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	// 设置相应头的内容类型为 JSON
	w.Header().Set("Content-Type", "application/json")

	// 设置 HTTP 的状态码为 200 OK
	w.WriteHeader(http.StatusOK)

	// 创建一个 H 类型的实例 h
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret)
}

func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}

func RespOK(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0, data, msg)
}

func RespOKList(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, 0, data, total)
}
