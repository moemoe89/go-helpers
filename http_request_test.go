package helpers

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestDoRequestWithoutResp(t *testing.T) {
	ping := PingModel{"pong"}
	srv := HttpMock("/ping", http.StatusOK, ping)
	defer srv.Close()

	api := API{URL: srv.URL}

	headers := map[string]string{"key": "value"}
	queries := map[string]string{"key": "value"}

	to := time.Duration(10)
	opt := &HttpOptions{
		Ctx:     context.Background(),
		Url:     api.URL + "/ping",
		TO:      &to,
		Headers: headers,
		Queries: queries,
		Method:  http.MethodGet,
	}

	_, err := DoRequest(opt, nil)

	if err != nil {
		t.Error("expected", nil, "got", err.Error())
	}
}

func TestDoRequestErrCtx(t *testing.T) {
	ping := PingModel{"pong"}
	srv := HttpMock("/ping", http.StatusOK, ping)
	defer srv.Close()

	api := API{URL: srv.URL}

	opt := &HttpOptions{
		Url:    api.URL + "/ping",
		Method: http.MethodGet,
	}

	_, err := DoRequest(opt, nil)

	expt := "net/http: nil Context"
	if err.Error() != expt {
		t.Error("expected", expt, "got", err.Error())
	}
}

func TestDoRequestErrUnmarshal(t *testing.T) {
	ping := PingModel{"pong"}
	srv := HttpMock("/ping", http.StatusOK, ping)
	defer srv.Close()

	api := API{URL: srv.URL}
	rs := ""

	opt := &HttpOptions{
		Ctx:    context.Background(),
		Url:    api.URL + "/ping",
		Method: http.MethodGet,
	}

	_, err := DoRequest(opt, rs)

	expt := "json: Unmarshal(non-pointer string)"
	if err.Error() != expt {
		t.Error("expected", expt, "got", err.Error())
	}
}

func TestDoRequestErrEmptyURL(t *testing.T) {
	opt := &HttpOptions{
		Ctx:    context.Background(),
		Url:    "",
		Method: http.MethodGet,
	}

	_, err := DoRequest(opt, nil)

	expt := `Get : unsupported protocol scheme ""`
	if err.Error() != expt {
		t.Error("expected", expt, "got", err.Error())
	}
}
