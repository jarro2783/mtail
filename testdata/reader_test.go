package testdata

import (
	"os"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/google/mtail/metrics"
	"github.com/google/mtail/metrics/datum"
)

var expectedMetrics = map[string][]*metrics.Metric{
	"bytes_total": {
		&metrics.Metric{
			Name:    "bytes_total",
			Program: "reader_test",
			Kind:    metrics.Counter,
			Keys:    []string{"operation"},
			LabelValues: []*metrics.LabelValue{
				&metrics.LabelValue{
					Labels: []string{"sent"},
					Value:  datum.MakeInt(62793673, time.Date(2011, 2, 23, 5, 54, 10, 0, time.UTC))},
				&metrics.LabelValue{
					Labels: []string{"received"},
					Value:  datum.MakeInt(975017, time.Date(2011, 2, 23, 5, 54, 10, 0, time.UTC))}}},
	},
	"connections_total": {
		&metrics.Metric{
			Name:    "connections_total",
			Program: "reader_test",
			Kind:    metrics.Counter,
			Keys:    []string{},
			LabelValues: []*metrics.LabelValue{
				&metrics.LabelValue{
					Value: datum.MakeInt(52, time.Date(2011, 2, 22, 21, 54, 13, 0, time.UTC))}}},
	},
	"connection-time_total": {
		&metrics.Metric{
			Name:    "connection-time_total",
			Program: "reader_test",
			Kind:    metrics.Counter,
			Keys:    []string{},
			LabelValues: []*metrics.LabelValue{
				&metrics.LabelValue{
					Value: datum.MakeInt(1181011, time.Date(2011, 2, 23, 5, 54, 10, 0, time.UTC))}}},
	},
	"transfers_total": {
		&metrics.Metric{
			Name:    "transfers_total",
			Program: "reader_test",
			Kind:    metrics.Counter,
			Keys:    []string{"operation", "module"},
			LabelValues: []*metrics.LabelValue{
				&metrics.LabelValue{
					Labels: []string{"send", "module"},
					Value:  datum.MakeInt(2, time.Date(2011, 2, 23, 5, 50, 32, 0, time.UTC))},
				&metrics.LabelValue{
					Labels: []string{"send", "repo"},
					Value:  datum.MakeInt(25, time.Date(2011, 2, 23, 5, 51, 14, 0, time.UTC))}}},
	},
	"foo": {
		&metrics.Metric{
			Name:        "foo",
			Program:     "reader_test",
			Kind:        metrics.Gauge,
			Keys:        []string{"label"},
			LabelValues: []*metrics.LabelValue{},
		},
	},
	"bar": {
		&metrics.Metric{
			Name:    "bar",
			Program: "reader_test",
			Kind:    metrics.Counter,
			Keys:    []string{},
			LabelValues: []*metrics.LabelValue{
				&metrics.LabelValue{
					Value: datum.MakeInt(0, time.Unix(0, 0)),
				},
			},
		},
	},
	"floaty": {
		&metrics.Metric{
			Name:    "floaty",
			Program: "reader_test",
			Kind:    metrics.Gauge,
			Type:    datum.Float,
			Keys:    []string{},
			LabelValues: []*metrics.LabelValue{
				&metrics.LabelValue{
					Labels: []string{},
					Value:  datum.MakeFloat(37.0, time.Date(2017, 6, 15, 18, 9, 37, 0, time.UTC)),
				},
			},
		},
	},
}

func TestReadTestData(t *testing.T) {
	f, err := os.Open("reader_test.golden")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	store := metrics.NewStore()
	ReadTestData(f, "reader_test", store)
	diff := deep.Equal(expectedMetrics, store.Metrics)
	if diff != nil {
		t.Error(diff)
		t.Logf("store contains %s", store.Metrics)
	}
}
