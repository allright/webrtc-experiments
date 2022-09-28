package main

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"net/http"
	"time"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().Local().Format("2006-01-02 15:04:05.000 ") + string(bytes))
}

type Room struct {
	clients map[net.Addr]net.Conn
}

func newRoom(name string) *Room {
	log.Printf("created room: %v \n", name)

	room := &Room{clients: make(map[net.Addr]net.Conn)}
	return room
}

func (c *Room) start(conn net.Conn) {
	c.clients[conn.RemoteAddr()] = conn
	log.Printf("start %v\n", conn.RemoteAddr())
	go func(conn net.Conn) {
		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				// handle error
				//	log.Printf("%v:%v error: %v\n", ip, port, err)
				err = conn.Close()
				log.Printf("%v closing on message read, %v\n", conn.RemoteAddr(), err)
				return
			}

			for addr, cn := range c.clients {
				if addr != conn.RemoteAddr() {
					log.Printf("%v -> %v msg: %v, op: %v", conn.RemoteAddr(), cn.RemoteAddr(), msg, op)
					err = wsutil.WriteServerMessage(cn, op, msg)
					if err != nil {
						log.Printf("%v closing on write: %v\n", conn.RemoteAddr(), err)
						break
					}
				}
			}

		}

	}(conn)
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
	//killSignal := make(chan os.Signal, 1)
	//signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)

	m := make(map[string]*Room)

	log.Printf("Starting server!")
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		req := r.RequestURI
		room, ok := m[req]
		if !ok {
			room = newRoom(req)
			m[req] = room
		}

		//ip := r.Header["X-Real-Ip"]
		//port := r.Header["X-Real-Port"]
		//log.Printf("http: %v:%v \n", ip, port, r.RequestURI)
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err == nil {
			room.start(conn)
		}
		//if err != nil {
		//	// handle error
		//}
		//t1 := time.Now()
		//go func() {
		//
		//	msg, _, err := wsutil.ReadClientData(conn)
		//	if err != nil {
		//		// handle error
		//		//	log.Printf("%v:%v error: %v\n", ip, port, err)
		//		err := conn.Close()
		//		log.Printf("%v:%v closing of first message read: %v\n", ip, port, err)
		//		return
		//	}
		//	key := string(msg)
		//	s := strings.Split(key, "@")
		//	if len(s) >= 2 {
		//		key = s[0] + "@" + s[1]
		//	}
		//
		//	delta := time.Duration(0)
		//	{
		//		mLock.Lock()
		//		_, found := m[key]
		//		if !found {
		//			m[key] = t1
		//		} else {
		//			t := m[key]
		//			m[key] = t1
		//			delta = t1.Sub(t)
		//			log.Printf("%v:%v reconnect: %v diff: %v\n", ip, port, key, delta)
		//		}
		//		mLock.Unlock()
		//
		//	}
		//
		//	for {
		//		msg, op, err := wsutil.ReadClientData(conn)
		//		if err != nil {
		//			// handle error
		//			//					log.Printf("%v:%v rd error: %v msg: %v diff: %v\n", ip, port, err, key, delta)
		//			log.Printf("%v:%v rd close: %v %v\n", ip, port, err, key)
		//			break
		//		}
		//		newTime := time.Now()
		//		delta := newTime.Sub(t1)
		//		t1 = newTime
		//
		//		log.Printf("%v:%v msg: %v diff: %v\n", ip, port, string(msg), delta)
		//		err = wsutil.WriteServerMessage(conn, op, msg)
		//		if err != nil {
		//			// handle error
		//			//					log.Printf("%v:%v wr error: %v msg: %v diff: %v\n", ip, port, err, msg, delta)
		//			log.Printf("%v:%v wr close: %v %v\n", ip, port, err, key)
		//			break
		//		}
		//	}
		//	err = conn.Close()
		//}()
	}))
	//
	//<-killSignal
	//fmt.Println("Thanks for using Golang!")
}
