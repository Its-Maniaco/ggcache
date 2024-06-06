package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Command string

const (
	CMDSet Command = "SET"
	CMDGet Command = "GET"
)

type Message struct {
	Cmd Command
	Key []byte
	Val []byte
	TTL time.Duration
}

func ParseMessage(raw []byte) (*Message, error) {
	var (
		rawStr = string(raw)
		parts  = strings.Split(rawStr, " ")
	)

	if len(parts) < 2 {
		return nil, errors.New("invalid protocol")
	}

	msg := &Message{
		Cmd: Command(parts[0]),
		Key: []byte(parts[1]),
	}

	// set values for MSGSet
	if msg.Cmd == CMDSet {
		if len(parts) < 4 {
			return nil, errors.New("invalid SET command")
		}

		// convert ttl to set
		ttl, err := strconv.Atoi(parts[3])
		fmt.Println("ttl: ", ttl)
		fmt.Println("err: ", err)
		if err != nil {
			return nil, errors.New("invalid SET TTL")
		}

		msg.Val = []byte(parts[2])
		msg.TTL = time.Duration(ttl)

		return msg, nil
	}

	return nil, nil
}
