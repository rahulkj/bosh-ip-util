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
			totalReservedIps, isBoshIPReserved := getTotalReservedIPs(subnet, ips, boshIP)
			availableIPs := computeAvailableIPS(boshVMS, ipv4Net, cloudConfig, network.Name, totalIps,
				totalReservedIps, boshIP, isBoshIPReserved)

			o := Result{
				Network:           network.Name,
				ReservedIPRange:   subnet.Reserved,
				StaticIPRange:     subnet.Static,
				CIDR:              subnet.Range,
				TotalIPs:          totalIps,
				TotalAvailableIPs: availableIPs,
				TotalReservedIPs:  totalReservedIps,
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
func getTotalReservedIPs(subnet Subnet, ips []string, boshIP string) (int, bool) {
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

			if (boshIP == startIP || boshIP == endIP) && (boshIPIndex >= startIPIndex && boshIPIndex <= endIPIndex) {
				isboshIPReserved = true
			}

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

// If bosh director IP is part of reserved IP Ranges in the right network, then don't reduce the available IP's count
// If compilation VM's are part of the network, then subtract that number from available IP's count
// availableIPs = totalIps - totalReservedIps - ipsinuse - (above conditions if true)
func computeAvailableIPS(boshVMS BoshVMs, ipv4Net *net.IPNet, cloudConfig CloudConfig,
	network string, totalIps int, totalReservedIps int, boshIP string, isBoshIPReserved bool) int {
	availableIPs := totalIps - totalReservedIps
	for _, table := range boshVMS.Tables {
		for _, row := range table.Rows {
			if ipv4Net.Contains(net.ParseIP(row.IPS)) {
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
