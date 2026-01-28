package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	graphql "github.com/hasura/go-graphql-client"
)

const (
	// LinearAPIEndpoint is the global Linear GraphQL API endpoint
	LinearAPIEndpoint = "https://api.linear.app/graphql"
)

// Client wraps the Linear GraphQL client
type Client struct {
	graphql *graphql.Client
	apiKey  string
	http    *http.Client
}

// New creates a new Linear API client
func New(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	c := &Client{
		apiKey: apiKey,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Create GraphQL client with auth
	c.graphql = graphql.NewClient(LinearAPIEndpoint, c.http).
		WithRequestModifier(func(req *http.Request) {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
			req.Header.Set("User-Agent", "lirt/0.1.0")
		})

	return c, nil
}

// Option is a functional option for configuring the client
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.http = httpClient
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.http.Timeout = timeout
	}
}

// Query executes a GraphQL query
func (c *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return c.graphql.Query(ctx, q, variables)
}

// Mutate executes a GraphQL mutation
func (c *Client) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error {
	return c.graphql.Mutate(ctx, m, variables)
}

// GetAPIKey returns the configured API key
func (c *Client) GetAPIKey() string {
	return c.apiKey
}

// MaskAPIKey returns a masked version of an API key
func MaskAPIKey(key string) string {
	if len(key) < 12 {
		return "***"
	}
	return key[:12] + "..."
}
