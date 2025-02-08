package util

import "time"

func slotTimeMainnet(slot uint64) time.Time {
	anchorTimeUTC := time.Date(2020, time.December, 1, 12, 00, 23, 0, time.UTC)
	slotDuration := 12 * time.Second
	diffSeconds := time.Duration(int64(slot)) * slotDuration
	return anchorTimeUTC.Add(diffSeconds)
}

func slotTimeHolesky(slot uint64) time.Time {
	anchorTimeUTC := time.Date(2023, time.September, 28, 12, 00, 00, 0, time.UTC)
	slotDuration := 12 * time.Second
	diffSeconds := time.Duration(int64(slot)) * slotDuration
	return anchorTimeUTC.Add(diffSeconds)
}
