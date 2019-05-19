package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/likexian/whois-go"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type DomainInfo struct {
	Servers          []ServerDesc `json:"servers"`
	ServersChanged   string       `json:"servers_changed"`
	ServersSSLGrade  string       `json:"ssl_grade"`
	PreviousSSLGrade string       `json:"previous_ssl_grade"`
	Logo             string       `json:"logo"`
	Title            string       `json:"title"`
	IsDown           bool         `json:"is_down"`
}

type ServerDesc struct {
	ServerAddress string `json:"address"`
	SSLGrade      string `json:"ssl-grade"`
	Country       string `json:"country"`
	Owner         string `json:"owner"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		var domainInfo = retrieveDomainInfo()
		respondwithJSON(w, 200, domainInfo)
	})
	http.ListenAndServe(":3334", r)
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func retrieveDomainInfo() DomainInfo {
	resp, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=truora.com")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var domainInfo = createServerInfo(string(body))
	return domainInfo
}

func createServerInfo(serverBody string) DomainInfo {

	var result map[string]interface{}
	json.Unmarshal([]byte(serverBody), &result)

	servers := result["endpoints"].([]interface{})

	var domainInfo DomainInfo

	for _, value := range servers {

		var serverTemp = value.(map[string]interface{})
		var server ServerDesc

		server.ServerAddress = serverTemp["ipAddress"].(string)
		server.SSLGrade = serverTemp["grade"].(string)
		obtainWhoIsInfo(&server)
		domainInfo.Servers = append(domainInfo.Servers, server)

	}
	obtainHeaderInfo(&domainInfo)
	return domainInfo

}

func obtainWhoIsInfo(server *ServerDesc) {
	whoisInfo, e := whois.Whois(server.ServerAddress)

	if e != nil {
		log.Fatalln(e)
	} else {
		scanner := bufio.NewScanner(strings.NewReader(whoisInfo))
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			line := scanner.Bytes()
			if bytes.Contains(line, []byte{'O', 'r', 'g', 'N', 'a', 'm', 'e'}) {
				organizationLine := scanner.Text()
				server.Owner = strings.TrimSpace(strings.Split(organizationLine, ":")[1])
			}
			if bytes.Contains(line, []byte{'C', 'o', 'u', 'n', 't', 'r', 'y'}) {
				countryLine := scanner.Text()
				server.Country = strings.TrimSpace(strings.Split(countryLine, ":")[1])
			}
		}
	}
}

func obtainHeaderInfo(domainInfo *DomainInfo) {
	resp, err := http.Get("http://truora.com")
	if err != nil {
		domainInfo.IsDown = true
		log.Fatalln(err)
	} else {
		//fmt.Println( "Header: "+ resp.Header.Get("Title"))
		domainInfo.Logo = resp.Header.Get("og:image")
		domainInfo.Title = resp.Header.Get("Title")
		domainInfo.IsDown = false
	}

}
