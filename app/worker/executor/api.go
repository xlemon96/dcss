package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	_ TaskExecutor = APIData{}
)

type APIData struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	PayLoad string            `json:"payload"`
	Header  map[string]string `json:"header"`
}

func (a APIData) Type() string {
	return "api"
}

func (a APIData) Run(ctx context.Context) io.ReadCloser {
	pr, pw := io.Pipe()
	go func() {
		var exitCode = DefaultExitCode
		defer pw.Close()
		defer func() {
			now := time.Now().Local().Format("2006-01-02 15:04:05: ")
			pw.Write([]byte(fmt.Sprintf("\n%sRun Finished,Return Code:%5d", now, exitCode)))
		}()
		req, err := http.NewRequestWithContext(ctx, a.Method, a.URL, bytes.NewReader([]byte(a.PayLoad)))
		if err != nil {
			pw.Write([]byte(err.Error()))
			return
		}
		for k, v := range a.Header {
			req.Header.Add(k, v)
		}
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			log.Error("client Do failed", zap.Error(err))
			var customerr bytes.Buffer
			switch ctx.Err() {
			case context.DeadlineExceeded:
				customerr.WriteString(resp.GetMsg(resp.ErrCtxDeadlineExceeded))
			case context.Canceled:
				customerr.WriteString(resp.GetMsg(resp.ErrCtxCanceled))
			default:
				customerr.WriteString(err.Error())
			}
			pw.Write(customerr.Bytes())
			return
		}
		defer resp.Body.Close()
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		pw.Write(bs)
		if resp.StatusCode > 0 {
			exitCode = resp.StatusCode
		}
	}()
	return pr
}
