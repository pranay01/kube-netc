package collector

import (
	"fmt"

	"github.com/nirmata/kube-netc/pkg/cluster"
	"github.com/nirmata/kube-netc/pkg/tracker"
	"github.com/prometheus/client_golang/prometheus"
)

func StartCollector(tr *tracker.Tracker, ci *cluster.ClusterInfo) {
	for {
		select {
		case update := <-tr.NodeUpdateChan:
			ActiveConnections.Set(float64(update.NumConnections))

		case update := <-tr.ConnUpdateChan:

			var labels prometheus.Labels
			conn := update.Connection

			sourceFoundName := ci.PodIPMap[conn.SAddr]
			destinationFoundName := ci.PodIPMap[conn.DAddr]
			
			labels = prometheus.Labels{
				"source_pod_name":            sourceFoundName.Name,
				"destination_pod_name":       destinationFoundName.Name,
				"source_address":      tracker.IPPort(conn.SAddr, conn.SPort),
				"destination_address": tracker.IPPort(conn.DAddr, conn.DPort),
			}

			BytesSent.With(labels).Set(float64(update.Data.BytesSent))
			BytesRecv.With(labels).Set(float64(update.Data.BytesRecv))
			BytesSentPerSecond.With(labels).Set(float64(update.Data.BytesSentPerSecond))
			BytesRecvPerSecond.With(labels).Set(float64(update.Data.BytesRecvPerSecond))

		case update := <-ci.MapUpdateChan:

			if update.Info != nil {
				fmt.Printf("Pod Update: %s -> %s\n", update.IP, update.Info.Name)

			} else {
				fmt.Printf("Pod Deleted: %s\n", update.IP)
			}

		}
	}
}
