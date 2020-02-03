package helper

import "strings"

func StringJoin(one, two string) string {
	return strings.Join([]string{one, two}, "")
}

func StringJoinSeperated(one, two, sep string) string {
	return strings.Join([]string{one, two}, sep)
}
