package main

type Outputs struct {
	Result []Output
}

type Output struct {
	Network           string
	ReservedIPRange   []string
	StaticIPRange     []string
	CIDR              string
	TotalIPs          int
	TotalAvailableIPs int
	TotalReservedIPs  int
}
