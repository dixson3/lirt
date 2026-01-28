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

// ProjectsQuery represents the GraphQL projects query
type ProjectsQuery struct {
	Projects struct {
		Nodes []struct {
			ID          string `graphql:"id"`
			Name        string `graphql:"name"`
			Description string `graphql:"description"`
			State       string `graphql:"state"`
			Priority    int    `graphql:"priority"`
			Lead        *struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"lead"`
			CreatedAt string `graphql:"createdAt"`
			UpdatedAt string `graphql:"updatedAt"`
			URL       string `graphql:"url"`
		} `graphql:"nodes"`
	} `graphql:"projects"`
}

// ListProjects fetches all projects
func (c *Client) ListProjects(ctx context.Context) ([]model.Project, error) {
	var query ProjectsQuery
	if err := c.Query(ctx, &query, nil); err != nil {
		return nil, err
	}

	projects := make([]model.Project, 0, len(query.Projects.Nodes))
	for _, node := range query.Projects.Nodes {
		project := model.Project{
			ID:          node.ID,
			Name:        node.Name,
			Description: node.Description,
			State:       node.State,
			Priority:    node.Priority,
			URL:         node.URL,
		}

		if node.Lead != nil {
			project.Lead = &model.User{
				ID:   node.Lead.ID,
				Name: node.Lead.Name,
			}
		}

		projects = append(projects, project)
	}

	return projects, nil
}

// ProjectQuery represents a single project query
type ProjectQuery struct {
	Project struct {
		ID          string `graphql:"id"`
		Name        string `graphql:"name"`
		Description string `graphql:"description"`
		State       string `graphql:"state"`
		Priority    int    `graphql:"priority"`
		Lead        *struct {
			ID          string `graphql:"id"`
			Name        string `graphql:"name"`
			Email       string `graphql:"email"`
			DisplayName string `graphql:"displayName"`
		} `graphql:"lead"`
		Members struct {
			Nodes []struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"nodes"`
		} `graphql:"members"`
		CreatedAt string `graphql:"createdAt"`
		UpdatedAt string `graphql:"updatedAt"`
		URL       string `graphql:"url"`
	} `graphql:"project(id: $id)"`
}

// GetProject fetches a single project by ID
func (c *Client) GetProject(ctx context.Context, id string) (*model.Project, error) {
	variables := map[string]interface{}{
		"id": id,
	}

	var query ProjectQuery
	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	project := &model.Project{
		ID:          query.Project.ID,
		Name:        query.Project.Name,
		Description: query.Project.Description,
		State:       query.Project.State,
		Priority:    query.Project.Priority,
		URL:         query.Project.URL,
	}

	if query.Project.Lead != nil {
		project.Lead = &model.User{
			ID:          query.Project.Lead.ID,
			Name:        query.Project.Lead.Name,
			Email:       query.Project.Lead.Email,
			DisplayName: query.Project.Lead.DisplayName,
		}
	}

	return project, nil
}

// CreateProjectMutation represents the project creation mutation
type CreateProjectMutation struct {
	ProjectCreate struct {
		Success bool `graphql:"success"`
		Project struct {
			ID   string `graphql:"id"`
			Name string `graphql:"name"`
			URL  string `graphql:"url"`
		} `graphql:"project"`
	} `graphql:"projectCreate(input: $input)"`
}

// CreateProjectInput represents input for creating a project
type CreateProjectInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	State       *string `json:"state,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
	LeadID      *string `json:"leadId,omitempty"`
	TeamIDs     *[]string `json:"teamIds,omitempty"`
}

// CreateProject creates a new project
func (c *Client) CreateProject(ctx context.Context, input *CreateProjectInput) (*model.Project, error) {
	variables := map[string]interface{}{
		"input": input,
	}

	var mutation CreateProjectMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, err
	}

	if !mutation.ProjectCreate.Success {
		return nil, fmt.Errorf("failed to create project")
	}

	return &model.Project{
		ID:   mutation.ProjectCreate.Project.ID,
		Name: mutation.ProjectCreate.Project.Name,
		URL:  mutation.ProjectCreate.Project.URL,
	}, nil
}

// UpdateProjectMutation represents the project update mutation
type UpdateProjectMutation struct {
	ProjectUpdate struct {
		Success bool `graphql:"success"`
		Project struct {
			ID   string `graphql:"id"`
			Name string `graphql:"name"`
		} `graphql:"project"`
	} `graphql:"projectUpdate(id: $id, input: $input)"`
}

// UpdateProjectInput represents input for updating a project
type UpdateProjectInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	State       *string `json:"state,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
	LeadID      *string `json:"leadId,omitempty"`
}

// UpdateProject updates an existing project
func (c *Client) UpdateProject(ctx context.Context, id string, input *UpdateProjectInput) error {
	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var mutation UpdateProjectMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}

	if !mutation.ProjectUpdate.Success {
		return fmt.Errorf("failed to update project")
	}

	return nil
}

// ArchiveProjectMutation represents the project archive mutation
type ArchiveProjectMutation struct {
	ProjectArchive struct {
		Success bool `graphql:"success"`
	} `graphql:"projectArchive(id: $id)"`
}

// ArchiveProject archives a project
func (c *Client) ArchiveProject(ctx context.Context, id string) error {
	variables := map[string]interface{}{
		"id": id,
	}

	var mutation ArchiveProjectMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}

	if !mutation.ProjectArchive.Success {
		return fmt.Errorf("failed to archive project")
	}

	return nil
}

// DeleteProjectMutation represents the project deletion mutation
type DeleteProjectMutation struct {
	ProjectDelete struct {
		Success bool `graphql:"success"`
	} `graphql:"projectDelete(id: $id)"`
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(ctx context.Context, id string) error {
	variables := map[string]interface{}{
		"id": id,
	}

	var mutation DeleteProjectMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}

	if !mutation.ProjectDelete.Success {
		return fmt.Errorf("failed to delete project")
	}

	return nil
}

// ProjectIssuesQuery represents issues for a project
type ProjectIssuesQuery struct {
	Project struct {
		Issues struct {
			Nodes []struct {
				ID         string `graphql:"id"`
				Identifier string `graphql:"identifier"`
				Title      string `graphql:"title"`
				State      struct {
					Name string `graphql:"name"`
					Type string `graphql:"type"`
				} `graphql:"state"`
				Assignee *struct {
					Name string `graphql:"name"`
				} `graphql:"assignee"`
			} `graphql:"nodes"`
		} `graphql:"issues"`
	} `graphql:"project(id: $id)"`
}

// ListProjectIssues fetches issues for a project
func (c *Client) ListProjectIssues(ctx context.Context, projectID string) ([]model.Issue, error) {
	variables := map[string]interface{}{
		"id": projectID,
	}

	var query ProjectIssuesQuery
	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	issues := make([]model.Issue, 0, len(query.Project.Issues.Nodes))
	for _, node := range query.Project.Issues.Nodes {
		issue := model.Issue{
			ID:         node.ID,
			Identifier: node.Identifier,
			Title:      node.Title,
			State: &model.State{
				Name: node.State.Name,
				Type: node.State.Type,
			},
		}

		if node.Assignee != nil {
			issue.Assignee = &model.User{
				Name: node.Assignee.Name,
			}
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// MilestonesQuery represents the GraphQL milestones query
type MilestonesQuery struct {
	Milestones struct {
		Nodes []struct {
			ID          string  `graphql:"id"`
			Name        string  `graphql:"name"`
			Description string  `graphql:"description"`
			TargetDate  *string `graphql:"targetDate"`
			Project     struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"project"`
			CreatedAt string `graphql:"createdAt"`
		} `graphql:"nodes"`
	} `graphql:"milestones(filter: $filter)"`
}

// ListMilestones fetches milestones, optionally filtered by project
func (c *Client) ListMilestones(ctx context.Context, projectID string) ([]model.Milestone, error) {
	variables := map[string]interface{}{}

	if projectID != "" {
		variables["filter"] = map[string]interface{}{
			"project": map[string]interface{}{
				"id": map[string]interface{}{
					"eq": projectID,
				},
			},
		}
	}

	var query MilestonesQuery
	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	milestones := make([]model.Milestone, 0, len(query.Milestones.Nodes))
	for _, node := range query.Milestones.Nodes {
		milestone := model.Milestone{
			ID:          node.ID,
			Name:        node.Name,
			Description: node.Description,
			Project: &model.Project{
				ID:   node.Project.ID,
				Name: node.Project.Name,
			},
		}

		milestones = append(milestones, milestone)
	}

	return milestones, nil
}

// MilestoneQuery represents a single milestone query
type MilestoneQuery struct {
	Milestone struct {
		ID          string  `graphql:"id"`
		Name        string  `graphql:"name"`
		Description string  `graphql:"description"`
		TargetDate  *string `graphql:"targetDate"`
		Project     struct {
			ID   string `graphql:"id"`
			Name string `graphql:"name"`
		} `graphql:"project"`
		CreatedAt string `graphql:"createdAt"`
	} `graphql:"milestone(id: $id)"`
}

// GetMilestone fetches a single milestone by ID
func (c *Client) GetMilestone(ctx context.Context, id string) (*model.Milestone, error) {
	variables := map[string]interface{}{
		"id": id,
	}

	var query MilestoneQuery
	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	milestone := &model.Milestone{
		ID:          query.Milestone.ID,
		Name:        query.Milestone.Name,
		Description: query.Milestone.Description,
		Project: &model.Project{
			ID:   query.Milestone.Project.ID,
			Name: query.Milestone.Project.Name,
		},
	}

	return milestone, nil
}

// CreateMilestoneMutation represents the milestone creation mutation
type CreateMilestoneMutation struct {
	MilestoneCreate struct {
		Success   bool `graphql:"success"`
		Milestone struct {
			ID   string `graphql:"id"`
			Name string `graphql:"name"`
		} `graphql:"milestone"`
	} `graphql:"milestoneCreate(input: $input)"`
}

// CreateMilestoneInput represents input for creating a milestone
type CreateMilestoneInput struct {
	ProjectID   string  `json:"projectId"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	TargetDate  *string `json:"targetDate,omitempty"`
}

// CreateMilestone creates a new milestone
func (c *Client) CreateMilestone(ctx context.Context, input *CreateMilestoneInput) (*model.Milestone, error) {
	variables := map[string]interface{}{
		"input": input,
	}

	var mutation CreateMilestoneMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, err
	}

	if !mutation.MilestoneCreate.Success {
		return nil, fmt.Errorf("failed to create milestone")
	}

	return &model.Milestone{
		ID:   mutation.MilestoneCreate.Milestone.ID,
		Name: mutation.MilestoneCreate.Milestone.Name,
	}, nil
}

// UpdateMilestoneMutation represents the milestone update mutation
type UpdateMilestoneMutation struct {
	MilestoneUpdate struct {
		Success   bool `graphql:"success"`
		Milestone struct {
			ID   string `graphql:"id"`
			Name string `graphql:"name"`
		} `graphql:"milestone"`
	} `graphql:"milestoneUpdate(id: $id, input: $input)"`
}

// UpdateMilestoneInput represents input for updating a milestone
type UpdateMilestoneInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	TargetDate  *string `json:"targetDate,omitempty"`
}

// UpdateMilestone updates an existing milestone
func (c *Client) UpdateMilestone(ctx context.Context, id string, input *UpdateMilestoneInput) error {
	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var mutation UpdateMilestoneMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}

	if !mutation.MilestoneUpdate.Success {
		return fmt.Errorf("failed to update milestone")
	}

	return nil
}

// DeleteMilestoneMutation represents the milestone deletion mutation
type DeleteMilestoneMutation struct {
	MilestoneDelete struct {
		Success bool `graphql:"success"`
	} `graphql:"milestoneDelete(id: $id)"`
}

// DeleteMilestone deletes a milestone
func (c *Client) DeleteMilestone(ctx context.Context, id string) error {
	variables := map[string]interface{}{
		"id": id,
	}

	var mutation DeleteMilestoneMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}

	if !mutation.MilestoneDelete.Success {
		return fmt.Errorf("failed to delete milestone")
	}

	return nil
}

// MilestoneIssuesQuery represents issues for a milestone
type MilestoneIssuesQuery struct {
	Milestone struct {
		Issues struct {
			Nodes []struct {
				ID         string `graphql:"id"`
				Identifier string `graphql:"identifier"`
				Title      string `graphql:"title"`
				State      struct {
					Name string `graphql:"name"`
					Type string `graphql:"type"`
				} `graphql:"state"`
			} `graphql:"nodes"`
		} `graphql:"issues"`
	} `graphql:"milestone(id: $id)"`
}

// ListMilestoneIssues fetches issues for a milestone
func (c *Client) ListMilestoneIssues(ctx context.Context, milestoneID string) ([]model.Issue, error) {
	variables := map[string]interface{}{
		"id": milestoneID,
	}

	var query MilestoneIssuesQuery
	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	milestoneIssues := make([]model.Issue, 0, len(query.Milestone.Issues.Nodes))
	for _, node := range query.Milestone.Issues.Nodes {
		issue := model.Issue{
			ID:         node.ID,
			Identifier: node.Identifier,
			Title:      node.Title,
			State: &model.State{
				Name: node.State.Name,
				Type: node.State.Type,
			},
		}

		milestoneIssues = append(milestoneIssues, issue)
	}

	return milestoneIssues, nil
}

// InitiativesQuery represents the GraphQL initiatives query
type InitiativesQuery struct {
	Initiatives struct {
		Nodes []struct {
			ID          string `graphql:"id"`
			Name        string `graphql:"name"`
			Description string `graphql:"description"`
			CreatedAt   string `graphql:"createdAt"`
			UpdatedAt   string `graphql:"updatedAt"`
		} `graphql:"nodes"`
	} `graphql:"initiatives"`
}

// ListInitiatives fetches all initiatives
func (c *Client) ListInitiatives(ctx context.Context) ([]model.Initiative, error) {
	var query InitiativesQuery
	if err := c.Query(ctx, &query, nil); err != nil {
		return nil, err
	}

	initiatives := make([]model.Initiative, 0, len(query.Initiatives.Nodes))
	for _, node := range query.Initiatives.Nodes {
		initiatives = append(initiatives, model.Initiative{
			ID:          node.ID,
			Name:        node.Name,
			Description: node.Description,
		})
	}

	return initiatives, nil
}

// InitiativeQuery represents a single initiative query
type InitiativeQuery struct {
	Initiative struct {
		ID          string `graphql:"id"`
		Name        string `graphql:"name"`
		Description string `graphql:"description"`
		Projects    struct {
			Nodes []struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"nodes"`
		} `graphql:"projects"`
		CreatedAt string `graphql:"createdAt"`
		UpdatedAt string `graphql:"updatedAt"`
	} `graphql:"initiative(id: $id)"`
}

// GetInitiative fetches a single initiative by ID
func (c *Client) GetInitiative(ctx context.Context, id string) (*model.Initiative, error) {
	variables := map[string]interface{}{
		"id": id,
	}

	var query InitiativeQuery
	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	initiative := &model.Initiative{
		ID:          query.Initiative.ID,
		Name:        query.Initiative.Name,
		Description: query.Initiative.Description,
	}

	return initiative, nil
}

// CreateInitiativeMutation represents the initiative creation mutation
type CreateInitiativeMutation struct {
	InitiativeCreate struct {
		Success    bool `graphql:"success"`
		Initiative struct {
			ID   string `graphql:"id"`
			Name string `graphql:"name"`
		} `graphql:"initiative"`
	} `graphql:"initiativeCreate(input: $input)"`
}

// CreateInitiativeInput represents input for creating an initiative
type CreateInitiativeInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// CreateInitiative creates a new initiative
func (c *Client) CreateInitiative(ctx context.Context, input *CreateInitiativeInput) (*model.Initiative, error) {
	variables := map[string]interface{}{
		"input": input,
	}

	var mutation CreateInitiativeMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return nil, err
	}

	if !mutation.InitiativeCreate.Success {
		return nil, fmt.Errorf("failed to create initiative")
	}

	return &model.Initiative{
		ID:   mutation.InitiativeCreate.Initiative.ID,
		Name: mutation.InitiativeCreate.Initiative.Name,
	}, nil
}

// UpdateInitiativeMutation represents the initiative update mutation
type UpdateInitiativeMutation struct {
	InitiativeUpdate struct {
		Success    bool `graphql:"success"`
		Initiative struct {
			ID   string `graphql:"id"`
			Name string `graphql:"name"`
		} `graphql:"initiative"`
	} `graphql:"initiativeUpdate(id: $id, input: $input)"`
}

// UpdateInitiativeInput represents input for updating an initiative
type UpdateInitiativeInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateInitiative updates an existing initiative
func (c *Client) UpdateInitiative(ctx context.Context, id string, input *UpdateInitiativeInput) error {
	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var mutation UpdateInitiativeMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}

	if !mutation.InitiativeUpdate.Success {
		return fmt.Errorf("failed to update initiative")
	}

	return nil
}

// ArchiveInitiativeMutation represents the initiative archive mutation
type ArchiveInitiativeMutation struct {
	InitiativeArchive struct {
		Success bool `graphql:"success"`
	} `graphql:"initiativeArchive(id: $id)"`
}

// ArchiveInitiative archives an initiative
func (c *Client) ArchiveInitiative(ctx context.Context, id string) error {
	variables := map[string]interface{}{
		"id": id,
	}

	var mutation ArchiveInitiativeMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}

	if !mutation.InitiativeArchive.Success {
		return fmt.Errorf("failed to archive initiative")
	}

	return nil
}

// DeleteInitiativeMutation represents the initiative deletion mutation
type DeleteInitiativeMutation struct {
	InitiativeDelete struct {
		Success bool `graphql:"success"`
	} `graphql:"initiativeDelete(id: $id)"`
}

// DeleteInitiative deletes an initiative
func (c *Client) DeleteInitiative(ctx context.Context, id string) error {
	variables := map[string]interface{}{
		"id": id,
	}

	var mutation DeleteInitiativeMutation
	if err := c.Mutate(ctx, &mutation, variables); err != nil {
		return err
	}

	if !mutation.InitiativeDelete.Success {
		return fmt.Errorf("failed to delete initiative")
	}

	return nil
}

// InitiativeProjectsQuery represents projects for an initiative
type InitiativeProjectsQuery struct {
	Initiative struct {
		Projects struct {
			Nodes []struct {
				ID    string `graphql:"id"`
				Name  string `graphql:"name"`
				State string `graphql:"state"`
			} `graphql:"nodes"`
		} `graphql:"projects"`
	} `graphql:"initiative(id: $id)"`
}

// ListInitiativeProjects fetches projects for an initiative
func (c *Client) ListInitiativeProjects(ctx context.Context, initiativeID string) ([]model.Project, error) {
	variables := map[string]interface{}{
		"id": initiativeID,
	}

	var query InitiativeProjectsQuery
	if err := c.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	projects := make([]model.Project, 0, len(query.Initiative.Projects.Nodes))
	for _, node := range query.Initiative.Projects.Nodes {
		projects = append(projects, model.Project{
			ID:    node.ID,
			Name:  node.Name,
			State: node.State,
		})
	}

	return projects, nil
}
