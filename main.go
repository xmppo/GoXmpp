package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	my := MyXMPP{
		ServerAddess: "192.168.1.1:5222",
		Name:         "id@xmpp-host",
		Password:     "1234"
		SSL:          false, //if setting with default when you install jabberd2
	}

	go func() {
		for {
			from, chat, err := my.GetMessage()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(from, chat)
		}
	}()

	for {
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			continue
		}
		line = strings.TrimRight(line, "\n")

		tokens := strings.SplitN(line, " ", 2)
		fmt.Printf("input name=%s, msg=%s \n", tokens[0], tokens[1])
		if len(tokens) == 2 {
			my.SendMessage(tokens[0], tokens[1])
		}
	}
}
