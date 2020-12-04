package lancher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

type Task struct {
	Cmd  string
	Args string
}
type rusult struct {
	Code int
	Msg  string
}
//TODO chan需要使用带锁的队列,或更换这块的实现
var pids = make(chan int, 1)

func executor(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error starting  command: %s\n", arg)
		return
	}

	pids <- cmd.Process.Pid

	log.Printf("Just ran subprocess %d \n", cmd.Process.Pid)

}

func handel(w http.ResponseWriter, r *http.Request) {
	var p Task
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	go executor(p.Cmd, p.Args)
	w.Header().Set("Content-Type", "application/json")
	pid := <-pids
	rv := rusult{
		Code: 200,
		Msg:  strconv.Itoa(pid),
	}

	err = json.NewEncoder(w).Encode(rv)
	if err != nil {
		//... handle error
	}
}

func LancherLocal() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handel)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
