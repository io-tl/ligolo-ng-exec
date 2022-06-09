package main

import (
	"flag"
	"github.com/desertbit/grumble"
	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"ligolo-ng/cmd/proxy/app"
	"ligolo-ng/pkg/proxy"
	"os"
	"strings"
)

func main() {
	var allowDomains []string
	var verboseFlag = flag.Bool("v", false, "enable verbose mode")
	var tunInterface = flag.String("tun", "ligolo", "tun interface name")
	var listenInterface = flag.String("laddr", "0.0.0.0:11601", "listening address ")
	var enableAutocert = flag.Bool("autocert", false, "automatically request letsencrypt certificates, requires port 80 to be accessible")
	var enableSelfcert = flag.Bool("selfcert", false, "dynamically generate self-signed certificates")
	var certFile = flag.String("certfile", "certs/cert.pem", "TLS server certificate")
	var keyFile = flag.String("keyfile", "certs/key.pem", "TLS server key")
	var domainWhitelist = flag.String("allow-domains", "", "autocert authorised domains, if empty, allow all domains, multiple domains should be comma-separated.")

	flag.Parse()

	if *verboseFlag {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// If verbose is set, include method and line in log messages
	logrus.SetReportCaller(*verboseFlag)

	if *domainWhitelist != "" {
		allowDomains = strings.Split(*domainWhitelist, ",")
	}

	app.Run(*tunInterface)

	proxyController := proxy.New(proxy.ControllerConfig{
		EnableAutocert:  *enableAutocert,
		EnableSelfcert:  *enableSelfcert,
		Address:         *listenInterface,
		Certfile:        *certFile,
		Keyfile:         *keyFile,
		DomainWhitelist: allowDomains,
	})

	go proxyController.ListenAndServe()

	// Wait for listener
	proxyController.WaitForReady()

	// Agent registration goroutine
	go func() {
		for {
			remoteConn := <-proxyController.Connection

			yamuxConn, err := yamux.Client(remoteConn, nil)
			if err != nil {
				panic(err)
			}

			agent, err := proxy.NewAgent(yamuxConn)
			if err != nil {
				logrus.Errorf("could not register agent, error: %v", err)
				continue
			}

			logrus.WithFields(logrus.Fields{"remote": remoteConn.RemoteAddr(), "name": agent.Name}).Info("Agent joined.")

			if err := app.RegisterAgent(agent); err != nil {
				logrus.Errorf("could not register agent: %s", err.Error())
			}
		}
	}()

	// Grumble doesn't like cli args
	os.Args = []string{}
	grumble.Main(app.App)
}
