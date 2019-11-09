package main

type Results struct {
	Result []Result
}

type Result struct {
	Network                         string
	SubnetName                      string
	AZs                             []string
	ReservedIPRange                 []string
	StaticIPRange                   []string
	CIDR                            string
	TotalIPs                        int
	TotalReservedIPs                int
	TotalIPsInUse                   int
	TotalAvailableIPs               int
	TotalIPsNeededForCompilationVMs int
}
