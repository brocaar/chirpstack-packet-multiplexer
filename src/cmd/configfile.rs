use handlebars::Handlebars;

use crate::config::Configuration;

pub fn run(config: &Configuration) {
    let template = r#"
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
  level="{{ logging.level }}"


# Multiplexer configuration.
[multiplexer]

  # Interface:port of UDP bind.
  #
  # This this is the interface:port on which the Multiplexer will receive
  # data from the gateways.
  bind="{{ multiplexer.bind }}"

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
  {{#each multiplexer.servers}}
  [[multiplexer.server]]
    server="{{this.server}}"
    uplink_only={{this.uplink_only}}
    gateway_id_prefixes=[
      {{#each this.gateway_id_prefixes}}
      "{{this}}",
      {{/each}}
    ]

  {{/each}}


# Monitoring configuration.
[monitoring]

  # Interface:port.
  #
  # If set, this will enable the monitoring endpoints. If not set, the endpoint
  # will be disabled. Endpoints:
  #
  # * /metrics: Exposes Prometheus metrics.
  bind="{{ monitoring.bind }}"
"#;

    let reg = Handlebars::new();
    println!(
        "{}",
        reg.render_template(template, config)
            .expect("Render configfile error")
    );
}
