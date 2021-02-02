package main


import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/keepcalmist/Deanonimazer/pkg/handlers"
)

func main() {
	initCheckHandler := handlers.MakeCheckHandler()
	initRootHandler := handlers.MakeRootHandler()
	initPathHandler := handlers.MakeGetVars()
	router := initRouter(initCheckHandler,initRootHandler,initPathHandler)
	srv := initServer(router)
	go func() {
		log.Println("Server started on port:", os.Getenv("PORT"))
		if err := srv.ListenAndServe(); err != nil {
			log.Println("Server error: ", err)
		}
	}()
	channelToStopClicent := make (chan bool,1)
	go createClient(channelToStopClicent)
	c := make(chan os.Signal, 1)
	signal.Notify(c,os.Interrupt)
	<-c
	channelToStopClicent<-true
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	<-channelToStopClicent
	log.Println("Shutting down Client")
	log.Println("Shutting down Server")
	os.Exit(0)

}

func initRouter(checkHandle,RootHandler,MakeGetVars http.Handler) *mux.Router{
	r := mux.NewRouter()
	r.PathPrefix("/check").Handler(checkHandle)
	r.PathPrefix("/setIPs").Handler(MakeGetVars)
	r.PathPrefix("/").Handler(RootHandler)
	return r
}

func initServer(r *mux.Router) *http.Server{
	return &http.Server{
		Addr:              ":"+os.Getenv("PORT"),
		Handler:           r,
	}
}

func createClient(stop chan bool) {
	_, err := http.Get("http://127.0.0.1:"+os.Getenv("PORT")+"/setIPs/vpn")
	if err != nil {
		log.Println(err)
	}
	_, err = http.Get("http://127.0.0.1:"+os.Getenv("PORT")+"/setIPs/tor")
	if err != nil {
		log.Println(err)
	}
	ticker := time.NewTicker(time.Hour*24)
	go func() {
		for _ = range ticker.C{
			_, err = http.Get(":"+os.Getenv("PORT")+"/setIPs/vpn")
			if err != nil {
				log.Println(err)
			}
			_, err = http.Get(":"+os.Getenv("PORT")+"/setIPs/tor")
			if err != nil {
				log.Println(err)
			}
		}
	}()
	<-stop
	ticker.Stop()
}