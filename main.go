package main

import (
	irc "github.com/fluffle/goirc/client"
	"flag"
	"fmt"
	"os"
)

func main() {
	// check flags
	var serverAddress string
	var port int
	var targetNick string
	var nick string
	var message string

	flag.StringVar(&serverAddress, "s", "", "The server to log in to")
	flag.IntVar(&port, "p", 6667, "The port to connect on")
	flag.StringVar(&targetNick, "t", "", "The user to send the message to")
	flag.StringVar(&nick, "n", "dropaline", "The nick to log into the server with")
	flag.StringVar(&message, "m", "Alert!", "The message to be sent to the target")
	flag.Parse()

	quit := make(chan bool)

	// connect and send the message
	ircConn := irc.SimpleClient(nick)

	ircConn.HandleFunc("connected", func(conn *irc.Conn, line *irc.Line) {
		conn.Privmsg(targetNick, message)
		defer os.Exit(0)
		quit <- true
	})
	ircConn.HandleFunc("disconnected", func(conn *irc.Conn, line *irc.Line) {
		defer os.Exit(1)
		quit <- true
	})

	err := ircConn.ConnectTo(serverAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		defer os.Exit(1)
		return
	}

	<-quit
	ircConn.Quit()
}

