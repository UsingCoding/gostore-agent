# gostore-agent

Expose gostore ssh keys via ssh-agent

## Configuration

Configure agent
```toml
# custom path to gostore
gostore = "${HOME}/go/bin/gostore

[stores.store-id] # define store
ids = [
   "github.com/ssh-key", # define keys
]
```

## Installation

### Go

```shell
go install github.com/UsingCoding/gostore-agent/cmd/gostore-agent@latest
```

### Homebrew

```shell
# Add tap
brew tap UsingCoding/public
# Install agent
brew install gostore-agent
# Start service
brew services start usingcoding/public/gostore-agent
```

Use agent
```shell
export SSH_AUTH_SOCK="$(brew --prefix)/var/run/gostore-agent.sock"
```

## Run manually

Serve SSH agent
```shell
gostore-agent ssh -s ${HOME}/.gostore-agent/gostore-agent.sock
```
