package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/golang/protobuf/jsonpb"
	pb "github.com/webmocha/lumberman/pb"
	"google.golang.org/grpc"
)

var (
	port         = flag.Int("port", 80, "Listen address")
	lmServerAddr = flag.String("server_addr", "127.0.0.1:9090", "The Lumberman server address in the format of host:port")
)

type UnsubscribeFunc func() error

type Subscriber interface {
	Subscribe(c chan []byte) (UnsubscribeFunc, error)
}

type Notifier interface {
	Notify(b []byte) error
}

type NotificationCenter struct {
	subscribers   map[chan []byte]struct{}
	subscribersMu *sync.Mutex
}

func NewNotificationCenter() *NotificationCenter {
	return &NotificationCenter{
		subscribers:   map[chan []byte]struct{}{},
		subscribersMu: &sync.Mutex{},
	}
}

func (nc *NotificationCenter) Subscribe(c chan []byte) (UnsubscribeFunc, error) {
	nc.subscribersMu.Lock()
	nc.subscribers[c] = struct{}{}
	nc.subscribersMu.Unlock()

	unsubscribeFn := func() error {
		nc.subscribersMu.Lock()
		delete(nc.subscribers, c)
		nc.subscribersMu.Unlock()

		return nil
	}

	return unsubscribeFn, nil
}

func (nc *NotificationCenter) Notify(b []byte) error {
	nc.subscribersMu.Lock()
	defer nc.subscribersMu.Unlock()

	for c := range nc.subscribers {
		select {
		case c <- b:
		default:
		}
	}

	return nil
}

func handleListPrefixes(lmc *lmClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, lmc.ListPrefixes())
	}
}

func handleListKeys(lmc *lmClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params, ok := r.URL.Query()["prefix"]

		if !ok || len(params[0]) < 1 {
			log.Println("Url Param 'prefix' is missing")
			http.Error(w, `{"error": "Url Param 'prefix' is missing"}`, http.StatusBadRequest)
			return
		}

		prefix := string(params[0])

		fmt.Fprintf(w, lmc.ListKeys(prefix))
	}
}

func handleGetLog(lmc *lmClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params, ok := r.URL.Query()["key"]

		if !ok || len(params[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			http.Error(w, `{"error": "Url Param 'key' is missing"}`, http.StatusBadRequest)
			return
		}

		key := string(params[0])

		fmt.Fprintf(w, lmc.GetLog(key))
	}
}

func handleGetLogsStream(lmc *lmClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")

		params, ok := r.URL.Query()["prefix"]

		if !ok || len(params[0]) < 1 {
			log.Println("Url Param 'prefix' is missing")
			http.Error(w, `{"error": "Url Param 'prefix' is missing"}`, http.StatusBadRequest)
			return
		}

		prefix := string(params[0])

		lmc.GetLogsStream(prefix, w)
	}
}

func handleTailLogsStream(lmc *lmClient, s Subscriber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Subscribe
		c := make(chan []byte)
		unsubscribeFn, err := s.Subscribe(c)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
			return
		}

		params, ok := r.URL.Query()["prefix"]

		if !ok || len(params[0]) < 1 {
			log.Println("Url Param 'prefix' is missing")
			http.Error(w, `{"error": "Url Param 'prefix' is missing"}`, http.StatusBadRequest)
			return
		}

		prefix := string(params[0])

		// SSE Support
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		lmc.TailLogsStream(prefix, s, unsubscribeFn, w)
	}
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*lmServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %s\n%v", *lmServerAddr, err)
	}
	defer conn.Close()

	lmc := &lmClient{
		client: pb.NewLoggerClient(conn),
		m:      &jsonpb.Marshaler{},
	}

	nc := NewNotificationCenter()

	go func() {
		for {
			b := []byte(time.Now().Format(time.RFC3339))
			if err := nc.Notify(b); err != nil {
				log.Fatal(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.HandleFunc("/api/list-prefixes", handleListPrefixes(lmc))
	http.HandleFunc("/api/list-keys", handleListKeys(lmc))
	http.HandleFunc("/api/get-log", handleGetLog(lmc))
	http.HandleFunc("/api/get-logs-stream", handleGetLogsStream(lmc))
	http.HandleFunc("/api/tail-logs-stream", handleTailLogsStream(lmc, nc))

	fmt.Printf("Listening on port %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
