package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang/protobuf/jsonpb"
	pb "github.com/webmocha/lumberman/pb"
	"google.golang.org/grpc"
)

var (
	port          = flag.Int("port", 80, "Listen address")
	lmServerAddr  = flag.String("server_addr", "127.0.0.1:9090", "The Lumberman server address in the format of host:port")
	sb            *Switchboard
	basicAuthUser = os.Getenv("AUTH_USER")
	basicAuthPass = os.Getenv("AUTH_PASS")
)

func main() {
	flag.Parse()

	authHandler := basicAuth

	if basicAuthUser == "" && basicAuthPass == "" {
		authHandler = func(h http.HandlerFunc) http.HandlerFunc {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.ServeHTTP(w, r)
			})
		}
		log.Println("WARNING: Basic auth is disabled.\nEnable auth by setting envs AUTH_USER and AUTH_PASS")
	}

	conn, err := grpc.Dial(*lmServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %s\n%v", *lmServerAddr, err)
	}
	defer conn.Close()

	lmc := &lmClient{
		client: pb.NewLoggerClient(conn),
		m:      &jsonpb.Marshaler{},
	}

	sb = NewSwitchboard(lmc)

	http.HandleFunc("/", authHandler(func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir("ui/public")).ServeHTTP(w, r)
	}))
	http.HandleFunc("/api/list-prefixes", authHandler(handleListPrefixes(lmc)))
	http.HandleFunc("/api/list-keys", authHandler(handleListKeys(lmc)))
	http.HandleFunc("/api/get-log", authHandler(handleGetLog(lmc)))
	http.HandleFunc("/api/get-logs-stream", authHandler(handleGetLogsStream(lmc)))
	http.HandleFunc("/api/tail-logs-stream", authHandler(handleTailLogsStream(lmc, sb)))

	fmt.Printf("Listening on port %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
