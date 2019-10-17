package main

type Results struct {
	Result []Result
}

type Result struct {
	Network           string
	ReservedIPRange   []string
	StaticIPRange     []string
	CIDR              string
	TotalIPs          int
	TotalAvailableIPs int
	TotalReservedIPs  int
}
