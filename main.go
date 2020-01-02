package main

import(
	"net"
	"fmt"
	"time"
	"os"
	"os/signal"
)

type serverStatus struct {
	Path string
	Status bool
}

func checkServer(path string) bool {
	c, err := net.Dial("unix", path)
	if err != nil {
		return false
	}
	c.Close()
	return true
}

func checkServers(sS *[]serverStatus) {
	for i, s := range *sS {
		if checkServer(s.Path) {
			if !s.Status {
				fmt.Println("Server:", s.Path, " Up")
			}
			(*sS)[i].Status = true
		} else {
			if s.Status {
				fmt.Println("Server:", s.Path, " Down")
			}
			(*sS)[i].Status = false
		}
	}
}

func startCheckTask(servers []string, interval time.Duration) {
	serverstatuses := make([]serverStatus, len(servers))
	for i, p := range servers {
		serverstatuses[i] = serverStatus{Path: p, Status: true}
	}
	tick := time.NewTicker(time.Second * interval)
	for _ = range tick.C {
		checkServers(&serverstatuses)
	}
}

func main(){
	qt := make(chan os.Signal)
	signal.Notify(qt, os.Interrupt, os.Kill)
	servers := []string{
		"/home/mastodon/run/puma.sock",
	}
	go startCheckTask(servers, 30)
	<-qt
}