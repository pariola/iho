package main

import (
	"flag"
	"log"

	"iho/internal/pkg/client"
	"iho/internal/pkg/server"
)

func main() {
	var listenAddr string
	var toAddr, remoteAddr string

	flag.StringVar(&listenAddr, "listen", ":3333", "port for server to listen on")

	flag.StringVar(&toAddr, "to", "", "local address to forward requests to.")
	flag.StringVar(&remoteAddr, "remote", ":3333", "remote address to iho instance")
	flag.Parse()

	switch flag.Arg(0) {
	default:
		log.Fatalf("unknown mode")

	case "client":
		err := client.Connect(toAddr, remoteAddr)
		if err != nil {
			log.Fatalf("Client failed to start %v", err)
		}

	case "server":
		log.Printf("Tunnel server starting on %s", listenAddr)
		err := server.ListenAndServe(listenAddr)
		if err != nil {
			log.Fatalf("Tunnel failed to start %v", err)
		}
	}

	// TODO: proper shutdown
	// stop := make(chan os.Signal)
	// signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	// <-stop
}
