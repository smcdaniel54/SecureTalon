// Package agent implements the core behavior (A):
//   - ToolIntent parsing (from POST body intents or last message content as JSON array)
//   - Policy Engine evaluation and capability token issuance
//   - Tool Broker execution for allowed intents
//   - Audit events: policy.intent.received, policy.decision, capability.issued, tool.executed, run.finished
package agent

import (
	"encoding/json"
	"fmt"
	"time"

	"securetalon/internal/audit"
	"securetalon/internal/broker"
	"securetalon/internal/core"
	"securetalon/internal/policy"
)

// Agent runs the loop: intents → policy eval → broker execute → steps + audit.
type Agent struct {
	Store      *core.Store
	Policy     *policy.Engine
	Broker     *broker.Broker
	AuditStore *audit.Store
}

// NewAgent returns an agent with the given dependencies.
func NewAgent(store *core.Store, policyEngine *policy.Engine, b *broker.Broker, auditStore *audit.Store) *Agent {
	return &Agent{
		Store:      store,
		Policy:     policyEngine,
		Broker:     b,
		AuditStore: auditStore,
	}
}

// Run processes the run: resolve intents (from list or parse last message), then for each intent
// evaluate policy, optionally execute via broker, append steps and audit events. Marks run completed/failed.
// On panic, run is marked failed and run.finished is still emitted.
func (a *Agent) Run(sessionID, runID string, intents []core.ToolIntent) {
	run := a.Store.GetRun(runID)
	if run == nil {
		return
	}
	if a.Policy == nil || a.Broker == nil {
		a.finishRun(runID, sessionID, "failed", 0)
		return
	}

	a.Store.UpdateRunStatus(runID, "running", nil, nil)
	var finalStatus = "completed"
	var stepCount int

	defer func() {
		if r := recover(); r != nil {
			finalStatus = "failed"
		}
		a.finishRun(runID, sessionID, finalStatus, stepCount)
	}()

	if len(intents) == 0 {
		intents = a.parseIntentsFromLastMessage(sessionID)
	}

	for i, intent := range intents {
		stepID := core.NewStepID(i + 1)
		a.emitAudit(runID, sessionID, "policy.intent.received", map[string]interface{}{
			"tool": intent.Tool, "step_id": stepID,
		})

		result := a.Policy.Evaluate(intent, sessionID)

		a.emitAudit(runID, sessionID, "policy.decision", map[string]interface{}{
			"decision": string(result.Decision),
			"tool":     intent.Tool,
			"reason":   result.Reason,
			"step_id":  stepID,
		})

		if result.Decision == core.DecisionAllow && result.Token != nil {
			a.Store.AppendRunStep(runID, core.Step{
				StepID:  stepID,
				Type:    "policy_eval",
				Status:  "allow",
				Tool:    intent.Tool,
				Details: map[string]interface{}{"reason": result.Reason},
			})
			stepCount++
			a.emitAudit(runID, sessionID, "capability.issued", map[string]interface{}{
				"token_hash": result.Token.Signature,
				"tool":       intent.Tool,
			})

			out, err := a.Broker.Execute(intent, result.Token)
			step := core.Step{
				StepID:  stepID,
				Type:    "tool_exec",
				Tool:    intent.Tool,
				Status:  "ok",
				Details: map[string]interface{}{"result": out},
			}
			if err != nil {
				step.Status = "error"
				if step.Details == nil {
					step.Details = make(map[string]interface{})
				}
				step.Details["error"] = err.Error()
				finalStatus = "failed"
			}
			a.Store.AppendRunStep(runID, step)
			stepCount++
			a.emitAudit(runID, sessionID, "tool.executed", map[string]interface{}{
				"tool": intent.Tool, "step_id": stepID, "status": step.Status,
			})
		} else {
			a.Store.AppendRunStep(runID, core.Step{
				StepID:  stepID,
				Type:    "policy_eval",
				Status:  "denied",
				Tool:    intent.Tool,
				Details: map[string]interface{}{"reason": result.Reason},
			})
			stepCount++
			finalStatus = "failed"
		}
	}
}

// finishRun sets run status, emits run.finished, and appends an assistant summary message to the session.
func (a *Agent) finishRun(runID, sessionID, status string, stepCount int) {
	ended := time.Now().UTC()
	a.Store.UpdateRunStatus(runID, status, &ended, nil)
	a.emitAudit(runID, sessionID, "run.finished", map[string]interface{}{"status": status})

	summary := fmt.Sprintf("Run %s %s. Steps: %d.", runID, status, stepCount)
	a.Store.AppendMessage(sessionID, "assistant", summary, map[string]string{"run_id": runID})
}

func (a *Agent) emitAudit(runID, sessionID, evType string, data map[string]interface{}) {
	if a.AuditStore == nil {
		return
	}
	ev := &core.AuditEvent{
		SessionID: sessionID,
		RunID:     runID,
		Type:      evType,
		Data:      data,
	}
	_ = a.AuditStore.Append(ev)
}

// parseIntentsFromLastMessage returns intents parsed from the last message content.
// Expected format: JSON array of {"tool": "...", "params": {...}}. Invalid or non-array returns nil.
func (a *Agent) parseIntentsFromLastMessage(sessionID string) []core.ToolIntent {
	msgs, ok := a.Store.GetMessages(sessionID, 1)
	if !ok || len(msgs) == 0 {
		return nil
	}
	content := msgs[len(msgs)-1].Content
	var raw []map[string]interface{}
	if err := json.Unmarshal([]byte(content), &raw); err != nil {
		return nil
	}
	var intents []core.ToolIntent
	for _, m := range raw {
		tool, _ := m["tool"].(string)
		if tool == "" {
			continue
		}
		params, _ := m["params"].(map[string]interface{})
		if params == nil {
			params = make(map[string]interface{})
		}
		intents = append(intents, core.ToolIntent{Tool: tool, Params: params})
	}
	return intents
}
