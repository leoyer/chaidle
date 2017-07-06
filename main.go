package main

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"zrsf.com/hbase_client/router"
	"zrsf.com/hbase_client/helper"
)

func main() {
	defer helper.GetSeelog().Flush()
	helper.GetSeelog().Info("开始启动服务.....")
	helper.GetSeelog().Info("加载app.properties配置start.....")
	config := helper.GetAppProperties()
	helper.GetSeelog().Info("加载app.properties配置end.....")
	helper.GetSeelog().Info("配置路由start.....")
	ws := new(restful.WebService)
	router.Route(ws,config)
	restful.Add(ws)
	helper.GetSeelog().Info("配置路由end.....")
	port := config.Read("project.port")
	helper.GetSeelog().Info("开启服务端口.....==>"+port)
	http.ListenAndServe(":"+port, nil)
}
