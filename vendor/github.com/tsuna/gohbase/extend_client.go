package gohbase
import (
	"github.com/tsuna/gohbase/hrpc"
)

type ExtendClient interface{
	Client
	ExecService(request *hrpc.Coprocessor,response hrpc.ServiceResponse) hrpc.CoprocessorIte
}


type extendClient struct {
	*client
}

// NewClient creates a new HBase client.
func NewExtendClient(zkquorum string, options ...Option) ExtendClient {
	return  &extendClient{
		client:newClient(zkquorum, options...),
	}
}



func (c *extendClient) ExecService(request *hrpc.Coprocessor,response hrpc.ServiceResponse) hrpc.CoprocessorIte {
	return newCoprocessorIte(c,request,response)
}
