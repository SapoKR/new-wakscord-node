package discord

import (
	// "bytes"
	"encoding/json"
	"fmt"

	// "io"
	// "net/http"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

func RequestFastHTTP(key string, data any, retry int) Response {
	req := fasthttp.AcquireRequest()

	req.SetRequestURI(fmt.Sprintf("%s%s", baseURL, key))

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")

	body, _ := json.Marshal(data)
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()

	err := fasthttpClient.Do(req, resp)
	if err != nil {
		return Response{
			Key:   key,
			Error: err,
		}
	}

	code := resp.StatusCode()
	if code == fasthttp.StatusTooManyRequests {
		retryAfter := fastjson.GetFloat64(resp.Body(), "retry_after")
		fmt.Printf("Webhook (%s) is being rate limited. Retrying in %.2f seconds.\n", key[:35], retryAfter)
		time.Sleep(time.Duration(float64(time.Second) * retryAfter))
		if retry != 0 {
			return RequestFastHTTP(key, data, retry-1)
		}
	}

	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	return Response{
		Key:  key,
		Code: code,
		Body: string(resp.Body()),
	}
}

// func RequestHTTP(key string, data any, retry int) int {
// 	body, _ := json.Marshal(data)
// 	bodyReader := bytes.NewReader(body)

// 	req, _ := http.NewRequest(http.MethodPost, baseURL+key, bodyReader)

// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		fmt.Printf("Uncaught error occurred. Error: %s", err.Error())
// 		return 0
// 	}
// 	defer resp.Body.Close()

// 	code := resp.StatusCode
// 	respBody, _ := io.ReadAll(resp.Body)
// 	if code >= 400 && code != http.StatusNotFound {
// 		if code == http.StatusTooManyRequests {
// 			retryAfter := fastjson.GetFloat64(respBody, "retry_after")
// 			fmt.Printf("Webhook (%s) is being rate limited. Retrying in %.2f seconds.\n", key[:35], retryAfter)
// 			time.Sleep(time.Duration(float64(time.Second) * retryAfter))
// 			if retry != 0 {
// 				return RequestHTTP(key, data, retry-1)
// 			}
// 		} else {
// 			fmt.Printf("Uncaught error occurred. Status Code: %d and Body: %s\n", code, string(respBody))
// 		}
// 	}

// 	return resp.StatusCode
// }
