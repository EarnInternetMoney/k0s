package hub

import (
	"net/http"
)

func NewServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/agent/", static)
	mux.HandleFunc("/api/agents/", getAgents)
	mux.HandleFunc("/api/new", newAgentSlave)
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
