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
		// numCritical int
		// numWarning  int
		// numOK       int
		// numUnknown  int
		numTotal        int
		numConnectivity int
		numHost         int
		numService      int
		numExternal     int
		numExpression   int
	)

	for _, a := range resp.Alerts {
		if time.Unix(a.OpenedAt, 0).Month() != time.Now().AddDate(0, -1, 0).Month() {
			continue
		}

		numTotal++

		switch a.Type {
		case "connectivity":
			numConnectivity++
		case "host":
			numHost++
		case "service":
			numService++
		case "external":
			numExternal++
		case "expression":
			numExpression++
		}
	}

	nextID := resp.NextID
	for nextID != "" {
		respNext, err2 := client.FindWithClosedAlertsByNextID(nextID)
		if err2 != nil {
			log.Fatalln("Fetch alerts failed: ", err)
		}
		nextID = respNext.NextID
		for _, n := range respNext.Alerts {
			if time.Unix(n.OpenedAt, 0).Month() == time.Now().Month() {
				continue
			}
			if time.Unix(n.OpenedAt, 0).Month() >= time.Now().AddDate(0, -2, 0).Month() {
				nextID = ""
				break
			}
			numTotal++

			switch n.Type {
			case "connectivity":
				numConnectivity++
			case "host":
				numHost++
			case "service":
				numService++
			case "external":
				numExternal++
			case "expression":
				numExpression++
			}
		}
	}

	metricValues := []*mackerel.MetricValue{
		{
			Name:  "alerts-last-month.total_num",
			Time:  time.Now().Unix(),
			Value: numTotal,
		},
		{
			Name:  "alerts-last-month.connectivity_num",
			Time:  time.Now().Unix(),
			Value: numConnectivity,
		},
		{
			Name:  "alerts-last-month.host_num",
			Time:  time.Now().Unix(),
			Value: numHost,
		},
		{
			Name:  "alerts-last-month.service_num",
			Time:  time.Now().Unix(),
			Value: numService,
		},
		{
			Name:  "alerts-last-month.external_num",
			Time:  time.Now().Unix(),
			Value: numExternal,
		},
		{
			Name:  "alerts-last-month.expression_num",
			Time:  time.Now().Unix(),
			Value: numExpression,
		},
	}

	err = client.PostServiceMetricValues(serviceName, metricValues)

	if err != nil {
		log.Fatalln("Post service metric failed: ", err)
	}

}
