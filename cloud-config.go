package main

type CloudConfig struct {
  Networks []Network `json:"networks"`
  Compilation Compilation
}

type Network struct {
  Name string `json:"name"`
  Subnets []Subnet `json:"subnets"`
}

type Subnet struct {
  Range string `json:"range"`
  Reserved []string `json:"reserved"`
  Static []string `json:"static"`
}

type Compilation struct {
  Network string `json:network`
  Workers int `json:"workers"`
}
