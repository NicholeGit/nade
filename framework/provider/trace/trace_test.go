package trace_test

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/NicholeGit/nade/framework/contract"
	"github.com/NicholeGit/nade/tests"
)

func trace(msg string) func() {
	log.Printf("enter %s", msg)
	start := time.Now()
	return func() { log.Printf("exit %s (%s)", msg, time.Since(start)) }
}

func TestTrace(t *testing.T) {
	// 注意后面需要加 `()`
	defer trace("TestTrace run ")()
	container := tests.InitBaseContainer()
	trace := container.MustMake(contract.TraceKey).(contract.ITrace)
	tc := trace.NewTrace()
	t.Log(trace.ToMap(tc))

	t.Log("context 中使用")
	ctx := trace.WithTrace(context.Background(), tc)
	tc = trace.GetTrace(ctx)
	t.Log(trace.ToMap(tc))

	// =====
	t.Log("rpc 调用使用")
	request, _ := http.NewRequest("POST", "/", bytes.NewBufferString("foo=bar&bar=foo"))
	// 发送之前设置 cspan_id 子节点的span
	tc = trace.StartSpan(tc)
	request = trace.InjectHTTP(request, tc)
	t.Logf("发送前:\t%s\n", trace.ToMap(tc))
	// =====
	tc = trace.ExtractHTTP(request)
	t.Logf("接收后:\t%s\n", trace.ToMap(tc))

}
