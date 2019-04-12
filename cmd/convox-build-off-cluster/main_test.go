package main_test

import (
	"os"

	. "github.com/GetTerminus/convox-off-cluster-builder/cmd/convox-build-off-cluster"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {

	It("should call RunCommand", func() {
		Expect(RunCommand("ls -l", true)).NotTo(Equal(""))
	})

	It("should call GetGitHash", func() {
		Expect(GetGitHash()).NotTo(Equal(""))
	})

	It("should call GetRepo", func() {
		account := "1234567890"
		region := "us-east-2"
		Expect(GetRepo(&account, &region)).NotTo(Equal(""))
	})

	It("should call GetRegion", func() {
		os.Setenv("AWS_REGION", "us-east-1")
		Expect(GetRegion()).To(Equal("us-east-1"))
	})

	Describe("GenerateBuildJSONFile", func() {
		It("should error when called with a nil manifest", func() {
			_, err := GenerateBuildJSONFile(nil, "appname", "buildid", "description")
			Expect(err).NotTo(BeNil())
		})
	})

})
