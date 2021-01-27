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
	initCheckHadnle := handlers.MakeCheckHandler()
	initRootHandler := handlers.MakeRootHandler()
	router := initRouter(initCheckHadnle,initRootHandler)
	srv := initServer(router)
	//http.ListenAndServe()
	go func() {
		log.Println("Server started on port:", os.Getenv("PORT"))
		if err := srv.ListenAndServe(); err != nil {
			log.Println("Server error: ", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c,os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Shutting down")
	os.Exit(0)

}

func initRouter(checkHandle,RootHandler http.Handler) *mux.Router{
	r := mux.NewRouter()
	r.PathPrefix("/check").Handler(checkHandle)
	r.PathPrefix("/").Handler(RootHandler)
	return r
}

func initServer(r *mux.Router) *http.Server{
	return &http.Server{
		Addr:              ":"+os.Getenv("PORT"),
		Handler:           r,
	}
}