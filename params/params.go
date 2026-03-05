package params

import (
	"strconv"
	"strings"
)

func ParseParams(params string) []int {
	params = strings.TrimSpace(params)
	if params == "" {
		return []int{}
	}
	parts := strings.Split(params, ";")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if n, err := strconv.Atoi(p); err == nil {
			result = append(result, n)
		}
	}
	return result
}

func GetArg(args []int, index int, def int) int {
	if index >= len(args) || args[index] == 0 {
		return def
	}
	return args[index]
}
