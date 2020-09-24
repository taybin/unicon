package unicon_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/taybin/unicon"
)

var _ = Describe("URLConfig", func() {
	var (
		cfg Config
		err error
	)
	BeforeEach(func() {
		cfg = NewConfig(nil)
	})
	JustBeforeEach(func() {
		cfg.Use("url", NewURLConfig(fmt.Sprintf("http://127.0.0.1:%d", HttpPort)))
	})

	It("Should load config from URL", func() {
		Expect(cfg).ToNot(BeNil())
		Expect(err).ToNot(HaveOccurred())
		err := cfg.Load()
		Expect(err).ToNot(HaveOccurred())
		Expect(cfg.Get("test")).To(Equal("abc"))
	})
})
