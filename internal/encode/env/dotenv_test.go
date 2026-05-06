package env

import (
	"strings"
	"testing"

	"github.com/hydronica/trial"
)

func TestReadDotenvMap(t *testing.T) {
	fn := func(in string) (map[string]string, error) {
		return readDotenvMap(strings.NewReader(in))
	}
	cases := trial.Cases[string, map[string]string]{
		"default": {
			Input: `NAME=apply
VALUE=23
ENABLE=true
TIME="2010-08-10T00:00:00Z"
FLOAT64=99.9
DURA="10s"`,
			Expected: map[string]string{
				"NAME":    "apply",
				"VALUE":   "23",
				"ENABLE":  "true",
				"TIME":    "2010-08-10T00:00:00Z",
				"FLOAT64": "99.9",
				"DURA":    "10s",
			},
		},
		"empty": {
			Input:    "",
			Expected: map[string]string{},
		},
		"whitespace_only": {
			Input:    "\n\n  \t  \n",
			Expected: map[string]string{},
		},
		"comments_and_blank_lines": {
			Input: `# leading comment

KEY1=a
  # indented comment
KEY2=b

# trailing section
KEY3=c`,
			Expected: map[string]string{
				"KEY1": "a",
				"KEY2": "b",
				"KEY3": "c",
			},
		},
		"single_quoted_value": {
			Input:    `MSG='hello world'`,
			Expected: map[string]string{"MSG": "hello world"},
		},
		"double_quoted_with_escapes": {
			Input: `LINE="a\nb\tc\""
PATH="C:\\temp"`,
			Expected: map[string]string{
				"LINE": "a\nb\tc\"",
				"PATH": `C:\temp`,
			},
		},
		"inline_space_hash_strips_comment": {
			Input:    `FOO=bar baz # not used`,
			Expected: map[string]string{"FOO": "bar baz"},
		},
		"quoted_value_preserves_hash": {
			Input:    `FOO="bar # still inside"`,
			Expected: map[string]string{"FOO": "bar # still inside"},
		},
		"line_without_equals_skipped": {
			Input: `NOT_A_VAR_LINE
OK=yes`,
			Expected: map[string]string{"OK": "yes"},
		},
		"empty_key_skipped": {
			Input: `=nope
A=ok`,
			Expected: map[string]string{"A": "ok"},
		},
		"first_equals_splits_key_and_value": {
			Input:    `URL=https://x.example?q=a=b`,
			Expected: map[string]string{"URL": "https://x.example?q=a=b"},
		},
		"trimmed_key_and_value": {
			Input:    "  KEY  =  spaced  ",
			Expected: map[string]string{"KEY": "spaced"},
		},
		"duplicate_key_last_wins": {
			Input: `X=first
X=second`,
			Expected: map[string]string{"X": "second"},
		},
		"unquoted_value_trimmed": {
			Input:    "PORT=8080",
			Expected: map[string]string{"PORT": "8080"},
		},
	}
	trial.New(fn, cases).SubTest(t)
}
