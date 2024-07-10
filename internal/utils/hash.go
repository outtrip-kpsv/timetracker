package utils

import "hash/fnv"

func HashStringToInt(s string) uint64 {
  h := fnv.New64a()
  h.Write([]byte(s))
  return h.Sum64()
}
