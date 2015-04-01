package strings

import "fmt"
import "strings"

func FormatMultiColumn(words []string, columns int) (lines []string) {
	avg := AverageWorldlen(words) + 2
	for ncol := (columns / avg); ncol > 1; ncol -= 1 {
		wcols := make([]int, ncol)
		rows := make([][]string, 0)

		row, y := make([]string, 0), 0
		for _, w := range words {
			if len(row)+1 == ncol {
				rows = append(rows, row)
				row, y = make([]string, 0), 0
			}
			if len(w) > wcols[y] {
				wcols[y] = len(w)
			}
			row, y = append(row, w), y+1
		}
		if len(row) > 0 {
			rows = append(rows, row)
		}

		for y := 0; y < ncol; y++ {
			for _, row := range rows {
				if y < len(row) {
					row[y] = SuffixAlign(row[y], wcols[y]+2)
				}
			}
		}

		lines = make([]string, 0)
		for _, row := range rows {
			line := strings.Join(row, "")
			if len(line) < columns {
				lines = append(lines, line)
			} else {
				break
			}
		}

		if len(lines) == len(rows) {
			break
		}
	}
	return
}

func LargestWord(words []string) string {
	if len(words) == 0 {
		return ""
	}

	lword := words[0]
	for _, w := range words {
		if len(w) > len(lword) {
			lword = w
		}
	}
	return lword
}

func AverageWorldlen(words []string) int {
	sum := 0
	for _, w := range words {
		sum += len(w)
	}
	return (sum / len(words))
}

func SuffixAlign(word string, width int) string {
	f := fmt.Sprintf("%%-%vs", width)
	return fmt.Sprintf(f, word)
}
