# Rosenbridge CLI

Rosenbridge CLI provides an easy way to connect with and use a Rosenbridge cluster.

## Installation

Simply execute:
```shell
go install github.com/shivanshkc/rosenbridge-cli@latest
```

To make the commands concise, an alias can be added in your `~/.bashrc`:
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
rosen send -s anakin -r obiwan,quigon,yoda -m 'when master'
```
Here, the sender is `anakin`, the receivers are `obiwan`, `quigon` and `yoda`, and the message is `when master`.
This command sends the message and exits immediately.

If the users need to send multiple messages (more like a chat), then the `-m` flag can be skipped.
The following command will start a shell where messages can be written continuously.
```shell
rosen send -s anakin -r obiwan,quigon,yoda
```

The output will look something like this:
```
$ rosen send -s anakin -r obiwan
>> You: <write here>
```

Execute `rosen --help` for more information.

## Configurations

Rosenbridge CLI accepts a few configuration parameters. They can be provided by creating a `.rosen.yaml` file in the
home directory of the user. The yaml file must have the following structure:

```yaml
---
backend:
  # Base URL of the rosenbridge deployment WITHOUT protcol (http, https, ws, wss etc)
  base_url: rosenbridge.ledgerkeep.com
  # Flag to specify if the target rosenbridge deployment is using TLS.
  is_tls_enabled: true

general:
  # Since the default Rosenbridge cluster (rosenbridge.ledgerkeep.com) runs on GCP free-tier, it occasionally 
  # experiences server cold-start problems. The CLI automatically retries the operation if that's the case. So, we need
  # a max retry count.
  cold_start_retry_count: 10
```

This yaml example is also the default configuration used by the CLI. If users want to specify their own Rosenbridge
deployment, it can be done through the `~/.rosen.yaml` file.
