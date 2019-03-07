# LoRa Packet Multiplexer

The LoRa Packet Multiplexer utility forwards the [Semtech packet-forwarder](https://github.com/lora-net/packet_forwarder)
UDP data to multiple endpoints. It makes it possible to connect a single
LoRa gateway to multiple networks. It is part of the [LoRa Server project](https://www.loraserver.io).

## Building

### Binary

It is recommended to run the commands below inside a [Docker Compose](https://docs.docker.com/compose/)
environment.

```bash
docker-compose run --rm packetmultiplexer bash
```

```bash
# build binary
make

# create snapshot release
make snapshot

# run tests
make test
```

### Docker image

```bash
docker build -t IMAGENAME .
```

## Usage

Run `lora-packet-multiplexer --help` for usage information.

## Example configuration

Executing `lora-packet-multiplexer configfile` returns the following configuration
template:

```toml
[general]
# Log level
#
# debug=5, info=4, warning=3, error=2, fatal=1, panic=0
log_level=4


[packet_multiplexer]
# Bind
#
# The interface:port on which the packet-multiplexer will bind for receiving
# data from the packet-forwarder (UDP data).
bind="0.0.0.0:1700"


# Backends
#
# The backends to which the packet-multiplexer will forward the
# packet-forwarder UDP data.
#
# Example:
# [[packet_multiplexer.backend]]
# # Host
# #
# # The host:IP of the backend.
# host="192.16.1.5:1700"
# 
# # Gateway IDs
# #
# # The Gateway IDs to forward data for.
# gateway_ids = [
#   "0101010101010101",
#   "0202020202020202",
# ]
```
