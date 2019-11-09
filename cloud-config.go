package main

type CloudConfig struct {
	Networks    []Network `json:"networks"`
	Compilation Compilation
}

type Network struct {
	Name    string   `json:"name"`
	Subnets []Subnet `json:"subnets"`
}

type Subnet struct {
	AZs             []string        `json:"azs"`
	Range           string          `json:"range"`
	Reserved        []string        `json:"reserved"`
	Static          []string        `json:"static"`
	CloudProperties CloudProperties `json:"cloud_properties"`
}

type Compilation struct {
	Network string `json:network`
	Workers int    `json:"workers"`
}

type CloudProperties struct {
	Name string `json:"name"`
}
