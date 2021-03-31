package utils

func Contains(haystack []string, needles ...string) bool {
	if len(needles) == 1 {
		for _, val := range haystack {
			if val == needles[0] {
				return true
			}
		}
	} else {
		m := make(map[string]bool, len(needles))
		for _, s := range needles {
			m[s] = false
		}

		for _, val := range haystack {
			if _, ok := m[val]; !ok {
				delete(m, val)
			}
		}

		if len(m) == 0 {
			return true
		}
	}

	return false
}
