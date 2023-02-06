#!/usr/bin/env python3
import os
import json
from hm_pyhelper.miner_param import get_ethernet_addresses
from jinja2 import Template


CONF_FILE = '/etc/chirpstack-packet-multiplexer/chirpstack-packet-multiplexer.toml'
TTN_CONF_FILE = '/var/nebra/ttn_conf.json'


def read_ttn_config():
    try:
        with open(TTN_CONF_FILE) as file_:
            ttn_config = json.loads(file_.read())

        return ttn_config
    except FileNotFoundError:
        # Config file doesn't exist yet, assume TTN not enabled.
        return {
            'ttn_enabled': False,
            'ttn_cluster': 'eu'
        }


def populate_conf_file():
    mac_addrs = {'E0': '', 'W0': ''}
    get_ethernet_addresses(mac_addrs)

    mac_address = mac_addrs.get('E0')
    if not mac_address:
        mac_address = mac_addrs.get('W0')

    gateway_id = mac_address.replace(':', '')
    with open(CONF_FILE) as file_:
        template = Template(file_.read())

    ttn_config = read_ttn_config()
    rendered = template.render(gateway_id=gateway_id, ttn_config=ttn_config)
    with open(CONF_FILE, "w") as file_:
        file_.write(rendered)


def run_multiplexer():
    os.system('./chirpstack-packet-multiplexer')


def main():
    populate_conf_file()
    run_multiplexer()


if __name__ == '__main__':
    main()
