package handlers

import (
	"fmt"
	"github.com/avct/uasurfer"
	"github.com/gorilla/mux"
	"github.com/ipinfo/go/v2/ipinfo"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

func MakeCheckHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/check", checkHandler()).Methods(http.MethodGet)
	return r
}

var (
	usersList []User
)

type User struct {
	UserAgentList uasurfer.UserAgent
	IP            string
}

func MakeRootHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler()).Methods(http.MethodGet)
	return r
}

func rootHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HELLO FROM ", r.RemoteAddr)
		w.Write([]byte("LOLKEKROOTLINK"))
		return
	}
}

func checkHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		user := User{}
		user.IP, _,_= net.SplitHostPort(r.RemoteAddr)
		fmt.Println(r.RemoteAddr)
		ua := uasurfer.Parse(r.UserAgent())
		user.UserAgentList = *ua
		fmt.Println(user)
		curlForProxy(user.IP)
		w.Write([]byte("Hello"))
		wg.Wait()
		return
	}
}

//ipinfo.io/8.8.8.8/json?token=d1f08163ebabbc
func curlForProxy(ip string) {
	client := ipinfo.NewClient(nil, nil, "d1f08163ebabbc")
	info, err := client.GetIPInfo(net.ParseIP(ip))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info.Privacy)
	fmt.Println(info)
}

func checkBrowser(user User, group sync.WaitGroup) (browser bool) {
	for _, r := range usersList {
		if r.UserAgentList.Browser == user.UserAgentList.Browser {
			group.Done()
			return true
		}
	}
	group.Done()
	return false
}

func checkOS(user User, group sync.WaitGroup) (browser bool) {
	for _, r := range usersList {
		if r.UserAgentList.OS == user.UserAgentList.OS {
			group.Done()
			return true
		}
	}
	group.Done()
	return false
}

func checkIP(user User, group sync.WaitGroup) (browser bool) {
	userIP := strings.Split(user.IP, ":")
	for _, r := range usersList {
		IPs := strings.Split(r.IP, ":")
		if userIP[0] == IPs[0] {
			group.Done()
			return true
		}
	}
	group.Done()
	return false
}
