package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/likexian/whois-go"
	"golang.org/x/net/html"
	"io"
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
			line := scanner.Text()
			if strings.Contains(line, "OrgName") {
				server.Owner = strings.TrimSpace(strings.Split(line, ":")[1])
			}
			if strings.Contains(line, "Country") {
				server.Country = strings.TrimSpace(strings.Split(line, ":")[1])
			}
		}
	}
}

//func obtainHeaderInfo(domainInfo *DomainInfo) {
//	resp, err := http.Get("https://truora.com")
//	if err != nil {
//		domainInfo.IsDown = true
//		log.Fatalln(err)
//	} else {
//		defer resp.Body.Close()
//
//		responseData, err := ioutil.ReadAll(resp.Body)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		scanner := bufio.NewScanner(strings.NewReader(string(responseData)))
//		scanner.Split(bufio.ScanLines)
//		for scanner.Scan() {
//			//fmt.Println( scanner.Text())
//			line := scanner.Text()
//
//			if strings.Contains(line, "og:image"){
//			  tags := strings.Split(line, "<meta")
//			  for _,v := range tags{
//			  	  fmt.Println( v )
//			  	  if strings.Contains( v,"og:image"){
//			  	  	 domainInfo.Logo = strings.Split(v,"content=")[1]
//				  }
//			  }
//			}
//			if strings.Contains(line, "<title>") {
//				titleLine := strings.TrimSpace( line )
//				titleLine = strings.TrimLeft(titleLine, "<title>")
//				titleLine = strings.Split(titleLine,"<")[0]
//				domainInfo.Title = titleLine
//			}
//		}
//		domainInfo.IsDown = false
//	}

func obtainHeaderInfo(domainInfo *DomainInfo) {
	resp, err := http.Get("https://truora.com")
	if err != nil {
		domainInfo.IsDown = true
		log.Fatalln(err)
	} else {
		tokenizer := html.NewTokenizer(resp.Body)
		titleSet := false
		for {
			tokenType := tokenizer.Next()

			if tokenType == html.ErrorToken {
				err := tokenizer.Err()
				if err == io.EOF {
					break
				}
				log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
			}
			token := tokenizer.Token()
			if "title" == token.Data && !titleSet {
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					domainInfo.Title = tokenizer.Token().Data
					fmt.Println("content: " + token.DataAtom.String() + "," + token.Type.String() + ",token: " + token.String())
					titleSet = true
					continue
				}
			}
			if "meta" == token.Data && strings.Contains(token.String(), "og:image") {

				for k, v := range token.Attr {
					if v.Key == "content" && k == 0 {
						domainInfo.Logo = v.Val
						fmt.Println("content: " + v.Key + ", " + v.Val + ", " + v.Namespace + ", " + token.DataAtom.String() + "," + token.Type.String() + ",token: " + token.String())
						break
					}
				}
			}
		}
	}
}
