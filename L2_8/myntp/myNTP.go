// Package myntp implements a simple NTP-client
package myntp

import (
    "fmt"
    "github.com/beevik/ntp"
)

// GetCurrentTime gets the current time from an NTP-server
func GetCurrentTime() (string, error) {
    ntpTime, err := ntp.Time("pool.ntp.org")
    if err != nil {
        return "", fmt.Errorf("NTP error: %w", err)
    }

    return ntpTime.Local().Format("Mon Jan 2 15:04:05 MST 2006"), nil
}