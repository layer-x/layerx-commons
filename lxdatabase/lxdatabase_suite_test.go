package lxdatabase_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLxdatabase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lxdatabase Suite")
}
