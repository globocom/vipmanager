package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"text/tabwriter"
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
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Business    string `json:"business"`
	Environment int    `json:"environmentvip"`
	Ipv4Id      int    `json:"ipv4"`
	Ports       []Port `json:"ports"`
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

	vips.Vips[0].print()
	return vips.Vips[0], nil
}

func (v *Vip) print() {
	portNumbers := []int16{}
	for _, x := range v.Ports {
		portNumbers = append(portNumbers, x.Port)
	}
	fmt.Println()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "ID\tBusiness\tEnvironment\tPorts\t\n")
	fmt.Fprintf(w, "%d\t%s\t%d\t%v\t\n\n", v.ID, v.Business, v.Environment, portNumbers)
	w.Flush()

}
