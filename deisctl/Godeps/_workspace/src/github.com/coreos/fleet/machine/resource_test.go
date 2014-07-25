package machine

import (
	"testing"
)

const sampleMeminfo = `MemTotal:        1022228 kB
MemFree:          933544 kB
Buffers:            6464 kB
Cached:            45108 kB
SwapCached:            0 kB
Active:            26820 kB
Inactive:          35280 kB
Active(anon):      10648 kB
Inactive(anon):      380 kB
Active(file):      16172 kB
Inactive(file):    34900 kB
Unevictable:           0 kB
Mlocked:               0 kB
SwapTotal:             0 kB
SwapFree:              0 kB
Dirty:              2488 kB
Writeback:             0 kB
AnonPages:         10520 kB
Mapped:            18172 kB
Shmem:               512 kB
Slab:              14464 kB
SReclaimable:       8852 kB
SUnreclaim:         5612 kB
KernelStack:         512 kB
PageTables:         1104 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:      511112 kB
Committed_AS:      75624 kB
VmallocTotal:   34359738367 kB
VmallocUsed:        4224 kB
VmallocChunk:   34359731228 kB
HardwareCorrupted:     0 kB
AnonHugePages:      2048 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:       42944 kB
DirectMap2M:     1005568 kB
`

func TestParseMeminfo(t *testing.T) {
	mem, err := parseMeminfo([]byte(sampleMeminfo))
	if err != nil {
		t.Fatalf("Couldn't parse command output: %v", err)
	}

	if mem != 998 {
		t.Errorf("expected 998 MB, got %d", mem)
	}
}
