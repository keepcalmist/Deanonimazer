package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ipinfo/go/v2/ipinfo"
	"net"
	"net/http"
	"github.com/avct/uasurfer"
	"strings"
	"sync"
)

func MakeCheckHandler() http.Handler{
	r := mux.NewRouter()
	r.HandleFunc("/check", checkHandler()).Methods(http.MethodGet)
	return r
}

var (
	usersList []User
)

type User struct {
	UserAgentList uasurfer.UserAgent
	IP string
}

func checkHandler() func (w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		user := User{}
		user.IP = r.RemoteAddr
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
func curlForProxy(ip string){
client := ipinfo.NewClient(nil,nil,"d1f08163ebabbc")
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
	userIP := strings.Split(user.IP,":")
	for _, r := range usersList {
		IPs := strings.Split(r.IP,":")
		if userIP[0] == IPs[0] {
			group.Done()
			return true
		}
	}
	group.Done()
	return false
}