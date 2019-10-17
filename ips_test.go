package main

import (
	"encoding/json"
	"net"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestHandler(t *testing.T) {
	spec.Run(t, "Handler", testAvailableIPs, spec.Report(report.Terminal{}))
}

func testAvailableIPs(t *testing.T, when spec.G, it spec.S) {

	it.Before(func() {
		RegisterTestingT(t)
	})

	when("computeAvailableIPS", func() {
		it("reports the correct ips", func() {

			cloudConfig, boshVMS := getTestData()
			_, ipv4Net, _ := net.ParseCIDR(cloudConfig.Networks[0].Subnets[0].Range)
			networkName := cloudConfig.Networks[0].Name
			totalIps := 61
			totalReservedIps := 10
			boshIP := "192.168.10.11"
			isBoshIPReserved := false

			availableIPs := computeAvailableIPS(boshVMS, ipv4Net, cloudConfig, networkName, totalIps,
				totalReservedIps, boshIP, isBoshIPReserved)
			Expect(availableIPs).To(Equal(46))
		})
	})

	when("getTotalReservedIPs", func() {
		it("reports the correct reserved ips", func() {

			cloudConfig, _ := getTestData()
			ip, ipv4Net, _ := net.ParseCIDR(cloudConfig.Networks[0].Subnets[0].Range)
			boshIP := "192.168.10.11"
			ips := getAllIPsInCIDR(ip, ipv4Net)

			reservedIPs, isBoshInReservedRange := getTotalReservedIPs(cloudConfig.Networks[0].Subnets[0], ips, boshIP)
			Expect(reservedIPs).To(Equal(10))
			Expect(isBoshInReservedRange).To(Equal(false))
		})
	})

	when("getTotalReservedIPs", func() {
		it("reports the correct reserved ips with director IP in reserved range", func() {

			cloudConfig, _ := getTestData()
			ip, ipv4Net, _ := net.ParseCIDR(cloudConfig.Networks[0].Subnets[0].Range)
			boshIP := "192.168.10.10"
			ips := getAllIPsInCIDR(ip, ipv4Net)

			reservedIPs, isBoshInReservedRange := getTotalReservedIPs(cloudConfig.Networks[0].Subnets[0], ips, boshIP)
			Expect(reservedIPs).To(Equal(10))
			Expect(isBoshInReservedRange).To(Equal(true))
		})
	})

	when("getTotalReservedIPs", func() {
		it("reports the correct reserved ips with director IP in reserved range", func() {

			cloudConfig, _ := getTestData()
			ip, ipv4Net, _ := net.ParseCIDR(cloudConfig.Networks[0].Subnets[0].Range)
			boshIP := "192.168.10.10"
			ips := getAllIPsInCIDR(ip, ipv4Net)

			reservedIPs, isBoshInReservedRange := getTotalReservedIPs(cloudConfig.Networks[0].Subnets[0], ips, boshIP)
			Expect(reservedIPs).To(Equal(10))
			Expect(isBoshInReservedRange).To(Equal(true))
		})
	})

	when("getTotalReservedIPs", func() {
		it("reports the correct reserved ips multiple reserved ranges", func() {

			cloudConfig, _ := getTestData()
			ip, ipv4Net, _ := net.ParseCIDR(cloudConfig.Networks[2].Subnets[0].Range)
			boshIP := "192.168.10.11"
			ips := getAllIPsInCIDR(ip, ipv4Net)

			reservedIPs, isBoshInReservedRange := getTotalReservedIPs(cloudConfig.Networks[2].Subnets[0], ips, boshIP)
			Expect(reservedIPs).To(Equal(34))
			Expect(isBoshInReservedRange).To(Equal(false))
		})
	})

	when("getTotalReservedIPs", func() {
		it("reports the correct reserved ips multiple reserved ranges and includes BOSH Director IP", func() {

			cloudConfig, _ := getTestData()
			ip, ipv4Net, _ := net.ParseCIDR(cloudConfig.Networks[2].Subnets[0].Range)
			boshIP := "192.168.14.11"
			ips := getAllIPsInCIDR(ip, ipv4Net)

			reservedIPs, isBoshInReservedRange := getTotalReservedIPs(cloudConfig.Networks[2].Subnets[0], ips, boshIP)
			Expect(reservedIPs).To(Equal(34))
			Expect(isBoshInReservedRange).To(Equal(true))
		})
	})

	when("compute", func() {
		it("reports a json string with all the details", func() {

			boshIP := "192.168.10.11"

			result, _ := compute(getCloudConfigTestData(), getBoshVMSTestData(), boshIP)
			var results Results
			json.Unmarshal([]byte(result), &results)

			Expect(result).ShouldNot(BeEmpty())
			Expect(len(results.Result)).To(Equal(4))
			Expect(results.Result[0].Network).To(Equal("INFRASTRUCTURE"))
			Expect(results.Result[0].CIDR).To(Equal("192.168.10.0/26"))
			Expect(results.Result[0].TotalIPs).To(Equal(61))
			Expect(results.Result[0].TotalReservedIPs).To(Equal(10))
			Expect(results.Result[0].TotalAvailableIPs).To(Equal(46))

			Expect(results.Result[1].Network).To(Equal("DEPLOYMENT"))
			Expect(results.Result[1].CIDR).To(Equal("192.168.12.0/23"))
			Expect(results.Result[1].TotalIPs).To(Equal(508))
			Expect(results.Result[1].TotalReservedIPs).To(Equal(10))
			Expect(results.Result[1].TotalAvailableIPs).To(Equal(497))

			Expect(results.Result[2].Network).To(Equal("SERVICES"))
			Expect(results.Result[2].CIDR).To(Equal("192.168.14.0/23"))
			Expect(results.Result[2].TotalIPs).To(Equal(508))
			Expect(results.Result[2].TotalReservedIPs).To(Equal(34))
			Expect(results.Result[2].TotalAvailableIPs).To(Equal(474))

			Expect(results.Result[3].Network).To(Equal("PKS"))
			Expect(results.Result[3].CIDR).To(Equal("192.168.16.0/23"))
			Expect(results.Result[3].TotalIPs).To(Equal(508))
			Expect(results.Result[3].TotalReservedIPs).To(Equal(10))
			Expect(results.Result[3].TotalAvailableIPs).To(Equal(498))
		})
	})

	when("compute", func() {
		it("reports a error when empty cloud config is passed", func() {

			boshIP := "192.168.10.11"

			result, err := compute("", getBoshVMSTestData(), boshIP)

			Expect(err).ShouldNot(Equal(nil))
			Expect(result).Should(Equal(""))
		})
	})

	when("compute", func() {
		it("reports a error when empty bosh vms json is passed", func() {

			boshIP := "192.168.10.11"

			result, err := compute(getCloudConfigTestData(), "", boshIP)

			Expect(err).ShouldNot(Equal(nil))
			Expect(result).Should(Equal(""))
		})
	})

}

func getTestData() (CloudConfig, BoshVMs) {
	// Generic interface to read the file into
	cloudConfigRawData := getCloudConfigTestData()
	boshVMSRawData := getBoshVMSTestData()

	var cloudConfig CloudConfig
	json.Unmarshal([]byte(cloudConfigRawData), &cloudConfig)

	var boshVMS BoshVMs
	json.Unmarshal([]byte(boshVMSRawData), &boshVMS)

	return cloudConfig, boshVMS
}

func getCloudConfigTestData() string {
	cloudConfigFile, _ := filepath.Abs("./sample-data/cloud-config.json")

	cloudConfigRawData := GetRaw(cloudConfigFile)

	return string(cloudConfigRawData)
}

func getBoshVMSTestData() string {
	boshVMSOutputFile, _ := filepath.Abs("./sample-data/vms.json")

	boshVMSRawData := GetRaw(boshVMSOutputFile)

	return string(boshVMSRawData)
}
