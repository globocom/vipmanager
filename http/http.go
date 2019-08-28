package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var napiURLQA = "http://networkapi.qa01.globoi.com/api/v3/"
var napiURLPROD = "https://networkapi.globoi.com/api/v3/"

//HTTP Wrapper
type HTTP struct {
	napiUser string
	napiPass string
	napiURL  string
}

//New Creates a new instance of HTTP
func New(env string, user string, pass string) HTTP {
	http := HTTP{napiUser: user, napiPass: pass}

	if env == "prod" {
		http.napiURL = napiURLPROD
	} else {
		http.napiURL = napiURLQA
	}

	return http
}

//Call Makes one http call, encode and decode the request body and response body
func (h *HTTP) Call(method string, suffix string, body interface{}, resp interface{}) error {
	var mBody []byte
	var err error
	if body != nil {
		mBody, err = json.Marshal(body)

		fmt.Println(string(mBody))
		if err != nil {
			return err
		}
	}

	rq, err := http.NewRequest(method, h.napiURL+suffix, bytes.NewBuffer(mBody))
	if err != nil {
		return err
	}
	rq.Header.Set("Content-Type", "application/json")
	rq.SetBasicAuth(h.napiUser, h.napiPass)

	client := &http.Client{}
	rs, err := client.Do(rq)

	if err != nil {
		return err
	}

	if rs.StatusCode != 200 && rs.StatusCode != 201 {
		s, _ := ioutil.ReadAll(rs.Body)
		log.Printf("Status code error: %d \n %s\n", rs.StatusCode, string(s))
		return fmt.Errorf("Status code: %d", rs.StatusCode)
	}

	dec := json.NewDecoder(rs.Body)
	return dec.Decode(resp)
}
