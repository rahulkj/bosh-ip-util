BOSH IP Utility
---

This tool allows one to inspect the cloud config, bosh vms output, to determine:
* Network names
* Subnet Names
* Availiability Zones
* CIDR defined
* Reserved IP Range/s
* State IP Range/s
* Total IP's in the CIDR
* Total Reserved IPs
* Total IP's in use
* Total Available IPs
* Total IP's needed for compilation vms


### How to Build

`go get -u github.com/rahulkj/bosh-ip-util`

For linux `GOOS=linux GOARCH=amd64 go build -o releases/bosh-ip-util-linux-amd64 github.com/rahulkj/bosh-ip-util`

For mac `GOOS=darwin GOARCH=amd64 go build -o releases/bosh-ip-util-darwin-amd64 github.com/rahulkj/bosh-ip-util`

For windows `GOOS=windows GOARCH=386 go build -o releases/bosh-ip-util-windows-amd64.exe github.com/rahulkj/bosh-ip-util`

### Usage

Export all the BOSH environment variables
```
export BOSH_ENVIRONMENT=
export BOSH_CLIENT=
export BOSH_CLIENT_SECRET=
export BOSH_CA_CERT=
```

Now run the tool
```
./bosh-ip-util | jq .
```

You will get a JSON output, that you can parse. Sample Data:
```
{
  "Result": [
    {
      "Network": "INFRASTRUCTURE",
      "SubnetName": "INFRASTRUCTURE",
      "AZs": [
        "MGMT-AZ",
        "AZ-1",
        "AZ-2"
      ],
      "ReservedIPRange": [
        "192.168.10.1-192.168.10.10"
      ],
      "StaticIPRange": [],
      "CIDR": "192.168.10.0/26",
      "TotalIPs": 62,
      "TotalReservedIPs": 10,
      "TotalIPsInUse": 1,
      "TotalAvailableIPs": 51,
      "TotalIPsNeededForCompilationVMs": 4
    },
    {
      "Network": "DEPLOYMENT",
      "SubnetName": "DEPLOYMENT",
      "AZs": [
        "AZ-1",
        "AZ-2"
      ],
      "ReservedIPRange": [
        "192.168.12.1-192.168.12.10"
      ],
      "StaticIPRange": [
        "192.168.12.245",
        "192.168.12.250",
        "192.168.12.252",
        "192.168.12.240"
      ],
      "CIDR": "192.168.12.0/23",
      "TotalIPs": 510,
      "TotalReservedIPs": 10,
      "TotalIPsInUse": 22,
      "TotalAvailableIPs": 478,
      "TotalIPsNeededForCompilationVMs": 0
    },
    {
      "Network": "SERVICES",
      "SubnetName": "SERVICES",
      "AZs": [
        "AZ-1",
        "AZ-2"
      ],
      "ReservedIPRange": [
        "192.168.14.1-192.168.14.10"
      ],
      "StaticIPRange": [],
      "CIDR": "192.168.14.0/23",
      "TotalIPs": 510,
      "TotalReservedIPs": 10,
      "TotalIPsInUse": 0,
      "TotalAvailableIPs": 500,
      "TotalIPsNeededForCompilationVMs": 0
    },
    {
      "Network": "PKS",
      "SubnetName": "PKS",
      "AZs": [
        "AZ-1",
        "AZ-2"
      ],
      "ReservedIPRange": [
        "192.168.16.1-192.168.16.10"
      ],
      "StaticIPRange": [
        "192.168.16.12"
      ],
      "CIDR": "192.168.16.0/23",
      "TotalIPs": 510,
      "TotalReservedIPs": 10,
      "TotalIPsInUse": 2,
      "TotalAvailableIPs": 498,
      "TotalIPsNeededForCompilationVMs": 0
    }
  ]
}
```

### Logic

```
Total IPs = Total IP's in CIDR - 2 (Gateway + Broadcast IP)
Reserved IPs = Summation of all the Reserved IP ranges defined for each subnet
Used IPs = Summation of all the IPs in a given subnet, by looking at the `bosh vms` output
Available IP = Total IPs - Reserved IPs - Used IPs - _(Bosh IP + Compilation IP)_
```

Hope you enjoy this utility!
