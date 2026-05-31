//go:build !solution

package hotelbusiness

import "slices"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) (result []Load) {
	if len(guests) == 0 {
		return
	}

	events := make(map[int]int)
	for _, g := range guests {
		events[g.CheckInDate]++
		events[g.CheckOutDate]--
	}

	uniqDays := []int{}
	for d := range events {
		uniqDays = append(uniqDays, d)
	}
	slices.Sort(uniqDays)

	currentCnt := 0
	for _, dt := range uniqDays {
		if events[dt] == 0 {
			continue
		}
		currentCnt += events[dt]
		result = append(result, Load{
			StartDate:  dt,
			GuestCount: currentCnt,
		})
	}
	return
}
