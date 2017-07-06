package hrpc

import (
	"github.com/tsuna/gohbase/pb"
	"github.com/golang/protobuf/proto"
	"context"
	"github.com/tsuna/gohbase/filter"
)

type Coprocessor struct {
	base
	corRquest ServiceRequest
}

func baseCoprocessor(ctx context.Context, table []byte,request ServiceRequest,
	options ...func(Call) error) (*Coprocessor, error) {
	s := &Coprocessor{
		base: base{
			table: table,
			ctx:   ctx,
		},
		corRquest:request,
	}
	err := applyOptions(s, options...)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func NewCoprocessor(ctx context.Context, table string,request ServiceRequest, options ...func(Call) error)(*Coprocessor, error){
	return baseCoprocessor(ctx,[]byte(table),request,options... )
}

func NewRangeCoprocessor(ctx context.Context, table string,startKey []byte ,request ServiceRequest, options ...func(Call) error)(*Coprocessor, error){
	cor,err  := baseCoprocessor(ctx,[]byte(table),request,options... )
	if err  != nil{
		return nil,err
	}
	cor.key =startKey

	return cor,nil
}


type CoprocessorIte interface {

	Next() (ServiceResponse, error)

	Close() error
}

type ServiceRequest interface{
	GetServiceName() string
	GetMethodName() string
	GetRow() []byte
	GetProtoRequest() ([]byte,error)
	GetStartRow() string
	GetEndRow() string
	SetRow([]byte)
	GetTableName() string

}
type ServiceResponse interface{
	GetResponse(responseByte []byte) ServiceResponse

}





// Name returns the name of this RPC call.
func (s *Coprocessor) Name() string {
	return "ExecService"
}

// ToProto converts this Scan into a protobuf message
func (s *Coprocessor) ToProto() (proto.Message, error) {
	coprocessor := &pb.CoprocessorServiceRequest{
		Region:       s.regionSpecifier(),
		Call: &pb.CoprocessorServiceCall{},
	}
	coprocessor.Call.Row = s.corRquest.GetRow()
	serviceName := s.corRquest.GetServiceName()
	coprocessor.Call.ServiceName = &serviceName
	methodName := s.corRquest.GetMethodName()
	coprocessor.Call.MethodName = &methodName


	if  requstBytes,err :=  s.corRquest.GetProtoRequest(); err == nil{
		coprocessor.Call.Request = requstBytes
	}else{
		return nil, err
	}
	return coprocessor, nil
}

// NewResponse creates an empty protobuf message to read the response
// of this RPC.
func (s *Coprocessor) NewResponse() proto.Message {
	return &pb.CoprocessorServiceResponse{}
}

// SetFamilies sets the families covered by this scanner.
func (s *Coprocessor) SetFamilies(fam map[string][]string) error {
	return nil
}

// SetFilter sets the request's filter.
func (s *Coprocessor) SetFilter(ft filter.Filter) error {
	return nil
}


// SetFilter sets the request's filter.
func (s *Coprocessor) GetCorRequest() ServiceRequest {
	return s.corRquest
}


// SetFilter sets the request's filter.
func (s *Coprocessor) SetCorRequest(request ServiceRequest) error {
	s.corRquest = request
	return nil
}

func (s *Coprocessor) StopRow() []byte {
	return []byte(s.corRquest.GetEndRow())
}

// StartRow returns the start key (inclusive) of this scanner.
func (s *Coprocessor) StartRow() []byte {
	return []byte(s.corRquest.GetStartRow())
}