package collectors

import (
	"context"
	"fmt"
	"strings"

	"github.com/hydn-co/mesh-sdk/pkg/catalog/entities"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/types"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
	"github.com/hydn-co/mesh-slack/internal/options"
	slackapi "github.com/hydn-co/mesh-slack/internal/slack_api"
	"github.com/hydn-co/mesh-slack/internal/users"
)

// SlackUsersCollector collects users from Slack workspaces and emits them as
// catalog entities.
type SlackUsersCollector struct {
	*connector.TypedFeatureContext[*options.SlackUsersCollectorOptions, *connector.NoPayload]
	token string
	state connectorutil.FeatureState
}

// NewSlackUsersCollector constructs a SlackUsersCollector.
func NewSlackUsersCollector(ctx *connector.TypedFeatureContext[*options.SlackUsersCollectorOptions, *connector.NoPayload]) runner.Feature {
	return &SlackUsersCollector{TypedFeatureContext: ctx}
}

// Init prepares the collector for operation.
func (c *SlackUsersCollector) Init(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	token, err := connectorutil.ExtractAPIKey(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("failed to extract bot token: %w", err)
	}

	c.token = token
	c.state.MarkReady()
	return nil
}

// Start begins collecting users from the Slack workspace.
func (c *SlackUsersCollector) Start(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := c.state.RequireReady(); err != nil {
		return err
	}

	cursor := ""
	for {
		if err := slackapi.EnsureContextActive(ctx); err != nil {
			return err
		}

		result, err := users.ListUsers(ctx, c.token, cursor)
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		for _, user := range result.Members {
			if err := slackapi.EnsureContextActive(ctx); err != nil {
				return err
			}

			account := toAccountEntity(user)
			if err := c.Emit(ctx, account); err != nil {
				return fmt.Errorf("failed to emit account %s: %w", user.ID, err)
			}

			if err := c.emitUserChannelAccounts(ctx, user.ID); err != nil {
				return err
			}
		}

		cursor = result.ResponseMetadata.NextCursor
		if cursor == "" {
			break
		}
	}

	return nil
}

// Stop halts user collection and releases resources.
func (c *SlackUsersCollector) Stop(ctx context.Context) error {
	if err := slackapi.EnsureContextActive(ctx); err != nil {
		return err
	}

	if err := c.state.RequireReady(); err != nil {
		return err
	}

	c.state.Reset()
	c.token = ""
	return nil
}

// ============= Private helpers below =============

func (c *SlackUsersCollector) emitUserChannelAccounts(ctx context.Context, userID string) error {
	seen := make(map[string]struct{})
	cursor := ""

	for {
		if err := slackapi.EnsureContextActive(ctx); err != nil {
			return err
		}

		result, err := users.ListUserConversations(ctx, c.token, userID, cursor)
		if err != nil {
			return fmt.Errorf("failed to list conversations for user %s: %w", userID, err)
		}

		for _, channel := range result.Channels {
			if channel.ID == "" {
				continue
			}

			if _, exists := seen[channel.ID]; exists {
				continue
			}
			seen[channel.ID] = struct{}{}

			channelAccount := &entities.ChannelAccount{
				Metadata:   types.EntityMetadata{Space: spaces.ChannelAccounts},
				ChannelRef: channel.ID,
				AccountRef: userID,
			}

			if err := c.Emit(ctx, channelAccount); err != nil {
				return fmt.Errorf("failed to emit channel account %s/%s: %w", channel.ID, userID, err)
			}
		}

		cursor = result.ResponseMetadata.NextCursor
		if cursor == "" {
			return nil
		}
	}
}

func toAccountEntity(user users.SlackUser) *entities.Account {
	entity := &entities.Account{
		Metadata:    types.EntityMetadata{Space: spaces.Accounts},
		AccountRef:  user.ID,
		AccountType: toAccountType(user),
		Name: firstNonEmpty(
			user.RealName,
			user.Profile.RealName,
			user.Profile.DisplayName,
			user.Name,
		),
		DisplayName: firstNonEmpty(
			user.Profile.DisplayName,
			user.RealName,
			user.Name,
		),
		FirstName:   user.Profile.FirstName,
		LastName:    user.Profile.LastName,
		Description: strings.TrimSpace(user.Profile.Title),
		Enabled:     !user.Deleted,
	}

	email := strings.TrimSpace(user.Profile.Email)
	if email != "" {
		entity.PrimaryEmail = &types.Email{Address: email}
	}

	return entity
}

func toAccountType(user users.SlackUser) types.AccountType {
	if user.IsRestricted || user.IsUltraRestricted {
		return types.AccountTypeGuest
	}

	if user.IsBot || user.IsAppUser {
		return types.AccountTypeServicePrincipal
	}

	return types.AccountTypeUser
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}

	return ""
}
