package loadbalancer

func roundRobinStrategy(ld *LoadBalancer) *ServerEntry {
	ld.currIdx = ld.nextIdx
	ld.nextIdx = ld.currIdx + 1
	if ld.nextIdx >= ld.serverCount {
		ld.nextIdx = 0
	}

	currIdx := ld.currIdx
	s := ld.ServerList[ld.currIdx]
	for !s.Healthy {
		s = *ld.chooseServer()
		if ld.currIdx == currIdx {
			return nil
		}
	}

	return &s
}

func weightedRoundRobinStrategy(ld *LoadBalancer) *ServerEntry {
	ld.currIdx = ld.nextIdx
	ld.nextIdx = ld.currIdx + 1
	if ld.nextIdx >= ld.serverCount {
		ld.nextIdx = 0
	}

	currIdx := ld.currIdx
	s := ld.ServerList[ld.currIdx]
	for !s.Healthy {
		s = *ld.chooseServer()
		if ld.currIdx == currIdx {
			return nil
		}
	}

	return &s
}
