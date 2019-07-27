package core

import (
	"testing"
)

func TestFreeIndex(t *testing.T) {
	indeces := []uint64{1, 2, 4, 5, 6}
	freeIndex := FreeIndex(indeces)
	var reference uint64 = 3
	if freeIndex != reference {
		t.Errorf("Free index is %d (should be %d)\n", freeIndex, reference)
	}

	indeces = []uint64{4, 5, 6, 8, 9}
	freeIndex = FreeIndex(indeces)
	reference = 1
	if freeIndex != reference {
		t.Errorf("Free index is %d (should be %d)\n", freeIndex, reference)
	}

	indeces = []uint64{}
	freeIndex = FreeIndex(indeces)
	reference = 1
	if freeIndex != reference {
		t.Errorf("Free index is %d (should be %d)\n", freeIndex, reference)
	}

	indeces = []uint64{1, 2, 3, 10, 11}
	freeIndex = FreeIndex(indeces)
	reference = 4
	if freeIndex != reference {
		t.Errorf("Free index is %d (should be %d)\n", freeIndex, reference)
	}

	indeces = []uint64{1, 2, 3, 4, 5}
	freeIndex = FreeIndex(indeces)
	reference = 6
	if freeIndex != reference {
		t.Errorf("Free index is %d (should be %d)\n", freeIndex, reference)
	}

}
