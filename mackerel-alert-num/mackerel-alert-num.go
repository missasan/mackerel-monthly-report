package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mackerelio/mackerel-client-go"
)

func main() {
	var (
		apikey = flag.String("a", "", "apikey")
	)

	flag.Parse()
	if *apikey == "" {
		fmt.Fprint(os.Stderr, "missing apikey/n")
		os.Exit(1)
	}

	client := mackerel.NewClient(*apikey)
	serviceName := "mackerel-container-agent"

	resp, err := client.FindWithClosedAlerts()
	if err != nil {
		log.Fatalln("Fetch alerts failed: ", err)
	}

	var (
		numCritical int
		numWarning  int
		numOK       int
		numUnknown  int
	)

	for _, a := range resp.Alerts {
		if time.Unix(a.OpenedAt, 0).Month() != time.Now().AddDate(0, -1, 0).Month() {
			continue
		}

		switch a.Status {
		case "CRITICAL":
			numCritical++
		case "WARNING":
			numWarning++
		case "OK":
			numOK++
		case "UNKNOWN":
			numUnknown++
		}
	}

	metricValues := []*mackerel.MetricValue{
		{
			Name:  "alerts.total_num",
			Time:  time.Now().Unix(),
			Value: len(resp.Alerts),
		},
		{
			Name:  "alerts.critical_num",
			Time:  time.Now().Unix(),
			Value: numCritical,
		},
		{
			Name:  "alerts.warning_num",
			Time:  time.Now().Unix(),
			Value: numWarning,
		},
		{
			Name:  "alerts.ok_num",
			Time:  time.Now().Unix(),
			Value: numOK,
		},
		{
			Name:  "alerts.unknown_num",
			Time:  time.Now().Unix(),
			Value: numUnknown,
		},
	}
	err = client.PostServiceMetricValues(serviceName, metricValues)

	if err != nil {
		log.Fatalln("Post service metric failed: ", err)
	}

}
