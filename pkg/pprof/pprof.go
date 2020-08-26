package pprof

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const (
	Addr    = ":38888"
	PidFile = "server.pid"
)

func Run() {
	ch := make(chan os.Signal, 1)
	ch1 := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR1)
	signal.Notify(ch1, syscall.SIGUSR2)
	// 写入文件
	f, err := os.Create(PidFile)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	pid := os.Getpid()
	_, err = f.WriteString(strconv.Itoa(pid))
	if err != nil {
		log.Printf("%v", err)
	}
	f.Close()
	log.Print("进程id：", pid)
	var server *http.Server
	for {
		select {
		case <-ch:
			go func() {
				server = &http.Server{
					Addr: Addr,
				}
				log.Print("Listen addr:", Addr)
				err := server.ListenAndServe()
				if err != nil {
					log.Printf("Listen%v", err)
				}
			}()
		case <-ch1:
			if server != nil {
				err := server.Close()
				if err != nil {
					log.Printf("server close : %v", err)
				}
			}
		}
	}

}
