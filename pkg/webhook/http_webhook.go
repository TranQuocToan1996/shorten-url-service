package webhook

import (
	"context"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

type httpWebhookClient struct {
	client *resty.Request
}

func NewRestyRequest() *resty.Request {
	return resty.New().
		SetTimeout(10 * time.Second).
		SetRetryCount(10).
		SetRetryWaitTime(100 * time.Millisecond).
		SetRetryMaxWaitTime(2 * time.Second).
		SetContentLength(true).
		SetDebug(false).
		NewRequest()
}

func NewHTTPWebhookClient() Client {
	client := NewRestyRequest()
	return &httpWebhookClient{
		client: client,
	}
}

func (w *httpWebhookClient) Notify(ctx context.Context, callbackURL string, payload any) error {
	if callbackURL == "" {
		return nil
	}

	resp, err := w.client.
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(payload).
		Post(callbackURL)
	if err != nil {
		log.Printf("httpWebhookClient, callbackURL [%v], payload [%v],err [%v]", callbackURL, payload, err)
	}

	log.Printf("httpWebhookClient, callbackURL [%v], payload [%v],callbackURL, payloadbody [%v], status [%v]", callbackURL, payload, resp.String(), resp.Status())
	return nil
}
