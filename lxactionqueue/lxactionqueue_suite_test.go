package lxactionqueue_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLxactionqueue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lxactionqueue Suite")
}
