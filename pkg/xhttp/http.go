package xhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"regexp"
	"time"

	"miniapp/pkg/log"
	"miniapp/pkg/metrics"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var digitsRE = regexp.MustCompile(`\b\d+\b`)

type myClient struct {
	cli *http.Client
}

// Client http客户端接口
type Client interface {
	// DoGet Get方法, req表示query
	DoGet(ctx context.Context, url string, header map[string]string, query map[string]string) (resp *http.Response, err error)

	// DoPost Post方法，req表示body
	DoPost(ctx context.Context, url string, header map[string]string, body interface{}) (resp *http.Response, err error)
}

// NewClient 创建一个HTTP client
func NewClient(timeout time.Duration) Client {
	return &myClient{
		cli: &http.Client{
			Timeout: timeout,
		},
	}
}

// DoGet 发起Get请求
// 注意: 上层应用读取response之后，需要close
func (c *myClient) DoGet(ctx context.Context, url string, header map[string]string, query map[string]string) (resp *http.Response, err error) {
	tr := otel.Tracer("HTTP-Call")
	ctx, span := tr.Start(ctx, "do-get-http")
	defer span.End()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		onSpanError(span, err)
		return
	}

	// query参数
	for k, v := range query {
		req.URL.Query().Add(k, v)
	}

	// header头设置
	for k, v := range header {
		req.Header.Add(k, v)
	}

	// TODO(@x-lambda): 这里是在干啥?
	tracer := NewClientTrace(ctx)
	ctx = httptrace.WithClientTrace(ctx, tracer)
	req = req.WithContext(ctx)

	// trace.InjectHeader(span.Context(), req)

	start := time.Now()
	resp, err = c.cli.Do(req)
	duration := time.Since(start)

	// 只有err==nil时 resp不为nil
	status := http.StatusOK
	if err != nil {
		onSpanError(span, err)
		status = http.StatusGatewayTimeout
	} else {
		status = resp.StatusCode
	}

	log.Get(ctx).Debugf("[HTTP] Get url: %s, status: %d query: %s", url, status, req.URL.RawQuery)

	//span.SetAttributes(string(ext.Component), "http")
	//span.SetAttributes(string(ext.HTTPUrl), url)
	//span.SetAttributes(string(ext.HTTPMethod), req.Method)
	//span.SetAttributes(string(ext.HTTPStatusCode), status)

	// url 中带有的纯数字替换成 %d，不然 prometheus 就炸了
	// /v123/4/56/foo => /v123/%d/%d/foo
	url = digitsRE.ReplaceAllString(url, "%d")
	metrics.HTTPDurationSeconds.WithLabelValues(url, fmt.Sprint(status)).Observe(duration.Seconds())
	return
}

// DoPost 发起Post请求
// 注意: 上层应用读取response之后，需要close
func (c *myClient) DoPost(ctx context.Context, url string, header map[string]string, body interface{}) (resp *http.Response, err error) {
	tr := otel.Tracer("HTTP-Call")
	ctx, span := tr.Start(ctx, "do-post-http")
	defer span.End()

	b, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		onSpanError(span, err)
		return
	}

	// header头设置
	for k, v := range header {
		req.Header.Add(k, v)
	}

	req = req.WithContext(ctx)

	// trace.InjectHeader(span.Context(), req)

	start := time.Now()
	resp, err = c.cli.Do(req)
	duration := time.Since(start)

	// 只有err==nil时 response不为nil
	status := http.StatusOK
	if err != nil {
		onSpanError(span, err)
		status = http.StatusGatewayTimeout
	} else {
		status = resp.StatusCode
	}

	log.Get(ctx).Debugf("[HTTP] Get url: %s, status: %d query: %s", url, status, req.URL.RawQuery)

	//span.SetTag(string(ext.Component), "http")
	//span.SetTag(string(ext.HTTPUrl), url)
	//span.SetTag(string(ext.HTTPMethod), req.Method)
	//span.SetTag(string(ext.HTTPStatusCode), status)

	// url 中带有的纯数字替换成 %d，不然 prometheus 就炸了
	// /v123/4/56/foo => /v123/%d/%d/foo
	url = digitsRE.ReplaceAllString(url, "%d")
	metrics.HTTPDurationSeconds.WithLabelValues(url, fmt.Sprint(status)).Observe(duration.Seconds())
	return
}

// NewCommonHeader
func NewCommonHeader() (header map[string]string) {
	header = map[string]string{
		"Content-Type": "text/plain",
	}
	return
}

// onSpanError error信息记录到span上
func onSpanError(span oteltrace.Span, err error) {
	// span.SetTag(string(ext.Error), true)
	// span.LogKV(xlog.Error(err))
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
