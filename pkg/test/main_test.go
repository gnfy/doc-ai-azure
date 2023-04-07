package test

import (
	"doc-ai-azure/pkg/api"
	"doc-ai-azure/pkg/common"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	go func() {
		r := api.SetupRouter()
		if err := r.Run(":" + common.GlobalObject.Server.Port); err != nil {
			fmt.Printf("startup service failed, err:%v\n", err)
		}
	}()
	<-time.After(time.Second)
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
