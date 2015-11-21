package lxdatabase_test

import (
	"github.com/layer-x/layerx-commons/lxdatabase"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var binaryUrl string
var fileName string
var extract *exec.Cmd
var run *exec.Cmd

//not working, must run etcd manually first
func startETCD() {
	//start a test etcd server (will not work if etcd already running on host
	if runtime.GOOS == "darwin" {
		binaryUrl = "https://github.com/coreos/etcd/releases/download/v2.2.2/etcd-v2.2.2-darwin-amd64.zip"
		fileName = "etcd-v2.2.2-darwin-amd64.zip"
		extract = exec.Command("unzip", fileName, "-d", "etcd")
		run = exec.Command("etcd/etcd-v2.2.2-darwin-amd64/etcd")
	}
	if runtime.GOOS == "linux" {
		binaryUrl = "https://github.com/coreos/etcd/releases/download/v2.2.2/etcd-v2.2.2-linux-amd64.tar.gz"
		fileName = "etcd-v2.2.2-linux-amd64.tar.gz"
		exec.Command("mkdir", "etcd").Run()
		extract = exec.Command("tar", "xzvf", fileName, "-C", "etcd")
		run = exec.Command("etcd/etcd-v2.2.2-linux-amd64/etcd")
	}
	exec.Command("curl", "-L", binaryUrl, "-o", fileName).Run()
	extract.Run()
	go func() {
		run.Run()
	}()
	//5 seconds to initialize etcd
	time.Sleep(5 * time.Second)
}

var _ = Describe("Lxdb", func() {
	startETCD()

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
			err := lxdatabase.Mkdir("foo_dir")
			Expect(err).To(BeNil())
		})
	})
	Describe("lxdatabase.Ls(key)", func() {
		It("lists directories", func() {
			keys, err := lxdatabase.Ls("foo_dir")
			Expect(err).To(BeNil())
			Expect(keys).To(BeEmpty())
			err = lxdatabase.Set("foo_dir/foo", "bar")
			Expect(err).To(BeNil())
			keys, err = lxdatabase.Ls("foo_dir")
			Expect(err).To(BeNil())
			Expect(keys).To(ContainElement("bar"))
		})
	})
	Describe("cleanup", func() {
		It("cleans up etcd", func() {
			os.RemoveAll(fileName)
			os.RemoveAll("etcd")
			os.RemoveAll("default.etcd")
			exec.Command("pkill", "etcd")
		})
	})
})
