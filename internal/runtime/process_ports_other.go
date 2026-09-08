//go:build !windows

package runtime

import "fmt"

// ListeningTCPPorts is currently implemented for the Windows dogfood target.
// Other platforms return no facts rather than falling back to a host-wide
// generic process scanner.
func ListeningTCPPorts(pid int) ([]int, error) {
	if pid <= 0 {
		return nil, fmt.Errorf("process pid must be positive")
	}
	return nil, nil
}
