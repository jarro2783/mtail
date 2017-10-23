# Introduction

mtail is intended to run one per machine, and serve as monitoring glue for multiple applications running on that machine.  It runs one or more programs in a 1:1 mapping to those client applications.

## Configuration Overview

mtail is configured through commandline flags.

The `--help` flag will print a list of flags for configuring `mtail`.

(Flags may be prefixed with either `-` or `--`)

Basic flags necessary to start `mtail`:

  * `--logs` is a comma separated list of filenames to extract from, but can also be used multiple times, and each filename can be a [glob pattern](http://godoc.org/path/filepath#Match).
  * `--progs` is a directory path containing [mtail programs](Language). Programs must have the `.mtail` suffix.

mtail runs an HTTP server on port 3903, which can be changed with the `--port` flag.

# Details

## Launching mtail

```
mtail --progs /etc/mtail --logs /var/log/syslog --logs /var/log/ntp/peerstats
```

## Getting the Metrics Out

### Pull based collection

Point your collection tool at `localhost:3903/json` for JSON format metrics.

Prometheus can be directed to the /metrics endpoint for Prometheus text-based format.

### Push based collection

Use the `collectd_socketpath` or `graphite_host_port` flags to enable pushing to a collectd or graphite instance.

Configure collectd on the same machine to use the unixsock plugin, and set `collectd_socketpath` to that unix socket.

```
mtail --progs /etc/mtail --logs /var/log/syslog,/var/log/rsyncd.log --collectd_socketpath=/var/run/collectd-unixsock
```

Set `graphite_host_port` to be the host:port of the carbon server.

```
mtail --progs /etc/mtail --logs /var/log/syslog,/var/log/rsyncd.log --graphite_host_port=localhost:9999
```

Likewise, set `statsd_hostport` to the host:port of the statsd server.

Additionally, the flag `metric_push_interval_seconds` can be used to configure the push frequency.  It defaults to 60, i.e. a push every minute.

## Troubleshooting

Lots of state is logged to the log file, by default in `/tmp/mtail.INFO`.  See [Troubleshooting](Troubleshooting) for more information.
