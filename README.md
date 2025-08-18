# Rosenbridge CLI

A command-line interface for connecting to and communicating with Rosenbridge clusters. Rosenbridge CLI provides real-time messaging capabilities between clients through WebSocket connections.

## Installation

Install the latest version using Go:

```shell
go install github.com/shivanshkc/rosenbridge-cli@latest
```

For convenience, add an alias to your shell configuration file (e.g., `~/.bashrc` or `~/.zshrc`):

```shell
alias rosen='rosenbridge-cli'
```

## Usage

### Connect and Listen for Messages

Start listening for incoming messages on a specific client ID:

```shell
rosen connect --client-id obiwan
```

Or using the short flag:

```shell
rosen connect -c obiwan
```

All messages sent to `obiwan` will be displayed on the console with timestamps.

### Send Messages

#### Send a Single Message

Send a message and exit immediately:

```shell
rosen send --sender anakin --receivers obiwan,quigon,yoda --message "when master"
```

Using short flags:

```shell
rosen send -s anakin -r obiwan,quigon,yoda -m "when master"
```

#### Interactive Message Mode

Start an interactive session to send multiple messages:

```shell
rosen send --sender anakin --receivers obiwan,quigon,yoda
```

This opens a console where you can type messages continuously:

```
>> You: <type your message here>
```

## Configuration

### Custom Rosenbridge URL

By default, the CLI connects to `wss://rosenbridge.rudraksh.dev`. To use a different Rosenbridge instance:

```shell
rosen connect --url wss://your-rosenbridge.example.com --client-id obiwan
```

Or using the short flag:

```shell
rosen connect -u wss://your-rosenbridge.example.com -c obiwan
```

The `--url` flag is available for all commands and can be used to specify the Rosenbridge cluster endpoint.

## Commands

### connect

Establishes a connection with Rosenbridge and streams incoming messages.

**Flags:**
- `--client-id, -c` (required): ID of the client making the connection
- `--url, -u` (optional): URL of the Rosenbridge cluster

### send

Sends messages to specified clients.

**Flags:**
- `--sender, -s` (required): ID of the client sending the message
- `--receivers, -r` (required): Comma-separated list of client IDs to receive the message
- `--message, -m` (optional): Message to send. If omitted, opens interactive mode
- `--url, -u` (optional): URL of the Rosenbridge cluster

## Examples

### Basic Communication

Terminal 1 (listener):
```shell
rosen connect -c alice
```

Terminal 2 (sender):
```shell
rosen send -s bob -r alice -m "Hello Alice!"
```

### Multi-recipient Message

```shell
rosen send -s coordinator -r worker1,worker2,worker3 -m "Start processing"
```

### Custom Rosenbridge Instance

```shell
rosen connect -u wss://private-rosenbridge.company.com -c employee123
```

### Interactive Chat Session

```shell
rosen send -s alice -r bob
>> You: How are you?
>> You: Did you finish the project?
```

## Error Handling

The CLI automatically handles common scenarios:

- **Connection failures**: Clear error messages with connection details
- **Server overload**: Automatic retry with exponential backoff (up to 10 retries)
- **Invalid client IDs**: Validation with helpful error messages
- **Malformed URLs**: URL validation before attempting connection

## Getting Help

```shell
rosen --help           # General help
rosen connect --help   # Help for connect command
rosen send --help      # Help for send command
```

## Client ID Requirements

Client IDs must:
- Be between 1 and 100 characters long
- Contain only alphanumeric characters, hyphens, underscores, dots, and @ symbols
- Match the pattern: `^[a-zA-Z0-9-@._]*$`

## Troubleshooting

### Connection Issues

If you see "Failed to connect" errors:

1. Verify the Rosenbridge URL is correct and accessible
2. Check if the server supports WebSocket connections
3. Ensure your network allows WebSocket traffic

### Server Overload

If you see "server is under load" messages:

- The CLI will automatically retry up to 10 times
- This is common with cloud-hosted instances during cold starts
- Wait a moment and try again if retries fail

### Invalid Client ID

Ensure your client ID follows the requirements listed above. Common issues:
- Contains spaces or special characters
- Too long (over 100 characters)
- Empty or too short