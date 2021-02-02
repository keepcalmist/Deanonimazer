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
)

var (
	usersList    []User
	PROXYHEADERS = []string{"HTTP_VIA", "HTTP_X_FORWARDED_FOR", "HTTP_FORWARDED_FOR", "HTTP_X_FORWARDED",
		"HTTP_FORWARDED", "HTTP_CLIENT_IP", "HTTP_FORWARDED_FOR_IP", "VIA", "X_FORWARDED_FOR", "FORWARDED_FOR",
		"X_FORWARDED", "FORWARDED", "CLIENT_IP", "FORWARDED_FOR_IP", "HTTP_PROXY_CONNECTION"}
)

func MakeCheckHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/check", checkHandler()).Methods(http.MethodGet)
	return r
}

type User struct {
	UserAgentList uasurfer.UserAgent
	IP            string
}

type trues struct {
	Browser chan bool
	OS      chan bool
	IP      chan bool
	Headers chan []string
	VPN     chan bool
	Tor     chan bool
}

func MakeRootHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler()).Methods(http.MethodGet)
	return r
}

func rootHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, port, _ := net.SplitHostPort(r.RemoteAddr)
		log.Println("HELLO FROM ",ip,":", port )
		w.Write([]byte("LOLKEKROOTLINK"))
		return
	}
}

func checkHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		channelForCheck := trues{
			Browser: make(chan bool, 1),
			OS:      make(chan bool, 1),
			IP:      make(chan bool, 1),
			Headers: make(chan []string, 1),
			VPN:     make(chan bool, 1),
			Tor:     make(chan bool, 1),
		}
		user := User{}
		ua := uasurfer.Parse(r.UserAgent())
		user.UserAgentList = *ua
		user.IP, _, _ = net.SplitHostPort(r.RemoteAddr)
		////

		go checkIP(user.IP, channelForCheck.IP)
		go checkOS(user, channelForCheck.OS)
		go checkBrowser(user, channelForCheck.Browser)
		go checkForProxyHeaders(r, channelForCheck.Headers)
		go checkVPN(user.IP, channelForCheck.VPN)
		go checkTOR(user.IP, channelForCheck.Tor)

		////Example of output response
		/*
		Do is user exist? - Y/N		(If [NO] then add new user to UserList)
		Tor - Y/N
		VPN - Y/N
		Proxy - Y/N
		[if Proxy == Y] then Println(Headers)
						else Nothing
 		*/
		////
		//check for exists
		args := map[string]string {
			"Exists": "",
			"Tor": "",
			"VPN": "",
			"Proxy Headers": "",
			"Browser and OS - ": "",
		}
		func () {
			//check ip
			 if trues := <- channelForCheck.IP; trues{
				//user exists
			 	args["Exists"] = "Do is user exist? - Yes"
			 	args["Tor"] = " No"
			 	args["VPN"] = " No"
			 	args["Proxy Headers"] = " No"
			 } else {
			 	//Check for VPN and TOR and Proxy
				 vpn := <-channelForCheck.VPN
				 if vpn {
					 args["VPN"] = " Yes"
				 } else {
					 args["VPN"] = " No"
				 }
				 tor := <-channelForCheck.Tor
				 if tor {
					 args["Tor"] = " Yes"
				 } else {
					 args["Tor"] = " No"
				 }
				 headers := <- channelForCheck.Headers
				 fmt.Println(headers, len(headers))
				 if headers == nil {
					 args["Proxy Headers"] = " No"
				 } else {
				 	args["Proxy Headers"] = strings.Join(headers,", ")
				 }
					if tor||vpn||(headers!=nil) {
						os := <-channelForCheck.OS
						browser := <-channelForCheck.Browser
						if os && browser {
							args["Browser and OS - "]=  "Probably we saw this person"
						}
					} else {
						args["Exists"] = "Do is user exist? - Yes\n"
						usersList = append(usersList,user)
					}
			 	//
			 }

		}()


		curlForProxy(user.IP)
		key, _ := args["Exists"]
		fmt.Println(key)
		w.Write([]byte(key))
		for key, value := range args{
			if key != "Exists" {
				w.Write([]byte(key + value + "\n"))
			}
		}

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

func checkForProxyHeaders(r *http.Request, channel chan []string) {
	Headers := make([]string, 0, 2)
	for _, header := range PROXYHEADERS {
		if r.Header.Get(header) != "" {
			Headers = append(Headers, header)
		}
	}
	channel <- Headers
}

func checkBrowser(user User, channel chan bool) {
	for _, r := range usersList {
		if r.UserAgentList.Browser == user.UserAgentList.Browser {
			channel <- true
			return
		}
	}
	channel <- false
	return
}

func checkOS(user User, channel chan bool) {
	for _, r := range usersList {
		if r.UserAgentList.OS == user.UserAgentList.OS {
			channel <- true
			return
		}
	}
	channel <- false
	return
}

func checkIP(remoteAddr string, YeaOrNo chan bool) {
	userIP, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Println(err)
		YeaOrNo <- false
		return
	}
	for _, r := range usersList {
		IPs, _, _ := net.SplitHostPort(r.IP)
		if userIP == IPs {
			YeaOrNo <- true
			return
		}
	}
	YeaOrNo <- false
	return
}

func checkVPN(ip string, channel chan bool) {
	for key, _ := range VPNIPS {
		if ip == key {
			channel <- true
			return
		}
	}
	channel <- false
	return
}
func checkTOR(ip string, channel chan bool) {
	for key, _ := range TorIPS {
		if key == ip {
			channel <- true
			return
		}
	}
	channel <- false
	return
}
