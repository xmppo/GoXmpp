package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/mattn/go-xmpp"
	"log"
	"os"
	"strings"
)

var talk *xmpp.Client

type MyXMPP struct {
	ServerAddess string //include port
	Name         string //include host/resource
	Password     string
	SSL          bool
}

type MyXMPPAction interface {
	Connect() error
	GetMessage() (string, string, error)
	SendMessage(to string, msg string) error
}

func (my MyXMPP) Connect() error {
	var server = flag.String("server", my.ServerAddess, "server")
	var username = flag.String("username", my.Name, "username")
	var password = flag.String("password", my.Password, "password")
	var status = flag.String("status", "xa", "status")
	var statusMessage = flag.String("status-msg", "I for one welcome our new codebot overlords.", "status message")
	var notls = flag.Bool("notls", !my.SSL, "No TLS")
	var debug = flag.Bool("debug", true, "debug output")
	var session = flag.Bool("session", false, "use server session")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: example [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	if *username == "" || *password == "" {
		flag.Usage()
	}

	if !*notls {
		xmpp.DefaultConfig = tls.Config{
			ServerName:         serverName(*server),
			InsecureSkipVerify: true,
		}
	}

	var err error
	options := xmpp.Options{Host: *server,
		User:                         *username,
		Password:                     *password,
		NoTLS:                        *notls,
		Debug:                        *debug,
		Session:                      *session,
		Status:                       *status,
		StatusMessage:                *statusMessage,
		InsecureAllowUnencryptedAuth: true,
	}

	talk, err = options.NewClient()

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (my MyXMPP) GetMessage() (string, string, error) {
	chat, err := talk.Recv()
	if err != nil {
		log.Fatal(err)
	}
	switch v := chat.(type) {
	case xmpp.Chat:
		//fmt.Println(v.Remote, v.Text)
		return v.Remote, v.Text, nil
	case xmpp.Presence:
		fmt.Printf("presence .. %s is %s \n", v.From, v.Show)
	}
	return "", "", nil
}

func (my MyXMPP) SendMessage(to string, msg string) error {
	talk.Send(xmpp.Chat{Remote: to, Type: "chat", Text: msg})
	return nil
}

func serverName(host string) string {
	return strings.Split(host, ":")[0]
}
