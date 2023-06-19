package internal

import "bytes"

func StringJoin(in []string) *string {
	var m2 = bytes.NewBuffer(make([]byte, 0, 100))
	for _, v := range in {
		m2.WriteString(v)
	}
	str := m2.String()
	return &str
}
