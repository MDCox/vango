package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
)

type Bot struct {
	server        string
	port          string
	nick          string
	user          string
	channel       string
	pass          string
	pread, pwrite chan string
	conn          net.Conn
}

func MakeBot() *Bot {
	return &Bot{server: "irc.freenode.net",
		port:    "6667",
		nick:    "VincentVanGoBot",
		channel: "#nearpdx",
		pass:    "",
		conn:    nil,
		user:    "Tester"}
}

func (bot *Bot) Connect() (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", bot.server+":"+bot.port)
	if err != nil {
		log.Fatal("unable to connect to IRC server", err)
	}

	bot.conn = conn
	log.Printf("Connected to IRC server %s (%s)\r\n", bot.server, bot.conn.RemoteAddr())
	return bot.conn, nil
}

func main() {

	ircbot := MakeBot()

	conn, _ := ircbot.Connect()
	conn.Write([]byte("NICK " + ircbot.nick + "\r\n"))
	conn.Write([]byte("USER " + ircbot.user + " 8 * :" + ircbot.user + "\r\n"))
	conn.Write([]byte("JOIN " + ircbot.channel + "\r\n"))
	defer conn.Close()

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()

		// Features and outside functionality
		if strings.Contains(line, "PRIVMSG") {

			if strings.Contains(line, "!help") {
				conn.Write([]byte("PRIVMSG #nearpdx :I actually don't really do anything yet.\r\n"))
				conn.Write([]byte("PRIVMSG #nearpdx :Soon I'll have a little webserver of my very own where I'll display fun facts that I gather about this channel.\r\n"))
				conn.Write([]byte("PRIVMSG #nearpdx :You know, most active users, most said words, popularity of tacos vs pizza, that sort of thing.\r\n"))
			} else if strings.Contains(line, "Vincent") {
				conn.Write([]byte("PRIVMSG #nearpdx :Did someone call me? I'm a bit hard of hearing\r\n"))
			} else if strings.Contains(line, "!pizza") {
				conn.Write([]byte("PIZZA?!?"))
			}

		} else if strings.Contains(line, "JOIN "+ircbot.channel) {
			if strings.Contains(line, "Vincent") == false {
				user := line[1:strings.Index(line, "!")]
				fmt.Printf("%s", user)
				conn.Write([]byte("PRIVMSG #nearpdx :Welcome " + user + "!\r\n"))
			}

		} else if strings.Contains(line, "PING") {
			conn.Write([]byte("PONG " + line[5:] + "\r\n"))
		}

		if err != nil {
			fmt.Printf("%s", err)
			break
		}

		fmt.Printf("%s\n", line)
	}
}
