package echo

import "testing"

func TestGenTable(t *testing.T) {
	table := make([]byte, 128)
	for i := range table {
		table[i] = noescchr
	}
	escapechars := "\"\\/\b\f\n\r\t"
	escapeto := `"\/bfnrt`
	for i := range escapechars {
		table[escapechars[i]] = escapeto[i]
	}
	t.Log(string(table))
}
