package cerbos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"opita-sync-framework/internal/engine/policy"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type checkRequest struct {
	RequestID string         `json:"request_id"`
	Principal map[string]any `json:"principal"`
	Resource  map[string]any `json:"resource"`
	Actions   []string       `json:"actions"`
}

type checkResponse struct {
	Results map[string]struct {
		Effect string `json:"effect"`
	} `json:"results"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Evaluate(input policy.Input) (policy.DecisionRecord, error) {
	reqBody := checkRequest{
		RequestID: input.ExecutionID,
		Principal: map[string]any{
			"id":         input.ExecutionID,
			"roles":      []string{"system"},
			"tenant_id":  input.TenantID,
			"attributes": map[string]any{"classification_level": input.ClassificationLevel},
		},
		Resource: map[string]any{
			"kind": input.ResourceKind,
			"id":   input.ContractID,
			"attributes": map[string]any{
				"classification_level":    input.ClassificationLevel,
				"approval_mode_effective": input.ApprovalModeEffective,
				"risk_level":              input.RiskLevel,
			},
		},
		Actions: []string{input.Action},
	}
	raw, err := json.Marshal(reqBody)
	if err != nil {
		return policy.DecisionRecord{}, fmt.Errorf("marshal cerbos request: %w", err)
	}
	httpReq, err := http.NewRequest(http.MethodPost, c.BaseURL+"/api/check/resources", bytes.NewReader(raw))
	if err != nil {
		return policy.DecisionRecord{}, fmt.Errorf("build cerbos request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return policy.DecisionRecord{}, fmt.Errorf("call cerbos: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return policy.DecisionRecord{}, fmt.Errorf("cerbos status %d", resp.StatusCode)
	}
	var parsed checkResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return policy.DecisionRecord{}, fmt.Errorf("decode cerbos response: %w", err)
	}
	decision := policy.DecisionDenyBlock
	reasonCodes := []string{"policy.deny.cerbos_default"}
	if result, ok := parsed.Results[input.Action]; ok {
		switch strings.ToUpper(result.Effect) {
		case "EFFECT_ALLOW", "ALLOW":
			decision = mapAllowDecision(input)
			reasonCodes = []string{"policy.allow.cerbos"}
		default:
			decision = policy.DecisionDenyBlock
			reasonCodes = []string{"policy.deny.cerbos"}
		}
	}
	return policy.DecisionRecord{
		PolicyDecisionID: fmt.Sprintf("policy-%d", time.Now().UTC().UnixNano()),
		Decision:         decision,
		ReasonCodes:      reasonCodes,
		PolicyVersion:    "cerbos-live",
		TraceID:          input.ExecutionID,
	}, nil
}

func mapAllowDecision(input policy.Input) policy.Decision {
	if input.Action == "view_restricted" && input.ClassificationLevel == "restricted" {
		return policy.DecisionRestrictedView
	}
	if input.Action == "execute" && input.ApprovalModeEffective != "" && input.ApprovalModeEffective != "auto" {
		return policy.DecisionRequireApproval
	}
	if input.Action == "execute" && input.RiskLevel == "critical" {
		return policy.DecisionRequireEscalation
	}
	return policy.DecisionAllow
}
