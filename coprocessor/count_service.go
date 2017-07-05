package coprocessor
import (
	"zrsf.com/hbase_client/protob"
	"github.com/tsuna/gohbase/hrpc"
	"github.com/golang/protobuf/proto"
)

type CountServiceResponse struct {
	Count int64
}

func (c *CountServiceResponse) GetResponse(responseByte []byte) hrpc.ServiceResponse {
	res := protob.CountResponse{}
	rsp := &CountServiceResponse{
	}
	if err := proto.Unmarshal(responseByte,&res);err== nil{
		rsp.Count = *res.Count
		return rsp
	}
	return nil
}

func NewCountServiceResponse() *CountServiceResponse{
	return  &CountServiceResponse{}
}



type countServiceRequest struct {
	startKey   string
	endKey           string
	params           []*Params
	column           string
	defaultQualifier string
}

func NewCountServiceRequest(start ,end ,family,defaultQualifier  string , params    []*Params) *countServiceRequest{
	return  &countServiceRequest{
		startKey:start,
		endKey:end,
		params:params,
		column:family,
		defaultQualifier:defaultQualifier,
	}
}

// Name returns the name of this RPC call.
func (s *countServiceRequest) GetServiceName() string {
	return "Count"
}


func (s *countServiceRequest) GetMethodName() string {
	return "sendCountRequest"
}

func (s *countServiceRequest) GetRow() []byte{
	return []byte(s.startKey)
}

func (s *countServiceRequest) GetProtoRequest() ([]byte,error) {
	var cusParams []*protob.Params
	for _,param := range s.params{
		pa := &protob.Params{
			Key: &(param.Key),
			Value: &(param.Value),
			Type: &(param.Type),
		}
		cusParams = append(cusParams,pa)
	}
	requstMessage := protob.CountRequest{
		StartKey:&s.startKey,
		EndKey:&s.endKey,
		Params:cusParams,
		Column:&s.column,
		DefaultQualifier:&s.defaultQualifier,
	}

	return proto.Marshal(&requstMessage)
}

