package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

type ServerPoolReq struct {
	Ids         []int         `json:"-"`
	ServerPools []*ServerPool `json:"server_pools"`
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
	UsersPermissions  map[string]interface{} `json:"users_permissions,omitempty"`
	GroupsPermissions map[string]interface{} `json:"groups_permissions,omitempty"`
}

type ServerPoolMember struct {
	ID           *int      `json:"id"`
	Identifier   string    `json:"identifier,omitempty"`
	PortReal     int       `json:"port_real"`
	Weight       int       `json:"weight"`
	Priority     int       `json:"priority"`
	Limit        int       `json:"limit"`
	MemberStatus int       `json:"member_status"`
	IP           MemberIp  `json:"ip"`
	IPV6         *MemberIp `json:"ipv6"`
}

type MemberIp struct {
	ID         int    `json:"id,omitempty"`
	IPFormated string `json:"ip_formated"`
}

func (s *ServerPoolReq) build(get bool) string {
	pURL := url.URL{}
	sIds := []string{}
	for _, id := range s.Ids {
		sIds = append(sIds, strconv.Itoa(id))
	}

	if get {
		pURL.Path += "pool/" + strings.Join(sIds, ";")
	} else {
		pURL.Path += "pool/deploy/" + strings.Join(sIds, ";") + "/"
	}
	return pURL.String()
}

func (s *ServerPoolReq) GET(vip Vip) error {
	log.Println("fetching Pools info")

	serverPoolIDs := []int{}
	for _, port := range vip.Ports {
		for _, pool := range port.Pools {
			serverPoolIDs = append(serverPoolIDs, pool.ServerPoolID)
		}
	}
	s.Ids = serverPoolIDs
	pURL := s.build(true)

	err := HTTP.Call("GET", pURL, nil, s)
	if err != nil {
		return err
	}
	if len(s.ServerPools) == 0 {
		return errors.New("No server pools were found")
	}

	if len(s.ServerPools) != len(serverPoolIDs) {
		return fmt.Errorf(
			"the number of ServerPoolIds: %d is differente from the returned in the response: %d",
			len(serverPoolIDs),
			len(s.ServerPools),
		)
	}
	s.print()
	return nil
}

func (s *ServerPoolReq) AddMember(ipID int, ipFormated string, dry bool) error {
	for _, p := range s.ServerPools {
		p.addMember(p.DefaultPort, ipID, ipFormated)
	}

	log.Print("Pool after adding")
	s.print()

	if dry {
		return nil
	}

	return s.store()
}

func (s *ServerPoolReq) RemMember(ipID int, dry bool) error {
	for _, p := range s.ServerPools {
		p.remMember(ipID)
	}

	log.Print("Pool after removing")
	s.print()

	if dry {
		return nil
	}

	return s.store()
}

func (s *ServerPoolReq) store() error {
	resp := &ServerPoolReq{Ids: s.Ids}
	pURL := s.build(false)
	j, _ := json.Marshal(s)
	fmt.Println(string(j))

	return HTTP.Call("PUT", pURL, s, resp)
}

// test member not found
// test member found
func (s *ServerPool) addMember(portReal int, ipID int, ipFormated string) {
	for _, m := range s.Members {
		if ipID == m.IP.ID {
			return
		}
	}

	m := ServerPoolMember{
		PortReal:     portReal,
		Identifier:   ipFormated,
		Weight:       0,
		Priority:     0,
		Limit:        0,
		MemberStatus: 7,
		IP:           MemberIp{ipID, ipFormated},
	}
	s.Members = append(s.Members, m)
}

// test member not found
// test member found
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

func (s *ServerPoolReq) print() {
	fmt.Println()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "Pool ID\tName\tDefault port\tMembers (ID/IP)\t\n")
	for _, p := range s.ServerPools {
		var members []string
		for _, p := range p.Members {
			members = append(members, strconv.Itoa(p.IP.ID)+"/"+p.IP.IPFormated)
		}
		fmt.Fprintf(w, "%d\t%s\t%d\t%s\t\n", p.ID, p.Identifier, p.DefaultPort, strings.Join(members, "; "))
	}
	w.Flush()
	fmt.Println()
}
