package cinit

import "github.com/nats-io/nats.go"

var Natscon *nats.Conn

func Natsinit() {
	c, _ := nats.Connect(Config.Nats.Addr)
	Natscon = c
}
func Natsclose() {
	Natscon.Close()
}
