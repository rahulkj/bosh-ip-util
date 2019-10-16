package main

type BoshVMs struct {
	Tables []Table `json:"Tables"`
}

type Table struct {
	Rows []Row `json:"Rows"`
}

type Row struct {
	IPS string `json:"ips"`
}
