package gohbase
import (
	"github.com/tsuna/gohbase/hrpc"
	"sync"
	"context"
	"io"
	"github.com/tsuna/gohbase/pb"
	"errors"
	"bytes"
)

type coprocessorRe struct{
	r *pb.NameBytesPair
	e error
}


type coprocessorIte struct {
	// fetcher's fileds shouldn't be accessed by scanner
	// TODO: maybe separate fetcher into a different package
	f       coprocessorFetcher
	once    sync.Once
	resultsCh chan coprocessorRe
	resultsM sync.Mutex
	results   *pb.NameBytesPair
	cancel  context.CancelFunc
}

func (s *coprocessorIte) peek() (*pb.NameBytesPair, error) {
	if s.f.ctx.Err() != nil {
		return nil, io.EOF
	}
	if s.results == nil   {
		re, ok := <-s.resultsCh
		if !ok {
			return nil, io.EOF
		}
		if re.e != nil {
			return nil, re.e
		}
		// fetcher never returns empty results
		s.results = re.r
	}
	return s.results, nil
}

func (s *coprocessorIte) shift() {
	if len(s.results.Value) == 0 {
		return
	}
	// set to nil so that GC isn't blocked to clean up the result
	s.results = nil
}

// coalesce combines result with partial if they belong to the same row
// and returns the coalesed result and whether coalescing happened
func (s *coprocessorIte) coalesce(result, partial *pb.NameBytesPair) (*pb.NameBytesPair, bool) {


	return partial, true
}


func newCoprocessorIte(c RPCClient, rpc *hrpc.Coprocessor,result hrpc.ServiceResponse) *coprocessorIte {
	ctx, cancel := context.WithCancel(rpc.Context())
	results := make(chan coprocessorRe)
	return &coprocessorIte{
		resultsCh: results,
		cancel:  cancel,
		f: coprocessorFetcher{
			RPCClient: c,
			rpc:       rpc,
			ctx:       ctx,
			resultsCh:   results,
			result:result,
			startRow: rpc.StartRow(),
		},
	}
}

func toLocalServiceResult(r *pb.NameBytesPair,respon hrpc.ServiceResponse) hrpc.ServiceResponse {
	if r == nil {
		return nil
	}
	return respon.GetResponse(r.Value)
}




func (s *coprocessorIte) Next() (hrpc.ServiceResponse, error) {
	s.once.Do(func() {
		go s.f.fetch()
	})
	s.resultsM.Lock()
	println("track Next() 1")

	var result, partial *pb.NameBytesPair
	var err error
	for {
		partial, err = s.peek()
		println("track partial ",partial)
		println("track partial ",err)
		if err == io.EOF && result != nil {
			// no more results, return what we have. Next call to the Next() will get EOF
			result.Value = []byte("")
			s.resultsM.Unlock()
			println("track Next() 2")
			return toLocalServiceResult(result,s.f.result), nil
		}
		if err != nil {

			// return whatever we have so far and the error
			s.resultsM.Unlock()
			println("track Next() 3",err.Error())
			return toLocalServiceResult(result,s.f.result), err
		}

		var done bool
		result, done = s.coalesce(result, partial)
		if done {
			s.shift()
		}

		s.resultsM.Unlock()
		println("track Next() 4")
		return  toLocalServiceResult(result,s.f.result), nil

	}

}

func (s *coprocessorIte) Close() error {
	s.cancel()
	return nil
}

type coprocessorFetcher struct {
	RPCClient
	resultsCh chan<- coprocessorRe
	// rpc is original scan query
	rpc *hrpc.Coprocessor
	ctx context.Context
	// result current result we are adding partials to
	startRow []byte
	result hrpc.ServiceResponse
}



func (f *coprocessorFetcher) trySend(rs  *pb.NameBytesPair, err error) bool {
	if err == nil && len(rs.Value) == 0 {
		return true
	}
	select {
	case <-f.ctx.Done():
		return true
	case f.resultsCh <- coprocessorRe{r: rs, e: err}:
		return false
	}
}


// fetch scans results from appropriate region, sends them to client and updates
// the fetcher for the next scan
func (f *coprocessorFetcher) fetch() {
	for {
		resp, region, err := f.next()
		if err != nil {
			if err != ErrDeadline {
				// if the context of the scan rpc wasn't cancelled (same as calling Close()),
				// return the error to client
				f.trySend(nil, err)
			}
			break
		}
		f.update( region)
		if f.trySend(resp.Value,nil ){
			break
		}
		// check whether we should close the scanner before making next request
		if f.shouldClose(resp, region) {
			break
		}
	}

	close(f.resultsCh)

}

func (f *coprocessorFetcher) next() (*pb.CoprocessorServiceResponse, hrpc.RegionInfo, error) {
	var rpc *hrpc.Coprocessor
	var err error

	rpc,err  = hrpc.NewRangeCoprocessor(f.ctx, string(f.rpc.Table()),f.startRow, f.rpc.GetCorRequest())


	res, err := f.SendRPC(rpc)
	if err != nil {
		return nil, nil, err
	}
	coprocessors, ok := res.(*pb.CoprocessorServiceResponse)
	if !ok {
		return nil, nil, errors.New("got non-ScanResponse for scan request")
	}
	return coprocessors, rpc.Region(), nil
}


// update updates the fetcher for the next scan request
func (f *coprocessorFetcher) update(region hrpc.RegionInfo) {
	f.rpc.GetCorRequest()

	f.startRow = region.StopKey()
}


// shouldClose check if this scanner should be closed and should stop fetching new results
func (f *coprocessorFetcher) shouldClose(resp *pb.CoprocessorServiceResponse, region hrpc.RegionInfo) bool {

	select {
	case <-f.ctx.Done():
		// scanner has been asked to close
		return true
	default:
	}
	// Check to see if this region is the last we should scan because:
	// (1) it's the last region
	if len(region.StopKey()) == 0 {
		return true
	}
	// (3) because its stop_key is greater than or equal to the stop_key of this scanner,
	// provided that (2) we're not trying to scan until the end of the table.
	return len(f.rpc.StopRow()) != 0 && // (2)
		bytes.Compare(f.rpc.StopRow(), region.StopKey()) <= 0 // (3)
}
