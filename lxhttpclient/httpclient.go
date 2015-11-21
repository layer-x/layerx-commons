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
	if err != nil {
		return resp, emptyBytes, lxerrors.New("error reading get response", err)
	}

	return resp, respBytes, nil
}

func Post(url string, path string, headers map[string]string, message interface{}) (*http.Response, []byte, error) {
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
}

func postPB(url string, path string, headers map[string]string, pb proto.Message) (*http.Response, []byte, error) {
	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, emptyBytes, lxerrors.New("could not proto.Marshal mesasge", err)
	}
	fmt.Printf("posting pb: %s", data)
	return postData(url, path, headers, data)
}

func postJson(url string, path string, headers map[string]string, jsonStruct interface{}) (*http.Response, []byte, error) {
	//err has already been caught
	data, _ := json.Marshal(jsonStruct)
	fmt.Printf("posting json: %s", data)
	return postData(url, path, headers, data)
}

func postData(url string, path string, headers map[string]string, data []byte) (*http.Response, []byte, error) {
	completeURL := parseURL(url, path)
	request, err := http.NewRequest("POST", completeURL, bytes.NewReader(data))
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
	if err != nil {
		return resp, emptyBytes, lxerrors.New("error reading get response", err)
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
