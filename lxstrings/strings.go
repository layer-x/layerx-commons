package lxstrings
import "github.com/layer-x/layerx-commons/lxmath"

func FirstN(str string, n int) string {
	if str == "" {
		return str
	}
	return str[:lxmath.Min(len(str)-1, n)]
}

func LastN(str string, n int) string {
	if str == "" {
		return str
	}
	return str[len(str)-lxmath.Min(len(str)-1, n):len(str)-1]
}