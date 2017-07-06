package router

import (
	"github.com/emicklei/go-restful"
	"zrsf.com/hbase_client/helper"
	"io"
	"time"
	"zrsf.com/hbase_client/types"

)

func Route(ws *restful.WebService,config *helper.Config){
	queryControllerPath := config.Read("project.name") + "/" + helper.CONTROLLER_HQUERY + "/{method}"
	delControllerPath := config.Read("project.name") + "/" +helper.CONTROLLER_HDEL  + "/{method}"
	saveControllerPath := config.Read("project.name") + "/" +helper.CONTROLLER_HSAVE  + "/{method}"
	monitorControllerPath := config.Read("project.name") + "/" +helper.CONTROLLER_MONITOR

	defer  helper.GetSeelog().Flush()

	helper.GetSeelog().Info("加载GET服务==>"+queryControllerPath)
	ws.Route(ws.GET(queryControllerPath).Filter(filterRequest).To(queryController))
	helper.GetSeelog().Info("加载POST服务==>"+queryControllerPath)
	ws.Route(ws.POST(queryControllerPath).Filter(filterRequest).To(queryController))
	helper.GetSeelog().Info("{method}==>"+helper.METHOD_QUERY_ONE + "," +
		helper.METHOD_QUERY_MORE+ "," + helper.METHOD_QUERY_PDF+ "," +
		helper.METHOD_QUERY_PAGEDATA+ "," + helper.METHOD_QUERY_COUNT+ "," + helper.METHOD_QUERY_MONEY_COUNT+ "," +
		helper.METHOD_QUERY_STS_COUNT+ "," + helper.METHOD_QUERY_SUBACCOUNT_STS)

	helper.GetSeelog().Info("加载GET服务==>"+delControllerPath)
	ws.Route(ws.GET(delControllerPath).Filter(filterRequest).To(delController))
	helper.GetSeelog().Info("加载POST服务==>"+delControllerPath)
	ws.Route(ws.POST(delControllerPath).Filter(filterRequest).To(delController))
	helper.GetSeelog().Info("{method}==>"+helper.METHOD_DEL_INCOICE_PDF_IMG + "," +
		helper.METHOD_DEL_ENT_UPLOAD_RES+ "," + helper.METHOD_DEL_BY_KEYTABLE+ "," +helper.METHOD_DEL_MEMBER_ENT_APPLY)

	helper.GetSeelog().Info("加载GET服务==>"+saveControllerPath)
	ws.Route(ws.GET(saveControllerPath).Filter(filterRequest).To(saveController))
	helper.GetSeelog().Info("加载POST服务==>"+saveControllerPath)
	ws.Route(ws.POST(saveControllerPath).Filter(filterRequest).To(saveController))
	helper.GetSeelog().Info("{method}==>"+helper.METHOD_SAVE + "," +
		helper.METHOD_SAVE_INVOICE_INDEXS+ "," + helper.METHOD_SAVE_PDF)

	helper.GetSeelog().Info("加载GET服务==>"+monitorControllerPath)
	ws.Route(ws.GET(monitorControllerPath).Filter(filterRequest).To(monitorController))
	helper.GetSeelog().Info("加载POST服务==>"+monitorControllerPath)
	ws.Route(ws.POST(monitorControllerPath).Filter(filterRequest).To(monitorController))
}

func filterRequest(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	//TODO 时间限定 ip打印
	defer  helper.GetSeelog().Flush()

	helper.GetSeelog().Info("Request.RemoteAddr ==>"+req.Request.RemoteAddr)
	helper.GetSeelog().Info("Request.RequestURI ==>"+req.Request.RequestURI)
	helper.GetSeelog().Info("Request.Method  ==>"+req.Request.Method)
	helper.GetSeelog().Info("Request.Host  ==>"+req.Request.Host)
	helper.GetSeelog().Info("Request.Header.Referer  ==>"+req.Request.Header.Get("Referer"))
	helper.GetSeelog().Info("Request.Header.Proxy-Client-IP  ==>"+req.Request.Header.Get("Proxy-Client-IP"))
	helper.GetSeelog().Info("Request.Header.WL-Proxy-Client-IP  ==>"+req.Request.Header.Get("WL-Proxy-Client-IP"))
	helper.GetSeelog().Info("Request.Header.X-Forwarded-For  ==>"+req.Request.Header.Get("X-Forwarded-For"))
	helper.GetSeelog().Info("Request.Header.X-Real-IP  ==>"+req.Request.Header.Get("X-Real-IP"))
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	helper.GetSeelog().Info("当前系统时间 ==>"+nowTime)
	if  nowTime > "2018-06-30 00:00:00" {
		helper.GetSeelog().Error("产品已过试用期请联系服务公司")
		resp.WriteAsXml(&types.BaseResponse{
			Msg:"产品已过试用期请联系服务公司",
			RsCode:helper.CODE_SERVER_EXPIRED,
		})
		return
	}
	chain.ProcessFilter(req, resp)
}


func queryController(req *restful.Request, resp *restful.Response) {

	switch req.PathParameter("method") {
		case helper.METHOD_QUERY_ONE:
			result := types.QueryOneResponse{
				BaseResponse:types.BaseResponse{
					Msg:"查询成功",
					RsCode:helper.CODE_SUCCESS,
				},
				Qualifiers: types.StringMap{"a":"s"},
			}
		    resp.WriteAsXml(result)
		case helper.METHOD_QUERY_MORE:
			result := types.QueryMoreResponse{
				BaseResponse:types.BaseResponse{
					Msg:"查询成功",
					RsCode:helper.CODE_SUCCESS,
				},
				Qualifiers: []types.StringMap{{"a":"s"},{"sss":"ddsd"}},
			}
			resp.WriteAsXml(result)
		case helper.METHOD_QUERY_PDF:
			result := types.QueryOneResponse{
				BaseResponse:types.BaseResponse{
					Msg:"查询成功",
					RsCode:helper.CODE_SUCCESS,
				},
				Qualifiers: types.StringMap{"a":"s"},
			}
			resp.WriteAsXml(result)
		case helper.METHOD_QUERY_COUNT:
			result := types.QueryCountResponse{
				BaseResponse:types.BaseResponse{
					Msg:"查询成功",
					RsCode:helper.CODE_SUCCESS,
				},
				Count: 0,
			}
			resp.WriteAsXml(result)
		case helper.METHOD_QUERY_MONEY_COUNT:
			result := types.QueryMoneyCountResponse{
				BaseResponse:types.BaseResponse{
					Msg:"查询成功",
					RsCode:helper.CODE_SUCCESS,
				},
				Count: types.StringMap{"a":"s"},
			}
			resp.WriteAsXml(result)
		case helper.METHOD_QUERY_PAGEDATA:
			result := types.QueryPageDataResponse{
				QueryMoreResponse: types.QueryMoreResponse{
						BaseResponse: types.BaseResponse{
							Msg:    "查询成功",
							RsCode: helper.CODE_SUCCESS,
						},
						Qualifiers: []types.StringMap{{"a":"s"},{"sss":"ddsd"}},
				},
			}
			resp.WriteAsXml(result)
		case helper.METHOD_QUERY_STS_COUNT:
			result := types.QueryStsCountResponse{
				QueryMoreResponse: types.QueryMoreResponse{
					BaseResponse: types.BaseResponse{
						Msg:    "查询成功",
						RsCode: helper.CODE_SUCCESS,
					},
					Qualifiers: []types.StringMap{{"a":"s"},{"sss":"ddsd"}},
				},
			}
			resp.WriteAsXml(result)
		case helper.METHOD_QUERY_SUBACCOUNT_STS:
			result := types.QuerySubAccountStsResponse{
				QueryMoreResponse: types.QueryMoreResponse{
					BaseResponse: types.BaseResponse{
						Msg:    "查询成功",
						RsCode: helper.CODE_SUCCESS,
					},
					Qualifiers: []types.StringMap{{"a":"s"},{"sss":"ddsd"}},
				},
				Total :  types.StringMap{"a":"s"},
			}
			resp.WriteAsXml(result)
		default:
			helper.GetSeelog().Info("queryController---未知服务类型==>"+req.PathParameter("method"))
			resp.WriteAsXml(&types.BaseResponse{
				Msg:"未知服务类型",
				RsCode:helper.CODE_CSERROR,
			})
	}



}
func delController(req *restful.Request, resp *restful.Response) {
	result := &types.BaseResponse{
		Msg:"删除成功",
		RsCode:helper.CODE_SUCCESS,
	}

	switch req.PathParameter("method") {
		case helper.METHOD_DEL_BY_KEYTABLE:
		case helper.METHOD_DEL_ENT_UPLOAD_RES:
		case helper.METHOD_DEL_INCOICE_PDF_IMG:
		case helper.METHOD_DEL_MEMBER_ENT_APPLY:
		default:
			helper.GetSeelog().Info("delController---未知服务类型==>"+req.PathParameter("method"))
			result.Msg = "未知服务类型"
			result.RsCode = helper.CODE_CSERROR
	}

	resp.WriteAsXml(result)
}
func saveController(req *restful.Request, resp *restful.Response) {
	result := &types.BaseResponse{
		Msg:"保存成功",
		RsCode:helper.CODE_SUCCESS,
	}

	switch req.PathParameter("method") {
		case helper.METHOD_SAVE:

		case helper.METHOD_SAVE_PDF:
		case helper.METHOD_SAVE_INVOICE_INDEXS:
		default:
			helper.GetSeelog().Info("saveController---未知服务类型==>"+req.PathParameter("method"))
			result.Msg = "未知服务类型"
			result.RsCode = helper.CODE_CSERROR
	}
	resp.WriteAsXml(result)
}
func monitorController(req *restful.Request, resp *restful.Response) {

	switch helper.GetAppProperties().Read("monitor.type") {
		case helper.MONITOR_HSAVE_TYPE:
		case helper.MONITOR_HQUERY_TYPE:
	}
	io.WriteString(resp, "42")
}