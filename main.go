package uncorrupt

import (
	"encoding/base64"
	"math"
)

func populateKey(key []byte, keylength int) []byte {
	sum := 0
	for loc, k := range key {
		sum = (sum + ((loc + 1) * int(k))) % 255
	}

	for i := len(key); i < keylength; i++ {
		key = append(key, byte(sum^((i+1)%255)))
		sum = (sum + ((i + 1) * int(key[i]))) % 255
	}

	for i := range key {
		key[i] = ((key[i] ^ byte(sum)) * byte((i+1)%255)) % 255
	}

	return key
}

func updateKey(key []byte, keylength int) []byte {
	sum := 0
	occurence := [256]int{}

	key = populateKey(key, keylength)

	weight := make([]int, len(key))

	for _, k := range key {
		occurence[k]++
	}
	for loc, k := range key {
		weight[loc] = int(math.Round(float64(occurence[k]+(int(k)*(loc+1)))) / (float64(occurence[k]) / float64(len(key))))
	}

	for i := range key {
		sum += weight[i]
	}

	weightedsum := byte(sum % 127)
	for i := range key {
		key[i] = byte(key[i]^(weightedsum-byte(weight[i]))) % 127
	}

	return key
}

func Run(input []byte, keystring string) []byte {
	keylength := int(math.Max(float64(len(input)), float64(len(keystring))))

	key := make([]byte, 0, keylength)
	key = append(key, []byte(keystring)...)

	key = updateKey(key, keylength)

	keyIdx := 0

	out := make([]byte, 0, len(input))

	for _, x := range input {
		k := key[keyIdx]
		out_byte := (x ^ k)

		keyIdx += 1
		if keyIdx >= len(key) {
			keyIdx = 0
		}

		out = append(out, out_byte)
	}

	return out
}

func Corrupt(input []byte, key string) []byte {
	out := Run(input, key)
	return []byte(base64.StdEncoding.EncodeToString(out))
}

func Uncorrupt(input []byte, key string) []byte {
	decoded, _ := base64.StdEncoding.DecodeString(string(input))
	return Run(decoded, key)
}
