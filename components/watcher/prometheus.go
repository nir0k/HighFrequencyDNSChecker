package watcher

import (
	sqldb "HighFrequencyDNSChecker/components/db"
	log "HighFrequencyDNSChecker/components/log"
	"context"
	"encoding/base64"
	"strconv"
	"sync"
	"time"

	"github.com/castai/promwrite"
	"github.com/miekg/dns"
)

var (
    Buffer []promwrite.TimeSeries
    Config sqldb.Config
    Mu sync.Mutex
)


func basicAuth() string {
    auth := Config.Prometheus.Username + ":" + Config.Prometheus.Password
    return base64.StdEncoding.EncodeToString([]byte(auth))
}


func collectLabels(server sqldb.Resolver, r_header dns.MsgHdr, promLabels sqldb.PrometheusLabelConfiguration) []promwrite.Label {
    var label promwrite.Label

    labels := []promwrite.Label{
        {
            Name:  "__name__",
            Value: Config.Prometheus.MetricName,
        },
        {
            Name: "server",
            Value: server.Server,
        },
        {
            Name: "server_ip",
            Value: server.IPAddress,
        },
        {
            Name: "domain",
            Value: server.Domain,
        },
        {
            Name: "location",
            Value: server.Location,
        },
        {
            Name: "site",
            Value: server.Site,
        },
        {
            Name: "watcher",
            Value: Config.General.Hostname,
        },
        {
            Name: "watcher_ip",
            Value: Config.General.IPAddress,
        },
        {
            Name: "watcher_security_zone",
            Value: Config.Watcher.SecurityZone,
        },
        {
            Name: "watcher_location",
            Value: Config.Watcher.Location,
        },
        {
            Name: "protocol",
            Value: server.Protocol,
        },
        {
            Name: "server_security_zone",
            Value: server.ServerSecurityZone,
        },
        {
            Name: "service_mode",
            Value: strconv.FormatBool(server.ServiceMode),
        },
    }

    label.Name = "zonename"
    label.Value = server.Zonename
    labels = append(labels, label)

    if promLabels.AuthenticatedData {
        label.Name = "authenticated_data"
        label.Value = strconv.FormatBool(r_header.AuthenticatedData)
        labels = append(labels, label)
    }
    if promLabels.Authoritative {
        label.Name = "authoritative"
        label.Value = strconv.FormatBool(r_header.Authoritative)
        labels = append(labels, label)
    }
    if promLabels.CheckingDisabled {
        label.Name = "checking_disabled"
        label.Value = strconv.FormatBool(r_header.CheckingDisabled)
        labels = append(labels, label)
    }
    if promLabels.Opcode {
        label.Name = "opscodes"
        label.Value = strconv.Itoa(r_header.Opcode)
        labels = append(labels, label)
    }
    if promLabels.Rcode {
        label.Name = "rcode"
        label.Value = strconv.Itoa(r_header.Rcode)
        labels = append(labels, label)
    }
    if promLabels.RecursionAvailable {
        label.Name = "recursion_available"
        label.Value = strconv.FormatBool(r_header.RecursionAvailable)
        labels = append(labels, label)
    }
    if promLabels.RecursionDesired {
        label.Name = "recursion_desired"
        label.Value = strconv.FormatBool(r_header.RecursionDesired)
        labels = append(labels, label)
    }
    if promLabels.Truncated {
        label.Name = "truncated"
        label.Value = strconv.FormatBool(r_header.Truncated)
        labels = append(labels, label)
    }
    if promLabels.PollingRate {
        label.Name = "polling_rate"
        label.Value = strconv.Itoa(server.QueryCount)
        labels = append(labels, label)
    }
    if promLabels.Recursion {
        label.Name = "recursion"
        label.Value = strconv.FormatBool(server.Recursion)
        labels = append(labels, label)
    }
    return labels
}


func BufferTimeSeries(server sqldb.Resolver, tm time.Time, value float64, response_header dns.MsgHdr) {
    Mu.Lock()
	defer Mu.Unlock()
    if len(Buffer) >= Config.Prometheus.BuferSize {
        go sendVM(Buffer)
        Buffer = nil
        return
    }
    instance := promwrite.TimeSeries{
        Labels: collectLabels(server, response_header, Config.PrometheusLabels),
        Sample: promwrite.Sample{
            Time:  tm,
            Value: value,
        },
    }
    Buffer = append(Buffer, instance)
}


func sendVM(items []promwrite.TimeSeries) bool {
    client := promwrite.NewClient(Config.Prometheus.Url)
    
    req := &promwrite.WriteRequest{
        TimeSeries: items,
    }
    log.AppLog.Debug("TimeSeries:", items)
    for i := 0; i < Config.Prometheus.RetriesCount; i++ {
        _, err := client.Write(context.Background(), req, promwrite.WriteHeaders(map[string]string{"Authorization": "Basic " + basicAuth()}))
        if err == nil {
            log.AppLog.Debug("Remote write to VM succesfull. URL:", Config.Prometheus.Url ,", timestamp:", time.Now().Format("2006/01/02 03:04:05.000"))
            return true
        }
        log.AppLog.Warn("Remote write to VM failed. Retry ", i+1, " of ", Config.Prometheus.RetriesCount, ". URL:", Config.Prometheus.Url, ", timestamp:", time.Now().Format("2006/01/02 03:04:05.000"), ", error:", err)
    }
    log.AppLog.Error("Remote write to VM failed. URL:", Config.Prometheus.Url ,", timestamp:", time.Now().Format("2006/01/02 03:04:05.000"))
    log.AppLog.Debug("Request:", req)
    return false
}
