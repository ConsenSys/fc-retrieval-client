# fc-retrieval-client
Filecoin Secondary Retrieval Market client library.

See also Filecoin Secondary Retrieval Market [gateway](https://github.com/ConsenSys/fc-retrieval-gateway) and [retrieval provider](https://github.com/ConsenSys/fc-retrieval-provider) repositories.

## Start the service

1. Create a config file
Create a `.env` file, using [.env.example](./.env.example) as a reference:
```
cp .env.example .env
```

2. Run command
Run integration tests
```
go run cmd/client/integration-tests.go
```

Available arguments:
- **client-id**: ID of the client. Generated if no value provided
```
go run cmd/client/integration-tests.go --client-id=0123456789
```
- **ttl**: Time to live for the establishment message between client and gateway. Default value is *100*
```
go run cmd/client/integration-tests.go --ttl=100
```
- **log-level**: Logging level. Default value is *info*
```
go run cmd/client/integration-tests.go --log-level=info
```

## Config

Config variables description:

| name              | description               | options       | default                     |
| ----------------- | ------------------------- | ------------- | --------------------------- |
| LOG_LEVEL         | logging level             |               | info                        |
| LOG_TARGET        | logging target            | STDOUT / FILE | STDOUT                      |
| CLIENT_ID         | ID of the client          |               |                             |
| ESTABLISHMENT_TTL | Gateway establishment TTL |               |                             |

## POC1

As an initial experiment, this project implemented a ping client. Usage:
```
bash ./run.sh # runs the ping client in a Docker container
```
or
```
sudo go run ping.go # to ping example.com
```
or
```
sudo go run ping.go <servername> # to ping a given server
```
NB: Raw socket ICMP pings require superuser privileges on most modern operating systems.

TODO: Extend for demonstrating Filecoin client to Gateway protocol.

## POC2

TODO

## MVP

TODO

