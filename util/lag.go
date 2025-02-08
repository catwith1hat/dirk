package util

import (
	"fmt"
	"slices"
	"sort"
	"time"
)

// FIXME: Dirk must know the chain timing for the lag mechanism to
// work effectively. ATM hard code all my node names with the
// respective networks.
//
// FIXME: Does the signing domain (given on the signing RPC) contain
// something from which we can infer the chain and therefore the slot
// timings? I do think so...

func slotTimeForClient(cname string, slot uint64) time.Time {
	if cname == "n4-i0" {
		return slotTimeHolesky(slot)
	}
	return slotTimeMainnet(slot)
}


// FIXME: action should probably be an enum instead of a string.
func Delay(action string, cname string, pubkey []byte, slot uint64) {
	var delay time.Duration

	if action == "SignBeaconAttestation" || action == "SignBeaconAttestations" {
		// Attestations happen 4 seconds into the slot
		delay = 4 * time.Second
	} else if action == "SignBeaconProposal" {
		// Block proposal happens at the slot boundary but we
		// configured a 1 second proposer-delay on vouch. So
		// account for that.
		delay = time.Second
	} else {
		delay = 0
	}

	// Refresh last seen names and find the index of cname.
	cNameMemory := GetCNameMemory()
	cNameMemory.Refresh(cname)
	activeCnames := cNameMemory.GetAllActive()
	// fmt.Printf("Presort: activeCnames=%v\n", activeCnames)
	sort.Slice(activeCnames, func(i, j int) bool {
		// FIXME: Should be configurable somehow, but
		// all Dirk instances must arrive at the same
		// ordering.
		//
		// n3-i0 is my best node. Always prefer its
		// signatures.
		return activeCnames[i] == "n3-i0" || activeCnames[i] < activeCnames[j] && activeCnames[j] != "n3-i0"
	})
	index := slices.Index(activeCnames, cname)

	now := time.Now().UTC()
	slotT := slotTimeForClient(cname, slot)
	delta := now.Sub(slotT)

	if index > 0 {
		// Sign request from backup nodes 4 seconds
		// into the slot, plus 200ms extra for each
		// index position (which is equal to the
		// backup node priority).
		networkLag := time.Duration(200*index) * time.Millisecond
		timeout := slotT.Add(delay + networkLag)
		sleep := timeout.Sub(now)
		fmt.Printf("cname: %s, activeCnames=%v, index=%d, current_time: %s, slot_time(%d)=%s, delta=%v, sleep=%v\n", cname, activeCnames, index, now.Format(time.RFC3339), slot, slotT.Format(time.RFC3339), delta, sleep)
		if sleep > 0 {
			time.Sleep(sleep)
		}
	} else {
		fmt.Printf("cname: %s, activeCnames=%v, index=%d, current_time: %s, slot_time(%d)=%s, delta=%v\n", cname, activeCnames, index, now.Format(time.RFC3339), slot, slotT.Format(time.RFC3339), delta)
	}
}
