package main

import (
	"io"
	"log"
	"net"

	"github.com/minio/cli"
)

// Handles forwarding incoming connections to forward address.
func connForward(forwardAddress string, incoming net.Conn) {
	// Attempt a tcp client connection to forward address.
	client, err := net.Dial("tcp", forwardAddress)
	if err != nil {
		log.Fatalf("Failed to establish client connection to %s, upon %s", forwardAddress, err)
	}
	// Redirect all the incoming to the established client connection
	// in a routine.
	go func() {
		defer client.Close()
		defer incoming.Close()
		io.Copy(client, incoming)
	}()
	// Get all the replies from the client back to the incoming connection.
	go func() {
		defer client.Close()
		defer incoming.Close()
		io.Copy(incoming, client)
	}()
}

func main() {
	app := cli.NewApp()
	app.Usage = "Simple port forwarding tool"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen",
			Value: ":80",
		},
		cli.StringFlag{
			Name:  "forward",
			Value: ":9000",
		},
	}
	app.Action = func(c *cli.Context) {
		// Listen address to accept new connections.
		listenAddr := c.String("listen")
		// Initialize a listener at listen address.
		listener, err := net.Listen("tcp", listenAddr)
		if err != nil {
			log.Fatalf("Listening on %s failed with %s\n", listenAddr, err)
		}
		// Start the loop to accept incoming tcp connections on listen address.
		for {
			incoming, err := listener.Accept()
			if err != nil {
				log.Fatalf("Failed to accept connection %s\n", err)
			}
			forwardAddress := c.String("forward")
			// Forward all incoming connections to forward address.
			log.Println("Accepted and successfully forwarded the connection to", forwardAddress)
			go connForward(forwardAddress, incoming)
		}
	}
	app.Version = "0.0.1"
	app.RunAndExitOnError()
}
