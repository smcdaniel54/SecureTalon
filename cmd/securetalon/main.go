// Command securetalon runs the SecureTalon API server.
// Auth: Authorization: Bearer <ADMIN_TOKEN> for all /v1/* endpoints.
package main

import (
	"log"
	"net/http"
	"os"

	"securetalon/internal/agent"
	"securetalon/internal/api"
	"securetalon/internal/audit"
	"securetalon/internal/auth"
	"securetalon/internal/broker"
	"securetalon/internal/config"
	"securetalon/internal/core"
	"securetalon/internal/policy"
)

func main() {
	cfg := config.DefaultConfig()
	if cfg.AdminToken == "" {
		log.Fatal("ADMIN_TOKEN is required (env or config)")
	}
	if err := cfg.EnsureDataDirs(); err != nil {
		log.Fatalf("data dirs: %v", err)
	}

	store := core.NewStore()
	auditStore, err := audit.NewStore(cfg.AuditDir())
	if err != nil {
		log.Fatalf("audit store: %v", err)
	}
	tokenSecret := cfg.TokenSecret
	if tokenSecret == "" {
		tokenSecret = cfg.AdminToken
	}
	issuer := policy.NewIssuer(tokenSecret)
	verifier := policy.NewVerifier(tokenSecret)
	policyEngine := policy.NewEngine(issuer)
	brokerSvc := broker.NewBroker(verifier)
	agentLoop := agent.NewAgent(store, policyEngine, brokerSvc, auditStore)
	handlers := &api.Handlers{
		Store:      store,
		Policy:     policyEngine,
		AuditStore: auditStore,
		Agent:      agentLoop,
	}
	router := api.NewRouter(handlers)
	handler := auth.Middleware(cfg.AdminToken)(router)

	log.Printf("SecureTalon listening on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, handler); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
