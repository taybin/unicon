package unicon_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
	. "github.com/taybin/unicon"
)

var _ = Describe("FlagSetConfig", func() {
	var (
		err error
		cfg ReadableConfig
	)
	BeforeEach(func() {
		fs := pflag.NewFlagSet("arguments", pflag.ContinueOnError)
		cfg = NewFlagSetConfig(fs, "")
		err = cfg.Load()
	})
	It("Should load variables from commandline", func() {
		Expect(len(cfg.All()) >= 0).To(BeTrue())
		fs := pflag.NewFlagSet("arguments", pflag.ContinueOnError)
		fs.Int("test", 1, "")
		cfg2 := NewFlagSetConfig(fs, "")
		cfg2.Load()
		Expect(len(cfg2.All()) >= len(cfg.All())).To(BeTrue())
	})
	It("Should remove prefix from arguments", func() {
		fs := pflag.NewFlagSet("arguments", pflag.ContinueOnError)
		fs.Int("test-a", 1, "")
		cfg2 := NewFlagSetConfig(fs, "test-")
		cfg2.Load()
		Expect(cfg2.GetInt("a")).To(Equal(1))
	})
	It("Should namespace arguments appropriately", func() {
		fs := pflag.NewFlagSet("arguments", pflag.ContinueOnError)
		fs.String("postgres-host", "localhost", "")
		fs.Int("postgres-port", 5432, "")
		cfg2 := NewFlagSetConfig(fs, "", "postgres")
		cfg2.Load()
		Expect(cfg2.Get("postgres.host")).To(Equal("localhost"))
		Expect(cfg2.GetInt("postgres.port")).To(Equal(5432))
	})
})
