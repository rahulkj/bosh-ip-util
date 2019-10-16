package main

import (
	"fmt"
	"net"
  "os"
  "encoding/json"
  "log"
  "strings"
)

func compute(cloudConfigFile string, boshVMSOutputFile string, boshIP string) {
  cloudConfigRawData := GetRaw(cloudConfigFile)
  boshVMSRawData := GetRaw(boshVMSOutputFile)

  // Generic interface to read the file into
  var cloudConfig CloudConfig
  err := json.Unmarshal(cloudConfigRawData, &cloudConfig)
  if err != nil {
    fmt.Println("Error parsing JSON: ", err)
  }

  var boshVMS BoshVMs
  err = json.Unmarshal(boshVMSRawData, &boshVMS)
  if err != nil {
    fmt.Println("Error parsing JSON: ", err)
  }

	var outputs []Output

  for _,network := range cloudConfig.Networks {

    for _,subnet := range network.Subnets {
      ip, ipv4Net, err := net.ParseCIDR(subnet.Range)
      if err != nil {
        log.Fatal(err)
        os.Exit(1)
      }
      // fmt.Println(ip)
      // fmt.Println(ipv4Net)

      ips := getAllIPsInCIDR(ip, ipv4Net)

      totalIps := (len(ips) - 2)
      totalReservedIps := getTotalReservedIPs(subnet, ips)
      availableIPs := computeAvailableIPS(boshVMS, ipv4Net, cloudConfig, network.Name, totalIps, totalReservedIps, boshIP)

			o := Output {
				Network: network.Name,
				ReservedIPRange: subnet.Reserved,
				StaticIPRange: subnet.Static,
				CIDR: subnet.Range,
				TotalIPs: totalIps,
				TotalAvailableIPs: availableIPs,
				TotalReservedIPs: totalReservedIps,
			}

			outputs = append(outputs, o)
    }
  }

	output := Outputs {
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
    if ! strings.HasSuffix(ip.String(),".0") {
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

func getTotalReservedIPs(subnet Subnet, ips []string) int {
  totalReservedIps := 0
  for _,reserved := range subnet.Reserved {
    reservedIps := strings.Split(reserved, "-")
    if len(reservedIps) == 2 {
      startIP := reservedIps[0]
      endIP := reservedIps[1]
      // fmt.Println("Start of the Reserved IP is: ", startIP)
      // fmt.Println("End of the Reserved IP is: ", endIP)

      startIPIndex :=  0
      endIPIndex := 0
      for i := range ips {
        if startIP == ips[i] {
          startIPIndex = i
          } else if endIP == ips[i] {
            endIPIndex = i
          }
        }

        // fmt.Println("Start Index of the Reserved IP is: ", startIPIndex)
        // fmt.Println("End Index of the Reserved IP is: ", endIPIndex)
        // fmt.Println("reserved ips length: ", len(ips[startIPIndex:endIPIndex+1]))
        totalReservedIps += len(ips[startIPIndex:endIPIndex+1])
    } else if len(reservedIps) == 1 {
      totalReservedIps += 1
    }
  }

  return totalReservedIps
}

func computeAvailableIPS(boshVMS BoshVMs, ipv4Net *net.IPNet, cloudConfig CloudConfig, network string, totalIps int, totalReservedIps int, boshIP string) int {
  availableIPs := totalIps-totalReservedIps
  for _,table := range boshVMS.Tables {
    for _,row := range table.Rows {
      if ipv4Net.Contains(net.ParseIP(row.IPS)) {
        // fmt.Println("IP %s Belongs to subnet: %s", row.IPS, ipv4Net)
        availableIPs =  availableIPs - len(table.Rows)
        break
      }
    }
  }

  if ipv4Net.Contains(net.ParseIP(boshIP)) {
    availableIPs = availableIPs - 1
  }

  if network == cloudConfig.Compilation.Network {
    availableIPs = availableIPs - cloudConfig.Compilation.Workers
  }

  return availableIPs
}
