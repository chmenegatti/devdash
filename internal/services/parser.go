// Package services — parser.go provides output parsing utilities.
package services

// ParseLines splits raw output into lines, stripping empties.
func ParseLines(raw string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(raw); i++ {
		if raw[i] == '\n' {
			line := raw[start:i]
			if len(line) > 0 && line != "\r" {
				// trim trailing \r
				if line[len(line)-1] == '\r' {
					line = line[:len(line)-1]
				}
				lines = append(lines, line)
			}
			start = i + 1
		}
	}
	// Trailing content without newline
	if start < len(raw) {
		line := raw[start:]
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}
	return lines
}
