package api

import "time"

// Bead represents a devdash issue/task.
type Bead struct {
	ID           string `json:"id"`
	LocalBeadID  string `json:"localBeadId,omitempty"`
	ProjectID    string `json:"projectId,omitempty"`
	ProjectName  string `json:"projectName,omitempty"`
	Subject      string `json:"subject"`
	Description  string `json:"description,omitempty"`
	Status       string `json:"status"`
	Priority     int    `json:"priority"`
	BeadType     string `json:"beadType"`
	AssignedTo   string `json:"assignedTo,omitempty"`
	ParentBeadID string `json:"parentBeadId,omitempty"`

	BlockedBy []string `json:"blockedBy,omitempty"`
	Blocks    []string `json:"blocks,omitempty"`

	PreInstructions  string           `json:"preInstructions,omitempty"`
	CompletionResult *CompletionResult `json:"completionResult,omitempty"`
	BurnIntelligence *BurnIntelligence `json:"burnIntelligence,omitempty"`

	DueDate          string `json:"dueDate,omitempty"`
	EstimatedMinutes int    `json:"estimatedMinutes,omitempty"`
	StaleSince       string `json:"staleSince,omitempty"`
	StaleMinutes     int    `json:"staleMinutes,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CompletionResult holds metadata from closing an issue.
type CompletionResult struct {
	Summary string `json:"summary,omitempty"`
	PR      string `json:"pr,omitempty"`
	Commit  string `json:"commit,omitempty"`
}

// BurnIntelligence holds scoring data.
type BurnIntelligence struct {
	ComplexityScore    int      `json:"complexityScore"`
	AutomabilityScore  int      `json:"automabilityScore"`
	AutomabilityGrade  string   `json:"automabilityGrade"`
	Factors            []string `json:"factors,omitempty"`
	ScoredAt           string   `json:"scoredAt,omitempty"`
}

// Project represents a devdash project.
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	GithubRepo  string `json:"githubRepo,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

// Job represents an async job.
type Job struct {
	ID              string          `json:"id"`
	BeadID          string          `json:"beadId,omitempty"`
	ProjectID       string          `json:"projectId,omitempty"`
	Status          string          `json:"status"`
	WorkerType      string          `json:"workerType,omitempty"`
	Prompt          string          `json:"prompt,omitempty"`
	OutputLog       string          `json:"output_log,omitempty"`
	Error           string          `json:"error,omitempty"`
	FailureAnalysis *FailureAnalysis `json:"failureAnalysis,omitempty"`
	Result          interface{}     `json:"result,omitempty"`
	CreatedAt       string          `json:"createdAt,omitempty"`
	StartedAt       string          `json:"startedAt,omitempty"`
	CompletedAt     string          `json:"completedAt,omitempty"`
}

// FailureAnalysis holds diagnostic info for failed jobs.
type FailureAnalysis struct {
	Summary    string `json:"summary,omitempty"`
	RootCause  string `json:"rootCause,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
}

// TeamMember represents a project member.
type TeamMember struct {
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status,omitempty"`
}

// CreateBeadRequest is the body for POST /beads.
type CreateBeadRequest struct {
	ProjectID        string `json:"projectId"`
	Subject          string `json:"subject"`
	Description      string `json:"description,omitempty"`
	BeadType         string `json:"beadType,omitempty"`
	Priority         *int   `json:"priority,omitempty"`
	ParentBeadID     string `json:"parentBeadId,omitempty"`
	DueDate          string `json:"dueDate,omitempty"`
	EstimatedMinutes *int   `json:"estimatedMinutes,omitempty"`
}

// UpdateBeadRequest is the body for PATCH /beads/{id}.
type UpdateBeadRequest struct {
	Subject          *string `json:"subject,omitempty"`
	Description      *string `json:"description,omitempty"`
	Status           *string `json:"status,omitempty"`
	Priority         *int    `json:"priority,omitempty"`
	AssignedTo       *string `json:"assignedTo,omitempty"`
	ParentBeadID     *string `json:"parentBeadId,omitempty"`
	PreInstructions  *string `json:"preInstructions,omitempty"`
	DueDate          *string `json:"dueDate,omitempty"`
	EstimatedMinutes *int    `json:"estimatedMinutes,omitempty"`
}

// CloseBeadRequest is the body for closing a single bead.
type CloseBeadRequest struct {
	Status           string           `json:"status"`
	CompletionResult *CompletionResult `json:"completionResult,omitempty"`
}

// BulkCloseRequest is the body for POST /beads/bulk/close.
type BulkCloseRequest struct {
	Beads []BulkCloseItem `json:"beads"`
}

// BulkCloseItem is a single bead in a bulk close.
type BulkCloseItem struct {
	ID               string           `json:"id"`
	CompletionResult *CompletionResult `json:"completionResult,omitempty"`
}

// ReportRequest is the body for POST /beads/{id}/report.
type ReportRequest struct {
	Status       string `json:"status"`
	Summary      string `json:"summary,omitempty"`
	FilesChanged *int   `json:"filesChanged,omitempty"`
	Branch       string `json:"branch,omitempty"`
	Commit       string `json:"commit,omitempty"`
	Error        string `json:"error,omitempty"`
}

// AddDependencyRequest is the body for POST /beads/{id}/dependencies.
type AddDependencyRequest struct {
	BlockedBy string `json:"blockedBy"`
}

// CommentRequest is the body for POST /beads/{id}/comments.
type CommentRequest struct {
	Body string `json:"body"`
}
