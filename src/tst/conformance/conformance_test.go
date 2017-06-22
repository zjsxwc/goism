package conformance

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var home = os.Getenv("GOISM_DIR")

func init() {
	// Loads "emacs/conformance" package into Emacs daemon.
	eval(`(goism-load "conformance")`)
}

func evalCall(call string) string {
	return eval("(goism-conformance." + call + ")")
}

func evalVar(name string) string {
	return eval("goism-conformance." + name)
}

func eval(expr string) string {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := exec.Command(home+"/script/tst/daemon_eval", expr)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("error: %s (%v)", stderr.String(), err)
	}
	return strings.TrimRight(string(stdout.String()), "\t\n\r")
}

// Enclose each passed string value into double quotes (unconditionnaly).
func q(xs ...string) []string {
	for i, x := range xs {
		xs[i] = `"` + x + `"`
	}
	return xs
}

type callTest struct {
	call           string
	outputExpected string
}

func testCalls(t *testing.T, table []callTest) {
	for _, row := range table {
		res := evalCall(row.call)
		if res != row.outputExpected {
			t.Errorf("(%s)=>%s (want %s)", row.call, res, row.outputExpected)
		}
	}
}

func Test1Ops(t *testing.T) {
	testCalls(t, []callTest{
		// Int ops.
		{"add1Int 1", "2"},
		{"addInt 1 2 3", "6"},
		{"sub1Int 1", "0"},
		{"subInt 3 2 1", "0"},
		{"mulInt 2 2 20", "80"},
		{"quoInt 20 2 2", "5"},
		{"gtInt 2 1", "t"},
		{"ltInt 2 1", "nil"},
		// Float ops.
		{"add1Float 1.0", "2.0"},
		{"addFloat 1.1 2.2 3.3", "6.6"},
		{"sub1Float 1.5", "0.5"},
		{"subFloat 3.5 2.5 1.5", "-0.5"},
		{"mulFloat 2.0 2.0 20.0", "80.0"},
		{"quoFloat 20.0 2.0 2.0", "5.0"},
		{"gtFloat 2.0 1.0", "t"},
		{"ltFloat 2.0 1.0", "nil"},
		// String ops.
		{`concatStr "a" "b" "c"`, `"abc"`},
		{`ltStr "abc" "foo"`, "t"},
		{`ltStr "foo" "abc"`, "nil"},
	})
}

func Test2GlobalVars(t *testing.T) {
	table := []struct {
		name          string
		valueExpected string
	}{
		{"var1", "1"},
		{"var2", "2"},
		{"var3", "3"},
		{"var4", "4"},
		{"var5", "5"},
		{"var6", "6"},
	}

	for _, row := range table {
		res := evalVar(row.name)
		if res != row.valueExpected {
			t.Errorf("%s=>%s (want %s)", row.name, res, row.valueExpected)
		}
	}
}

func Test3MultiResult(t *testing.T) {
	table := []struct {
		call           string
		names          []string
		valuesExpected []string
	}{
		{"return2", []string{"r2_1", "r2_2"}, q("a", "b")},
		{"return3", []string{"r3_2", "r3_3"}, q("b", "c")},
		{"return4", []string{"r4_1", "r4_3", "r4_4"}, q("a", "c", "d")},
	}

	for _, row := range table {
		evalCall(row.call) // For side effects
		for i, name := range row.names {
			res := evalVar(name)
			if res != row.valuesExpected[i] {
				t.Errorf("%s=>%s (want %s)", name, res, row.valuesExpected[i])
			}
		}
	}
}

func Test4If(t *testing.T) {
	testCalls(t, []callTest{
		{"alwaysZero", "0"},
		{"neverZero", "1"},
		{"isZero 0", "t"},
		{"isZero 1", "nil"},
		{"stringifyInt1 0", `"0"`},
		{"stringifyInt1 1", `"1"`},
		{"stringifyInt1 2", `"2"`},
		{"stringifyInt1 3", `"x"`},
		{"stringifyInt2 0", `"0"`},
		{"stringifyInt2 1", `"1"`},
		{"stringifyInt2 2", `"2"`},
		{"stringifyInt2 3", `"x"`},
	})
}

func Test5Switch(t *testing.T) {
	testCalls(t, []callTest{
		{"stringifyInt3 0", `"0"`},
		{"stringifyInt3 1", `"1"`},
		{"stringifyInt3 2", `"2"`},
		{"stringifyInt3 3", `"x"`},
		{"stringifyInt4 0", `"0"`},
		{"stringifyInt4 1", `"1"`},
		{"stringifyInt4 2", `"2"`},
		{"stringifyInt4 3", `"x"`},
	})
}