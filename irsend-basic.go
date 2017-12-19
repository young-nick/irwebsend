package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/chbmuc/lirc"
)

func parseKeyNames(reply []string) []string {

	keyNames := []string{}
	for i := 0; i < len(reply); i++ {
		keyString := strings.Split(reply[i], " ")
		keyNames = append(keyNames, keyString[1])
	}
	return keyNames
}

// Remote represents a single lircd remote and all its available commands.
type Remote struct {
	Name     string
	Commands []string
}

func main() {
	// Initialize with path to lirc socket
	ir, err := lirc.Init("/var/run/lirc/lircd")
	if err != nil {
		panic(err)
	}

	remotesReply := ir.Command(`LIST`)
	// the ir object only keeps one Data object across replies, it seems
	// so, copy the list of remotes out to a new slice
	remotes := make([]string, len(remotesReply.Data))
	remoteCommands := make([]Remote, 0)
	copy(remotes, remotesReply.Data)

	fmt.Printf("%+v\n", remotes)

	for j := 0; j < len(remotes); j++ {
		currentRemote := remotes[j]
		log.Printf("Getting commands for %v\n", currentRemote)
		reply := ir.Command(fmt.Sprintf("LIST %v", currentRemote))
		keyNames := parseKeyNames(reply.Data)
		newRemote := Remote{Name: currentRemote, Commands: keyNames}
		remoteCommands = append(remoteCommands, newRemote)
	}

	tmpl, err := template.New("remotelist").ParseFiles("remotes.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.ExecuteTemplate(os.Stdout, "remotes.tmpl", remoteCommands)
	if err != nil {
		panic(err)
	}
	// pretty.PrettyPrint(remoteCommands)
	// Send Commands
	// reply := ir.Command(`LIST Samsung_TV`)
	// keyNames := parseKeyNames(reply.Data)
	// fmt.Printf("%+v\n", keyNames)

	// err = ir.Send("Samsung_TV KEY_POWER")
	// if err != nil {
	// 	log.Println(err)
	// }
}
