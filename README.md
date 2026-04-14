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

Posts messages to Slack channels based on catalog events.

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