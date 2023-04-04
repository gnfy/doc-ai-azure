package openai

import (
	"bytes"
	"doc-ai-azure/pkg/common"
	"io"
	"net/http"
)

// 用户花费token的数据
type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens int `json:"total_tokens"`
}

func Send(modelType string, reqBytes []byte) (body []byte, err error) {
	requrl := common.GlobalObject.OpenAi.Endpoint + "openai/deployments/" + modelType + "?api-version=" + common.GlobalObject.OpenAi.Apiversion  
	req, err := http.NewRequest(http.MethodPost, requrl, bytes.NewBuffer(reqBytes))
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", common.GlobalObject.OpenAi.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	return
}
