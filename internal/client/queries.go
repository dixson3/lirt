package client

import (
	"context"

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
