# mosquitto-exporter

Prometheus exporter for the [Mosquitto MQTT message broker](https://mosquitto.org/).

There is a docker image available:
```
docker run \
  -p 9234:9234 jnovack/mosquitto-exporter \
  -broker tcp://mosquitto:1883
```

## Features

### TLS and Websocket Endpoints

MQTT endpoints can be plaintext (`tcp://mosquitto:1883`) or TLS-enabled
(`tls://mosquitto:8883`). You can also utilize websockets
(`ws://mosquitto:8080`) or TLS-enabled websockets (`wss://mosquitto:8081`).


## Related Projects

This code is heavily based off of [sapcc/mosquitto-exporter](https://github.com/sapcc/mosquitto-exporter/).

I am unhappy with the speed at which the upstream makes changes, and I wanted
another project in Go.  Additionally, rewriting and refactoring helps someone
understand more about a project and technology.

So should you use my fork?  It's your choice!  At the time of this writing, the
upstream project had a number of pull-requests open for months for what I'll
consider are "simple necessary fixes" and they have staled out.  I needed those
fixes today, so I forked it and rewrote it, today.