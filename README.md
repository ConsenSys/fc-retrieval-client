# fc-retrieval-client
Filecoin Secondary Retrieval Market client library.

See also Filecoin Secondary Retrieval Market [gateway](https://github.com/ConsenSys/fc-retrieval-gateway) and [retrieval provider](https://github.com/ConsenSys/fc-retrieval-provider) repositories.

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

