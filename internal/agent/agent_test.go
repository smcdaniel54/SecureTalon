package agent

import (
	"testing"

	"securetalon/internal/audit"
	"securetalon/internal/broker"
	"securetalon/internal/core"
	"securetalon/internal/policy"
)

func TestRun_NoIntents_CompletesWithZeroSteps(t *testing.T) {
	store := core.NewStore()
	sess := store.CreateSession("test", nil)
	store.AppendMessage(sess.ID, "user", "hello", nil) // content not JSON array
	run := store.CreateRun(sess.ID)
	store.SetMessageRunID(sess.ID, run.ID)

	issuer := policy.NewIssuer("secret")
	verifier := policy.NewVerifier("secret")
	policyEngine := policy.NewEngine(issuer)
	brokerSvc := broker.NewBroker(verifier)
	auditStore, _ := audit.NewStore(t.TempDir())

	agent := NewAgent(store, policyEngine, brokerSvc, auditStore)
	agent.Run(sess.ID, run.ID, nil)

	r := store.GetRun(run.ID)
	if r == nil {
		t.Fatal("run not found")
	}
	if r.Status != "completed" {
		t.Fatalf("expected status completed, got %s", r.Status)
	}
	if len(r.Steps) != 0 {
		t.Fatalf("expected 0 steps, got %d", len(r.Steps))
	}
	if r.EndedAt == nil {
		t.Fatal("expected ended_at set")
	}
	msgs, _ := store.GetMessages(sess.ID, 10)
	var foundAssistant bool
	for _, m := range msgs {
		if m.Role == "assistant" && len(m.Content) > 0 {
			foundAssistant = true
			break
		}
	}
	if !foundAssistant {
		t.Fatal("expected assistant summary message")
	}
}

func TestRun_IntentDenied_AppendsStepAndMarksFailed(t *testing.T) {
	store := core.NewStore()
	sess := store.CreateSession("test", nil)
	store.AppendMessage(sess.ID, "user", "run shell", nil)
	run := store.CreateRun(sess.ID)
	store.SetMessageRunID(sess.ID, run.ID)

	issuer := policy.NewIssuer("secret")
	verifier := policy.NewVerifier("secret")
	policyEngine := policy.NewEngine(issuer)
	brokerSvc := broker.NewBroker(verifier)
	auditStore, _ := audit.NewStore(t.TempDir())

	agent := NewAgent(store, policyEngine, brokerSvc, auditStore)
	intents := []core.ToolIntent{{Tool: "shell.exec", Params: map[string]interface{}{}}}
	agent.Run(sess.ID, run.ID, intents)

	r := store.GetRun(run.ID)
	if r == nil {
		t.Fatal("run not found")
	}
	if r.Status != "failed" {
		t.Fatalf("expected status failed, got %s", r.Status)
	}
	if len(r.Steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(r.Steps))
	}
	if r.Steps[0].Type != "policy_eval" || r.Steps[0].Status != "denied" {
		t.Fatalf("expected policy_eval denied step, got %s %s", r.Steps[0].Type, r.Steps[0].Status)
	}
}

func TestRun_IntentFromMessageContent(t *testing.T) {
	store := core.NewStore()
	sess := store.CreateSession("test", nil)
	// Content as JSON array of intents (will be denied without policy override)
	store.AppendMessage(sess.ID, "user", `[{"tool":"file.read","params":{"path":"/work/foo"}}]`, nil)
	run := store.CreateRun(sess.ID)
	store.SetMessageRunID(sess.ID, run.ID)

	issuer := policy.NewIssuer("secret")
	verifier := policy.NewVerifier("secret")
	policyEngine := policy.NewEngine(issuer)
	policyEngine.SetSessionPolicy(sess.ID, &policy.SessionPolicy{
		Overrides: []policy.RuleOverride{
			{Tool: "file.read", Allow: true, Constraints: map[string]interface{}{
				"roots": []string{"/work"},
				"max_bytes": 1024.0,
			}},
		},
	})
	brokerSvc := broker.NewBroker(verifier)
	auditStore, _ := audit.NewStore(t.TempDir())

	agent := NewAgent(store, policyEngine, brokerSvc, auditStore)
	agent.Run(sess.ID, run.ID, nil) // no intents = parse from last message

	r := store.GetRun(run.ID)
	if r == nil {
		t.Fatal("run not found")
	}
	// file.read /work/foo is allowed by policy: we record policy_eval (allow) + tool_exec (ok or error)
	if len(r.Steps) != 2 {
		t.Fatalf("expected 2 steps (policy_eval + tool_exec), got %d", len(r.Steps))
	}
	if r.Steps[0].Type != "policy_eval" || r.Steps[0].Status != "allow" || r.Steps[0].Tool != "file.read" {
		t.Fatalf("expected step 0 policy_eval allow file.read, got %s %s %s", r.Steps[0].Type, r.Steps[0].Status, r.Steps[0].Tool)
	}
	if r.Steps[1].Type != "tool_exec" || r.Steps[1].Tool != "file.read" {
		t.Fatalf("expected step 1 tool_exec file.read, got %s %s", r.Steps[1].Type, r.Steps[1].Tool)
	}
}
