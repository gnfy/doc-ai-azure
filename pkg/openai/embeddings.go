package openai

import (
	"doc-ai-azure/pkg/common"
	"encoding/json"
)

const (
	// 这里需要自己的azure openai部署模型
	// todo 放进配置文件里
	embeddingsApi = "embedding-ada/embeddings"
)
const (
	TextEmbeddingAda002 = "text-embedding-ada-002"
)

type EmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

// embedding响应数据的结构体
type EmbeddingResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
}

type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// 获得embedding的数据
func SendEmbeddings(request EmbeddingRequest) (embeddingResponse EmbeddingResponse, err error) {
	var reqBytes []byte
	reqBytes, err = json.Marshal(request)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}

	body, err := Send(embeddingsApi, reqBytes)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}

	err = json.Unmarshal(body, &embeddingResponse)
	return
}
