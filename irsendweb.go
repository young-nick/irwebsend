package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/chbmuc/lirc"
	"github.com/gorilla/mux"
	"github.com/young-nick/lircdremotes"
)

func Index(w http.ResponseWriter, r *http.Request, remoteCommands []lircdremotes.Remote) {
	remotesTmpl, err := template.New("remotelist").ParseFiles("templates/base.tmpl", "templates/remotes.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = remotesTmpl.ExecuteTemplate(w, "irsendweb", remoteCommands)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getRemote(remoteCommands []lircdremotes.Remote, remotename string) (lircdremotes.Remote, error) {
	for i, v := range remoteCommands {
		if v.Name == remotename {
			return remoteCommands[i], nil
		}
	}
	var emptyremote lircdremotes.Remote

	return emptyremote, fmt.Errorf("Remote name not found: %s", remotename)
}

func verifyCommand(remote lircdremotes.Remote, operation string) bool {
	for _, v := range remote.Commands {
		if v == operation {
			return true
		}
	}
	return false
}

func Device(w http.ResponseWriter, r *http.Request, remoteCommands []lircdremotes.Remote) {
	vars := mux.Vars(r)
	remote, err := getRemote(remoteCommands, vars["device"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	remoteControlTmpl, err := template.New("remoteControl").ParseFiles("templates/base.tmpl", "templates/control.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = remoteControlTmpl.ExecuteTemplate(w, "irsendweb", remote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func Clicked(w http.ResponseWriter, r *http.Request, remoteCommands []lircdremotes.Remote, ir *lirc.Router) {
	vars := mux.Vars(r)

	device := vars["device"]
	remote, err := getRemote(remoteCommands, device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	operation := vars["operation"]

	if verifyCommand(remote, operation) {
		ir.Send(fmt.Sprintf("%s %s", device, operation))
	}

	http.Error(w, fmt.Sprintf("Couldn't find operation %s", operation), http.StatusInternalServerError)

}

func main() {
	ir, err := lirc.Init("/var/run/lirc/lircd")
	if err != nil {
		panic(err)
	}
	//remotesReply := ir.Command(`LIST`)
	// the ir object only keeps one Data object across replies, it seems
	// so, copy the list of remotes out to a new slice
	//remotes := make([]string, len(remotesReply.Data))
	remoteCommands := lircdremotes.RemoteCommands(ir)
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Index(w, r, remoteCommands)
	})

	router.HandleFunc("/device/{device}", func(w http.ResponseWriter, r *http.Request) {
		Device(w, r, remoteCommands)
	})

	router.HandleFunc("/device/{device}/clicked/{operation}", func(w http.ResponseWriter, r *http.Request) {
		Clicked(w, r, remoteCommands, ir)
	})

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	srv := &http.Server{
		Handler:      router,
		Addr:         ":5001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
