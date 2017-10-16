## MUTE Authentication Proxy

The aim of this project is to provide an OAUTH provider and a proxy for ConiksServer.
These two features are requirements for MUTE with end-to-end encryption.

## Cloning and installing

```sh
go get github.com/coast-team/mute_auth_proxy
```

## Configuration

Fill in the `config.toml` with the required information (OAUTH client ID, client secret ...)

## Launch it

```
mute-auth-proxy
```
