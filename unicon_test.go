package unicon_test

import (
	. "github.com/ndeanNovetta/m-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/taybin/unicon"
	"time"
)

var _ = Describe("Unicon", func() {
	Describe("Config struct", func() {
		var cfg *Unicon
		BeforeEach(func() {
			cfg = NewConfig(nil)
		})
		Describe("config.Default", func() {
			It("Should automatically create memory config for defaults", func() {
				memconf := NewMemoryConfig()
				memconf.Set("a", "b")
				cfg.ResetDefaults(memconf.All())
				Expect(cfg.GetString("a")).To(Equal("b"))
			})
		})
		It("Should use memory store to set and get by default", func() {
			cfg.Set("test_a", "10")
			Expect(cfg.GetString("test_a")).Should(Equal(cfg.GetString("test_a")))
		})
		It("Should return nil when key is non-existing", func() {
			Expect(cfg.Get("some-key")).To(BeNil())
			Expect(cfg.GetString("some-key")).To(Equal(""))
			Expect(cfg.GetBool("some-key")).To(Equal(false))
			Expect(cfg.GetInt("some-key")).To(Equal(0))
			Expect(cfg.GetInt64("some-key")).To(Equal(int64(0)))
			Expect(cfg.GetFloat64("some-key")).To(Equal(0.0))
			Expect(cfg.GetTime("some-key")).To(BeTemporally("==", time.Time{}))
			Expect(cfg.GetDuration("some-key")).To(Equal(time.Duration(0)))
		})
		It("Should return and use Defaults", func() {
			cfg.SetDefault("test_var", "abc")
			Expect(cfg.GetString("test_var")).Should(Equal("abc"))
			cfg.Set("test_var", "bca")
			Expect(cfg.GetDefault("test_var")).Should(Equal("abc"), "Setting to memory should not override defaults")
			Expect(cfg.Get("test_var")).Should(Equal("bca"), "Set to config should set in memory and use it over defaults")
		})
		It("Should get and set with case insensitivity", func() {
			cfg.Set("A", 1)
			cfg.Set("a", 2)
			Expect(cfg.Get("A")).To(Equal(2))
			Expect(cfg.Get("a")).To(Equal(2))
		})
		It("Should reset everything else but Defaults() on reset", func() {
			cfg.SetDefault("test_var", "abc")
			Expect(cfg.GetDefault("test_var")).Should(Equal("abc"))
			cfg.Set("test_var", "bca")
			Expect(cfg.GetDefault("test_var")).Should(Equal("abc"), "Setting to memory should not override defaults")
			Expect(cfg.Get("test_var")).Should(Equal("bca"), "Set to config should set in memory and use it over defaults")
			cfg.Reset()
			Expect(cfg.Get("test_var")).Should(Equal("abc"), "Set to config should set in memory and use it over defaults")
		})

		It("Should load & save all relevant sources", func() {
			cfg.Use("json1", NewJSONConfig("./config_test_1.json"))
			cfg.Use("json2", NewJSONConfig("./config_test_2.json"))
			cfg.Use("json2").Set("asd", "123")
			cfg.Use("json1").Set("asd", "321")
			err := cfg.Save()
			Expect(err).ToNot(HaveOccurred())
			cfg.Reset()
			Expect(cfg.Use("json1").All()).To(HaveLen(0))
			err = cfg.Load()
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg.Use("json1").Get("asd")).To(Equal("321"))
			Expect(cfg.Use("json2").Get("asd")).To(Equal("123"))
		})

		It("Should return all values from all storages", func() {
			cfg.Use("mem1", NewMemoryConfig())
			cfg.Use("mem2", NewMemoryConfig())
			cfg.Set("asd", 123456)
			cfg.Use("mem1").Set("das", 654321)
			cfg.Use("mem2").Set("sad", 654321)
			i := 0
			for key, value := range cfg.All() {
				Expect(cfg.GetInt(key)).To(Equal(value))
				i++
			}
			Expect(i).To(Equal(3))
		})
		It("Should be able to use Config objects in the hierarchy", func() {
			cfg.Use("test", NewConfig(nil))
			cfg.Set("test_123", "321test")
			Expect(cfg.Use("test").Get("test_123")).To(BeNil())
		})
		It("should prefer using defaults deeper in hierarchy (reverse order to normal fetch.)", func() {
			deeper := NewConfig(nil)
			deeper.ResetDefaults(M{
				"test":  "123",
				"testb": "321",
			})
			cfg.Use("test", deeper)
			cfg.ResetDefaults(M{
				"test": "333",
			})
			Expect(cfg.GetString("test")).To(Equal("123"))
			Expect(cfg.GetString("testb")).To(Equal("321"))
			cfg.Set("testb", "1")
			Expect(cfg.GetString("testb")).To(Equal("1"))
		})
		It("Should allow focusing on a specific namespace", func() {
			cfg.Set("test", 1)
			cfg.Set("a.b", 2)
			cfg.Set("a.b.c", 3)
			sub := cfg.Sub("a")
			Expect(sub.Get("test")).To(BeNil())
			Expect(sub.GetInt("b")).To(Equal(2))
			Expect(sub.GetInt("b.c")).To(Equal(3))
			sub.Set("b", 5)
			Expect(cfg.GetInt("a.b")).To(Equal(5))
		})
		It("should properly set/get defaults in Sub'd config", func() {
			sub := cfg.Sub("a")
			sub.SetDefault("b", 1)
			Expect(sub.Get("b")).To(Equal(1))
		})
		It("should be case insensitive when Sub'ing", func() {
			cfg.Set("A.B", 2)
			cfg.Set("A.B.C", 3)
			sub := cfg.Sub("a")
			Expect(sub.GetInt("b")).To(Equal(2))
			Expect(sub.GetInt("b.c")).To(Equal(3))
			sub.Set("b", 5)
			Expect(cfg.GetInt("a.b")).To(Equal(5))
		})
		It("should allow sub-subbing", func() {
			cfg.Set("A.B.C", 3)
			suba := cfg.Sub("A")
			subb := suba.Sub("B")
			Expect(subb.GetInt("C")).To(Equal(3))
		})
		It("should support bulk setting", func() {
			cfg.Set("foo", "oldvalue")
			cfg.Set("baz", "fuzz")

			bulk := make(map[string]interface{})
			bulk["foo"] = "bar"
			bulk["yes"] = true
			cfg.BulkSet(bulk)

			Expect(cfg.GetString("baz")).To(Equal("fuzz"))
			Expect(cfg.GetString("foo")).To(Equal("bar"))
			Expect(cfg.GetBool("yes")).To(Equal(true))
		})
		It("Should unmarshal maps when they are set", func() {
			dict := make(map[string]interface{})
			dict["test"] = "123"
			cfg.Set("A", dict)
			Expect(cfg.GetString("A.test")).To(Equal("123"))
		})
		It("Should unmarshal arrays when they are set", func() {
			cfg.Set("A", []interface{}{"123", "321"})
			Expect(cfg.GetInt("A.length")).To(Equal(2))
			Expect(cfg.GetInt("A[0]")).To(Equal(123))
			Expect(cfg.GetInt("A[1]")).To(Equal(321))
		})
	})
})
