package main

import (
	pb "github.com/webmocha/lumberman/pb"
	"sync"
)

type UnsubscribeFunc func() error

type Subscriber struct {
	C               chan *pb.LogDetail
	UnsubscribeFunc func(s *Subscriber) error
}

type Switchboard struct {
	lmc           *lmClient
	subscribers   map[string]map[*Subscriber]bool
	subscribersMu *sync.Mutex
}

func NewSwitchboard(lmc *lmClient) *Switchboard {
	return &Switchboard{
		lmc:           lmc,
		subscribers:   map[string]map[*Subscriber]bool{},
		subscribersMu: &sync.Mutex{},
	}
}

func (sb *Switchboard) Unsubscribe(prefix string, s *Subscriber) error {
	sb.subscribersMu.Lock()
	delete(sb.subscribers[prefix], s)
	sb.subscribersMu.Unlock()
	return nil
}

func (sb *Switchboard) Subscribe(prefix string) (*Subscriber, error) {
	s := &Subscriber{
		C: make(chan *pb.LogDetail),
		UnsubscribeFunc: func(s *Subscriber) error {
			sb.subscribersMu.Lock()
			delete(sb.subscribers[prefix], s)
			sb.subscribersMu.Unlock()
			return nil
		},
	}

	sb.subscribersMu.Lock()
	if _, ok := sb.subscribers[prefix]; !ok {
		sb.subscribers[prefix] = make(map[*Subscriber]bool)
		go sb.lmc.TailLogsStream(prefix)
	}
	sb.subscribers[prefix][s] = true
	sb.subscribersMu.Unlock()

	return s, nil
}

func (sb *Switchboard) Broadcast(prefix string, logReply *pb.LogDetail) bool {
	if len(sb.subscribers[prefix]) == 0 {
		return false
	}

	sb.subscribersMu.Lock()
	defer sb.subscribersMu.Unlock()

	for s := range sb.subscribers[prefix] {
		select {
		case s.C <- logReply:
		default:
		}
	}

	return true
}
