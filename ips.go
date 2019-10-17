package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func compute(cloudConfigOutputJSON string, boshVMsOutput string, boshIP string) {

	// Generic interface to read the file into
	var cloudConfig CloudConfig
	err := json.Unmarshal([]byte(cloudConfigOutputJSON), &cloudConfig)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
	}

	var boshVMS BoshVMs
	err = json.Unmarshal([]byte(boshVMsOutput), &boshVMS)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
	}

	var outputs []Output

	for _, network := range cloudConfig.Networks {

		for _, subnet := range network.Subnets {
			ip, ipv4Net, err := net.ParseCIDR(subnet.Range)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			// fmt.Println(ip)
			// fmt.Println(ipv4Net)

			ips := getAllIPsInCIDR(ip, ipv4Net)

			totalIps := (len(ips) - 2)
			totalReservedIps, isBoshIPReserved := getTotalReservedIPs(subnet, ips, boshIP)
			availableIPs := computeAvailableIPS(boshVMS, ipv4Net, cloudConfig, network.Name, totalIps,
				totalReservedIps, boshIP, isBoshIPReserved)

			o := Output{
				Network:           network.Name,
				ReservedIPRange:   subnet.Reserved,
				StaticIPRange:     subnet.Static,
				CIDR:              subnet.Range,
				TotalIPs:          totalIps,
				TotalAvailableIPs: availableIPs,
				TotalReservedIPs:  totalReservedIps,
			}

			outputs = append(outputs, o)
		}
	}

	output := Outputs{
		Result: outputs,
	}

	final, err1 := json.Marshal(output)
	if err1 != nil {
		fmt.Println("Error parsing JSON: ", err)
	}

	fmt.Println(string(final))
}

func getAllIPsInCIDR(ip net.IP, ipv4Net *net.IPNet) []string {

	var ips []string
	for ip := ip.Mask(ipv4Net.Mask); ipv4Net.Contains(ip); inc(ip) {
		if !strings.HasSuffix(ip.String(), ".0") {
			ips = append(ips, ip.String())
		}
	}

	// fmt.Println(ips)

	return ips
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func getTotalReservedIPs(subnet Subnet, ips []string, boshIP string) (int, bool) {
	totalReservedIps := 0
	isboshIPReserved := false
	for _, reserved := range subnet.Reserved {
		reservedIps := strings.Split(reserved, "-")
		if len(reservedIps) == 2 {
			startIP := reservedIps[0]
			endIP := reservedIps[1]
			// fmt.Println("Start of the Reserved IP is: ", startIP)
			// fmt.Println("End of the Reserved IP is: ", endIP)

			startIPIndex := 0
			endIPIndex := 0
			boshIPIndex := 0

			for i := range ips {
				if startIP == ips[i] {
					startIPIndex = i
				}

				if endIP == ips[i] {
					endIPIndex = i
				}

				if boshIP == ips[i] {
					boshIPIndex = i
				}
			}

			if (boshIP == startIP || boshIP == endIP) && (boshIPIndex >= startIPIndex && boshIPIndex <= endIPIndex) {
				isboshIPReserved = true
			}

			// fmt.Println("Start Index of the Reserved IP is: ", startIPIndex)
			// fmt.Println("End Index of the Reserved IP is: ", endIPIndex)
			// fmt.Println("reserved ips length: ", len(ips[startIPIndex:endIPIndex+1]))
			totalReservedIps += len(ips[startIPIndex : endIPIndex+1])
		} else if len(reservedIps) == 1 {
			totalReservedIps++

			if boshIP == reservedIps[0] {
				isboshIPReserved = true
			}
		}

	}

	return totalReservedIps, isboshIPReserved
}

func computeAvailableIPS(boshVMS BoshVMs, ipv4Net *net.IPNet, cloudConfig CloudConfig,
	network string, totalIps int, totalReservedIps int, boshIP string, isBoshIPReserved bool) int {
	availableIPs := totalIps - totalReservedIps
	for _, table := range boshVMS.Tables {
		for _, row := range table.Rows {
			if ipv4Net.Contains(net.ParseIP(row.IPS)) {
				// fmt.Println("IP %s Belongs to subnet: %s", row.IPS, ipv4Net)
				availableIPs = availableIPs - len(table.Rows)
				break
			}
		}
	}

	if network == cloudConfig.Compilation.Network {
		availableIPs = availableIPs - cloudConfig.Compilation.Workers
	}

	if ipv4Net.Contains(net.ParseIP(boshIP)) && !isBoshIPReserved {
		availableIPs = availableIPs - 1
	}

	return availableIPs
}
