package lxhttpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/layer-x/layerx-commons/lxerrors"
	"io/ioutil"
	"net/http"
	"strings"
)

var DefaultRetries = 5

type client struct {
	c *http.Client
}

func newClient() *client {
	return &client{
		c: http.DefaultClient,
	}
}

var emptyBytes []byte

func Get(url string, path string, headers map[string]string) (*http.Response, []byte, error) {
	return getWithRetries(url, path, headers, DefaultRetries)
}

func getWithRetries(url string, path string, headers map[string]string, retries int) (*http.Response, []byte, error) {
	resp, respBytes, err := func() (*http.Response, []byte, error) {
		completeURL := parseURL(url, path)
		request, err := http.NewRequest("GET", completeURL, nil)
		if err != nil {
			return nil, emptyBytes, lxerrors.New("error generating get request", err)
		}
		for key, value := range headers {
			request.Header.Add(key, value)
		}
		resp, err := newClient().c.Do(request)
		if err != nil {
			return resp, emptyBytes, lxerrors.New("error performing get request", err)
		}
		respBytes, err := ioutil.ReadAll(resp.Body)
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			return resp, emptyBytes, lxerrors.New("error reading get response", err)
		}

		return resp, respBytes, nil
	}()
	if err != nil && retries > 0 {
		return getWithRetries(url, path, headers, retries-1)
	}
	return resp, respBytes, err
}

func Post(url string, path string, headers map[string]string, message interface{}) (*http.Response, []byte, error) {
	return postWithRetries(url, path, headers, message, DefaultRetries)
}

func postWithRetries(url string, path string, headers map[string]string, message interface{}, retries int) (*http.Response, []byte, error) {
	resp, respBytes, err := func() (*http.Response, []byte, error) {
		switch message.(type) {
		case proto.Message:
			return postPB(url, path, headers, message.(proto.Message))
		default:
			_, err := json.Marshal(message)
			if err != nil {
				return nil, emptyBytes, lxerrors.New("message was not of expected type `json` or `protobuf`", err)
			}
			return postJson(url, path, headers, message)
		}
	}()
	if err != nil && retries > 0 {
		return postWithRetries(url, path, headers, message, retries-1)
	}
	return resp, respBytes, err
}

func postPB(url string, path string, headers map[string]string, pb proto.Message) (*http.Response, []byte, error) {
	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, emptyBytes, lxerrors.New("could not proto.Marshal mesasge", err)
	}
	return postData(url, path, headers, data)
}

func postJson(url string, path string, headers map[string]string, jsonStruct interface{}) (*http.Response, []byte, error) {
	//err has already been caught
	data, _ := json.Marshal(jsonStruct)
	return postData(url, path, headers, data)
}

func postData(url string, path string, headers map[string]string, data []byte) (*http.Response, []byte, error) {
	completeURL := parseURL(url, path)
	request, err := http.NewRequest("POST", completeURL, bytes.NewReader(data))
	if err != nil {
		return nil, emptyBytes, lxerrors.New("error generating post request", err)
	}
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	resp, err := newClient().c.Do(request)
	if err != nil {
		return resp, emptyBytes, lxerrors.New("error performing post request", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return resp, emptyBytes, lxerrors.New("error reading post response", err)
	}

	return resp, respBytes, nil
}

func parseURL(url string, path string) string {
	if !strings.HasPrefix(url, "http://") || !strings.HasPrefix(url, "https://") {
		url = fmt.Sprintf("http://%s", url)
	}
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	return fmt.Sprintf("%s/%s", url, path)
}
