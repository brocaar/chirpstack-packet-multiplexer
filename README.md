# ChirpStack Packet Multiplexer

[![CircleCI](https://circleci.com/gh/brocaar/chirpstack-packet-multiplexer.svg?style=svg)](https://circleci.com/gh/brocaar/chirpstack-packet-multiplexer)

The ChirpStack Packet Multiplexer utility forwards the [Semtech packet-forwarder](https://github.com/lora-net/packet_forwarder)
UDP data to multiple endpoints. It makes it possible to connect a single
LoRa gateway to multiple networks. It is part of [ChirpStack](https://www.chirpstack.io).

## Install

### Debian / Ubuntu

ChirpStack provides a repository that is compatible with the
Debian / Ubuntu apt package system. First make sure that both `dirmngr` and
`apt-transport-https` are installed:

```
sudo apt install apt-transport-https dirmngr
```

Set up the key for this new repository:

```
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 1CE2AFD36DBCCA00
```

Add the repository to the repository list by creating a new file:

```
sudo echo "deb https://artifacts.chirpstack.io/packages/3.x/deb stable main" | sudo tee /etc/apt/sources.list.d/chirpstack.list
```

Update the apt package cache and install `chirpstack-packet-multiplexer`:

```
sudo apt update
sudo apt install chirpstack-packet-multiplexer
```

To complete the installation, update the configuration file which is located
at `/etc/chirpstack-packet-multiplexer/chirpstack-packet-multiplexer.toml` and (re)start
the service:

```
sudo systemctl restart chirpstack-packet-multiplexer
```

## Building from source

### Binary

It is recommended to run the commands below inside a [Docker Compose](https://docs.docker.com/compose/)
environment.

```bash
docker-compose run --rm chirpstack-packet-multiplexer bash
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

Run `chirpstack-packet-multiplexer --help` for usage information.

## Example configuration

Executing `chirpstack-packet-multiplexer configfile` returns the following configuration
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
# # Uplink only
#
# # This backend is for uplink only. It is not able to send data
# # back to the gateways.
# uplink_only=false
# 
# # Gateway IDs
# #
# # The Gateway IDs to forward data for.
# gateway_ids = [
#   "0101010101010101",
#   "0202020202020202",
# ]
```

## Changelog

### v3.1.0

This release renames LoRa Packet Multiplexer to ChirpStack Packet Multiplexer.
See the [Rename Announcement](https://www.chirpstack.io/r/rename-announcement) for more information.

### v3.0.2

* Fix setting of configuration variable (used to resolve if backend allows downlink).

### v3.0.1

* Auto-lowercase configured gateway IDs.

### v3.0.0

* Initial release (part of LoRa Server v3 repository).
