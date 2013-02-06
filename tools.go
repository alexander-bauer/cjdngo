package cjdngo

import (
	"strings"
)

// Truncate is a utility function meant for shortening full-length
// IPv6 addresses to their smallest unique 4-character (16-bit)
// identifier. It does not make any DNS requests.
func Truncate(addrs []string) (t map[string]string) {
	t = make(map[string]string, len(addrs))

	// For every address in the list, check every other address to see
	// if it's okay to truncate.
iLoop:
	for i, addr := range addrs {
		iComp := strings.Split(addr, ":")
	jLoop:
		for j, addr2 := range addrs {
			// If i and j are the same, skip it.
			if i == j {
				continue jLoop
			}
			jComp := strings.Split(addr2, ":")
			if iComp[len(iComp)-1] == jComp[len(jComp)-1] {
				// If the final characters match, then map with the
				// last two components.
				// BUG(DuoNoxSol): If the last two parts of the IP
				// match, then they'll truncate to the same string.
				t[addr] = iComp[len(iComp)-2] + ":" +
					iComp[len(iComp)-1]
				continue iLoop
			}
		}
		t[addr] = iComp[len(iComp)-1]
	}
	return
}
