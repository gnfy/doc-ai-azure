package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/ztrue/tracerr"
)

const (
	DefaultRetryNum = 6 // 默认重试次数
)

const (
	LOGNONE_TYPE = iota // 不输出日志
	LOGIN_TYPE          // in 日志
	LOGOUT_TYPE         // out 日志
	LOGALL_TYPE         // all 日志
)

var (
	ErrServUnable = errors.New("service unavailable")
)

// NewhttpClient 新建http客户端
func NewhttpClient() (client *http.Client) {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:    50,
			IdleConnTimeout: 60 * time.Second,
		},
	}
}

type HTTPRequestParam struct {
	Client   *http.Client      // 共用的http client
	ReqURL   string            // 请求地址
	In       interface{}       // 入参；withJson 时为可序列化为JSON的对象，withString时为字符串
	Dest     interface{}       // 出参；指针类型
	LogType  int               // 日志输出类型
	RetryNum int               // 重复次数；当==0时为DefaultRetryNum，只有返回是503时才会触发重试
	Header   map[string]string // Header
}

// HTTPGetWithJSON get请求，并把结果以JSON解码到dest
func HTTPGetWithJSON(ctx context.Context, param HTTPRequestParam) (err error) {
	if param.Client == nil {
		param.Client = NewhttpClient()
	}
	if param.LogType == LOGIN_TYPE || param.LogType == LOGALL_TYPE {
		log.Infof("HTTPGetWithJSONWithLog in=%s", param.ReqURL)
	}
	request, err := http.NewRequest(http.MethodGet, param.ReqURL, nil)
	if err != nil {
		return stack(err)
	}
	if param.Header != nil {
		for k, v := range param.Header {
			request.Header.Set(k, v)
		}
	}
	resp, err := param.Client.Do(request)
	if err != nil {
		return stack(err)
	}
	if resp.StatusCode == http.StatusServiceUnavailable {
		return stack(ErrServUnable)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		return stackf(resp.Status)
	}
	buff := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buff, resp.Body)
	if err != nil {
		return stack(err)
	}
	if param.LogType == LOGOUT_TYPE || param.LogType == LOGALL_TYPE {
		log.Infof("HTTPGetWithJSONWithLog out=%s", buff.Bytes())
	}
	err = json.NewDecoder(buff).Decode(param.Dest)
	return stack(err)
}

// HTTPGetWithJSONTryAgain 重试HTTPGetWithJSON请求
func HTTPGetWithJSONRetry(ctx context.Context, param HTTPRequestParam) (err error) {
	retryNum := param.RetryNum
	if retryNum == 0 {
		retryNum = DefaultRetryNum
	}
	for i := 0; i < retryNum; i++ {
		err = HTTPGetWithJSON(ctx, param)
		if err == nil {
			break
		}
		if errors.Is(err, ErrServUnable) {
			log.Infof("HTTPGetWithJSONTryAgain failed: %s", err)
			<-time.After(50 * time.Millisecond)
			continue
		}
		break
	}
	return stack(err)
}

// HTTPGetWithString get请求，结果直接返回字符串
func HTTPGetWithString(ctx context.Context, param HTTPRequestParam) (res string, err error) {
	if param.Client == nil {
		param.Client = NewhttpClient()
	}
	request, err := http.NewRequest(http.MethodGet, param.ReqURL, nil)
	if err != nil {
		return "", stack(err)
	}
	if param.Header != nil {
		for k, v := range param.Header {
			request.Header.Set(k, v)
		}
	}
	resp, err := param.Client.Do(request)
	if err != nil {
		return "", stack(err)
	}
	if resp.StatusCode == http.StatusServiceUnavailable {
		return "", stack(ErrServUnable)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		return "", stackf(resp.Status)
	}
	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", stack(err)
	}
	return string(resBytes), nil
}

// HTTPGetWithStringRetry 重试HTTPGetWithString请求
func HTTPGetWithStringRetry(ctx context.Context, param HTTPRequestParam) (res string, err error) {
	retryNum := param.RetryNum
	if retryNum == 0 {
		retryNum = DefaultRetryNum
	}
	for i := 0; i < retryNum; i++ {
		res, err = HTTPGetWithString(ctx, param)
		if err == nil {
			break
		}
		if errors.Is(err, ErrServUnable) {
			log.Infof("HTTPGetWithStringTryAgain failed: %s", err)
			<-time.After(50 * time.Millisecond)
			continue
		}
		break
	}
	return res, stack(err)
}

// HTTPPostWithJSON post请求。in入参结构体，dest出参结构体指针
func HTTPPostWithJSON(ctx context.Context, param HTTPRequestParam) (err error) {
	if param.Client == nil {
		param.Client = NewhttpClient()
	}
	buff := bytes.NewBuffer([]byte{})
	err = json.NewEncoder(buff).Encode(param.In)
	if err != nil {
		return stack(err)
	}
	if param.LogType == LOGIN_TYPE || param.LogType == LOGALL_TYPE {
		log.Infof("HTTPPostWithJSONWithLog requestURL=%s", param.ReqURL)
		log.Infof("HTTPPostWithJSONWithLog in=%s", buff.Bytes())
	}
	request, err := http.NewRequest(http.MethodPost, param.ReqURL, buff)
	if err != nil {
		return stack(err)
	}
	request.Header.Set("Content-Type", "application/json")
	if param.Header != nil {
		for k, v := range param.Header {
			request.Header.Set(k, v)
		}
	}
	resp, err := param.Client.Do(request)
	if err != nil {
		return stack(err)
	}
	if resp.StatusCode == http.StatusServiceUnavailable {
		return stack(ErrServUnable)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		return stackf(resp.Status)
	}
	buff.Reset()
	_, err = io.Copy(buff, resp.Body)
	if err != nil {
		return stack(err)
	}
	if param.LogType == LOGOUT_TYPE || param.LogType == LOGALL_TYPE {
		log.Infof("HTTPPostWithJSONWithLog out=%s", buff.Bytes())
	}
	err = json.NewDecoder(buff).Decode(param.Dest)
	return stack(err)
}

// HTTPPostWithJSONRetry 重试HTTPPostWithJSON请求
func HTTPPostWithJSONRetry(ctx context.Context, param HTTPRequestParam) (err error) {
	retryNum := param.RetryNum
	if retryNum == 0 {
		retryNum = DefaultRetryNum
	}
	for i := 0; i < retryNum; i++ {
		err = HTTPPostWithJSON(ctx, param)
		if err == nil {
			break
		}
		if errors.Is(err, ErrServUnable) {
			log.Infof("HTTPPostWithJSONTryAgain failed: %s", err)
			<-time.After(50 * time.Millisecond)
			continue
		}
		break
	}
	return stack(err)
}

// HTTPPostWithString post请求，结果直接返回字符串
func HTTPPostWithString(ctx context.Context, param HTTPRequestParam) (res string, err error) {
	if param.Client == nil {
		param.Client = NewhttpClient()
	}
	reqData, ok := param.In.(string)
	if !ok {
		return "", stack(errors.New("入参格式错误"))
	}
	if param.LogType == LOGIN_TYPE || param.LogType == LOGALL_TYPE {
		log.Infof("HTTPPostWithString request=%s", param.ReqURL)
		log.Infof("HTTPPostWithString in=%s", reqData)
	}
	buff := bytes.NewBufferString(reqData)
	request, err := http.NewRequest(http.MethodPost, param.ReqURL, buff)
	if err != nil {
		return "", stack(err)
	}
	request.Header.Set("Content-Type", "application/json")
	if param.Header != nil {
		for k, v := range param.Header {
			request.Header.Set(k, v)
		}
	}
	resp, err := param.Client.Do(request)
	if err != nil {
		return "", stack(err)
	}
	if resp.StatusCode == http.StatusServiceUnavailable {
		return "", stack(ErrServUnable)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		return "", stackf(resp.Status)
	}
	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", stack(err)
	}
	if param.LogType == LOGOUT_TYPE || param.LogType == LOGALL_TYPE {
		log.Infof("HTTPPostWithString out=%s", resBytes)
	}
	return string(resBytes), nil
}

// HTTPPostWithStringRetry 重试HTTPPostWithString请求
func HTTPPostWithStringRetry(ctx context.Context, param HTTPRequestParam) (res string, err error) {
	retryNum := param.RetryNum
	if retryNum == 0 {
		retryNum = DefaultRetryNum
	}
	for i := 0; i < retryNum; i++ {
		res, err = HTTPPostWithString(ctx, param)
		if err == nil {
			break
		}
		if errors.Is(err, ErrServUnable) {
			log.Infof("HTTPPostWithStringTryAgain failed: %s", err)
			<-time.After(50 * time.Millisecond)
			continue
		}
		break
	}
	return res, stack(err)
}

// Stack 在err中添加stack信息
func stack(err error) error {
	if err == nil {
		return nil
	}
	return tracerr.Wrap(err)
}

func stackf(format string, v ...interface{}) error {
	if format == "" {
		return nil
	}
	return tracerr.Wrap(fmt.Errorf(format, v...))
}
