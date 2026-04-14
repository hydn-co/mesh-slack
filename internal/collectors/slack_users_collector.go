package collectors

import (
	"context"

	"github.com/fgrzl/json/polymorphic"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func init() {
	polymorphic.RegisterType[SlackUsersCollectorOptions]()
}

// SlackUsersCollectorOptions configures the Slack users collector.
type SlackUsersCollectorOptions struct {
	// IncludeWorkspaceUser optionally filters to include or exclude the workspace user.
	IncludeWorkspaceUser *bool `json:"include_workspace_user"`
	// IncludeDisabledUsers optionally filters to include or exclude disabled users.
	IncludeDisabledUsers *bool `json:"include_disabled_users"`
	// MaxResults optionally limits the number of users collected.
	MaxResults *int64 `json:"max_results"`
}

func (o *SlackUsersCollectorOptions) GetDiscriminator() string {
	return "mesh://slack/users_collector_options"
}

func (o *SlackUsersCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Activity}
}

func (o *SlackUsersCollectorOptions) GetRequirements() []string {
	return []string{"slack"}
}

// SlackUsersCollector collects users from Slack workspaces and emits them as
// catalog entities.
type SlackUsersCollector struct {
	ctx *connector.TypedFeatureContext[*SlackUsersCollectorOptions]
}

// NewSlackUsersCollector constructs a SlackUsersCollector.
func NewSlackUsersCollector(ctx *connector.TypedFeatureContext[*SlackUsersCollectorOptions]) runner.Feature {
	return &SlackUsersCollector{ctx: ctx}
}

// Init prepares the collector for operation.
func (c *SlackUsersCollector) Init(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return nil
}

// Start begins collecting users from the Slack workspace.
func (c *SlackUsersCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return nil
}

// Stop halts user collection and releases resources.
func (c *SlackUsersCollector) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return nil
}
