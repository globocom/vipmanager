package model

import (
	"encoding/json"
	"log"
	"net/url"
)

type VipRequest struct {
	ExtendsSearch []map[string]string `json:"extends_search"`
}

func (v *VipRequest) Build(base string) (string, error) {
	j, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	pURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("search", string(j))
	pURL.Path += "api/v3/vip-request"
	pURL.RawQuery = params.Encode()
	log.Printf("Encoded URL: %v", params.Encode())
	return pURL.String(), nil
}

type Vips struct {
	Vips []Vip `json:"vips"`
}

type Vip struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Ipv4Id int    `json:"ipv4"`
	Ports  []Port `json:"ports"`
}

type Port struct {
	ID    int `json:"id"`
	Port  int16
	Pools []Pool `json:"pools"`
}

type Pool struct {
	ID           int `json:"id"`
	ServerPoolID int `json:"server_pool"`
}
