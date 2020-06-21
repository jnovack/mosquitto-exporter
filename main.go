package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/jnovack/go-version"
	"github.com/mattn/go-isatty"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	endpoint = flag.String("endpoint", "tcp://mosquitto:1883", "mosquitto message broker endpoint")
	clientID = flag.String("clientID", "random", "mqtt client id")
	username = flag.String("username", "", "username for authentication")
	password = flag.String("password", "", "password for authentication")
	certFile = flag.String("certFile", "", "certificate (in pem format) for user authentication")
	keyFile  = flag.String("keyFile", "", "private key (in pem format) for user authentication")
	port     = flag.Int("port", 9344, "listen port for the http metrics endpoint")
)

func main() {
	flag.Parse()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Mosquitto Exporter</title></head>
			<body>
			<h1>Mosquitto Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})

	runMQTTClient()

	log.Info().Msg("Serving metrics on " + strconv.FormatInt(int64(*port), 10))
	log.Fatal().Err(http.ListenAndServe(":"+strconv.FormatInt(int64(*port), 10), nil))

}

func init() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		// Format using ConsoleWriter if running straight
		zerolog.TimestampFunc = func() time.Time {
			return time.Now().In(time.Local)
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		// Format using JSON if running as a service (or container)
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	prometheus.MustRegister(NewCollector())
}
