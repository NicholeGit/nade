package contract

import (
	"context"
	"net/http"
)

const TraceKey = "nade:trace"

const (
	TraceKeyTraceID  = "trace_id"
	TraceKeySpanID   = "span_id"
	TraceKeyCspanID  = "cspan_id"
	TraceKeyParentID = "parent_id"
	// TraceKeyMethod   = "method"
	// TraceKeyCaller   = "caller"
	TraceKeyTime = "span_time"
)

// TraceContext define struct according Google Dapper
type TraceContext struct {
	TraceID  string // traceID global unique
	ParentID string // 父节点SpanID
	SpanID   string // 当前节点SpanID
	CspanID  string // 子节点调用的SpanID, 由调用方指定

	Annotation map[string]string // 标记各种信息
}

type ITrace interface {
	// WithTrace register new trace to context
	WithTrace(c context.Context, trace *TraceContext) context.Context
	// GetTrace From trace context
	GetTrace(c context.Context) *TraceContext
	// NewTrace generate a new trace
	NewTrace() *TraceContext
	// StartSpan generate cspan for child call
	StartSpan(trace *TraceContext) *TraceContext

	// ExtractHTTP Get Trace By Http
	ExtractHTTP(req *http.Request) *TraceContext
	// InjectHTTP Set Trace to Http
	InjectHTTP(req *http.Request, trace *TraceContext) *http.Request

	// ToMap traceContext to map for logger
	ToMap(tc *TraceContext) map[string]string
}
