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

For linux `GOOS=linux GOARCH=amd64 go build -o releases/bosh-ip-util-linux-amd64 github.com/rahulkj/bosh-ip-util`

For mac `GOOS=darwin GOARCH=amd64 go build -o releases/bosh-ip-util-darwin-amd64 github.com/rahulkj/bosh-ip-util`

For windows `GOOS=windows GOARCH=386 go build -o releases/bosh-ip-util-windows-amd64.exe github.com/rahulkj/bosh-ip-util`

### Usage

```
./bosh-ip-util
   -b string
       Bosh Director IP
   -c string
       Cloud Config json file
   -v string
       Bosh VMS output json file
```

To fetch the cloud config in json format, please execute:
```
bosh cloud-config > cloud-config.yml

yq -r -f cloud-config.yml > cloud-config.json
```

**NOTE: Please download yq from https://github.com/mikefarah/yq/releases**

To fetch the bosh vms output in json format, please execute:
```
bosh vms --json >> bosh-vms-output.json
```

Now run the tool
```
/bosh-ip-util -c ~/some-folder/cloud-config.json -v ~/some-folder/bosh-vms-output.json -b BOSH_DIRECTOR_IP
```

You will get a JSON output, that you can parse. Sample Data:
```
/bosh-ip-util -c /Users/rjain/Documents/github/rahulkj/poc-bosh/cloud-config.json  -v /Users/rjain/Documents/github/rahulkj/poc-bosh/vms.json -b 192.168.10.11 | jq .
```

will produce:
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

Hope you enjoy this utility!
