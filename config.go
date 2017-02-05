package cachet

import (
	"net"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

type CachetMonitor struct {
	SystemName  string                   `json:"system_name"`
	API         CachetAPI                `json:"api"`
	RawMonitors []map[string]interface{} `json:"monitors" yaml:"monitors"`

	Monitors []MonitorInterface `json:"-" yaml:"-"`

	Immediate bool `json:"-" yaml:"-"`
}

// Validate configuration
func (cfg *CachetMonitor) Validate() bool {
	valid := true

	if len(cfg.SystemName) == 0 {
		// get hostname
		cfg.SystemName = getHostname()
	}

	if len(cfg.API.Token) == 0 || len(cfg.API.URL) == 0 {
		logrus.Warnf("API URL or API Token missing.\nGet help at https://github.com/castawaylabs/cachet-monitor")
		valid = false
	}

	if len(cfg.Monitors) == 0 {
		logrus.Warnf("No monitors defined!\nSee help for example configuration")
		valid = false
	}

	for index, monitor := range cfg.Monitors {
		if errs := monitor.Validate(); len(errs) > 0 {
			logrus.Warnf("Monitor validation errors (index %d): %v", index, "\n - "+strings.Join(errs, "\n - "))
			valid = false
		}
	}

	return valid
}

// getHostname returns id of the current system
func getHostname() string {
	hostname, err := os.Hostname()
	if err == nil && len(hostname) > 0 {
		return hostname
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil || len(addrs) == 0 {
		return "unknown"
	}

	return addrs[0].String()
}

func getMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetMonitorType(t string) string {
	if len(t) == 0 {
		return "http"
	}

	return t
}
