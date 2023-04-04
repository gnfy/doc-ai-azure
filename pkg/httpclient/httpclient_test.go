package httpclient

import (
	"encoding/json"
	"fmt"
	"testing"
)

type ReviewResult struct {
	Code   int          `json:"code"`
	Msg    string       `json:"msg"`
	Status int          `json:"status"`
	Count  int          `json:"count"`
	Data   []ReviewData `json:"data"`
}
type ReviewData struct {
	Slug       string `json:"slug"`
	Rid        int    `json:"rid"`
	Rating     int    `json:"rating"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	Page       int    `json:"page"`
	PostedTime int    `json:"posted_time"`
	EditedTime int    `json:"edited_time"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func TestHttpClient(t *testing.T) {
	url := "https://api.parcelpanel.com/v1/shopify/review-list?slug=self-faq"
	res := Get(url, 1)
	if res.StatusCode == 0 {
		return
	}
	out := &ReviewResult{}
	fmt.Println(res.Body)
	err := json.Unmarshal([]byte(res.Body), out)
	if err != nil {
		return
	}
	for k, v := range out.Data {
		fmt.Println(k, v)
	}
}
