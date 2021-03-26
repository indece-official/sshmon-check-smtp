package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/miekg/dns"
	"gopkg.in/mail.v2"
)

// Variables set during build
var (
	ProjectName  string
	BuildVersion string
	BuildDate    string
)

var statusMap = []string{
	"OK",
	"WARN",
	"CRIT",
	"UNKNOWN",
}

var (
	flagVersion  = flag.Bool("v", false, "Print the version info and exit")
	flagService  = flag.String("service", "", "Service name (defaults to SMTP_<host>)")
	flagHost     = flag.String("host", "", "SMTP-Server Host")
	flagPort     = flag.Int("port", 465, "SMTP-Port")
	flagSSL      = flag.Bool("ssl", true, "Use SSL for connection")
	flagTimeout  = flag.Int("timeout", 3, "Connect-Timeout in seconds")
	flagUsername = flag.String("username", "", "Username for SMTP-Server")
	flagPassword = flag.String("password", "", "Password for SMTP-Server")
	flagDNS      = flag.String("dns", "", "Use alternate dns server")
)

func resolveDNS(host string) (string, error) {
	c := dns.Client{}
	m := dns.Msg{}

	m.SetQuestion(host+".", dns.TypeA)

	r, _, err := c.Exchange(&m, *flagDNS)
	if err != nil {
		return "", fmt.Errorf("Can't resolve '%s' on %s: %s", host, *flagDNS, err)
	}

	if len(r.Answer) == 0 {
		return "", fmt.Errorf("Can't resolve '%s' on %s: No results", host, *flagDNS)
	}

	aRecord := r.Answer[0].(*dns.A)

	return aRecord.A.String(), nil
}

func main() {
	var err error

	flag.Parse()

	if *flagVersion {
		fmt.Printf("%s %s (Build %s)\n", ProjectName, BuildVersion, BuildDate)
		fmt.Printf("\n")
		fmt.Printf("https://github.com/indece-official/sshmon-check-smtp\n")
		fmt.Printf("\n")
		fmt.Printf("Copyright 2020 by indece UG (haftungsbeschränkt)\n")

		os.Exit(0)

		return
	}

	serviceName := *flagService
	if serviceName == "" {
		serviceName = fmt.Sprintf("SMTP_%s", *flagHost)
	}

	host := *flagHost
	if *flagDNS != "" {
		host, err = resolveDNS(host)
		if err != nil {
			fmt.Printf(
				"2 %s - %s - Error connecting via SMTP to '%s': %s\n",
				serviceName,
				statusMap[2],
				*flagHost,
				err,
			)

			os.Exit(0)

			return
		}
	}

	dialer := mail.NewDialer(host, *flagPort, *flagUsername, *flagPassword)
	dialer.TLSConfig = &tls.Config{
		ServerName: *flagHost,
	}
	dialer.Timeout = time.Duration(*flagTimeout) * time.Second
	dialer.RetryFailure = false
	dialer.SSL = *flagSSL

	closer, err := dialer.Dial()
	if err != nil {
		fmt.Printf(
			"2 %s - %s - Error connecting via SMTP to '%s': %s\n",
			serviceName,
			statusMap[2],
			*flagHost,
			err,
		)

		os.Exit(0)

		return
	}
	defer closer.Close()

	status := 0

	fmt.Printf(
		"%d %s - %s - SMTP-Server on %s:%d is healthy\n",
		status,
		serviceName,
		statusMap[status],
		*flagHost,
		*flagPort,
	)

	os.Exit(0)
}
