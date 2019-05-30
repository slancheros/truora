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
	"time"
	"truora/db"
	"truora/models"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", string('*'))
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
	domainInfo.ServersSSLGrade = getSSLGrade(domainInfo)
	UpdateDomain(domain, &domainInfo)
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

func UpdateDomain(domain string, info *models.DomainInfo) models.DomainInfo {
	db.Connect()

	var item models.Item
	domain_info := domain + " info"

	if db.GetDomain(domain_info, &item) {
		if item.SSLGrade == "" || item.SSLGrade == ":" {
			info.PreviousSSLGrade = info.ServersSSLGrade
		} else {
			info.PreviousSSLGrade = item.SSLGrade
		}

		if item.SSLGrade != info.ServersSSLGrade && item.QueryTime.Add(60*time.Minute).Before(time.Now()) {
			info.ServersChanged = "true"
		} else {
			info.ServersChanged = "false"
		}

	} else {
		info.ServersChanged = "false"
	}
	db.UpdateQueriedDomain(domain, *info)
	defer db.Close()
	return *info
}

func getSSLGrade(domainInfo models.DomainInfo) string {

	SSLGrade := "A"

	for i := range domainInfo.Servers {
		if domainInfo.Servers[i].SSLGrade > SSLGrade {
			SSLGrade = domainInfo.Servers[i].SSLGrade
		}
	}
	return SSLGrade
}

func Contains(a []models.Item, domain string) models.Item {
	var found models.Item
	for _, n := range a {
		if domain == n.Item {
			found = n
		}
	}
	return found
}

func ListDomainNamesQueried(w http.ResponseWriter, r *http.Request) {
	db.Connect()
	items := db.ListDomainItems()
	defer db.Close()
	var domainList models.ListDomains
	for i := range items.Domains {
		domainList.DomainNames = append(domainList.DomainNames, items.Domains[i].Item)
	}
	if len(domainList.DomainNames) == 0 {
		respondWithJSON(w, 200, make([]string, 0))
	} else {
		respondWithJSON(w, 200, domainList)

	}
}
