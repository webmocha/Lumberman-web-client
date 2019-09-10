package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
	empty "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/webmocha/lumberman/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type lmClient struct {
	client pb.LoggerClient
	m      *jsonpb.Marshaler
}

func (l *lmClient) ListPrefixes() string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prefixesReply, err := l.client.ListPrefixes(ctx, new(empty.Empty))

	if err != nil {
		log.Println(handleCallError("ListPrefixes", err))
	}

	out, mErr := l.m.MarshalToString(prefixesReply)
	if err != nil {
		log.Printf("Error Marshaling ListPrefixes()\n%v\n", mErr)
	}

	return out
}

func (l *lmClient) ListKeys(prefix string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	keysReply, err := l.client.ListKeys(ctx, &pb.PrefixRequest{
		Prefix: prefix,
	})

	if err != nil {
		log.Println(handleCallError("ListKeys", err))
	}

	out, mErr := l.m.MarshalToString(keysReply)
	if err != nil {
		log.Printf("Error Marshaling ListKeys()\n%v\n", mErr)
	}

	return out
}

func (l *lmClient) GetLog(key string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logReply, err := l.client.GetLog(ctx, &pb.KeyMessage{
		Key: key,
	})

	if err != nil {
		log.Println(handleCallError("GetLog", err))
	}

	out, mErr := l.m.MarshalToString(logReply)
	if err != nil {
		log.Printf("Error Marshaling GetLog()\n%v\n", mErr)
	}

	return out
}

func (l *lmClient) GetLogsStream(prefix string, w http.ResponseWriter) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := l.client.GetLogsStream(ctx, &pb.PrefixRequest{
		Prefix: prefix,
	})
	if err != nil {
		log.Println(handleCallError("GetLogsStream", err))
		return
	}

	for {
		logReply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(handleCallError("GetLogsStream.Recv", err))
			return
		}

		mErr := l.m.Marshal(w, logReply)
		if mErr != nil {
			log.Printf("Error Marshaling GetLog()\n%v\n", mErr)
		}
		w.Write([]byte("\n"))
	}
}

func (l *lmClient) TailLogsStream(prefix string, s Subscriber, unsubscribeFn UnsubscribeFunc, w http.ResponseWriter) {
	ctx := context.Background()

	stream, err := l.client.TailLogStream(ctx, &pb.PrefixRequest{
		Prefix: prefix,
	})
	if err != nil {
		log.Println(handleCallError("TailLogStream", err))
		return
	}
	for {
		logReply, err := stream.Recv()
		if err == io.EOF {
			if err := unsubscribeFn(); err != nil {
				log.Println(err.Error())
				http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
				return
			}
			break
		}
		if err != nil {
			log.Println(handleCallError("TailLogStream.Recv", err))
			return
		}
		mErr := l.m.Marshal(w, logReply)
		if mErr != nil {
			log.Printf("Error Marshaling GetLog()\n%v\n", mErr)
		}
		// w.Write([]byte("\n"))

		w.(http.Flusher).Flush()
	}
}

func handleCallError(rpcFunc string, err error) error {
	if s, ok := status.FromError(err); !ok {
		return status.Errorf(codes.Internal, "client.%s <- server Unknown Internal Error('%s')", rpcFunc, s.Message())
	} else {
		return status.Errorf(s.Code(), "client.%s<-server.Error('%s')", rpcFunc, s.Message())
	}
}
