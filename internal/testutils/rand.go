package testutils

import "math/rand"

// Overwrite rand.Int() to always return unique numbers
type UniqueRand struct {
	existing map[int]struct{}
}

func (r UniqueRand) Int() int {
	num := rand.Int()

	_, ok := r.existing[num]
	// do this until you create one that does not exist
	for ok {
		num = rand.Int()
		_, ok = r.existing[num]
	}
	r.existing[num] = struct{}{}
	return num
}

func NewUniqueRand() UniqueRand {
	return UniqueRand{
		existing: make(map[int]struct{}),
	}
}

func (r UniqueRand) Intn(n int) int {
	return rand.Intn(n)
}
