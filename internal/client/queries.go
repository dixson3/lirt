package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/dixson3/lirt/internal/model"
)

// ViewerQuery represents the GraphQL viewer query
type ViewerQuery struct {
	Viewer struct {
		ID    string `graphql:"id"`
		Name  string `graphql:"name"`
		Email string `graphql:"email"`
		Organization struct {
			ID     string `graphql:"id"`
			Name   string `graphql:"name"`
			URLKey string `graphql:"urlKey"`
		} `graphql:"organization"`
	} `graphql:"viewer"`
}

// GetViewer fetches the authenticated user and organization
func (c *Client) GetViewer(ctx context.Context) (*model.Viewer, error) {
	var query ViewerQuery
	if err := c.Query(ctx, &query, nil); err != nil {
		return nil, err
	}

	viewer := &model.Viewer{
		ID:    query.Viewer.ID,
		Name:  query.Viewer.Name,
		Email: query.Viewer.Email,
		Organization: &model.Organization{
			ID:     query.Viewer.Organization.ID,
			Name:   query.Viewer.Organization.Name,
			URLKey: query.Viewer.Organization.URLKey,
		},
	}

	return viewer, nil
}

// TeamsQuery represents the GraphQL teams query
type TeamsQuery struct {
	Teams struct {
		Nodes []struct {
			ID          string `graphql:"id"`
			Key         string `graphql:"key"`
			Name        string `graphql:"name"`
			Description string `graphql:"description"`
		} `graphql:"nodes"`
	} `graphql:"teams"`
}

// ListTeams fetches all teams
func (c *Client) ListTeams(ctx context.Context) ([]model.Team, error) {
	var query TeamsQuery
	if err := c.Query(ctx, &query, nil); err != nil {
		return nil, err
	}

	teams := make([]model.Team, 0, len(query.Teams.Nodes))
	for _, node := range query.Teams.Nodes {
		teams = append(teams, model.Team{
			ID:          node.ID,
			Key:         node.Key,
			Name:        node.Name,
			Description: node.Description,
		})
	}

	return teams, nil
}

// IssuesQuery represents the GraphQL issues query with filters
type IssuesQuery struct {
	Issues struct {
		Nodes []struct {
			ID          string `graphql:"id"`
			Identifier  string `graphql:"identifier"`
			Title       string `graphql:"title"`
			Description string `graphql:"description"`
			Priority    int    `graphql:"priority"`
			State       struct {
				ID    string `graphql:"id"`
				Name  string `graphql:"name"`
				Type  string `graphql:"type"`
				Color string `graphql:"color"`
			} `graphql:"state"`
			Assignee *struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"assignee"`
			Team struct {
				ID   string `graphql:"id"`
				Key  string `graphql:"key"`
				Name string `graphql:"name"`
			} `graphql:"team"`
			Labels struct {
				Nodes []struct {
					ID    string `graphql:"id"`
					Name  string `graphql:"name"`
					Color string `graphql:"color"`
				} `graphql:"nodes"`
			} `graphql:"labels"`
			CreatedAt string `graphql:"createdAt"`
			UpdatedAt string `graphql:"updatedAt"`
			URL       string `graphql:"url"`
		} `graphql:"nodes"`
		PageInfo struct {
			HasNextPage bool   `graphql:"hasNextPage"`
			EndCursor   string `graphql:"endCursor"`
		} `graphql:"pageInfo"`
	} `graphql:"issues(filter: $filter, first: $first, after: $after)"`
}

// IssueFilters represents filters for issue queries
type IssueFilters struct {
	TeamID     *string   `json:"team,omitempty"`
	StateID    *string   `json:"state,omitempty"`
	AssigneeID *string   `json:"assignee,omitempty"`
	LabelIDs   *[]string `json:"labels,omitempty"`
	ProjectID  *string   `json:"project,omitempty"`
	Priority   *int      `json:"priority,omitempty"`
	Search     *string   `json:"searchableContent,omitempty"`
}

// ListIssues fetches issues with optional filters
func (c *Client) ListIssues(ctx context.Context, filters *IssueFilters) ([]model.Issue, error) {
	variables := map[string]interface{}{
		"first": 50,
	}

	if filters != nil {
		filterMap := make(map[string]interface{})
		if filters.TeamID != nil {
			filterMap["team"] = map[string]interface{}{"id": map[string]interface{}{"eq": *filters.TeamID}}
		}
		if filters.StateID != nil {
			filterMap["state"] = map[string]interface{}{"id": map[string]interface{}{"eq": *filters.StateID}}
		}
		if filters.AssigneeID != nil {
			filterMap["assignee"] = map[string]interface{}{"id": map[string]interface{}{"eq": *filters.AssigneeID}}
		}
		if filters.Priority != nil {
			filterMap["priority"] = map[string]interface{}{"eq": *filters.Priority}
		}
		if filters.Search != nil && *filters.Search != "" {
			filterMap["searchableContent"] = map[string]interface{}{"containsIgnoreCase": *filters.Search}
		}
		if len(filterMap) > 0 {
			variables["filter"] = filterMap
		}
	}

	var query IssuesQuery
	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	issues := make([]model.Issue, 0, len(query.Issues.Nodes))
	for _, node := range query.Issues.Nodes {
		issue := model.Issue{
			ID:          node.ID,
			Identifier:  node.Identifier,
			Title:       node.Title,
			Description: node.Description,
			Priority:    node.Priority,
			State: &model.State{
				ID:    node.State.ID,
				Name:  node.State.Name,
				Type:  node.State.Type,
				Color: node.State.Color,
			},
			Team: &model.Team{
				ID:   node.Team.ID,
				Key:  node.Team.Key,
				Name: node.Team.Name,
			},
			URL: node.URL,
		}

		if node.Assignee != nil {
			issue.Assignee = &model.User{
				ID:   node.Assignee.ID,
				Name: node.Assignee.Name,
			}
		}

		if len(node.Labels.Nodes) > 0 {
			issue.Labels = make([]model.Label, len(node.Labels.Nodes))
			for i, label := range node.Labels.Nodes {
				issue.Labels[i] = model.Label{
					ID:    label.ID,
					Name:  label.Name,
					Color: label.Color,
				}
			}
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// IssueQuery represents a single issue query
type IssueQuery struct {
	Issue struct {
		ID          string `graphql:"id"`
		Identifier  string `graphql:"identifier"`
		Title       string `graphql:"title"`
		Description string `graphql:"description"`
		Priority    int    `graphql:"priority"`
		State       struct {
			ID    string `graphql:"id"`
			Name  string `graphql:"name"`
			Type  string `graphql:"type"`
			Color string `graphql:"color"`
		} `graphql:"state"`
		Assignee *struct {
			ID          string `graphql:"id"`
			Name        string `graphql:"name"`
			Email       string `graphql:"email"`
			DisplayName string `graphql:"displayName"`
		} `graphql:"assignee"`
		Team struct {
			ID   string `graphql:"id"`
			Key  string `graphql:"key"`
			Name string `graphql:"name"`
		} `graphql:"team"`
		Project *struct {
			ID   string `graphql:"id"`
			Name string `graphql:"name"`
		} `graphql:"project"`
		Labels struct {
			Nodes []struct {
				ID          string `graphql:"id"`
				Name        string `graphql:"name"`
				Color       string `graphql:"color"`
				Description string `graphql:"description"`
			} `graphql:"nodes"`
		} `graphql:"labels"`
		Parent *struct {
			ID         string `graphql:"id"`
			Identifier string `graphql:"identifier"`
			Title      string `graphql:"title"`
		} `graphql:"parent"`
		CreatedAt string `graphql:"createdAt"`
		UpdatedAt string `graphql:"updatedAt"`
		URL       string `graphql:"url"`
	} `graphql:"issue(id: $id)"`
}

// GetIssue fetches a single issue by ID
func (c *Client) GetIssue(ctx context.Context, id string) (*model.Issue, error) {
	variables := map[string]interface{}{
		"id": id,
	}

	var query IssueQuery
	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	issue := &model.Issue{
		ID:          query.Issue.ID,
		Identifier:  query.Issue.Identifier,
		Title:       query.Issue.Title,
		Description: query.Issue.Description,
		Priority:    query.Issue.Priority,
		State: &model.State{
			ID:    query.Issue.State.ID,
			Name:  query.Issue.State.Name,
			Type:  query.Issue.State.Type,
			Color: query.Issue.State.Color,
		},
		Team: &model.Team{
			ID:   query.Issue.Team.ID,
			Key:  query.Issue.Team.Key,
			Name: query.Issue.Team.Name,
		},
		URL: query.Issue.URL,
	}

	if query.Issue.Assignee != nil {
		issue.Assignee = &model.User{
			ID:          query.Issue.Assignee.ID,
			Name:        query.Issue.Assignee.Name,
			Email:       query.Issue.Assignee.Email,
			DisplayName: query.Issue.Assignee.DisplayName,
		}
	}

	if query.Issue.Project != nil {
		issue.Project = &model.Project{
			ID:   query.Issue.Project.ID,
			Name: query.Issue.Project.Name,
		}
	}

	if len(query.Issue.Labels.Nodes) > 0 {
		issue.Labels = make([]model.Label, len(query.Issue.Labels.Nodes))
		for i, label := range query.Issue.Labels.Nodes {
			issue.Labels[i] = model.Label{
				ID:          label.ID,
				Name:        label.Name,
				Color:       label.Color,
				Description: label.Description,
			}
		}
	}

	return issue, nil
}

// ResolveIssueID resolves an issue identifier (ENG-123 or UUID) to an ID
func (c *Client) ResolveIssueID(ctx context.Context, identifier string) (string, error) {
	// If it looks like a UUID, return as-is
	if len(identifier) > 20 && !strings.Contains(identifier, "-") {
		return identifier, nil
	}

	// Otherwise treat as identifier (e.g., ENG-123)
	type IdentifierQuery struct {
		Issue struct {
			ID string `graphql:"id"`
		} `graphql:"issue(filter: {identifier: {eq: $identifier}})"`
	}

	variables := map[string]interface{}{
		"identifier": identifier,
	}

	var query IdentifierQuery
	if err := c.Query(ctx, &query, variables); err != nil {
		return "", fmt.Errorf("failed to resolve issue %s: %w", identifier, err)
	}

	if query.Issue.ID == "" {
		return "", fmt.Errorf("issue not found: %s", identifier)
	}

	return query.Issue.ID, nil
}

// CreateIssueMutation represents the issue creation mutation
type CreateIssueMutation struct {
	IssueCreate struct {
		Success bool `graphql:"success"`
		Issue   struct {
			ID         string `graphql:"id"`
			Identifier string `graphql:"identifier"`
			Title      string `graphql:"title"`
			URL        string `graphql:"url"`
		} `graphql:"issue"`
	} `graphql:"issueCreate(input: $input)"`
}

// CreateIssueInput represents input for creating an issue
type CreateIssueInput struct {
	TeamID      string  `json:"teamId"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
	StateID     *string `json:"stateId,omitempty"`
	AssigneeID  *string `json:"assigneeId,omitempty"`
	ProjectID   *string `json:"projectId,omitempty"`
	ParentID    *string `json:"parentId,omitempty"`
	LabelIDs    *[]string `json:"labelIds,omitempty"`
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(ctx context.Context, input *CreateIssueInput) (*model.Issue, error) {
	variables := map[string]interface{}{
		"input": input,
	}

	var mutation CreateIssueMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, err
	}

	if !mutation.IssueCreate.Success {
		return nil, fmt.Errorf("failed to create issue")
	}

	return &model.Issue{
		ID:         mutation.IssueCreate.Issue.ID,
		Identifier: mutation.IssueCreate.Issue.Identifier,
		Title:      mutation.IssueCreate.Issue.Title,
		URL:        mutation.IssueCreate.Issue.URL,
	}, nil
}

// UpdateIssueMutation represents the issue update mutation
type UpdateIssueMutation struct {
	IssueUpdate struct {
		Success bool `graphql:"success"`
		Issue   struct {
			ID         string `graphql:"id"`
			Identifier string `graphql:"identifier"`
			Title      string `graphql:"title"`
		} `graphql:"issue"`
	} `graphql:"issueUpdate(id: $id, input: $input)"`
}

// UpdateIssueInput represents input for updating an issue
type UpdateIssueInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
	StateID     *string `json:"stateId,omitempty"`
	AssigneeID  *string `json:"assigneeId,omitempty"`
	ProjectID   *string `json:"projectId,omitempty"`
	ParentID    *string `json:"parentId,omitempty"`
	LabelIDs    *[]string `json:"labelIds,omitempty"`
}

// UpdateIssue updates an existing issue
func (c *Client) UpdateIssue(ctx context.Context, id string, input *UpdateIssueInput) error {
	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var mutation UpdateIssueMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}

	if !mutation.IssueUpdate.Success {
		return fmt.Errorf("failed to update issue")
	}

	return nil
}

// ArchiveIssueMutation represents the issue archive mutation
type ArchiveIssueMutation struct {
	IssueArchive struct {
		Success bool `graphql:"success"`
	} `graphql:"issueArchive(id: $id)"`
}

// ArchiveIssue archives an issue
func (c *Client) ArchiveIssue(ctx context.Context, id string) error {
	variables := map[string]interface{}{
		"id": id,
	}

	var mutation ArchiveIssueMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}

	if !mutation.IssueArchive.Success {
		return fmt.Errorf("failed to archive issue")
	}

	return nil
}

// DeleteIssueMutation represents the issue deletion mutation
type DeleteIssueMutation struct {
	IssueDelete struct {
		Success bool `graphql:"success"`
	} `graphql:"issueDelete(id: $id)"`
}

// DeleteIssue deletes an issue
func (c *Client) DeleteIssue(ctx context.Context, id string) error {
	variables := map[string]interface{}{
		"id": id,
	}

	var mutation DeleteIssueMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}

	if !mutation.IssueDelete.Success {
		return fmt.Errorf("failed to delete issue")
	}

	return nil
}

// WorkflowStatesQuery represents the workflow states query for a team
type WorkflowStatesQuery struct {
	WorkflowStates struct {
		Nodes []struct {
			ID       string  `graphql:"id"`
			Name     string  `graphql:"name"`
			Type     string  `graphql:"type"`
			Color    string  `graphql:"color"`
			Position float64 `graphql:"position"`
		} `graphql:"nodes"`
	} `graphql:"workflowStates(filter: {team: {id: {eq: $teamId}}})"`
}

// ListWorkflowStates fetches workflow states for a team
func (c *Client) ListWorkflowStates(ctx context.Context, teamID string) ([]model.State, error) {
	variables := map[string]interface{}{
		"teamId": teamID,
	}

	var query WorkflowStatesQuery
	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	states := make([]model.State, 0, len(query.WorkflowStates.Nodes))
	for _, node := range query.WorkflowStates.Nodes {
		states = append(states, model.State{
			ID:       node.ID,
			Name:     node.Name,
			Type:     node.Type,
			Color:    node.Color,
			Position: int(node.Position),
		})
	}

	return states, nil
}
