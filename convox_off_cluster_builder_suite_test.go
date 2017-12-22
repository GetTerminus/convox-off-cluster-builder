package convox_off_cluster_builder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestConvoxOffClusterBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ConvoxOffClusterBuilder Suite")
}
