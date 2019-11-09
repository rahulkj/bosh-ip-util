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
			totalIps := 61
			reservedIPs := []string{"192.168.10.1", "192.168.10.2", "192.168.10.3", "192.168.10.4", "192.168.10.5", "192.168.10.6", "192.168.10.7", "192.168.10.8", "192.168.10.9", "192.168.10.10"}
			boshIP := "192.168.10.11"
			isBoshIPReserved := false

			availableIPs, _ := computeAvailableIPS(boshVMS, ipv4Net, totalIps,
				reservedIPs, boshIP, isBoshIPReserved)
			Expect(availableIPs).To(Equal(50))
		})
	})

	when("getTotalReservedIPs", func() {
		it("reports the correct reserved ips", func() {

			cloudConfig, _ := getTestData()
			ip, ipv4Net, _ := net.ParseCIDR(cloudConfig.Networks[0].Subnets[0].Range)
			boshIP := "192.168.10.11"
			ips := getAllIPsInCIDR(ip, ipv4Net)

			reservedIPs, isBoshInReservedRange := getTotalReservedIPs(cloudConfig.Networks[0].Subnets[0], ips, boshIP)
			Expect(len(reservedIPs)).To(Equal(10))
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
			Expect(len(reservedIPs)).To(Equal(10))
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
			Expect(len(reservedIPs)).To(Equal(10))
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
			Expect(len(reservedIPs)).To(Equal(34))
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
			Expect(len(reservedIPs)).To(Equal(34))
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
			Expect(len(results.Result)).To(Equal(5))
			Expect(results.Result[0].Network).To(Equal("INFRASTRUCTURE"))
			Expect(results.Result[0].SubnetName).To(Equal("INFRASTRUCTURE"))
			Expect(len(results.Result[0].AZs)).To(Equal(3))
			Expect(results.Result[0].CIDR).To(Equal("192.168.10.0/26"))
			Expect(results.Result[0].TotalIPs).To(Equal(61))
			Expect(results.Result[0].TotalReservedIPs).To(Equal(10))
			Expect(results.Result[0].TotalIPsNeededForCompilationVMs).To(Equal(4))
			Expect(results.Result[0].TotalIPsInUse).To(Equal(1))
			Expect(results.Result[0].TotalAvailableIPs).To(Equal(50))

			Expect(results.Result[1].Network).To(Equal("DEPLOYMENT"))
			Expect(results.Result[1].SubnetName).To(Equal("DEPLOYMENT"))
			Expect(len(results.Result[1].AZs)).To(Equal(2))
			Expect(results.Result[1].CIDR).To(Equal("192.168.12.0/23"))
			Expect(results.Result[1].TotalIPs).To(Equal(508))
			Expect(results.Result[1].TotalReservedIPs).To(Equal(10))
			Expect(results.Result[1].TotalIPsNeededForCompilationVMs).To(Equal(0))
			Expect(results.Result[1].TotalIPsInUse).To(Equal(1))
			Expect(results.Result[1].TotalAvailableIPs).To(Equal(497))

			Expect(results.Result[2].Network).To(Equal("SERVICES"))
			Expect(results.Result[2].SubnetName).To(Equal("SERVICES"))
			Expect(len(results.Result[2].AZs)).To(Equal(2))
			Expect(results.Result[2].CIDR).To(Equal("192.168.14.0/23"))
			Expect(results.Result[2].TotalIPs).To(Equal(508))
			Expect(results.Result[2].TotalReservedIPs).To(Equal(34))
			Expect(results.Result[2].TotalIPsNeededForCompilationVMs).To(Equal(0))
			Expect(results.Result[2].TotalIPsInUse).To(Equal(0))
			Expect(results.Result[2].TotalAvailableIPs).To(Equal(474))

			Expect(results.Result[3].Network).To(Equal("PKS"))
			Expect(results.Result[3].SubnetName).To(Equal("PKS"))
			Expect(len(results.Result[3].AZs)).To(Equal(2))
			Expect(results.Result[3].CIDR).To(Equal("192.168.16.0/23"))
			Expect(results.Result[3].TotalIPs).To(Equal(508))
			Expect(results.Result[3].TotalReservedIPs).To(Equal(266))
			Expect(results.Result[3].TotalIPsNeededForCompilationVMs).To(Equal(0))
			Expect(results.Result[3].TotalIPsInUse).To(Equal(0))
			Expect(results.Result[3].TotalAvailableIPs).To(Equal(242))
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

	when("compute", func() {
		it("reports a json string with all the details", func() {

			boshIP := "192.168.10.11"

			result, _ := compute(getCloudConfigTestData2(), getBoshVMSTestData2(), boshIP)

			var results Results
			json.Unmarshal([]byte(result), &results)

			Expect(result).ShouldNot(BeEmpty())
			Expect(len(results.Result)).To(Equal(12))
			Expect(results.Result[0].Network).To(Equal("DEPLOYMENT"))
			Expect(results.Result[0].CIDR).To(Equal("10.135.133.0/24"))
			Expect(results.Result[0].TotalIPs).To(Equal(253))
			Expect(results.Result[0].TotalReservedIPs).To(Equal(137))
			Expect(results.Result[0].TotalAvailableIPs).To(Equal(85))
			Expect(results.Result[0].SubnetName).To(Equal("XY13711_HAT_10.135.133.0_24_v7"))
			Expect(results.Result[0].TotalIPsNeededForCompilationVMs).To(Equal(8))

			Expect(results.Result[3].Network).To(Equal("SERVICES"))
			Expect(results.Result[3].CIDR).To(Equal("10.135.133.0/24"))
			Expect(results.Result[3].TotalIPs).To(Equal(253))
			Expect(results.Result[3].TotalReservedIPs).To(Equal(184))
			Expect(results.Result[3].TotalAvailableIPs).To(Equal(69))
			Expect(results.Result[3].SubnetName).To(Equal("XY13711_HAT_10.135.133.0_24_v7"))
			Expect(results.Result[3].TotalIPsNeededForCompilationVMs).To(Equal(0))

			Expect(results.Result[6].Network).To(Equal("ABC-DEPLOYMENT"))
			Expect(results.Result[6].CIDR).To(Equal("10.135.133.0/24"))
			Expect(results.Result[6].TotalIPs).To(Equal(253))
			Expect(results.Result[6].TotalReservedIPs).To(Equal(236))
			Expect(results.Result[6].TotalAvailableIPs).To(Equal(16))
			Expect(results.Result[6].SubnetName).To(Equal("XY13711_HAT_10.135.133.0_24_v7"))
			Expect(results.Result[6].TotalIPsNeededForCompilationVMs).To(Equal(0))

			Expect(results.Result[9].Network).To(Equal("ABC-SERVICES"))
			Expect(results.Result[9].CIDR).To(Equal("10.135.133.0/24"))
			Expect(results.Result[9].TotalIPs).To(Equal(253))
			Expect(results.Result[9].TotalReservedIPs).To(Equal(235))
			Expect(results.Result[9].TotalAvailableIPs).To(Equal(18))
			Expect(results.Result[9].SubnetName).To(Equal("XY13711_HAT_10.135.133.0_24_v7"))
			Expect(results.Result[9].TotalIPsNeededForCompilationVMs).To(Equal(0))
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

func getCloudConfigTestData2() string {
	cloudConfigFile, _ := filepath.Abs("./sample-data/cloud-config2.json")

	cloudConfigRawData := GetRaw(cloudConfigFile)

	return string(cloudConfigRawData)
}

func getBoshVMSTestData2() string {
	boshVMSOutputFile, _ := filepath.Abs("./sample-data/vms2.json")

	boshVMSRawData := GetRaw(boshVMSOutputFile)

	return string(boshVMSRawData)
}
