package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/ncabatoff/go-vmguestlib/vmguestlib"
	"github.com/prometheus/client_golang/prometheus"
)

const metric_name_pfx = "vmwareguest_"

var (
	descs       []*prometheus.Desc
	isGuestDesc = prometheus.NewDesc(
		"vmwareguest_isguest",
		"1 if running on a vmware guest with tools installed, 0 otherwise",
		[]string{},
		nil)
	collecterrsDesc = prometheus.NewDesc(
		"vmwareguest_collecterrors",
		"errors harvesting metrics",
		[]string{},
		nil)
	eventsDesc = prometheus.NewDesc(
		"vmwareguest_events",
		"events e.g. snapshot, vmotion, etc",
		[]string{},
		nil)
)

func init() {
	for _, m := range metrics {
		descs = append(descs, prometheus.NewDesc(
			m.prometheus_name(),
			m.desc,
			[]string{},
			nil))
	}
}

type Collector struct {
	session *vmguestlib.Session
	errors  int
	events  int
}

// Return a new Collector.  If we can't initialize the vmguestlib session,
// we'll return both a collector and an error, but the collector will publish
// only the isguest metric (with a value of 0).
func NewCollector() (*Collector, error) {
	s, err := vmguestlib.NewSession()
	if err != nil {
		s = nil // just to be sure
	}
	return &Collector{session: s}, err
}

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":9263", "Address on which to expose metrics and web interface.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	)
	flag.Parse()

	c, err := NewCollector()
	if err != nil {
		log.Printf("Error creating collector: %v", err)
	}
	prometheus.MustRegister(c)

	http.Handle(*metricsPath, prometheus.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Vmware Guest Exporter</title></head>
			<body>
			<h1>Vmware Guest Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	http.ListenAndServe(*listenAddress, nil)
}

// Collect implements prometheus.Collector.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	isguest := 0
	if c.session != nil {
		isguest = 1
	}
	ch <- prometheus.MustNewConstMetric(isGuestDesc,
		prometheus.GaugeValue,
		float64(isguest))

	if c.session == nil {
		return
	}

	if event, err := c.session.RefreshInfo(); err != nil {
		if err != nil {
			log.Printf("An error occured: %v", err)
		}
		os.Exit(1)
	} else if event {
		c.events++
	}

	for i, m := range metrics {
		val, err := m.Get(c.session)
		if err != nil {
			// log.Printf("error reading %s: %v", m.name, err)
			c.errors++
		} else {
			ch <- prometheus.MustNewConstMetric(descs[i], m.ValueType, val)
		}
	}

	ch <- prometheus.MustNewConstMetric(eventsDesc,
		prometheus.CounterValue,
		float64(c.events))
	ch <- prometheus.MustNewConstMetric(collecterrsDesc,
		prometheus.CounterValue,
		float64(c.errors))
}

// Describe implements prometheus.Collector.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- isGuestDesc
	if c.session == nil {
		return
	}

	for _, d := range descs {
		ch <- d
	}
	ch <- eventsDesc
	ch <- collecterrsDesc
}

type (
	getf64 interface {
		Get(*vmguestlib.Session) (float64, error)
	}
	getu32 func(*vmguestlib.Session) (uint32, error)
	getu64 func(*vmguestlib.Session) (uint64, error)
	metric struct {
		getter     getf64
		name, desc string
		multiplier float64
		unit       string
		prometheus.ValueType
	}
)

func (m metric) Get(s *vmguestlib.Session) (float64, error) {
	val, err := m.getter.Get(s)
	if err != nil {
		return 0, err
	}
	return val * m.multiplier, nil
}

func (g getu32) Get(s *vmguestlib.Session) (float64, error) {
	u32, err := g(s)
	if err != nil {
		return 0, err
	}
	return float64(u32), nil
}

func (g getu64) Get(s *vmguestlib.Session) (float64, error) {
	u64, err := g(s)
	if err != nil {
		return 0, err
	}
	return float64(u64), nil
}

// Create a metric whose underlying value is given in MBs.
func mbmetric(getter getf64, name, desc string) metric {
	return metric{getter, name, desc, 1024 * 1024, "bytes", prometheus.GaugeValue}
}

// Create a metric whose underlying value is given in MHz.
func mzmetric(getter getf64, name, desc string) metric {
	return metric{getter, name, desc, 1024 * 1024, "hertz", prometheus.GaugeValue}
}

// Create a metric whose underlying value is given in some unspecified unit.
func unmetric(getter getf64, name, desc string) metric {
	return metric{getter, name, desc, 1, "", prometheus.GaugeValue}
}

// Create a metric whose underlying value is given in milliseconds.
// Hack: we're assuming that going forward, all millisecond-unit
// metrics the guest lib returns will be counters.
func msmetric(getter getf64, name, desc string) metric {
	return metric{getter, name, desc, 0.001, "seconds", prometheus.CounterValue}
}

func (gm metric) prometheus_name() string {
	name := metric_name_pfx + gm.name
	if gm.unit != "" {
		name += "_" + gm.unit
	}
	return name
}

var metrics = []metric{
	mbmetric(getu64((*vmguestlib.Session).GetHostMemUnmapped), "HostMemUnmapped",
		"total amount of unmapped memory on the host"),
	mbmetric(getu64((*vmguestlib.Session).GetHostMemMapped), "HostMemMapped",
		"total amount of mapped memory on the host"),
	mbmetric(getu64((*vmguestlib.Session).GetHostMemKernOvhd), "HostMemKernOvhd",
		"total amount of host kernel memory overhead"),
	mbmetric(getu64((*vmguestlib.Session).GetHostMemPhysFree), "HostMemPhysFree",
		"total amount of physical memory free on host"),
	mbmetric(getu64((*vmguestlib.Session).GetHostMemPhys), "HostMemPhys",
		"total amount of memory available to the host OS kernel"),
	mbmetric(getu64((*vmguestlib.Session).GetHostMemUsed), "HostMemUsed",
		"total amount of consumed memory on the host"),
	mbmetric(getu64((*vmguestlib.Session).GetHostMemShared), "HostMemShared",
		"total amount of COW (Copy-On-Write) memory on the host"),
	mbmetric(getu64((*vmguestlib.Session).GetHostMemSwapped), "HostMemSwapped",
		"total amount of memory swapped out on the host"),
	msmetric(getu64((*vmguestlib.Session).GetHostCPUUsed), "HostCPUUsed",
		"total CPU time used by host."),
	msmetric(getu64((*vmguestlib.Session).GetCPUUsed), "CPUUsed",
		"time during which the virtual machine has been using the CPU."),
	mbmetric(getu64((*vmguestlib.Session).GetMemTargetSize), "MemTargetSize",
		"memory target Size"),
	msmetric(getu64((*vmguestlib.Session).GetCPUStolen), "CPUStolen",
		"time that the VM was runnable but not scheduled to run."),
	msmetric(getu64((*vmguestlib.Session).GetTimeElapsed), "TimeElapsed",
		"real time passed since the virtual machine started running on the current host system."),
	unmetric(getu32((*vmguestlib.Session).GetHostNumCPUCores), "HostNumCPUCores",
		"number of physical CPU cores on the host machine."),
	mbmetric(getu32((*vmguestlib.Session).GetMemUsed), "MemUsed",
		"estimated amount of physical host memory currently consumed for this virtual machine's physical memory."),
	mbmetric(getu32((*vmguestlib.Session).GetMemSharedSaved), "MemSharedSaved",
		"estimated amount of physical memory on the host saved from copy-on-write (COW) shared guest physical memory."),
	mbmetric(getu32((*vmguestlib.Session).GetMemShared), "MemShared",
		"physical memory associated with this virtual machine that is copy-on-write (COW) shared on the host."),
	mbmetric(getu32((*vmguestlib.Session).GetMemSwapped), "MemSwapped",
		"memory associated with this virtual machine that has been swapped by the host system."),
	mbmetric(getu32((*vmguestlib.Session).GetMemBallooned), "MemBallooned",
		"memory that has been reclaimed from this virtual machine via the VMware Memory Balloon mechanism."),
	mbmetric(getu32((*vmguestlib.Session).GetMemOverhead), "MemOverhead",
		"overhead memory associated with this virtual machine consumed on the host system."),
	mbmetric(getu32((*vmguestlib.Session).GetMemActive), "MemActive",
		"estimated amount of memory the virtual machine is actively using."),
	mbmetric(getu32((*vmguestlib.Session).GetMemMapped), "MemMapped",
		"mapped memory size of this virtual machine."),
	unmetric(getu32((*vmguestlib.Session).GetMemShares), "MemShares",
		"number of memory shares allocated to the virtual machine."),
	mbmetric(getu32((*vmguestlib.Session).GetMemLimit), "MemLimit",
		"maximum amount of memory that is available to the virtual machine."),
	mbmetric(getu32((*vmguestlib.Session).GetMemReservation), "MemReservation",
		"minimum amount of memory that is available to the virtual machine."),
	mzmetric(getu32((*vmguestlib.Session).GetHostProcessorSpeed), "HostProcessorSpeed",
		"host processor speed."),
	unmetric(getu32((*vmguestlib.Session).GetCPUShares), "CPUShares",
		"number of CPU shares allocated to the virtual machine."),
	mzmetric(getu32((*vmguestlib.Session).GetCPULimit), "CPULimit",
		"maximum processing power available to the virtual machine."),
	mzmetric(getu32((*vmguestlib.Session).GetCPUReservation), "CPUReservation",
		"minimum processing power available to the virtual machine."),
}
