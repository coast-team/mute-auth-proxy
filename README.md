## MUTE Authentication Proxy

The aim of this project is to provide an OAUTH provider and a proxy for ConiksServer.
These two features are requirements for MUTE with end-to-end encryption.

## Cloning and installing

```sh
sudo apt install golang-go
go get github.com/coast-team/mute-auth-proxy
```

## Configuration

Generate an config file template :

```
mute-auth-proxy init
```

Fill in the `config.toml` with the required information (OAUTH client ID, client secret, keyserver DB file ...)

## Launch it

```
mute-auth-proxy run
```

## Compile it

```
go install # binary will be put in $GOPATH/bin
```

or

```
go build # in working directory
```

## Help ?

```
mute-auth-proxy help
```
