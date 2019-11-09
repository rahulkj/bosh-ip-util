package main

import (
	"encoding/json"
	"net"
	"strings"
)

// Get the Results json response, with all the relevant info
func compute(cloudConfigOutputJSON string, boshVMsOutput string, boshIP string) (string, error) {

	var cloudConfig CloudConfig
	err := json.Unmarshal([]byte(cloudConfigOutputJSON), &cloudConfig)
	if err != nil {
		return "", err
	}

	var boshVMS BoshVMs
	err = json.Unmarshal([]byte(boshVMsOutput), &boshVMS)
	if err != nil {
		return "", err
	}

	var results []Result

	for _, network := range cloudConfig.Networks {

		for _, subnet := range network.Subnets {
			ip, ipv4Net, _ := net.ParseCIDR(subnet.Range)

			ips := getAllIPsInCIDR(ip, ipv4Net)

			totalIps := (len(ips) - 2)
			reservedIPs, isBoshIPReserved := getTotalReservedIPs(subnet, ips, boshIP)
			availableIPs, totalIPsInUse := computeAvailableIPS(boshVMS, ipv4Net, totalIps,
				reservedIPs, boshIP, isBoshIPReserved)

			totalCompilationVMIPs := 0
			if network.Name == cloudConfig.Compilation.Network {
				totalCompilationVMIPs = cloudConfig.Compilation.Workers
			}

			o := Result{
				Network:                         network.Name,
				SubnetName:                      subnet.CloudProperties.Name,
				AZs:                             subnet.AZs,
				CIDR:                            subnet.Range,
				ReservedIPRange:                 subnet.Reserved,
				StaticIPRange:                   subnet.Static,
				TotalIPs:                        totalIps,
				TotalReservedIPs:                len(reservedIPs),
				TotalIPsInUse:                   totalIPsInUse,
				TotalAvailableIPs:               availableIPs,
				TotalIPsNeededForCompilationVMs: totalCompilationVMIPs,
			}

			results = append(results, o)
		}
	}

	result := Results{
		Result: results,
	}

	final, _ := json.Marshal(result)

	return string(final), nil
}

func getAllIPsInCIDR(ip net.IP, ipv4Net *net.IPNet) []string {

	var ips []string
	for ip := ip.Mask(ipv4Net.Mask); ipv4Net.Contains(ip); inc(ip) {
		if !strings.HasSuffix(ip.String(), ".0") {
			ips = append(ips, ip.String())
		}
	}

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

// Sum all the reserved ip/range that belong to the subnet of the given network
func getTotalReservedIPs(subnet Subnet, ips []string, boshIP string) ([]string, bool) {
	var reservedIPs []string
	totalReservedIps := 0
	isboshIPReserved := false
	for _, reserved := range subnet.Reserved {
		reservedIps := strings.Split(reserved, "-")
		if len(reservedIps) == 2 {
			startIP := reservedIps[0]
			endIP := reservedIps[1]

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

			slice1 := ips[startIPIndex : endIPIndex+1]

			reservedIPs = append(reservedIPs, slice1...)

			if (boshIP == startIP || boshIP == endIP) && (boshIPIndex >= startIPIndex && boshIPIndex <= endIPIndex) {
				isboshIPReserved = true
			}

			totalReservedIps += len(ips[startIPIndex : endIPIndex+1])
		} else if len(reservedIps) == 1 {
			totalReservedIps++
			reservedIPs = append(reservedIPs, reservedIps[0])

			if boshIP == reservedIps[0] {
				isboshIPReserved = true
			}
		}

	}

	return reservedIPs, isboshIPReserved
}

// If bosh director IP is part of reserved IP Ranges in the right network, then don't reduce the available IP's count
// If compilation VM's are part of the network, then subtract that number from available IP's count
// availableIPs = totalIps - totalReservedIps - ipsinuse - (above conditions if true)
func computeAvailableIPS(boshVMS BoshVMs, ipv4Net *net.IPNet, totalIps int, reservedIPs []string, boshIP string, isBoshIPReserved bool) (int, int) {
	usedIPsCount := 0
	availableIPs := totalIps - len(reservedIPs)
	for _, table := range boshVMS.Tables {
		for _, row := range table.Rows {
			if ipv4Net.Contains(net.ParseIP(row.IPS)) {
				if !contains(reservedIPs, row.IPS) {
					availableIPs = availableIPs - 1
					usedIPsCount++
				}
			}
		}
	}

	if ipv4Net.Contains(net.ParseIP(boshIP)) && !isBoshIPReserved {
		availableIPs = availableIPs - 1
		usedIPsCount++
	}

	return availableIPs, usedIPsCount
}

func contains(reservedIPs []string, ip string) bool {
	for _, reservedIP := range reservedIPs {
		if reservedIP == ip {
			return true
		}
	}
	return false
}
