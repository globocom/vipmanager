package model

import (
	"net/url"
	"strconv"
	"strings"
)

type ServerPools struct {
	Ids         []int         `json:"-"`
	ServerPools []*ServerPool `json:"server_pools"`
}

func (s *ServerPools) AddMember(ipID int, ipFormated string) {
	for _, p := range s.ServerPools {
		p.addMember(p.DefaultPort, ipID, ipFormated)
	}
}

func (s *ServerPools) RemMember(ipID int) {
	for _, p := range s.ServerPools {
		p.remMember(ipID)
	}
}

func (v *ServerPools) Build(base string) (string, error) {
	pURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	sIds := []string{}
	for _, id := range v.Ids {
		sIds = append(sIds, strconv.Itoa(id))
	}

	pURL.Path += "api/v3/pool/" + strings.Join(sIds, ";")
	return pURL.String(), nil
}

type ServerPool struct {
	ID                int                    `json:"id"`
	Identifier        string                 `json:"identifier"`
	DefaultPort       int                    `json:"default_port"`
	Environment       int                    `json:"environment"`
	ServiceDownAction map[string]interface{} `json:"servicedownaction"`
	LBMethod          string                 `json:"lb_method"`
	HealthCheck       map[string]interface{} `json:"healthcheck"`
	DefaultLimit      int                    `json:"default_limit"`
	Members           []ServerPoolMember     `json:"server_pool_members"`
	PoolCreated       bool                   `json:"pool_created"`
}

func (s *ServerPool) addMember(portReal int, ipID int, ipFormated string) {
	m := ServerPoolMember{
		PortReal:     portReal,
		Identifier:   ipFormated,
		Weight:       0,
		Priority:     0,
		Limit:        0,
		MemberStatus: 7,
		IP:           Ip{ipID, ipFormated},
	}
	s.Members = append(s.Members, m)
}

func (s *ServerPool) remMember(ipID int) {
	var remI int
	for i, m := range s.Members {
		if ipID == m.IP.ID {
			remI = i
			break
		}

	}

	s.Members = append(s.Members[0:remI], s.Members[remI+1:]...)
}

type ServerPoolMember struct {
	ID           *int   `json:"id"`
	Identifier   string `json:"identifier,omitempty"`
	PortReal     int    `json:"port_real"`
	Weight       int    `json:"weight"`
	Priority     int    `json:"priority"`
	Limit        int    `json:"limit"`
	MemberStatus int    `json:"member_status"`
	IP           Ip     `json:"ip"`
	IPV6         *Ip    `json:"ipv6"`
}

type Ip struct {
	ID         int    `json:"id,omitempty"`
	IPFormated string `json:"ip_formated"`
}

type ServerPoolPutRequest struct {
	IDs         []int
	ServerPools ServerPools
}

func (v *ServerPoolPutRequest) Build(base string) (string, error) {
	pURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	sIds := []string{}
	for _, id := range v.IDs {
		sIds = append(sIds, strconv.Itoa(id))
	}

	pURL.Path += "api/v3/pool/deploy/" + strings.Join(sIds, ";") + "/"
	return pURL.String(), nil
}
