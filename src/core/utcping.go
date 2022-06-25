package main

import (
	"core/config"
	"core/tcping"
	"core/utils"
	"fmt"
)

func main() {
	//全局配置初始化
	err := config.GetGlobalFlag().Init()
	if nil != err {
		fmt.Println(err)
		return
	}
	//协程池配置初始化
	utils.GetRoutine().Init()

	//任务启动
	err = utcping.GetTCPing().Do(config.GetGlobalFlag())
	if nil != err {
		fmt.Println(err)
		return
	}
}
