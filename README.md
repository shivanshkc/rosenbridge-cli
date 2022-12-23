# Rosenbridge CLI

Rosenbridge CLI provides an easy way to connect with and use a Rosenbridge cluster.

## Installation

Simply execute:
```shell
go install github.com/shivanshkc/rosenbridge-cli@latest
```

To make the commands concise, an alias can be made:
```shell
alias rosen='rosenbridge-cli'
```

## Getting started

#### Listen to messages
Execute the following to start listening to all incoming messages:
```shell
rosen connect -c obiwan
```
Now, all messages that are sent to `obiwan` will start getting printed on the console.

#### Send messages
To send a message, execute the following:
```shell
rosen send -s anakin -r obiwan
```
Here, the sender is `anakin` and the receiver is `obiwan`.
Also, this command will start a shell where messages can be written. It will look something like this:
```
$ rosen send -s anakin -r obiwan
>> You: <write here>
```

Execute `rosen --help` for more information.

## Configurations

Rosenbridge CLI requires only two configuration parameters. They can be provided by creating a `.rosen.yaml` file in the home directory of the user. The yaml file must have the following structure:

```yaml
---
backend:
  # Base URL of the rosenbridge deployment WITHOUT protcol (http, https, ws, wss etc)
  base_url: rosenbridge.ledgerkeep.com
  # Flag to specify if the target rosenbridge deployment is using TLS.
  is_tls_enabled: true
```

This yaml example is also the default configuration used by the CLI. If users want to specify their own Rosenbridge deployment, it can be done through the `~/.rosen.yaml` file.
