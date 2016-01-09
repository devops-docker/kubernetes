/*
Copyright 2014 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package options contains flags for initializing a proxy.
package options

import (
	"net"
	_ "net/http/pprof"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/qos"
	"k8s.io/kubernetes/pkg/util"

	"github.com/spf13/pflag"
)

const (
	ExperimentalProxyModeAnnotation = "net.experimental.kubernetes.io/proxy-mode"
)

// ProxyServerConfig contains configurations for a Kubernetes proxy server
type ProxyServerConfig struct {
	BindAddress                    net.IP
	HealthzPort                    int
	HealthzBindAddress             net.IP
	OOMScoreAdj                    int
	ResourceContainer              string
	Master                         string
	Kubeconfig                     string
	PortRange                      util.PortRange
	HostnameOverride               string
	ProxyMode                      string
	IptablesSyncPeriod             time.Duration
	ConfigSyncPeriod               time.Duration
	NodeRef                        *api.ObjectReference // Reference to this node.
	MasqueradeAll                  bool
	CleanupAndExit                 bool
	KubeAPIQPS                     float32
	KubeAPIBurst                   int
	UDPIdleTimeout                 time.Duration
	ConntrackMax                   int
	ConntrackTCPTimeoutEstablished int // seconds
}

func NewProxyConfig() *ProxyServerConfig {
	return &ProxyServerConfig{
		BindAddress:                    net.ParseIP("0.0.0.0"),
		HealthzPort:                    10249,
		HealthzBindAddress:             net.ParseIP("127.0.0.1"),
		OOMScoreAdj:                    qos.KubeProxyOOMScoreAdj,
		ResourceContainer:              "/kube-proxy",
		IptablesSyncPeriod:             30 * time.Second,
		ConfigSyncPeriod:               15 * time.Minute,
		KubeAPIQPS:                     5.0,
		KubeAPIBurst:                   10,
		UDPIdleTimeout:                 250 * time.Millisecond,
		ConntrackMax:                   256 * 1024, // 4x default (64k)
		ConntrackTCPTimeoutEstablished: 86400,      // 1 day (1/5 default)
	}
}

// AddFlags adds flags for a specific ProxyServer to the specified FlagSet
func (s *ProxyServerConfig) AddFlags(fs *pflag.FlagSet) {
	fs.IPVar(&s.BindAddress, "bind-address", s.BindAddress, "The IP address for the proxy server to serve on (set to 0.0.0.0 for all interfaces)")
	fs.StringVar(&s.Master, "master", s.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	fs.IntVar(&s.HealthzPort, "healthz-port", s.HealthzPort, "The port to bind the health check server. Use 0 to disable.")
	fs.IPVar(&s.HealthzBindAddress, "healthz-bind-address", s.HealthzBindAddress, "The IP address for the health check server to serve on, defaulting to 127.0.0.1 (set to 0.0.0.0 for all interfaces)")
	fs.IntVar(&s.OOMScoreAdj, "oom-score-adj", s.OOMScoreAdj, "The oom-score-adj value for kube-proxy process. Values must be within the range [-1000, 1000]")
	fs.StringVar(&s.ResourceContainer, "resource-container", s.ResourceContainer, "Absolute name of the resource-only container to create and run the Kube-proxy in (Default: /kube-proxy).")
	fs.MarkDeprecated("resource-container", "This feature will be removed in a later release.")
	fs.StringVar(&s.Kubeconfig, "kubeconfig", s.Kubeconfig, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	fs.Var(&s.PortRange, "proxy-port-range", "Range of host ports (beginPort-endPort, inclusive) that may be consumed in order to proxy service traffic. If unspecified (0-0) then ports will be randomly chosen.")
	fs.StringVar(&s.HostnameOverride, "hostname-override", s.HostnameOverride, "If non-empty, will use this string as identification instead of the actual hostname.")
	fs.StringVar(&s.ProxyMode, "proxy-mode", "", "Which proxy mode to use: 'userspace' (older) or 'iptables' (faster). If blank, look at the Node object on the Kubernetes API and respect the '"+ExperimentalProxyModeAnnotation+"' annotation if provided.  Otherwise use the best-available proxy (currently iptables).  If the iptables proxy is selected, regardless of how, but the system's kernel or iptables versions are insufficient, this always falls back to the userspace proxy.")
	fs.DurationVar(&s.IptablesSyncPeriod, "iptables-sync-period", s.IptablesSyncPeriod, "How often iptables rules are refreshed (e.g. '5s', '1m', '2h22m').  Must be greater than 0.")
	fs.DurationVar(&s.ConfigSyncPeriod, "config-sync-period", s.ConfigSyncPeriod, "How often configuration from the apiserver is refreshed.  Must be greater than 0.")
	fs.BoolVar(&s.MasqueradeAll, "masquerade-all", false, "If using the pure iptables proxy, SNAT everything")
	fs.BoolVar(&s.CleanupAndExit, "cleanup-iptables", false, "If true cleanup iptables rules and exit.")
	fs.Float32Var(&s.KubeAPIQPS, "kube-api-qps", s.KubeAPIQPS, "QPS to use while talking with kubernetes apiserver")
	fs.IntVar(&s.KubeAPIBurst, "kube-api-burst", s.KubeAPIBurst, "Burst to use while talking with kubernetes apiserver")
	fs.DurationVar(&s.UDPIdleTimeout, "udp-timeout", s.UDPIdleTimeout, "How long an idle UDP connection will be kept open (e.g. '250ms', '2s').  Must be greater than 0. Only applicable for proxy-mode=userspace")
	fs.IntVar(&s.ConntrackMax, "conntrack-max", s.ConntrackMax, "Maximum number of NAT connections to track (0 to leave as-is)")
	fs.IntVar(&s.ConntrackTCPTimeoutEstablished, "conntrack-tcp-timeout-established", s.ConntrackTCPTimeoutEstablished, "Idle timeout for established TCP connections (0 to leave as-is)")
}