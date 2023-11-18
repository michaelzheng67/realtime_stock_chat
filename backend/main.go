package main

import (
	"os"
	"os/signal"
    "encoding/json"
	"time"
	"sync"

	polygonws "github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
	"github.com/sirupsen/logrus"

    "golang.org/x/net/websocket"
	"fmt"
	"io"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
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


// getBusinessDaysBefore calculates the last 90 business days before the given date.
func getBusinessDaysBefore(inputDate string) ([]time.Time, error) {
    // Parse the input date
    date, err := time.Parse("2006-01-02", inputDate)
    if err != nil {
        return nil, err
    }

    var businessDays []time.Time
    for len(businessDays) < 90 {
        date = date.AddDate(0, 0, -1) // Go back one day
        weekday := date.Weekday()

        // Check if it's a weekend
        if weekday != time.Saturday && weekday != time.Sunday {
            businessDays = append(businessDays, date)
        }
    }

    return businessDays, nil
}


// Response object from Polygon API
type StockResponse struct {
    Status     string  `json:"status"`
    From       string  `json:"from"`
    Symbol     string  `json:"symbol"`
    Open       float64 `json:"open"`
    High       float64 `json:"high"`
    Low        float64 `json:"low"`
    Close      float64 `json:"close"`
    Volume     float64 `json:"volume"`
    AfterHours float64 `json:"afterHours"`
    PreMarket  float64 `json:"preMarket"`
}

// GET endpoint for stock chart data
func stockHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    date := vars["date"]
    ticker := vars["ticker"]

	days, _ := getBusinessDaysBefore(date)

	// Prepare the response object
    response := struct {
        Ticker string   `json:"ticker"`
        Date   string   `json:"date"`
        Days   []string `json:"days"`
    }{
        Ticker: ticker,
        Date:   date,
        Days:   make([]string, len(days)),
    }

	// Convert each day to a string
    var wg sync.WaitGroup

    for i, day := range days {
        wg.Add(1)
        go func(i int, day time.Time) {
            defer wg.Done()

            url := fmt.Sprintf("https://api.polygon.io/v1/open-close/%s/%s?adjusted=true&apiKey=%s", ticker, day.Format("2006-01-02"), os.Getenv("POLYGON_API_KEY"))

            resp, err := http.Get(url)
            if err != nil {
                fmt.Printf("Error making request: %v\n", err)
                return
            }
            defer resp.Body.Close()

            body, err := io.ReadAll(resp.Body)
            if err != nil {
                fmt.Printf("Error reading response: %v\n", err)
                return
            }

            var stockResponse StockResponse
            err = json.Unmarshal(body, &stockResponse)
            if err != nil {
                fmt.Printf("Error parsing JSON: %v\n", err)
                return
            }

            response.Days[i] = fmt.Sprintf("%.2f", stockResponse.Close)
        }(i, day)
    }

    wg.Wait()

	


	// Convert the response object to JSON
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	// Set Content-Type header
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    // Write the JSON to the response
    w.Write(jsonResponse)
}


func main() {

	// :8000 - stocks
	// :8080 - trader chat

    stock_server := NewServer()
	stockMux := http.NewServeMux()
	stockMux.Handle("/stock-ws", websocket.Handler(stock_server.handleWS))
	go func() {
		if err := http.ListenAndServe(":8000", stockMux); err != nil {
			fmt.Printf("Failed to start stock server: %v", err)
		}
	}()

	trader_server := NewServer()
	traderMux := http.NewServeMux()
	traderMux.Handle("/trader-ws", websocket.Handler(trader_server.handleWS))
	go func() {
		if err := http.ListenAndServe(":9000", traderMux); err != nil {
			fmt.Printf("Failed to start trader server: %v", err)
		}
	}()

	// http server
	r := mux.NewRouter()

    // Define a route with path parameters date and ticker
	go func() {
		r.HandleFunc("/stock/{ticker}/{date}", stockHandler)
    	http.ListenAndServe(":9090", r)
	}()
	
    sendData("ws://localhost:8000/stock-ws", "http://localhost/");

}
