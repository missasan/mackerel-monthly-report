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
		downtime float64
	)

	for _, a := range resp.Alerts {
		if time.Unix(a.OpenedAt, 0).Month() != time.Now().Month() {
			continue
		}
		if a.Type != "external" || a.ClosedAt == 0 {
			continue
		}
		t1 := time.Unix(a.OpenedAt, 0)
		t2 := time.Unix(a.ClosedAt, 0)
		duration := t2.Sub(t1)
		downtime += duration.Seconds()
	}

	nextID := resp.NextID
	for nextID != "" {
		respNext, err2 := client.FindWithClosedAlertsByNextID(nextID)
		if err2 != nil {
			log.Fatalln("Fetch alerts failed: ", err)
		}
		nextID = respNext.NextID
		for _, n := range respNext.Alerts {
			if time.Unix(n.OpenedAt, 0).Month() != time.Now().Month() {
				nextID = ""
				break
			}
			if n.Type != "external" || n.ClosedAt == 0 {
				continue
			}
			t1 := time.Unix(n.OpenedAt, 0)
			t2 := time.Unix(n.ClosedAt, 0)
			duration := t2.Sub(t1)
			downtime += duration.Seconds()

		}
	}

	metricValues := []*mackerel.MetricValue{
		{
			Name:  "donwtime.external",
			Time:  time.Now().Unix(),
			Value: downtime,
		},
	}
	err = client.PostServiceMetricValues(serviceName, metricValues)

	if err != nil {
		log.Fatalln("Post service metric failed: ", err)
	}

}
