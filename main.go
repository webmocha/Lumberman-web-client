package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	pb "github.com/webmocha/lumberman/pb"
	"google.golang.org/grpc"
)

var (
	port         = flag.Int("port", 80, "Listen address")
	lmServerAddr = flag.String("server_addr", "127.0.0.1:9090", "The Lumberman server address in the format of host:port")
	sb           *Switchboard
)

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

	sb = NewSwitchboard(lmc)

	http.Handle("/", http.FileServer(http.Dir("ui/public")))
	http.HandleFunc("/api/list-prefixes", handleListPrefixes(lmc))
	http.HandleFunc("/api/list-keys", handleListKeys(lmc))
	http.HandleFunc("/api/get-log", handleGetLog(lmc))
	http.HandleFunc("/api/get-logs-stream", handleGetLogsStream(lmc))
	http.HandleFunc("/api/tail-logs-stream", handleTailLogsStream(lmc, sb))

	fmt.Printf("Listening on port %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
