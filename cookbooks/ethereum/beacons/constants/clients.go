package client_consts

func GetAnyClientHTTP(clientName string) []string {
	switch clientName {
	case Lighthouse:
		return LighthouseBeaconPorts
	case Geth:
		return GethBeaconPorts
	}

	return []string{}
}
