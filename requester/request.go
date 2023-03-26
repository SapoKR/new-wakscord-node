package requester

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

func Request(key string, data any, retry int) int {
	req := fasthttp.AcquireRequest()

	req.SetRequestURI(fmt.Sprintf("%s%s", baseURL, key))

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")

	body, _ := json.Marshal(data)
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()

	err := client.Do(req, resp)
	if err != nil {
		if resp.StatusCode() == fasthttp.StatusTooManyRequests {
			retryAfter := fastjson.GetFloat64(resp.Body(), "retry_after")
			log.Printf("We are being rate limited. Webhook (%s) responded with 429. Retrying in %.2f seconds.", key[:35], retryAfter)
			time.Sleep(time.Duration(float64(time.Second) * retryAfter))
			if retry != 0 {
				return Request(key, data, retry-1)
			}
		}
	}

	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	return resp.StatusCode()
}
