package config

import "strconv"

func ParseInt(value *int, defaultValue int, arg string) {
	if v, e := strconv.Atoi(arg); e == nil {
		*value = v
	} else {
		*value = defaultValue
	}
}
