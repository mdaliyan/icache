package icache

const (
	// offset64 FNVa offset basis. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	offset64 = 14695981039346656037
	// prime64 FNVa prime value. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	prime64 = 1099511628211

	shardAndOpVal = 255
	shardsCount   = 256
)

// Sum64 gets the string and returns its uint64 hash value.
func keyGen(key string) (hashVal, shardID uint64) {
	var hash uint64 = offset64
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}
	return hash, hash & shardAndOpVal
}

// Sum64 gets the string and returns its uint64 hash value.
func TagKeyGen(tags []string) []uint64 {
	var count = len(tags)
	if count == 0 {
		return nil
	}
	keys := make([]uint64, count)
	for i, tag := range tags {
		var hash uint64 = offset64
		for i := 0; i < len(tag); i++ {
			hash ^= uint64(tag[i])
			hash *= prime64
		}
		keys[i] = hash
	}
	return keys
}

//
// func KeyShard(data []byte) (hashVal, segID uint64) {
//	hashVal = xxhash.Sum64(data)
//	segID = hashVal & shardAndOpVal
//	return
// }
