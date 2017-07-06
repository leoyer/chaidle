package helper

import (
	"github.com/tsuna/gohbase"
	"zrsf.com/hbase_client/types"
	"github.com/tsuna/gohbase/hrpc"
	"context"
	"strings"
	"github.com/tsuna/gohbase/filter"
	"io"
)


type HbaseClient interface{
	QueryByRowKey(HbaseQueryRequest)(types.StringMap,error)
	QueryRowKeys(HbaseQueryRequest) []types.StringMap
	QueryRowKeyExists(HbaseQueryRequest) bool
	QueryByRowKeyEncodeBASE64(HbaseQueryRequest) types.StringMap

	Scan()
	DoCoprocessor(request hrpc.ServiceRequest,response hrpc.ServiceResponse)[]hrpc.ServiceResponse

	SaveWithRowKey(HbaseSaveRequest) error
	SaveWithoutRowKey()
	SaveDecodeBASE64()
}

type HbaseRequest interface {
	GetTable() string
	GetRowKey() string
	GetFamily()  string
}

type HbaseQueryRequest interface{
	HbaseRequest
	SetRowKey(rowkey string )
	//只支持一次只操作一个族信息
	GetQulifier() []string
}

type HbaseSaveRequest interface {
	HbaseRequest
	GetQulifierValues() types.StringMap
}




var tsunaClient  gohbase.ExtendClient

type TsunaGoHbase struct {

}

func NewTsunaGoHbase () HbaseClient{
	return &TsunaGoHbase{}
}

func ( this TsunaGoHbase) getDefaultClient() gohbase.ExtendClient {
	if tsunaClient == nil {
		tsunaClient = gohbase.NewExtendClient(GetAppProperties().Read("hbase.zookeeper.quorum"))
	}
	return  tsunaClient
}

func ( this TsunaGoHbase)QueryByRowKey(request HbaseQueryRequest) (types.StringMap,error){
	defer  GetSeelog().Flush()
	GetSeelog().Info("@@@@@@@@@@@根据rowkey查询.[QueryByRowKey]....start....")
	GetSeelog().Trace("table="+request.GetTable())
	GetSeelog().Trace("rowkey="+request.GetRowKey())
	GetSeelog().Trace("family="+request.GetFamily())
	GetSeelog().Trace("qulifier=",request.GetQulifier())

	get,err := hrpc.NewGetStr(context.Background(),request.GetTable(),request.GetRowKey())
	if err != nil{
		GetSeelog().Error("QueryByRowKey.[NewGetStr] 根据rowkey查询错误,err="+err.Error())
		panic(err)
		return  nil,err
	}
	//设置列族以及列，这里每次只查询一个列族的信息
	get.SetFamilies(map[string][]string {request.GetFamily():request.GetQulifier()})
	result ,err := this.getDefaultClient().Get(get)
	if err != nil{
		GetSeelog().Error("QueryByRowKey.[Get]根据rowkey查询错误,err="+err.Error())
		panic(err)
		return  nil,err
	}
	resultMap := types.StringMap{}
	for _,cell := range result.Cells{
		resultMap[string(cell.Qualifier)] = string(cell.Value)
		GetSeelog().Trace("查询结果打印----->KEY="+string(cell.Qualifier)+",VAL="+string(cell.Value))
	}
	GetSeelog().Info("@@@@@@@@@@@根据rowkey查询.[QueryByRowKey]....end....")
	return  resultMap,nil
}

func ( this TsunaGoHbase)QueryRowKeys(request HbaseQueryRequest) []types.StringMap{
	defer  GetSeelog().Flush()
	GetSeelog().Info("@@@@@@@@@@@根据rowkeys查询.[QueryRowKeys]....start....")
	GetSeelog().Trace("rowkey="+request.GetRowKey())
	resultList := []types.StringMap{}
	for _,rowkey := range strings.Split(request.GetRowKey(),","){
		request.SetRowKey(rowkey)
		result ,err := this.QueryByRowKey(request)
		if err == nil {
			resultList = append(resultList,result)
		}
	}
	GetSeelog().Info("@@@@@@@@@@@根据rowkeys查询.[QueryRowKeys]....end....")
	return resultList
}
func ( this TsunaGoHbase)QueryRowKeyExists(request HbaseQueryRequest) bool{
	defer  GetSeelog().Flush()
	GetSeelog().Info("@@@@@@@@@@@根据rowkey查询.[QueryRowKeyExists]....start....")
	GetSeelog().Trace("table="+request.GetTable())
	GetSeelog().Trace("rowkey="+request.GetRowKey())
	GetSeelog().Trace("family="+request.GetFamily())
	GetSeelog().Trace("qulifier=",request.GetQulifier())

	get,err := hrpc.NewGetStr(context.Background(),request.GetTable(),request.GetRowKey())
	if err != nil{
		GetSeelog().Error("QueryRowKeyExists.[NewGetStr] 根据rowkey查询错误,err="+err.Error())
		panic(err)
		return  false
	}
	//设置列族以及列，这里每次只查询一个列族的信息
	get.SetFamilies(map[string][]string {request.GetFamily():request.GetQulifier()})
	get.SetFilter(filter.NewKeyOnlyFilter(true))
	result ,err := this.getDefaultClient().Get(get)
	if err != nil{
		GetSeelog().Error("QueryRowKeyExists.[Get]根据rowkey查询错误,err="+err.Error())
		panic(err)
		return  false
	}
	if len(result.Cells) > 0 {
		return  true
	}
	return false
}
func ( this TsunaGoHbase)QueryByRowKeyEncodeBASE64(HbaseQueryRequest) types.StringMap{
	return nil
}

func ( this TsunaGoHbase)Scan(){

}
func ( this TsunaGoHbase)DoCoprocessor(request hrpc.ServiceRequest,response hrpc.ServiceResponse)[]hrpc.ServiceResponse {
	defer  GetSeelog().Flush()
	GetSeelog().Info("@@@@@@@@@@@.[DoCoprocessor]....start....")
	GetSeelog().Trace("table="+request.GetTableName())
	GetSeelog().Trace("ServiceName="+request.GetServiceName())
	GetSeelog().Trace("MethodName="+request.GetMethodName())
	corp,err := hrpc.NewCoprocessor(context.Background(),request.GetTableName(),request)
	resultList := [] hrpc.ServiceResponse{}
	if err == nil{
		ite := this.getDefaultClient().ExecService(corp,response)
		result ,err :=   ite.Next()
		for ( err != io.EOF){
			resultList = append(resultList,result)
			result ,err  =   ite.Next()
		}
	}
	GetSeelog().Info("@@@@@@@@@@@.[DoCoprocessor]....end....")
	return  resultList
}

func ( this TsunaGoHbase)SaveWithRowKey(request HbaseSaveRequest) error{

	defer  GetSeelog().Flush()
	GetSeelog().Info("@@@@@@@@@@@插入.[SaveWithRowKey]....start....")
	GetSeelog().Trace("table="+request.GetTable())
	GetSeelog().Trace("rowkey="+request.GetRowKey())
	GetSeelog().Trace("family="+request.GetFamily())
	GetSeelog().Trace("GetQulifierValues=",request.GetQulifierValues())
	mapValue := map[string][]byte{}
	for key,val := range request.GetQulifierValues(){
		mapValue[key] = []byte(val)
	}
	saveValues := map[string]map[string][]byte{request.GetFamily():mapValue}
	put,err := hrpc.NewPutStr(context.Background(),request.GetTable(), request.GetRowKey(),saveValues)
	if err != nil{
		GetSeelog().Error("@@@@@@@@@@@插入SaveWithRowKey.[NewPutStr]异常,error="+err.Error())
		panic(err)
		return err
	}
	_ ,err = this.getDefaultClient().Put(put)

	if err != nil{
		GetSeelog().Error("@@@@@@@@@@@插入插入SaveWithRowKey.[Put]异常,error="+err.Error())
		panic(err)
		return err
	}
	return nil
}
func ( this TsunaGoHbase)SaveWithoutRowKey(){

}
func ( this TsunaGoHbase)SaveDecodeBASE64(){

}