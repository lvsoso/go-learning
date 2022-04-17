package base_demo

import "strings"

//func Split(s, sep string) (result []string) {
//	i := strings.Index(s, sep)
//
//	for i > -1 {
//		result = append(result, s[:i])
//		s = s[i+len(sep):]
//		i = strings.Index(s, sep)
//	}
//	result = append(result, s)
//	return
//}

func Split(s, sep string) (result []string) {
	result = make([]string, 0, strings.Count(s, sep)+1)
	i := strings.Index(s, sep)

	for i > -1 {
		result = append(result, s[:i])
		s = s[i+len(sep):]
		i = strings.Index(s, sep)
	}
	result = append(result, s)
	return
}
