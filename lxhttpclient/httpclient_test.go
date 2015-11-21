package lxhttpclient_test

import (
	"github.com/layer-x/layerx-commons/lxhttpclient"

	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/mesosproto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
)

const (
	EXPECTED_GET_RESULT  = "EXPECTED_GET_RESULT"
	EXPECTED_POST_RESULT = "EXPECTED_POST_RESULT"
	ERROR                = "ERROR"
)

type fakeJSON struct {
	Value string `json:"value"`
}

func runTestServer() {
	m := martini.Classic()
	m.Get("/get", func() (int, string) {
		return 200, EXPECTED_GET_RESULT
	})
	m.Post("/post_json", func(res http.ResponseWriter, req *http.Request) (int, string) {
		body, _ := ioutil.ReadAll(req.Body)
		var fake fakeJSON
		err := json.Unmarshal(body, &fake)
		if err != nil {
			return 500, ERROR
		}
		return 202, EXPECTED_POST_RESULT
	})
	m.Post("/post_pb", func(res http.ResponseWriter, req *http.Request) (int, string) {
		body, _ := ioutil.ReadAll(req.Body)
		fake := &mesosproto.FrameworkID{}
		err := proto.Unmarshal(body, fake)
		if err != nil {
			return 500, ERROR
		}
		return 202, EXPECTED_POST_RESULT
	})
	m.RunOnAddr(":22334")
}

var _ = Describe("LXHttpclient", func() {
	//start a test martini server
	go runTestServer()
	It("can get", func() {
		resp, body, errs := lxhttpclient.Get("127.0.0.1:22334", "/get", nil)
		Expect(errs).To(BeNil())
		Expect(resp.StatusCode).To(Equal(200))
		Expect(string(body)).To(Equal(EXPECTED_GET_RESULT))
	})
	It("can post json", func() {
		jsonStruct := fakeJSON{
			Value: "fake_val",
		}
		resp, body, errs := lxhttpclient.Post("127.0.0.1:22334", "/post_json", nil, jsonStruct)
		Expect(errs).To(BeNil())
		Expect(resp.StatusCode).To(Equal(202))
		Expect(string(body)).To(Equal(EXPECTED_POST_RESULT))
	})
	It("can post pb", func() {
		fakeMessage := &mesosproto.FrameworkID{
			Value: proto.String("fake_value"),
		}
		resp, body, errs := lxhttpclient.Post("127.0.0.1:22334", "/post_pb", nil, fakeMessage)
		Expect(errs).To(BeNil())
		Expect(resp.StatusCode).To(Equal(202))
		Expect(string(body)).To(Equal(EXPECTED_POST_RESULT))
	})
})
