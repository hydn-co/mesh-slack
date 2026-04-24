# mesh-slack

A mesh collector for Slack integration. Implements standardized mesh collectors
using [mesh-sdk](https://github.com/hydn-co/mesh-sdk) to receive commands and
emit catalog entities.

## Collectors

### `slack_users_collector`

Collects users from Slack workspaces and emits them as catalog entities in the
`activity` space.

### `slack_channels_collector`

Collects channels from Slack workspaces and emits them as catalog entities in the
`activity` space.

## Actions

### `slack_channel_message_post_action`

Posts messages to Slack channels. Configure the target channel in feature options and supply the message body as the runtime action payload.

## Slack app setup

All features require a Slack bot token (`xoxb-...`).

1. Go to [api.slack.com/apps](https://api.slack.com/apps) and click **Create New App → From scratch**
2. Name the app and select your workspace
3. Under **OAuth & Permissions → Scopes → Bot Token Scopes**, add the required scopes for the features you intend to use:
   - `chat:write` — post messages to channels
   - `channels:read` — validate public channel access
   - `groups:read` — validate private channel access
   - `users:read` — collect workspace users
4. Click **Install to Workspace** at the top of the **OAuth & Permissions** page and authorize the app
5. Copy the **Bot User OAuth Token** shown after installation

The credential payload expected by all features:

```json
{"token": "xoxb-your-token-here"}
```

> **Note:** The bot must be invited to any channel before it can post or read from it.
> Use `/invite @your-app-name` in the target channel.

## Requirements

- Go 1.25+

## Quick start

```bash
git clone https://github.com/hydn-co/mesh-slack.git
cd mesh-slack
go test ./... -v
go build ./...
```

## Usage

Generate the feature manifest:

```bash
go run ./cmd -describe
```

List registered features:

```bash
go run ./cmd -list
```

Run with a Unix socket transport:

```bash
go run ./cmd -transport-socket <socket-name>
```

## Contributing

Keep changes small and add unit tests for new behavior.