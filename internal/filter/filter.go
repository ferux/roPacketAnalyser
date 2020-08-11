package filter

type Packet interface {
	Allowed(pid string) (allowed bool)
}

func NewAllowList(pids []string) Packet {
	var mappedpids = make(map[string]struct{}, len(pids))
	for _, pid := range pids {
		mappedpids[pid] = struct{}{}
	}

	return packetAllowList{packetIDs: mappedpids}
}
func NewBlockList(pids []string) Packet {
	var mappedpids = make(map[string]struct{}, len(pids))
	for _, pid := range pids {
		mappedpids[pid] = struct{}{}
	}

	return packetBlockList{packetIDs: mappedpids}
}

type packetAllowList struct {
	packetIDs map[string]struct{}
}

// Allowed implementes Packet interace.
func (p packetAllowList) Allowed(pid string) (allowed bool) {
	_, allowed = p.packetIDs[pid]

	return allowed
}

type packetBlockList struct {
	packetIDs map[string]struct{}
}

func (p packetBlockList) Allowed(pid string) (allowed bool) {
	_, allowed = p.packetIDs[pid]

	return allowed
}
