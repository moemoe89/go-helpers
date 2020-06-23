package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpOptions struct {
	Ctx     context.Context
	Url     string
	TO      *time.Duration
	Headers map[string]string
	Queries map[string]string
	Data    []byte
	Method  string
}

func DoRequest(opt *HttpOptions, rs interface{}) (int, error) {
	if opt.TO != nil {
		timeout := *opt.TO
		ctx, cancel := context.WithTimeout(opt.Ctx, timeout*time.Second)
		defer cancel()

		opt.Ctx = ctx
	}

	body := bytes.NewBuffer(opt.Data)
	defer body.Reset()

	req, err := http.NewRequestWithContext(opt.Ctx, opt.Method, opt.Url, body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}

	queryValues := req.URL.Query()
	for key, val := range opt.Queries {
		queryValues.Set(key, val)
	}
	req.URL.RawQuery = queryValues.Encode()

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	if rs == nil {
		return resp.StatusCode, nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = json.Unmarshal(respBody, rs)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return resp.StatusCode, nil
}
