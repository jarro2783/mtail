// Copyright 2015 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package exporter

import (
	"expvar"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/google/mtail/metrics"
)

var (
	metricExportTotal = expvar.NewInt("metric_export_total")
)

const (
	prometheusFormat = "%s{%s} %s\n"
)

func noHyphens(s string) string {
	return strings.Replace(s, "-", "_", -1)
}

// HandlePrometheusMetrics exports the metrics in a format readable by
// Prometheus via HTTP.
func (e *Exporter) HandlePrometheusMetrics(w http.ResponseWriter, r *http.Request) {
	e.store.RLock()
	defer e.store.RUnlock()

	w.Header().Add("Content-type", "text/plain; version=0.0.4")

	for _, ml := range e.store.Metrics {
		emittype := true
		for _, m := range ml {
			m.RLock()
			metricExportTotal.Add(1)

			if emittype {
				fmt.Fprintf(w,
					"# TYPE %s %s\n",
					noHyphens(m.Name),
					kindToPrometheusType(m.Kind))
				emittype = false
			}

			lc := make(chan *metrics.LabelSet)
			go m.EmitLabelSets(lc)
			for l := range lc {
				if m.Source != "" {
					fmt.Fprintf(w, "# %s defined at %s\n", noHyphens(m.Name), m.Source)
				}
				line := metricToPrometheus(e.o, m, l)
				fmt.Fprint(w, line)
			}
			m.RUnlock()
		}
	}
}

func metricToPrometheus(options Options, m *metrics.Metric, l *metrics.LabelSet) string {
	var s []string
	for k, v := range l.Labels {
		// Prometheus quotes the value of each label=value pair.
		s = append(s, fmt.Sprintf("%s=%q", k, v))
	}
	sort.Strings(s)
	if !options.OmitProgLabel {
		s = append(s, fmt.Sprintf("prog=\"%s\"", m.Program))
	}
	return fmt.Sprintf(prometheusFormat,
		noHyphens(m.Name),
		strings.Join(s, ","),
		l.Datum.Value())
}

func kindToPrometheusType(kind metrics.Kind) string {
	if kind != metrics.Timer {
		return strings.ToLower(kind.String())
	}
	return "gauge"
}
