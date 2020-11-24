package mqttclient

import (
	"crypto/tls"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

var (
	gaugeMetrics   = map[string]prometheus.Gauge{}
	counterMetrics = map[string]prometheus.Counter{}
)

// Options contains configurable options for an Client.
type Options struct {
	Endpoint string
	ClientID string
	Username string
	Password string
	CertFile string
	KeyFile  string
	Port     int
}

// NewOptions will create a new ClientOptions type with some
// default values.
func NewOptions() *Options {
	o := &Options{
		Endpoint: "",
		ClientID: "",
		Username: "",
		Password: "",
		CertFile: "",
		KeyFile:  "",
		Port:     0,
	}
	return o
}

// RunMQTTClient sets up and runs a client connection
func RunMQTTClient(mqttOpts *Options) {
	opts := mqtt.NewClientOptions()
	opts.SetCleanSession(true)
	opts.AddBroker(mqttOpts.Endpoint)

	// set client id if it is not random
	if mqttOpts.ClientID != "random" {
		opts.SetClientID(mqttOpts.ClientID)
	} else {
		opts.SetClientID(fmt.Sprintf("mosquitto_exporter_%v", time.Now().Unix()))
	}

	// if you have a username you'll need a password with it
	if mqttOpts.Username != "" {
		opts.SetUsername(mqttOpts.Username)
		if mqttOpts.Password != "" {
			opts.SetPassword(mqttOpts.Password)
		}
	}
	// if you have a client certificate you want a key aswell
	if mqttOpts.CertFile != "" && mqttOpts.KeyFile != "" {
		keyPair, err := tls.LoadX509KeyPair(mqttOpts.CertFile, mqttOpts.KeyFile)
		if err != nil {
			log.Err(err).Msg("Failed to load certificate/keypair")
		}
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{keyPair},
			InsecureSkipVerify: true,
			ClientAuth:         tls.NoClientCert,
		}
		opts.SetTLSConfig(tlsConfig)
		if !strings.HasPrefix(mqttOpts.Endpoint, "ssl://") &&
			!strings.HasPrefix(mqttOpts.Endpoint, "tls://") {
			log.Warn().Msg("Warning: To use TLS the endpoint URL will have to begin with 'ssl://' or 'tls://'")
		}
	} else if (mqttOpts.CertFile != "" && mqttOpts.KeyFile == "") ||
		(mqttOpts.CertFile == "" && mqttOpts.KeyFile != "") {
		log.Warn().Msg("Warning: For TLS to work both certificate and private key are needed. Skipping TLS.")
	}

	opts.OnConnect = func(client mqtt.Client) {
		log.Info().Msgf("Connected to %s", mqttOpts.Endpoint)
		// subscribe on every (re)connect
		token := client.Subscribe("$SYS/#", 0, func(broker mqtt.Client, msg mqtt.Message) {
			processUpdate(msg.Topic(), string(msg.Payload()), broker.OptionsReader())
		})
		if !token.WaitTimeout(10 * time.Second) {
			log.Error().Msg("Error: Timeout subscribing to topic $SYS/#")
		}
		if err := token.Error(); err != nil {
			log.Error().Msgf("Failed to subscribe to topic $SYS/#: %s", err)
		}
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Warn().Msgf("Warning: Connection to %s lost: %s", mqttOpts.Endpoint, err)
	}

	client := mqtt.NewClient(opts)

	// launch the first connection in another thread so it is no blocking
	// and exporter can serve metrics in case of no connection
	go mqttConnect(client, mqttOpts)
}

// try to connect forever with the MQTT broker
func mqttConnect(client mqtt.Client, mqttOpts *Options) {
	// try to connect forever
	for {
		token := client.Connect()
		log.Info().Str("endpoint", mqttOpts.Endpoint).Msg("Attempting to connect to mosquitto endpoint")
		if token.WaitTimeout(5 * time.Second) {
			if token.Error() == nil {
				break
			}
			log.Error().Err(token.Error()).Str("endpoint", mqttOpts.Endpoint).Msg("Failed to connect to mosquitto endpoint")
		} else {
			log.Error().Str("endpoint", mqttOpts.Endpoint).Msg("Timeout connecting to mosquitto endpoint")
		}
		time.Sleep(5 * time.Second)
	}
}

// process the messages received in $SYS/
func processUpdate(topic, payload string, reader mqtt.ClientOptionsReader) {
	//log.Debugf("Got broker update with topic %s and data %s", topic, payload)
	labels := prometheus.Labels{"broker": reader.Servers()[0].Hostname()}

	if _, ok := ignoreKeyMetrics[topic]; !ok {
		if _, ok := counterKeyMetrics[topic]; ok {
			log.Debug().Str("topic", topic).Str("payload", payload).Msg("Processing counter metric")
			processCounterMetric(topic, payload, labels)
		} else {
			log.Debug().Str("topic", topic).Str("payload", payload).Msg("Processing gauge metric")
			processGaugeMetric(topic, payload, labels)
		}
		// restartSecondsSinceLastUpdate()
	} else {
		log.Debug().Str("topic", topic).Str("payload", payload).Msg("Ignoring metric")
	}
}

func processCounterMetric(topic, payload string, labels prometheus.Labels) {
	// if counterMetrics[topic] != nil {
	// 	value := parseValue(payload)
	// 	counterMetrics[topic].  .Set(value)
	// } else {
	// 	// create a mosquitto counter pointer
	// 	mCounter := NewMosquittoCounter(prometheus.NewDesc(
	// 		parseForPrometheus(topic),
	// 		topic,
	// 		[]string{},
	// 		prometheus.Labels{},
	// 	))

	// 	// register the metric
	// 	prometheus.MustRegister(mCounter)
	// 	// add the first value
	// 	value := parseValue(payload)
	// 	counterMetrics[topic].Set(value)
	// }
}

func processGaugeMetric(topic, payload string, labels prometheus.Labels) {
	if gaugeMetrics[topic] == nil {
		gaugeMetrics[topic] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        parseForPrometheus(topic),
			Help:        topic,
			ConstLabels: labels,
		})
		// register the metric
		prometheus.MustRegister(gaugeMetrics[topic])
		// 	// add the first value
	}
	value := parseValue(payload)
	gaugeMetrics[topic].Set(value)
}

func parseForPrometheus(incoming string) string {
	outgoing := strings.Replace(incoming, "$SYS", "mqtt", 1)
	outgoing = strings.Replace(outgoing, "/", "_", -1)
	outgoing = strings.Replace(outgoing, " ", "_", -1)
	outgoing = strings.Replace(outgoing, "-", "_", -1)
	outgoing = strings.Replace(outgoing, ".", "_", -1)
	return outgoing
}

func parseValue(payload string) float64 {
	// fmt.Printf("Payload %s \n", payload)
	var validValue = regexp.MustCompile(`-?\d{1,}[.]\d{1,}|\d{1,}`)
	// get the first value of the string
	strArray := validValue.FindAllString(payload, 1)
	if len(strArray) > 0 {
		// parse to float
		value, err := strconv.ParseFloat(strArray[0], 64)
		if err == nil {
			return value
		}
	}
	return 0
}
