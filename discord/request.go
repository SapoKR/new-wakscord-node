package discord

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

func RequestFastHTTP(key string, data WebhookParams, retry int) Response {
	req := fasthttp.AcquireRequest()

	req.SetRequestURI(fmt.Sprintf("%s%s", baseURL, key))

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")

	reqBody, _ := json.Marshal(data)
	req.SetBody(reqBody)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := fasthttpClient.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		if retry != 0 {
			return RequestFastHTTP(key, data, retry-1)
		}
		return Response{
			Key:   key,
			Error: err,
		}
	}

	code := resp.StatusCode()
	respBody := resp.Body()
	fasthttp.ReleaseResponse(resp)
	if code == fasthttp.StatusTooManyRequests {
		retryAfter := fastjson.GetFloat64(respBody, "retry_after")
		time.Sleep(time.Duration(float64(time.Second) * retryAfter))
	}
	if code != fasthttp.StatusNoContent {
		if retry != 0 {
			return RequestFastHTTP(key, data, retry-1)
		}
	}

	return Response{
		Key:  key,
		Code: code,
		Body: string(respBody),
	}
}
