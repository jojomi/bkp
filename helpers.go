package bkp

import (
	"os"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func mergeStringSlices(s1, s2 []string) []string {
	s := make([]string, len(s1)+len(s2))
	i := 0
	for _, elem := range s1 {
		s[i] = elem
		i++
	}
	for _, elem := range s2 {
		s[i] = elem
		i++
	}
	return s
}
