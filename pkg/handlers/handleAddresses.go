package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const(
	TOR_IPS_LIST = "https://check.torproject.org/torbulkexitlist"
)

var (
	TorIPS map[string]struct{}
)

func MakeGetVars() http.Handler{
	r := mux.NewRouter()
	r.HandleFunc("/setTors", setTorIPsHandler()).Methods(http.MethodGet)
	return r
}

func setTorIPsHandler() func(w http.ResponseWriter, r*http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
		TorIPS = make(map[string]struct{})
		req, err  := http.Get(TOR_IPS_LIST)
		if err != nil {
			log.Println(err)
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil{
			log.Println(err)
		}
		IPsString := string(body)
		listIPS := strings.Split(IPsString,"\n")
		for _, ip := range listIPS {
			TorIPS[ip] = struct{}{}
		}
		for key, _ := range TorIPS {
			fmt.Println(key)
		}
		_, found := m
		return
	}
}
