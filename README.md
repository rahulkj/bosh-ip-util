BOSH IP Utility
---

This tool allows one to inspect the cloud config, bosh vms output, to determine:
* Network names
* CIDR defined
* Reserved IP Range/s
* State IP Range/s
* Total IPs in the CIDR
* Total Reserved IPs
* Total Available IPs


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
      "ReservedIPRange": [
        "192.168.10.1-192.168.10.10"
      ],
      "StaticIPRange": [],
      "CIDR": "192.168.10.0/26",
      "TotalIPs": 61,
      "TotalAvailableIPs": 46,
      "TotalReservedIPs": 10
    },
    {
      "Network": "DEPLOYMENT",
      "ReservedIPRange": [
        "192.168.12.1-192.168.12.10"
      ],
      "StaticIPRange": [],
      "CIDR": "192.168.12.0/23",
      "TotalIPs": 508,
      "TotalAvailableIPs": 497,
      "TotalReservedIPs": 10
    },
    {
      "Network": "SERVICES",
      "ReservedIPRange": [
        "192.168.14.1-192.168.14.10",
        "192.168.14.100-192.168.14.110",
        "192.168.14.240-192.168.14.250",
        "192.168.15.255"
      ],
      "StaticIPRange": [],
      "CIDR": "192.168.14.0/23",
      "TotalIPs": 508,
      "TotalAvailableIPs": 475,
      "TotalReservedIPs": 33
    },
    {
      "Network": "PKS",
      "ReservedIPRange": [
        "192.168.16.1-192.168.16.10"
      ],
      "StaticIPRange": [],
      "CIDR": "192.168.16.0/23",
      "TotalIPs": 508,
      "TotalAvailableIPs": 498,
      "TotalReservedIPs": 10
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
