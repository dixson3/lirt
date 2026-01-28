package model

import "time"

// Issue represents a Linear issue
type Issue struct {
	ID          string    `json:"id"`
	Identifier  string    `json:"identifier"` // e.g., "ENG-123"
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Priority    int       `json:"priority"`              // 0-4
	State       *State    `json:"state,omitempty"`
	Assignee    *User     `json:"assignee,omitempty"`
	Team        *Team     `json:"team,omitempty"`
	Project     *Project  `json:"project,omitempty"`
	Labels      []Label   `json:"labels,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	URL         string    `json:"url,omitempty"`
}

// State represents a workflow state
type State struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"` // triage, backlog, unstarted, started, completed, canceled
	Color    string `json:"color"`
	Position int    `json:"position"`
}

// Team represents a Linear team
type Team struct {
	ID          string `json:"id"`
	Key         string `json:"key"` // e.g., "ENG"
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IssueCount  int    `json:"issueCount,omitempty"`
	MemberCount int    `json:"memberCount,omitempty"`
}

// User represents a Linear user
type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName,omitempty"`
	Active      bool   `json:"active"`
}

// Project represents a Linear project
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	State       string    `json:"state"` // backlog, planned, started, paused, completed, canceled
	Priority    int       `json:"priority,omitempty"`
	Lead        *User     `json:"lead,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	URL         string    `json:"url,omitempty"`
}

// Milestone represents a project milestone
type Milestone struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	TargetDate  *time.Time `json:"targetDate,omitempty"`
	Project     *Project   `json:"project,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// Initiative represents a Linear initiative
type Initiative struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Label represents an issue label
type Label struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description,omitempty"`
}

// Comment represents a comment on an issue, project, or initiative
type Comment struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Cycle represents a development cycle
type Cycle struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Number    int        `json:"number"`
	StartsAt  time.Time  `json:"startsAt"`
	EndsAt    time.Time  `json:"endsAt"`
	Team      *Team      `json:"team,omitempty"`
	Completed bool       `json:"completed"`
}

// Organization represents a Linear workspace/organization
type Organization struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	URLKey  string `json:"urlKey"`
}

// Viewer represents the authenticated user and their organization
type Viewer struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	Organization *Organization `json:"organization,omitempty"`
}

// PriorityLevel represents a priority value
type PriorityLevel struct {
	Value int    `json:"value"`
	Name  string `json:"name"` // urgent, high, medium, low, none
	Label string `json:"label"` // Urgent, High, Medium, Low, No Priority
}

// PageInfo represents cursor pagination information
type PageInfo struct {
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	StartCursor     string `json:"startCursor,omitempty"`
	EndCursor       string `json:"endCursor,omitempty"`
}
