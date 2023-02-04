#!/usr/bin/env python3
import os
from hm_pyhelper.miner_param import get_ethernet_addresses
from jinja2 import Template


CONF_FILE = '/etc/chirpstack-packet-multiplexer/chirpstack-packet-multiplexer.toml'


def populate_conf_file():
    mac_addrs = {'E0': '', 'W0': ''}
    get_ethernet_addresses(mac_addrs)

    mac_address = mac_addrs.get('E0')
    if not mac_address:
        mac_address = mac_addrs.get('W0')

    gateway_id = mac_address.replace(':', '')
    with open(CONF_FILE) as file_:
        template = Template(file_.read())

    rendered = template.render(gateway_id=gateway_id)
    with open(CONF_FILE, "w") as file_:
        file_.write(rendered)


def run_multiplexer():
    os.system('./chirpstack-packet-multiplexer')


def main():
    populate_conf_file()
    run_multiplexer()


if __name__ == '__main__':
    main()
