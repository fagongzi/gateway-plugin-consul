package consul

import (
	"time"

	api "github.com/hashicorp/consul/api"

	"fmt"

	"strings"

	"strconv"

	"github.com/fagongzi/gateway-plugin-consul/pkg/conf"
	"github.com/fagongzi/gateway/pkg/lb"
	"github.com/fagongzi/gateway/pkg/model"
)

const (
	// MaxQPSTag gateway max qps tag
	MaxQPSTag = "GATEWAY-MAX-QPS"
	// HalfToOpenTag gateway half to open tag
	HalfToOpenTag = "GATEWAY-HALF-TO-OPEN"
	// HalfTrafficRateTag gateway half traffic rate
	HalfTrafficRateTag = "GATEWAY-HALF-TRAFFIC-RATE"
	// CloseCountTag gateway close count
	CloseCountTag = "GATEWAY-CLOSE-COUNT"
)

// Backend consul backend
type Backend struct {
	cli   *api.Client
	agent *api.Agent
}

// NewBackend new consul backend
func NewBackend(conf *conf.Conf) (*Backend, error) {
	backend := &Backend{}

	cfg := api.DefaultConfig()
	cfg.Address = conf.ConsulAddr

	if conf.Timeout != 0 {
		cfg.HttpClient.Timeout = time.Duration(conf.Timeout) * time.Second
	}

	if conf.Token != "" {
		cfg.Token = conf.Token
	}

	if conf.AuthEnabled {
		cfg.HttpAuth = &api.HttpBasicAuth{
			Username: conf.AuthUserName,
			Password: conf.AuthPassword,
		}
	}

	cli, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	backend.cli = cli
	backend.agent = cli.Agent()

	return backend, nil
}

// GetClusters get clusters
func (b *Backend) GetClusters() (map[string]*model.Cluster, error) {
	services, err := b.agent.Services()
	if err != nil {
		return nil, err
	}

	names := make(map[string]string)
	for _, service := range services {
		if _, ok := names[service.Service]; !ok {
			names[service.Service] = service.Service
		}
	}

	clusters := make(map[string]*model.Cluster)
	for name := range names {
		clusters[name] = &model.Cluster{
			Name:     name,
			LbName:   lb.ROUNDROBIN,
			External: true,
		}
	}

	return clusters, nil
}

// GetServers get servers by cluster
func (b *Backend) GetServers(clusterName string) (map[string]*model.Server, error) {
	services, err := b.agent.Services()
	if err != nil {
		return nil, err
	}

	targetServices := make(map[string]*api.AgentService)
	for _, service := range services {
		if service.Service == clusterName {
			targetServices[service.ID] = service
		}
	}

	servers := make(map[string]*model.Server)
	for _, service := range targetServices {
		s := &model.Server{
			Schema: "http",
			Addr:   fmt.Sprintf("%s:%d", service.Address, service.Port),

			External: true,
		}

		if service.Tags != nil {
			for _, tag := range service.Tags {
				kv := strings.SplitN(tag, "=", 2)

				switch kv[0] {
				case MaxQPSTag:
					v, err := strconv.Atoi(kv[1])
					if err == nil {
						s.MaxQPS = v
					}
				case HalfToOpenTag:
					v, err := strconv.Atoi(kv[1])
					if err == nil {
						s.HalfToOpen = v
					}
				case HalfTrafficRateTag:
					v, err := strconv.Atoi(kv[1])
					if err == nil {
						s.HalfTrafficRate = v
					}
				case CloseCountTag:
					v, err := strconv.Atoi(kv[1])
					if err == nil {
						s.CloseCount = v
					}
				}
			}
		}

		servers[s.Addr] = s
	}

	return servers, nil
}
