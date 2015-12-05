package lxdatabase_test

import (
	"github.com/layer-x/layerx-commons/lxdatabase"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/layer-x/layerx-mesos-tpi/test_helpers"
)

var _ = Describe("Lxdb", func() {
	test_helpers.StartETCD()

	BeforeEach(func() {
		Describe("initializes without error", func() {
			err := lxdatabase.Init([]string{"http://127.0.0.1:2379"})
			Expect(err).To(BeNil())
		})
	})
	Describe("lxdatabase.Set(key, val)", func() {
		It("sets keys", func() {
			err := lxdatabase.Set("foo", "bar")
			Expect(err).To(BeNil())
		})
	})
	Describe("lxdatabase.Get(key)", func() {
		It("gets values", func() {
			val, err := lxdatabase.Get("foo")
			Expect(err).To(BeNil())
			Expect(val).To(Equal("bar"))
		})
	})
	Describe("lxdatabase.Rm(key)", func() {
		It("deletes keys", func() {
			err := lxdatabase.Rm("foo")
			Expect(err).To(BeNil())
			val, err := lxdatabase.Get("foo")
			Expect(val).To(Equal(""))
			Expect(err).To(Not(BeNil()))
			Expect(err.Error()).To(ContainSubstring("Key not found"))
		})
	})
	Describe("lxdatabase.Mkdir(key)", func() {
		It("makes directories", func() {
			err := lxdatabase.Rmdir("/foo_dir", true)
			Expect(err).To(BeNil())
			err = lxdatabase.Mkdir("foo_dir")
			Expect(err).To(BeNil())
		})
	})
	Describe("lxdatabase.GetKeys(key)", func() {
		It("lists keys in directory", func() {
			keys, err := lxdatabase.GetKeys("foo_dir")
			Expect(err).To(BeNil())
			Expect(keys).To(BeEmpty())
			err = lxdatabase.Set("foo_dir/foo", "bar")
			Expect(err).To(BeNil())
			keys, err = lxdatabase.GetKeys("foo_dir")
			Expect(err).To(BeNil())
			Expect(keys).To(ContainElement("bar"))
		})
	})
	Describe("lxdatabase.GetSubdirectories(dir)", func() {
		It("lists sbudirs in directory", func() {
			dirs, err := lxdatabase.GetSubdirectories("foo_dir")
			Expect(err).To(BeNil())
			Expect(dirs).To(BeEmpty())
			err = lxdatabase.Mkdir("foo_dir/bar_dir")
			Expect(err).To(BeNil())
			dirs, err = lxdatabase.GetSubdirectories("foo_dir")
			Expect(err).To(BeNil())
			Expect(dirs).To(ContainElement("/foo_dir/bar_dir"))
		})
	})
	Describe("cleanup", func() {
		It("cleans up etcd", func() {
			test_helpers.CleanupETCD()
		})
	})
})
