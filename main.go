package main

import (
	"fmt"
	"net/http"
	"nft_watcher/monitor"
	"os"
)

func main()  {
	go monitor.Task.StartMonitor()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! you've requested %s\n", r.URL.Path)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	bindAddr := fmt.Sprintf(":%s", port)
	fmt.Printf("==> Server listening at %s 🚀\n", bindAddr)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}
}