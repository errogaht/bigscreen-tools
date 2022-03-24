package bs

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type ApiError struct {
	Code    uint16
	Message string
}

func (bsRef *Bigscreen) request(url string, method string, headers map[string]string, body string) ([]byte, *http.Response) {
	context := *bsRef

	var body2 *strings.Reader
	if body == "" {
		body2 = &strings.Reader{}
	} else {
		body2 = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, body2)

	req.Header.Add("authorization", "Bearer "+context.Bearer)
	req.Header.Add("accept", "application/json")

	if context.JWT.Token != "" {
		req.Header.Add("x-access-token", context.JWT.Token)
	}
	for s, s2 := range headers {
		req.Header.Add(s, s2)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err.Error())
		}
	}(resp.Body)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 401 {
		panic(fmt.Sprintf("status code: %d, body: %s", resp.StatusCode, body))
	}

	return respBody, resp
}
