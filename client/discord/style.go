package discord

import (
	"strconv"
	"strings"
)

func StyleFieldValues(inputs ...any) string {
	outputs := make([]string, len(inputs))
	for _, i := range inputs {
		switch i := i.(type) {
		case string:
			outputs = append(outputs, i)
		case float64:
			outputs = append(outputs, strconv.FormatFloat(i, 'f', -1, 64))
		case int64:
			outputs = append(outputs, strconv.FormatInt(i, 10))
		case int:
			outputs = append(outputs, strconv.Itoa(i))
		}
	}

	s := strings.Join(outputs, "")

	if !strings.ContainsAny(s, "[]()") {
		return ">>> " + s
	}
	s = strings.ReplaceAll(s, "[", "**__")
	s = strings.ReplaceAll(s, "]", "__**")
	s = strings.ReplaceAll(s, "(", "**")
	s = strings.ReplaceAll(s, ")", "**")

	return ">>> " + s
}
