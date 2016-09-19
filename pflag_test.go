package unicon_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/taybin/unicon"
)

var _ = Describe("PflagConfig", func() {
	var (
		err error
		cfg ReadableConfig
	)
	BeforeEach(func() {
		cfg = NewPflagConfig("test")
		err = cfg.Load()
	})
	It("Should load variables from commandline", func() {
		Expect(len(cfg.All()) >= 0).To(BeTrue())
		cfg2 := NewArgvConfig("")
		cfg2.Load()
		Expect(len(cfg2.All()) >= len(cfg.All())).To(BeTrue())
	})
})
