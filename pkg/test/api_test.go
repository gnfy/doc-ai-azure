package test

import (
	"context"
	"doc-ai-azure/pkg/api"
	"doc-ai-azure/pkg/common"
	"doc-ai-azure/pkg/httpclient"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func Test_chat(t *testing.T) {
	ctx := context.Background()
	reqList := []api.Message{
		{
			Text: "成人头疼是什么原因",
		},
		{
			Text: "成人流鼻涕是什么原因",
		},
		{
			Text: "成人流鼻咳嗽是什么原因",
		},
	}
	for _, item := range reqList {
		reqData, err := json.Marshal(item)
		if err != nil {
			t.Fatal(err)
		}
		reqParams := httpclient.HTTPRequestParam{
			Client:  newTesthttpClient(),
			ReqURL:  fmt.Sprintf("http://127.0.0.1:%s/chat", common.GlobalObject.Server.Port),
			In:      string(reqData),
			LogType: httpclient.LOGNONE_TYPE,
		}
		res, err := httpclient.HTTPPostWithString(ctx, reqParams)
		if err != nil {
			t.Fatal(err.Error())
		}
		if err = checkHttpResp(ctx, res); err != nil {
			t.Fatal(err.Error())
		}
	}
	<-time.After(time.Hour)
}
