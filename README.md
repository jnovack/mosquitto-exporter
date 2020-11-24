# mosquitto-exporter

Prometheus exporter for the [Mosquitto MQTT message broker](https://mosquitto.org/).

There is a docker image available:
```
docker run \
  -p 9234:9234 jnovack/mosquitto-exporter \
  -endpoint tcp://mosquitto:1883
```

## Features

### TLS and Websocket Endpoints

MQTT endpoints can be plaintext (`tcp://mosquitto:1883`) or TLS-enabled
(`tls://mosquitto:8883`). You can also utilize websockets
(`ws://mosquitto:8080`) or TLS-enabled websockets (`wss://mosquitto:8081`).

## Sample

You can see all the metrics collected in [collector.go](internal/mqttclient/collector.go)

```bash
# HELP mosquitto_exporter github.com/jnovack/mosquitto-exporter
# TYPE mosquitto_exporter gauge
mosquitto_exporter{version="v0.0.2"} 1
# HELP mqtt_broker_clients_connected $SYS/broker/clients/connected
# TYPE mqtt_broker_clients_connected gauge
mqtt_broker_clients_connected{broker="mqtt.example.local:8883"} 6
# HELP mqtt_broker_messages_stored $SYS/broker/messages/stored
# TYPE mqtt_broker_messages_stored gauge
mqtt_broker_messages_stored{broker="mqtt.example.local:8883"} 70

  ...truncated...

# HELP mqtt_broker_retained_messages_count $SYS/broker/retained messages/count
# TYPE mqtt_broker_retained_messages_count gauge
mqtt_broker_retained_messages_count{broker="mqtt.example.local:8883"} 45
# HELP mqtt_broker_store_messages_bytes $SYS/broker/store/messages/bytes
# TYPE mqtt_broker_store_messages_bytes gauge
mqtt_broker_store_messages_bytes{broker="mqtt.example.local:8883"} 516
# HELP mqtt_broker_store_messages_count $SYS/broker/store/messages/count
# TYPE mqtt_broker_store_messages_count gauge
mqtt_broker_store_messages_count{broker="mqtt.example.local:8883"} 70
# HELP mqtt_broker_subscriptions_count $SYS/broker/subscriptions/count
# TYPE mqtt_broker_subscriptions_count gauge
mqtt_broker_subscriptions_count{broker="mqtt.example.local:8883"} 2
# HELP mqtt_connection_tls mqtt/connection/tls
# TYPE mqtt_connection_tls gauge
mqtt_connection_tls{broker="mqtt.example.local:8883"} 1
```

## Related Projects

This code is heavily based off of [sapcc/mosquitto-exporter](https://github.com/sapcc/mosquitto-exporter/).

I am unhappy with the speed at which the upstream makes changes, and I wanted
another project in Go.  Additionally, rewriting and refactoring helps someone
understand more about a project and technology.

So should you use my fork?  It's your choice!  At the time of this writing, the
upstream project had a number of pull-requests open for months for what I'll
consider are "simple necessary fixes" and they have staled out.  I needed those
fixes today, so I forked it and rewrote it, today.