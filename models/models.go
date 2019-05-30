package models

import "time"

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

type Item struct {
	Item      string
	QueryTime time.Time
	SSLGrade  string
}

type Items struct {
	Domains []Item `json:"items"`
}

type ListDomains struct {
	DomainNames []string `json:"items"`
}

type Config struct {
	Database database
}

type database struct {
	Server   string
	Port     string
	Database string
	User     string
	Password string
}
