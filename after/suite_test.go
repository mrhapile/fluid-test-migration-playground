package after_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAfter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "After Suite")
}
