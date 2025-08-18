# TODO

## Critical
- Update README.md to remove references to .rosen.yaml configuration
- Document new --url flag usage in README
- Add unit tests for validation functions
- Add unit tests for helper functions
- Add godoc comments to all exported functions

## High Priority
- Replace select{} with proper signal handling in connect command
- Make retry count configurable via CLI flag
- Add WebSocket connection timeout configuration
- Add integration tests for WebSocket connections
- Implement graceful shutdown on SIGINT/SIGTERM

## Medium Priority
- Add automatic reconnection for WebSocket connections
- Add verbose/debug logging mode
- Add connection status monitoring
- Add TLS certificate validation options
- Remove unused constants marked with nolint:unused

## Low Priority
- Add shell completion support
- Add Docker support
- Add version command
- Add connection testing command
- Add troubleshooting documentation