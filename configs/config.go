package configs

const (
	ROUND_ROBIN          = "round_robin"
	WEIGHTED_ROUND_ROBIN = "weighted_round_robin"
)

// general configs for the load balancer
type NetworkConfig struct {
	Port string
	Ip   string
}

func DefaultNeworkConfig() *NetworkConfig {
	return &NetworkConfig{
		Port: "5151",
		Ip:   "0.0.0.0",
	}
}

type LoadBalancerConfig struct {
	Strategy     string
	ServerList   []string
	MaxNodeCount int
}

func DefaultLoadBalancerConfig(uri ...string) *LoadBalancerConfig {
	return &LoadBalancerConfig{
		Strategy:     ROUND_ROBIN,
		ServerList:   uri,
		MaxNodeCount: len(uri),
	}
}
