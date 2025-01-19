# gostore-agent

Expose gostore ssh keys via ssh-agent

## Run

Configure agent
```toml
[stores.store-id] # define store
ids = [
   "github.com/ssh-key", # define keys
]
```

Serve SSH agent
```shell
gostore-agent ssh -s ${HOME}/.gostore-agent/gostore-agent.sock
```

## Installation

### Go

```shell
go install github.com/UsingCoding/gostore-agent/cmd/gostore-agent@latest
```
