package loadbalancer

type ServerEntry struct {
	Url     string
	Healthy bool
}

type Strategy interface {
	ChooseServer() *ServerEntry
}

func NewStrategy(name string, serverList []string) (Strategy, error) {
	sl := make(map[int]ServerEntry)
	for i, s := range serverList {
		sl[i] = ServerEntry{s, false}
	}

	var strategy Strategy = nil
	switch name {
	default:
		strategy = newRoundRobinStrategy()
		break
	}

	return strategy, nil
}

type RoundRobinStrategy struct {
	Name string
}

func (s *RoundRobinStrategy) ChooseServer() *ServerEntry {
	return nil
}

func newRoundRobinStrategy() *RoundRobinStrategy {
	s := &RoundRobinStrategy{}
	return s
}

type WeightedRoundRobin struct {
	Name string
}

func (s *WeightedRoundRobin) ChooseServer() *ServerEntry {
	return nil
}

func newWeightedRoundRobin() *WeightedRoundRobin {
	s := &WeightedRoundRobin{}
	return s
}
