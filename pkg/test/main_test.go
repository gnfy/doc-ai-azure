package test

import (
	"context"
	"doc-ai-azure/pkg/api"
	"doc-ai-azure/pkg/common"
	"encoding/json"
	"fmt"
	"middleproxy/pkg/errorx"
	"middleproxy/pkg/httpx"
	"middleproxy/pkg/netx"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/util/gconv"
)

func TestMain(m *testing.M) {
	go func() {
		r := api.SetupRouter()
		if err := r.Run(":" + common.GlobalObject.Server.Port); err != nil {
			fmt.Printf("startup service failed, err:%v\n", err)
		}
	}()
	if !netx.RetryCheckNetPort("tcp", "127.0.0.1", gconv.Int(common.GlobalObject.Server.Port), netx.DefaultCheckPortOption) {
		log.Fatal("server start failed")
	}
	// 运行测试用例
	exitCode := m.Run()
	// 退出
	os.Exit(exitCode)
}

// newTesthttpClient 新建http客户端
func newTesthttpClient() (client *http.Client) {
	return &http.Client{
		Timeout: 500 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 500 * time.Second,
			}).DialContext,
			MaxIdleConns:    50,
			IdleConnTimeout: 60 * time.Second,
		},
	}
}

// checkHttpResp 检查http返回结果
func checkHttpResp(ctx context.Context, in string) (err error) {
	if in == "" {
		return errorx.Stackf("resp data is empty")
	}
	resp := &httpx.Response{}
	err = json.Unmarshal([]byte(in), resp)
	if err != nil {
		return errorx.Stack(err)
	}
	if resp.Code != 0 {
		return errorx.Stackf(resp.Msg)
	}
	return nil
}
