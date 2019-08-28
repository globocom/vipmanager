package model

import (
	"encoding/json"
	"errors"
	"net/url"
)

type VipRequest struct {
	ExtendsSearch []map[string]string `json:"extends_search"`
}

func (v *VipRequest) build() (string, error) {
	j, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	pURL := url.URL{}
	params := url.Values{}
	params.Add("search", string(j))
	pURL.Path += "vip-request"
	pURL.RawQuery = params.Encode()
	return pURL.String(), nil
}

type VipResp struct {
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

//GET returns at most one vip from Napi accordingly to the vip url
func (req *VipRequest) GET(name string) (Vip, error) {
	req.ExtendsSearch = []map[string]string{
		{"name": name},
	}

	pURL, err := req.build()
	if err != nil {
		return Vip{}, err
	}

	var vips VipResp
	err = HTTP.Call("GET", pURL, nil, &vips)
	if err != nil {
		return Vip{}, err
	}
	if len(vips.Vips) == 0 {
		return Vip{}, errors.New("Vip not found")
	}

	if len(vips.Vips) > 1 {
		return Vip{}, errors.New("more than one Vip was returned in the query. Aborting")
	}

	return vips.Vips[0], nil
}
