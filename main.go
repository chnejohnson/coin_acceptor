package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tarm/serial"
	"golang.org/x/net/websocket"
)

var oneDollar = []byte{0x90, 0x06, 0x12, 0x01, 0x03, 0xAC}
var fiveDollar = []byte{0x90, 0x06, 0x12, 0x02, 0x03, 0xAD}
var tenDollar = []byte{0x90, 0x06, 0x12, 0x03, 0x03, 0xAE}

var amount = 0

func addCoin(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		for {
			// err := websocket.Message.Send(ws, "hello websocket")
			// if err != nil {
			// 	log.Println(err)
			// }

			go listenCoin(ws)

			// Read
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
				break
			}
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "./frontend/dist/")
	e.GET("/ws", addCoin)

	e.Logger.Fatal(e.Start(":8080"))
}

func listenCoin(ws *websocket.Conn) {
	c := &serial.Config{
		Name: "/dev/ttyUSB0",
		Baud: 9600,
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 128)

	msg := []byte{}

	for {
		n, err := s.Read(buf)

		if err != nil {
			log.Println(err)
			continue
		}

		// log.Printf("Received: %q\n", buf[:n])
		msg = append(msg, buf[:n]...)

		// 用長度判斷訊號終止點
		if len(msg) == 6 {
			// log.Printf("%x\n", msg)
			if bytes.Compare(msg, oneDollar) == 0 {
				log.Println("Get 1 dollar!")
				amount++

				err := websocket.Message.Send(ws, strconv.Itoa(amount))
				if err != nil {
					log.Println("Fail to Send message: ", err)
				}
			}

			if bytes.Compare(msg, fiveDollar) == 0 {
				log.Println("Get 5 dollar!")
				amount += 5

				err := websocket.Message.Send(ws, strconv.Itoa(amount))
				if err != nil {
					log.Println("Fail to Send message: ", err)
				}
			}

			if bytes.Compare(msg, tenDollar) == 0 {
				log.Println("Get 10 dollar!")
				amount += 10

				err := websocket.Message.Send(ws, strconv.Itoa(amount))
				if err != nil {
					log.Println("Fail to Send message: ", err)
				}
			}

			msg = []byte{}
		}
	}
}
