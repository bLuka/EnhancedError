package eerror

import "math/rand"

func generateUniqueID() uint {
	return uint(rand.Uint64())
}
