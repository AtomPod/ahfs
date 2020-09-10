package metrics

import (
	"github.com/czhj/ahfs/models"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "ahfs_"

type Collector struct {
	Users         *prometheus.Desc
	Files         *prometheus.Desc
	TotalFileSize *prometheus.Desc
}

func NewCollector() Collector {
	return Collector{
		Users: prometheus.NewDesc(
			namespace+"users",
			"Number of users",
			nil, nil,
		),
		Files: prometheus.NewDesc(
			namespace+"files",
			"Number of files",
			nil, nil,
		),
		TotalFileSize: prometheus.NewDesc(
			namespace+"total_file_size",
			"Numer of total file size",
			nil, nil,
		),
	}
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Users
	ch <- c.Files
	ch <- c.TotalFileSize
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	stats := models.GetStatistic()

	ch <- prometheus.MustNewConstMetric(
		c.Users,
		prometheus.GaugeValue,
		float64(stats.Counter.User),
	)

	ch <- prometheus.MustNewConstMetric(
		c.Files,
		prometheus.GaugeValue,
		float64(stats.Counter.File),
	)

	ch <- prometheus.MustNewConstMetric(
		c.TotalFileSize,
		prometheus.GaugeValue,
		float64(stats.Counter.TotalFileSize))
}
