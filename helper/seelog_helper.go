package helper

import "github.com/cihub/seelog"

var logs seelog.LoggerInterface



func GetSeelog() seelog.LoggerInterface{
	if logs == nil{
		var err error
		logs,err = seelog.LoggerFromConfigAsFile("conf/seelog.xml")
		if err != nil {
			println("日志初始化错误",err.Error())
			panic(err)
		}
	}
	return  logs
}