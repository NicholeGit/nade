package trace

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
)

type TraceKey string

var ContextKey = TraceKey("trace-key")

type NadeTraceService struct {
	idService contract.IDService

	traceIDGenerator contract.IDService
	spanIDGenerator  contract.IDService
}

func NewNadeTraceService(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.IContainer)
	idService := c.MustMake(contract.IDKey).(contract.IDService)
	return &NadeTraceService{idService: idService}, nil
}

// WithTrace register new trace to context
func (t *NadeTraceService) WithTrace(c context.Context, trace *contract.TraceContext) context.Context {
	if ginC, ok := c.(*gin.Context); ok {
		ginC.Set(string(ContextKey), trace)
		return ginC
	} else {
		newC := context.WithValue(c, ContextKey, trace)
		return newC
	}
}

func (t *NadeTraceService) ToMap(tc *contract.TraceContext) map[string]string {
	m := map[string]string{}
	if tc == nil {
		return m
	}
	m[contract.TraceKeyTraceID] = tc.TraceID
	m[contract.TraceKeySpanID] = tc.SpanID
	m[contract.TraceKeyCspanID] = tc.CspanID
	m[contract.TraceKeyParentID] = tc.ParentID

	if tc.Annotation != nil {
		for k, v := range tc.Annotation {
			m[k] = v
		}
	}
	return m
}

func (t *NadeTraceService) GetTrace(ctx context.Context) *contract.TraceContext {
	if ctx == nil {
		return nil
	}
	if ginC, ok := ctx.(*gin.Context); ok {
		if val, ok2 := ginC.Get(string(ContextKey)); ok2 {
			return val.(*contract.TraceContext)
		}
	}

	if tc, ok := ctx.Value(ContextKey).(*contract.TraceContext); ok {
		return tc
	}
	return nil
}

func (t *NadeTraceService) NewTrace() *contract.TraceContext {
	var traceID, spanID string
	if t.traceIDGenerator != nil {
		traceID = t.traceIDGenerator.NewID()
	} else {
		traceID = t.idService.NewID()
	}

	if t.spanIDGenerator != nil {
		spanID = t.spanIDGenerator.NewID()
	} else {
		spanID = t.idService.NewID()
	}

	tc := &contract.TraceContext{
		TraceID:    traceID,
		ParentID:   "",
		SpanID:     spanID,
		CspanID:    "",
		Annotation: map[string]string{},
	}
	return tc
}

func (t *NadeTraceService) StartSpan(tc *contract.TraceContext) *contract.TraceContext {
	var childSpanID string
	if t.spanIDGenerator != nil {
		childSpanID = t.spanIDGenerator.NewID()
	} else {
		childSpanID = t.idService.NewID()
	}
	childSpan := &contract.TraceContext{
		TraceID:  tc.TraceID,
		ParentID: "",
		SpanID:   tc.SpanID,
		CspanID:  childSpanID,
		Annotation: map[string]string{
			contract.TraceKeyTime: time.Now().String(),
		},
	}
	return childSpan
}

func (t *NadeTraceService) ExtractHTTP(req *http.Request) *contract.TraceContext {
	tc := &contract.TraceContext{}
	tc.TraceID = req.Header.Get(contract.TraceKeyTraceID)
	tc.ParentID = req.Header.Get(contract.TraceKeySpanID)
	tc.SpanID = req.Header.Get(contract.TraceKeyCspanID)
	tc.CspanID = ""

	if tc.TraceID == "" {
		tc.TraceID = t.idService.NewID()
	}

	if tc.SpanID == "" {
		tc.SpanID = t.idService.NewID()
	}

	return tc
}

func (t *NadeTraceService) InjectHTTP(req *http.Request, tc *contract.TraceContext) *http.Request {
	req.Header.Add(contract.TraceKeyTraceID, tc.TraceID)
	req.Header.Add(contract.TraceKeySpanID, tc.SpanID)
	req.Header.Add(contract.TraceKeyCspanID, tc.CspanID)
	req.Header.Add(contract.TraceKeyParentID, tc.ParentID)
	return req
}
