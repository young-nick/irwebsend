package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/chbmuc/lirc"
	"github.com/gobs/pretty"
)

func parseKeyNames(reply []string) []string {

	keyNames := []string{}
	for i := 0; i < len(reply); i++ {
		keyString := strings.Split(reply[i], " ")
		keyNames = append(keyNames, keyString[1])
	}
	return keyNames
}

func main() {
	// Initialize with path to lirc socket
	ir, err := lirc.Init("/var/run/lirc/lircd")
	if err != nil {
		panic(err)
	}

	remoteCommands := make(map[string][]string)

	remotesReply := ir.Command(`LIST`)
	// the ir object only keeps one Data object across replies, it seems
	// so, copy the list of remotes out to a new slice
	remotes := make([]string, len(remotesReply.Data))
	copy(remotes, remotesReply.Data)

	fmt.Printf("%+v\n", remotes)

	for j := 0; j < len(remotes); j++ {
		currentRemote := remotes[j]
		log.Printf("Getting commands for %v\n", currentRemote)
		reply := ir.Command(fmt.Sprintf("LIST %v", currentRemote))
		keyNames := parseKeyNames(reply.Data)
		remoteCommands[currentRemote] = keyNames
	}

	pretty.PrettyPrint(remoteCommands)
	// Send Commands
	// reply := ir.Command(`LIST Samsung_TV`)
	// keyNames := parseKeyNames(reply.Data)
	// fmt.Printf("%+v\n", keyNames)

	// err = ir.Send("Samsung_TV KEY_POWER")
	// if err != nil {
	// 	log.Println(err)
	// }
}
