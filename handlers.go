package main

import (
	"fmt"
	"log"
	"net/http"
)

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()

		if basicAuthUser != user || basicAuthPass != pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
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

func handleTailLogsStream(lmc *lmClient, sb *Switchboard) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, ok := r.URL.Query()["prefix"]

		if !ok || len(params[0]) < 1 {
			log.Println("Url Param 'prefix' is missing")
			http.Error(w, `{"error": "Url Param 'prefix' is missing"}`, http.StatusBadRequest)
			return
		}

		prefix := string(params[0])

		s, err := sb.Subscribe(prefix)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
			return
		}

		// SSE Support
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

	Looping:
		for {
			select {
			case <-r.Context().Done():
				if err := sb.Unsubscribe(prefix, s); err != nil {
					log.Println(err.Error())
					http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
					return
				}
				break Looping

			default:
				logReply := <-s.C

				w.Write([]byte("event: log\n"))
				w.Write([]byte("data: "))
				mErr := lmc.m.Marshal(w, logReply)
				if mErr != nil {
					http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
					log.Printf("Error Marshaling TailLogStream.Recv()\n%v\n", mErr)
				}
				w.Write([]byte("\n\n"))
				w.(http.Flusher).Flush()

			}
		}

	}
}
