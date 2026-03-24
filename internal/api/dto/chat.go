package dto

type ChatRequest struct {
	SessionID   string `json:"session_id"`
	Query       string `json:"query"`
	Environment string `json:"environment"`
	TimeRange   string `json:"time_range"`
}

type ChatResponse struct {
	RequestID string                 `json:"request_id"`
	Code      int                    `json:"code"`
	Message   string                 `json:"message"`
	Data      *DiagnosticResult      `json:"data,omitempty"`
}

type DiagnosticResult struct {
	Status      string       `json:"status"`
	Summary     string       `json:"summary"`
	Evidence    []Evidence   `json:"evidence"`
}

type Evidence struct {
	Type    string `json:"type"`
	Source  string `json:"source"`
	Summary string `json:"summary"`
}
