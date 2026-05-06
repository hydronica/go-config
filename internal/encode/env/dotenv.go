package env

import (
	"bufio"
	"io"
	"os"
	"strings"
	"unicode"
)

// LoadEnvFile opens path and calls LoadEnvReader. The caller can instead open the file and use
// LoadEnvReader directly (for example with strings.NewReader in tests).
func LoadEnvFile(path string, v interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	m, err := readDotenvMap(f)
	if err != nil {
		return err
	}
	d := &Decoder{
		GetVal: func(k string) string { return m[k] },
	}
	return d.Unmarshal(v)
}

// readDotenvMap reads r line by line. Each non-empty, non-comment line is split on the first
// '=' into key and value (trimmed).
func readDotenvMap(r io.Reader) (map[string]string, error) {
	sc := bufio.NewScanner(r)
	vars := make(map[string]string)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		eq := strings.Index(line, "=")
		if eq < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		if key == "" {
			continue
		}
		val := parseDotenvLineValue(line[eq+1:])
		vars[key] = val
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return vars, nil
}

func parseDotenvLineValue(raw string) string {
	raw = strings.TrimSpace(raw)
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		return unescapeEnvDoubleQuoted(raw[1 : len(raw)-1])
	}
	if len(raw) >= 2 && raw[0] == '\'' && raw[len(raw)-1] == '\'' {
		return raw[1 : len(raw)-1]
	}
	if i := strings.Index(raw, " #"); i >= 0 {
		raw = strings.TrimSpace(raw[:i])
	}
	return raw
}

func unescapeEnvDoubleQuoted(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			case '"', '\\':
				b.WriteByte(s[i+1])
			default:
				b.WriteByte(s[i+1])
			}
			i++
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func needsEnvQuotes(s string) bool {
	if s == "" {
		return true
	}
	for _, r := range s {
		if unicode.IsSpace(r) || r == '#' || r == '=' || r == '"' || r == '\'' || r == '\\' {
			return true
		}
	}
	return false
}

func quoteEnvString(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			b.WriteByte(s[i])
		}
	}
	b.WriteByte('"')
	return b.String()
}

func formatEnvScalar(s string) string {
	if needsEnvQuotes(s) {
		return quoteEnvString(s)
	}
	return s
}
