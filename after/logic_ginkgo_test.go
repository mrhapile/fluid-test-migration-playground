package after_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/fluid-cloudnative/fluid-test-migration-playground/after"
)

var _ = Describe("Logic", func() {
	var (
		cfg *after.Config
	)

	BeforeEach(func() {
		// Common setup for most tests
		cfg = &after.Config{
			Name:      "test-dataset",
			Namespace: "default",
			Replicas:  1,
			MountPath: "/var/data",
			Options:   make(map[string]string),
		}
	})

	Describe("ValidateAndNormalize", func() {
		Context("Validation Rules", func() {
			It("should error when Name is empty", func() {
				cfg.Name = ""
				err := after.ValidateAndNormalize(cfg)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("name cannot be empty"))
			})

			It("should error when Replicas is negative", func() {
				cfg.Replicas = -1
				err := after.ValidateAndNormalize(cfg)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("replicas must be non-negative"))
			})

			It("should error when using temporary storage with multiple replicas", func() {
				cfg.Replicas = 2
				cfg.Options["storage"] = "temporary"
				err := after.ValidateAndNormalize(cfg)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("temporary storage does not support replicas > 1"))
			})

			It("should NOT error when using temporary storage with single replica", func() {
				cfg.Replicas = 1
				cfg.Options["storage"] = "temporary"
				err := after.ValidateAndNormalize(cfg)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Normalization Logic", func() {
			It("should set default MountPath if empty", func() {
				cfg.MountPath = ""
				err := after.ValidateAndNormalize(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(cfg.MountPath).To(Equal("/mnt/data"))
			})

			It("should normalize MountPath starting with double slashes", func() {
				cfg.MountPath = "//data/path"
				err := after.ValidateAndNormalize(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(cfg.MountPath).To(Equal("/data/path"))
			})

			It("should keep explicit valid MountPath unchanged", func() {
				cfg.MountPath = "/custom/path"
				err := after.ValidateAndNormalize(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(cfg.MountPath).To(Equal("/custom/path"))
			})
		})
	})

	Describe("BuildConnectionURL", func() {
		It("should generate a standard TCP connection string by default", func() {
			url := after.BuildConnectionURL(cfg)
			// expected: tcp://test-dataset.default:8080
			Expect(url).To(Equal("tcp://test-dataset.default:8080"))
		})

		It("should generate an SSL connection string when ssl option is true", func() {
			cfg.Options["ssl"] = "true"
			url := after.BuildConnectionURL(cfg)
			Expect(url).To(Equal("ssl://test-dataset.default:8080"))
		})

		It("should ignore other ssl values", func() {
			cfg.Options["ssl"] = "false"
			url := after.BuildConnectionURL(cfg)
			Expect(url).To(Equal("tcp://test-dataset.default:8080"))
		})
	})
})
