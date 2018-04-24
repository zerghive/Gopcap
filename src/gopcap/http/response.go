package http

import (
	"reflect"
	"strconv"
)

var CLRF = []byte{13,10}
var SEP = []byte{58,32}

type HttpResponse struct {
	header 		*map[string]string
	version		*string
	stateCode	HttpStateCode
	contents	[]byte
	ContentType	string
}

// response 写入
func (rep *HttpResponse)Write(Data ...interface{}) {
	for argNum, arg := range Data {
		if argNum > 0 {
			rep.contents = append(rep.contents,' ')
		}
		if arg != nil {
			rep.contents = append(rep.contents, []byte(reflect.ValueOf(arg).String())...)
		} else {
			rep.contents = append(rep.contents,[]byte("nil")...)
		}
	}
}

// 设置response的首部
func (rep *HttpResponse)SetHeader(key string, value string) {
	(*rep.header)[key] = value
}

// response变成字节流
func (rep *HttpResponse)getBytes() []byte {
	(*rep.header)["Content-Length"] = strconv.Itoa(len(rep.contents))
	(*rep.header)["Conntion"] = "keep-alive"
	// 计算byte总共长度 防止append申请内存拷贝
	length := 38 + len(rep.contents)
	for key,value := range *rep.header {
		length += len(key) + len(value) + 4
	}
	// 申请固定capacity的内存
	buf := make([]byte, 0, length)
	buf = append(buf, []byte(*rep.version)...)
	buf = append(buf, 32)
	buf = append(buf, []byte(strconv.Itoa(int(rep.stateCode)))...)
	buf = append(buf, 32)
	buf = append(buf, []byte(getStateName(rep.stateCode))...)
	buf = append(buf, CLRF...)
	// header
	for key,value := range *rep.header {
		buf = append(buf, []byte(key)...)
		buf = append(buf, SEP...)
		buf = append(buf, []byte(value)...)
		buf = append(buf, CLRF...)
	}
	buf = append(buf, CLRF...)
	// content
	buf = append(buf, rep.contents...)
	return buf
}