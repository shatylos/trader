package utils

func ContainsInt(sliceToSearch []int64, value int64) bool {
	for _, item := range sliceToSearch {
		if value == item {
			return true
		}
	}
	return false
}
