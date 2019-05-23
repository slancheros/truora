package handlers

import (
	"bufio"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/likexian/whois-go"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"truora/db"
	"truora/models"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(code)
	w.Write(response)
}

func RetrieveDomainInfo(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	resp, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=" + domain)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	domainInfo := createServerInfo(string(body), domain)
	db.Connect()
	db.SaveQueriedDomain(domain)
	defer db.Close()
	respondWithJSON(w, 200, domainInfo)

}

func createServerInfo(serverBody string, domain string) models.DomainInfo {

	var result map[string]interface{}
	json.Unmarshal([]byte(serverBody), &result)

	servers := result["endpoints"].([]interface{})

	var domainInfo models.DomainInfo

	for _, value := range servers {
		var serverTemp = value.(map[string]interface{})
		var server models.ServerDesc
		if serverTemp["ipAddress"] != nil {
			server.ServerAddress = serverTemp["ipAddress"].(string)
		}
		if serverTemp["grade"] != nil {
			server.SSLGrade = serverTemp["grade"].(string)
		}
		obtainWhoIsInfo(&server, domain)
		domainInfo.Servers = append(domainInfo.Servers, server)
	}
	obtainHeaderInfo(&domainInfo, domain)
	return domainInfo
}

func obtainWhoIsInfo(server *models.ServerDesc, domain string) {
	whoisInfo, e := whois.Whois(domain)

	if e != nil {
		log.Fatalln(e)
	} else {
		scanner := bufio.NewScanner(strings.NewReader(whoisInfo))
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Registrant Organization") {
				server.Owner = strings.TrimSpace(strings.Split(line, ":")[1])
			}
			if strings.Contains(line, "Country") {
				server.Country = strings.TrimSpace(strings.Split(line, ":")[1])
			}
		}
	}

}

func obtainHeaderInfo(domainInfo *models.DomainInfo, domain string) {
	resp, err := http.Get("https://" + domain)

	if err != nil {
		domainInfo.IsDown = true
		log.Fatalln(err)
	} else {
		tokenizer := html.NewTokenizer(resp.Body)
		titleSet := false
		logoSet := false
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
					titleSet = true
					continue
				}
			}
			if "meta" == token.Data && strings.Contains(token.String(), "og:image") {
				for _, v := range token.Attr {
					if v.Key == "content" && !logoSet {
						domainInfo.Logo = v.Val
						logoSet = true
						break
					}
				}
			}
		}
	}
}

func ListDomainsQueried(w http.ResponseWriter, r *http.Request) {
	db.Connect()
	items := db.ListItems()
	defer db.Close()
	respondWithJSON(w, 200, items)

}
