package main

import (
	"encoding/json"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestCommand(t *testing.T) {
	spec.Run(t, "command.go", testConvertToJson, spec.Report(report.Terminal{}))
}

func testConvertToJson(t *testing.T, when spec.G, it spec.S) {

	it.Before(func() {
		RegisterTestingT(t)
	})

	when("convertYmlToJSON", func() {
		it("reports valid json", func() {

			cloudConfigYmlString := getData()

			cloudConfigOutputJSON := convertYmlToJSON(cloudConfigYmlString)
			var cloudConfig CloudConfig
			err := json.Unmarshal([]byte(cloudConfigOutputJSON), &cloudConfig)
			if err != nil {
				Expect(err).Should(Equal(nil))
			}

			Expect(len(cloudConfig.Networks)).To(Equal(1))
		})
	})
}

func getData() string {
	cloudConfigFile, _ := filepath.Abs("./sample-data/cloud-config.yml")

	cloudConfigRawData := GetRaw(cloudConfigFile)

	return string(cloudConfigRawData)
}
