package helper

const(

	CONTROLLER_HSAVE = "hsaveController"
	METHOD_SAVE = "save"
	METHOD_SAVE_INVOICE_INDEXS = "addInvoice"
	METHOD_SAVE_PDF = "savePdf"


	CONTROLLER_HDEL = "hdelController"
	METHOD_DEL_MEMBER_ENT_APPLY = "delMemberEntApply"
	METHOD_DEL_ENT_UPLOAD_RES = "delEntUploadRes"
	METHOD_DEL_INCOICE_PDF_IMG = "delInvoicePdfImg"
	METHOD_DEL_BY_KEYTABLE = "deleteByKeyAndTable"

	CONTROLLER_MONITOR = "monitorAction"

	CONTROLLER_HQUERY = "hqueryController"
	METHOD_QUERY_ONE = "queryOne"
	METHOD_QUERY_PDF = "queryPdf"
	METHOD_QUERY_MORE = "queryMore"
	METHOD_QUERY_COUNT = "queryCount"
	METHOD_QUERY_MONEY_COUNT = "queryMoneyCount"
	METHOD_QUERY_PAGEDATA = "queryPageData"
	METHOD_QUERY_STS_COUNT = "queryStatisticsCount"
	METHOD_QUERY_SUBACCOUNT_STS = "querySubAccountStatistics"

	//操作成功
	CODE_SUCCESS = "0000"
	//参数错误
	CODE_CSERROR = "1001"
	//网络错误
	CODE_WLERROR = "2001"
	//未查询到数据
	CODE_WCXDSJERROR = "3001"
	//不满足前置条件
	CODE_BMZQZTJERROR = "3002"
	//写入数据异常
	CODE_XRSJYCERROR = "3003"
	//查询数据异常
	CODE_CXSJYCERROR = "3004"
	//操作失败
	CODE_CZSBERROR = "3005"

	//服务过期
	CODE_SERVER_EXPIRED = "4001"

	//系统异常
	CODE_XTYCERROR = "9999"



	MONITOR_HSAVE_TYPE = "hsave"
	MONITOR_HQUERY_TYPE = "hquery"
)
