package model

type Ipv4Request struct {
	ExtendsSearch []Ipv4ExtendSearch `json:"extends_search"`
}

type Ipv4ExtendSearch struct {
	Oct1 int `json:"oct1"`
	Oct2 int `json:"oct2"`
	Oct3 int `json:"oct3"`
	Oct4 int `json:"oct4"`
}

type Ipv4Resp struct {
	Ips []Ipv4ResIp `json:"ips"`
}

type Ipv4ResIp struct {
	Id int `json:"id"`
}
