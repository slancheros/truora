package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"io/ioutil"
	"log"
	"net/http"
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
		//w.Write([]byte("hello world"))
		//respondwithJSON( w, 200,map[string]string{"message": "hello Sandra"})
		var domainInfo = retrieveDomainInfo()
		respondwithJSON(w, 200, domainInfo)
	})
	http.ListenAndServe(":3334", r)
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	fmt.Println(payload)
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

	for key, value := range servers {
		// Each value is an interface{} type, that is type asserted as a string
		fmt.Println(key, value.(map[string]interface{}))
		var serverTemp = value.(map[string]interface{})
		var server ServerDesc

		server.ServerAddress = serverTemp["ipAddress"].(string)
		server.Owner = serverTemp["serverName"].(string)
		server.SSLGrade = serverTemp["grade"].(string)
		domainInfo.Servers = append(domainInfo.Servers, server)
	}

	return domainInfo

}
