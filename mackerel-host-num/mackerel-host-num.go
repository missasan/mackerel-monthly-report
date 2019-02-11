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

	hosts, err := client.FindHosts(&mackerel.FindHostsParam{
		// Service:  serviceName,
		Statuses: []string{mackerel.HostStatusWorking},
	})

	if err != nil {
		log.Fatalln("Fetch alerts failed: ", err)
	}

	fmt.Println(len(hosts))

	metricValues := []*mackerel.MetricValue{
		{
			Name:  "hostnum.working",
			Time:  time.Now().Unix(),
			Value: len(hosts),
		},
	}
	err = client.PostServiceMetricValues(serviceName, metricValues)

	if err != nil {
		log.Fatalln("Post service metric failed: ", err)
	}

}
