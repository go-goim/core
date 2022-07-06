package util

func IsGroupUID(uid string) bool {
	return len(uid) > 1 && uid[:2] == "g_"
}
