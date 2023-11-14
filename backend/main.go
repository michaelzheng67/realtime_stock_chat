package main

import (
	"os"
	"os/signal"
    "encoding/json"

	polygonws "github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
	"github.com/sirupsen/logrus"

    "golang.org/x/net/websocket"
	"fmt"
	"io"
	"net/http"

	"github.com/joho/godotenv"
)


type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}
		msg := buf[:n]
		
		s.broadcast(msg)
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error: ", err)
			}
		}(ws)
	}
} 

func sendData(url string, origin string) {

	// load env variables
	err := godotenv.Load("variables.env")
	if err != nil {
        fmt.Println("Error loading .env file")
    }


	// create websocket client
	ws, err := websocket.Dial(url, "", origin)
    if err != nil {
        fmt.Println("Error connecting to WebSocket server:", err)
        os.Exit(1)
    }
    defer ws.Close()


    log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
	c, err := polygonws.New(polygonws.Config{
		APIKey: os.Getenv("POLYGON_API_KEY"),
		Feed:   polygonws.Delayed,
		Market: polygonws.Stocks,
		Log:    log,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// aggregates
	_ = c.Subscribe(polygonws.StocksSecAggs, "AAPL", "TSLA", "AMZN", "GOOG", "META", "BABA", "MSFT")

	if err := c.Connect(); err != nil {
		log.Error(err)
		return
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	for {
		select {
		case <-sigint:
			return
		case <-c.Error():
			return
		case out, more := <-c.Output():
			if !more {
				return
			}
			switch out.(type) {

			case models.EquityAgg:
				log.WithFields(logrus.Fields{"aggregate": out}).Info()

				jsonData, err := json.Marshal(out)
				if err != nil {
					fmt.Println("Error marshalling JSON:", err)
					return
				}

				if err := websocket.Message.Send(ws, jsonData); err != nil {
					fmt.Println("Error sending message:", err)
					return
				}

                

			case models.EquityTrade:
				log.WithFields(logrus.Fields{"trade": out}).Info()
			case models.EquityQuote:
				log.WithFields(logrus.Fields{"quote": out}).Info()
			}
		}
	}
}

func main() {

	// :8000 - stocks
	// :8080 - trader chat

    stock_server := NewServer()
	go func() {
		http.Handle("/stock-ws", websocket.Handler(stock_server.handleWS))
		if err := http.ListenAndServe(":8000", nil); err != nil {
			fmt.Printf("Failed to start stock server: %v", err)
		}
	}()

	trader_server := NewServer()
	go func() {
		http.Handle("/trader-ws", websocket.Handler(trader_server.handleWS))
		if err := http.ListenAndServe(":9000", nil); err != nil {
			fmt.Printf("Failed to start trader server: %v", err)
		}
	}()
	
    sendData("ws://localhost:8000/stock-ws", "http://localhost/");

}
