package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestConvoxBuildOffCluster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ConvoxBuildOffCluster Suite")
}
