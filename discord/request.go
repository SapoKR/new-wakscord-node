package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

func RequestFastHTTP(key string, data any, retry int) int {
	req := fasthttp.AcquireRequest()

	req.SetRequestURI(fmt.Sprintf("%s%s", baseURL, key))

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")

	body, _ := json.Marshal(data)
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()

	err := fasthttpClient.Do(req, resp)
	if err != nil {
		code := resp.StatusCode()
		if code == fasthttp.StatusTooManyRequests {
			retryAfter := fastjson.GetFloat64(resp.Body(), "retry_after")
			log.Printf("Webhook (%s) is being rate limited. Retrying in %.2f seconds.\n", key[:35], retryAfter)
			time.Sleep(time.Duration(float64(time.Second) * retryAfter))
			if retry != 0 {
				return RequestFastHTTP(key, data, retry-1)
			}
		} else {
			log.Printf("Uncaught error occurred. Status Code: %d and Detail: %v\n", code, err)
		}
	}

	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	return resp.StatusCode()
}

func RequestHTTP(key string, data any, retry int) int {
	body, _ := json.Marshal(data)

	req, _ := http.NewRequest(http.MethodPost, baseURL+key, bytes.NewReader(body))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		code := resp.StatusCode
		if code == http.StatusTooManyRequests {
			respBody, _ := io.ReadAll(resp.Body)

			retryAfter := fastjson.GetFloat64(respBody, "retry_after")
			log.Printf("Webhook (%s) is being rate limited. Retrying in %.2f seconds.\n", key[:35], retryAfter)
			time.Sleep(time.Duration(float64(time.Second) * retryAfter))
			if retry != 0 {
				return RequestHTTP(key, data, retry-1)
			}
		} else {
			log.Printf("Uncaught error occurred. Status Code: %d and Detail: %v\n", code, err)
		}
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
