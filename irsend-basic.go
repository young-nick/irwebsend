package main

import (
	"html/template"
	"os"

	"github.com/chbmuc/lirc"
	"github.com/young-nick/lircdremotes"
)

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
	remoteCommands := lircdremotes.RemoteCommands(ir)

	remotesTmpl, err := template.New("remotelist").ParseFiles("templates/base.tmpl", "templates/remotes.tmpl")
	if err != nil {
		panic(err)
	}
	err = remotesTmpl.ExecuteTemplate(os.Stdout, "irsendweb", remoteCommands)
	if err != nil {
		panic(err)
	}
	remoteControlTmpl, err := template.New("remoteControl").ParseFiles("templates/base.tmpl", "templates/control.tmpl")
	if err != nil {
		panic(err)
	}
	for k := 0; k < len(remotes); k++ {

		err = remoteControlTmpl.ExecuteTemplate(os.Stdout, "irsendweb", remoteCommands[k])
		if err != nil {
			panic(err)
		}
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
