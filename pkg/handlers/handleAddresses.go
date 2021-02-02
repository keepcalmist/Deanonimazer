package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	TOR_IPS_LINK = "https://check.torproject.org/torbulkexitlist"
	VPN_IPS_LINK = "https://hidemy.life/api/vpn.json"
)

var (
	TorIPS = make(map[string]struct{})
	VPNIPS = make(map[string]struct{})
)

func MakeGetVars() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/setIPs/vpn", getVPNIPsHandler()).Methods(http.MethodGet)
	r.HandleFunc("/setIPs/vpn", getVPNIPHandler()).Methods(http.MethodPost)
	r.HandleFunc("/setIPs/tor", setTorIPsHandler()).Methods(http.MethodGet)
	r.HandleFunc("/setIPs/tor", setTorIPHandler()).Methods(http.MethodPost)
	return r
}

func getVPNIPsHandler() func(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Time        string      `jsong:"time"`
		Source      string      `json:"source"`
		Title       string      `json:"title"`
		Description string      `json:"description"`
		Donate      interface{} `json:"donate"`
		ListHeader  []string    `json:"list-header"`
		List        [][]string  `json:"list"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vpn := req{}
		req, err := http.Get(VPN_IPS_LINK)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Problems with remote vpn server"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Cant read ips from remote server"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &vpn)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte("Cant unmarhal data from remote server"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, value := range vpn.List {
			if _, ok := VPNIPS[value[1]]; !ok {
				VPNIPS[value[1]] = struct{}{}
			}
		}
	}
}

func setTorIPsHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		TorIPS = make(map[string]struct{})
		req, err := http.Get(TOR_IPS_LINK)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Problems with remote tor server"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			w.Write([]byte("Cant read ips from remote server"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		IPsString := string(body)
		listIPS := strings.Split(IPsString, "\n")
		for _, ip := range listIPS { //local var with Tor nodes (every ip can be last ip address in tor chain)
			if _, ok := TorIPS[ip]; !ok {
				TorIPS[ip] = struct{}{}
			}
		}
		w.Write([]byte("Ok"))
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func setTorIPHandler() func(w http.ResponseWriter, r *http.Request) {
	type req struct {
		IP string `json:"ip"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		request := req{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		check, err := regexp.MatchString("^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$", request.IP)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !check {
			log.Println("Incorrect IP address")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad IP address\nExample 128.1.190.250\nPattern:[1-255].[0-255].[0-255].[0-255]"))
			return
		} else {
			_, found := TorIPS[request.IP]
			if found {
				w.Write([]byte("This ip address has already been added"))
				w.WriteHeader(http.StatusCreated)
				return
			} else {
				TorIPS[request.IP] = struct{}{}
				w.WriteHeader(http.StatusCreated)
				return
			}
		}

	}
}

func getVPNIPHandler() func(w http.ResponseWriter, r *http.Request) {
	type req struct {
		IP string `json:"ip"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		request := req{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		check, err := regexp.MatchString("^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$", request.IP)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !check {
			log.Println("Incorrect IP address")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Incorrect IP address\nExample 128.1.190.250\nPattern:[1-255].[0-255].[0-255].[0-255]"))
			return
		} else {
			_, found := VPNIPS[request.IP]
			if found {
				w.Write([]byte("This ip address has already been added"))
				w.WriteHeader(http.StatusCreated)
				return
			} else {
				VPNIPS[request.IP] = struct{}{}
				w.Write([]byte("Created"))
				w.WriteHeader(http.StatusCreated)
				return
			}
		}

	}
}
