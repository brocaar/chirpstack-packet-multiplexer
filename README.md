# ChirpStack Packet Multiplexer

The ChirpStack Packet Multiplexer makes it possible to connect gateways using
the [Semtech UDP packet-forwarder protocol](https://github.com/Lora-net/packet_forwarder/blob/master/PROTOCOL.TXT)
to multiple servers, with the option to mark servers as uplink only.

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
sudo echo "deb https://artifacts.chirpstack.io/packages/4.x/deb stable main" | sudo tee /etc/apt/sources.list.d/chirpstack.list
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

### Requirements

Building ChirpStack Packet Multiplexer requires:

* [Nix](https://nixos.org/download.html)
* [Docker](https://www.docker.com/)

#### Nix

Nix is used for setting up the development environment which is used for local
development and for creating the binaries.

If you do not have Nix installed and do not wish to install it, then you can
also replicate the development environment by installing the packages listed
in `shell.nix` manually.

#### Docker

Docker is used by [cross-rs](https://github.com/cross-rs/cross) for cross-compiling,
as well as some of the `make` commands.

### Starting the development shell

Run the following command to start the development shell:

```bash
nix-shell
```

### Running tests

Execute the following command to run the tests:

```bash
make test
```

### Building binaries

Execute the following commands to build the ChirpStack Packet Multiplexer binaries
and packages:

```bash
# Only build binaries
make build

# Build binaries + distributable packages.
make dist
```

## Usage

Run `chirpstack-packet-multiplexer --help` for usage information.

## Example configuration

Executing `chirpstack-packet-multiplexer configfile` returns the following configuration
template:

```toml
# Logging settings.
[logging]

  # Log level.
  #
  # Valid options are:
  #   * TRACE
  #   * DEBUG
  #   * INFO
  #   * WARN
  #   * ERROR
  level = "info"


# Multiplexer configuration.
[multiplexer]

  # Interface:port of UDP bind.
  #
  # This this is the interface:port on which the Multiplexer will receive
  # data from the gateways.
  bind = "0.0.0.0:1700"

  # Servers to forward gateway data to.
  #
  # Example configuration:
  # [[multiplexer.server]]

  #   # Hostname:port of the server.
  #   server="example.com:1700"

  #   # Only allow uplink.
  #   #
  #   # If set to true, any downlink will be discarded.
  #   uplink_only=false

  #   # Gateway ID prefix filters.
  #   #
  #   # If not set, data of all gateways will be forwarded. If set, only data
  #   # from gateways with a matching Gateway ID will be forwarded.
  #   #
  #   # Examplex:
  #   # * "0102030405060708/32": Exact match (all 32 bits of the filter must match)
  #   # * "0102030400000000/16": All gateway IDs starting with "01020304" (filter on 16 most significant bits)
  #   gateway_id_prefixes=[]


# Monitoring configuration.
[monitoring]

  # Interface:port.
  #
  # If set, this will enable the monitoring endpoints. If not set, the endpoint
  # will be disabled. Endpoints:
  #
  # * /metrics: Exposes Prometheus metrics.
  bind = ""
```

## Docker Compose example

```
services:
  chirpstack-packet-multiplexer:
    image: chirpstack-packet-multiplexer:4
    command: -c /etc/chirpstack-packet-multiplexer/chirpstack-packet-multiplexer.toml
    ports:
      - 1700:1700/udp
    volumes:
      - ./config:/etc/chirpstack-packet-multiplexer
```

The above example assumes that you have a local configuration directory named
`config` which contains a `chirpstack-packet-multiplexer.toml` file.

## Changelog

### v4.0.0

* Refactor code from Go to Rust.
* Allow Gateway ID prefix filtering.
* Forward all gateways in case `gateway_id_prefixes` is empty.
* Expose Prometheus metrics.

### v3.1.0

This release renames LoRa Packet Multiplexer to ChirpStack Packet Multiplexer.
See the [Rename Announcement](https://www.chirpstack.io/r/rename-announcement) for more information.

### v3.0.2

* Fix setting of configuration variable (used to resolve if backend allows downlink).

### v3.0.1

* Auto-lowercase configured gateway IDs.

### v3.0.0

* Initial release (part of LoRa Server v3 repository).

## License

ChirpStack Packet Multiplexer is distributed under the MIT license. See also
[LICENSE](https://github.com/chirpstack/chirpstack-packet-multiplexer/blob/master/LICENSE).
