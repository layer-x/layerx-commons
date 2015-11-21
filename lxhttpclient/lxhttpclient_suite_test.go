package lxhttpclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLxhttpclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lxhttpclient Suite")
}
