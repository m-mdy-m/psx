package logger

import (
	"bytes"
	"strings"
	"testing"
)

// capture redirects both streams and resets the level flags for one call.
func capture(t *testing.T, verbose, quiet bool, fn func()) (out, err string) {
	t.Helper()

	var o, e bytes.Buffer
	prevOut, prevErr := Out, Err
	Out, Err = &o, &e
	Configure(verbose, quiet, true)
	defer func() {
		Out, Err = prevOut, prevErr
		Configure(false, false, false)
	}()

	fn()
	return o.String(), e.String()
}

// A check report must be the only thing on stdout, which is what makes
// `psx check -o json | jq` work. Step is deliberately excluded: it reports what
// `fix` created, which is that command's output rather than a diagnostic.
func TestCheckDiagnosticsGoToStderr(t *testing.T) {
	out, errOut := capture(t, false, false, func() {
		Info("informational")
		Success("done")
		Warning("careful")
		Error("broken")
		Plain("plain")
	})

	if out != "plain\n" {
		t.Errorf("stdout = %q, want only the Plain line", out)
	}
	for _, want := range []string{"informational", "done", "careful", "broken"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("stderr is missing %q: %q", want, errOut)
		}
	}
}

// Step is what `fix` prints per created file, so it belongs on stdout where it
// can be read or piped like any other result.
func TestStepGoesToStdout(t *testing.T) {
	out, errOut := capture(t, false, false, func() { Step("created README.md") })

	if !strings.Contains(out, "created README.md") {
		t.Errorf("stdout = %q, want the step line", out)
	}
	if strings.Contains(errOut, "created README.md") {
		t.Errorf("the step line leaked to stderr: %q", errOut)
	}
}

func TestVerboseIsOffByDefault(t *testing.T) {
	_, errOut := capture(t, false, false, func() { Verbose("chatter") })
	if strings.Contains(errOut, "chatter") {
		t.Errorf("verbose output appeared without --verbose: %q", errOut)
	}
	if VerboseEnabled() {
		t.Error("VerboseEnabled should be false by default")
	}
}

func TestQuietSuppressesEverythingButErrors(t *testing.T) {
	out, errOut := capture(t, true, true, func() {
		Info("informational")
		Success("done")
		Warning("careful")
		Verbose("chatter")
		Error("broken")
	})

	if out != "" {
		t.Errorf("quiet mode wrote to stdout: %q", out)
	}
	if QuietEnabled() {
		t.Error("QuietEnabled should be true while quiet")
	}
	for _, unwanted := range []string{"informational", "done", "careful", "chatter"} {
		if strings.Contains(errOut, unwanted) {
			t.Errorf("quiet mode printed %q: %q", unwanted, errOut)
		}
	}
	if !strings.Contains(errOut, "broken") {
		t.Errorf("quiet mode swallowed the error: %q", errOut)
	}
}

// Plain carries machine-readable output: three of its four callers are `--json`
// bodies. Quiet must not silence a report, or `-o json --quiet` would print
// nothing at all.
func TestPlainSurvivesQuiet(t *testing.T) {
	out, errOut := capture(t, true, true, func() { Plain(`{"summary":{}}`) })

	if !strings.Contains(out, `{"summary":{}}`) {
		t.Errorf("quiet mode suppressed the report body: %q", out)
	}
	if errOut != "" {
		t.Errorf("the report body leaked to stderr: %q", errOut)
	}
}

// --quiet has to win over --verbose, or a script passing both still gets chatter.
func TestQuietBeatsVerbose(t *testing.T) {
	capture(t, true, true, func() {
		Verbose("chatter")
		if VerboseEnabled() {
			t.Error("VerboseEnabled should be false when quiet is set")
		}
	})
}

func TestVerboseEnabledWithVerboseAlone(t *testing.T) {
	_, errOut := capture(t, true, false, func() {
		Verbose("chatter")
		if !VerboseEnabled() {
			t.Error("VerboseEnabled should be true")
		}
	})
	if !strings.Contains(errOut, "chatter") {
		t.Errorf("--verbose printed nothing: %q", errOut)
	}
}

func TestFormattedHelpersFormat(t *testing.T) {
	_, errOut := capture(t, false, false, func() {
		Infof("count %d", 3)
		Successf("%s done", "all")
		Warningf("%d warnings", 2)
		Errorf("cannot read %s", "go.mod")
	})

	for _, want := range []string{"count 3", "all done", "2 warnings", "cannot read go.mod"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("stderr is missing %q: %q", want, errOut)
		}
	}
}

func TestErrorfReturnsAUsableError(t *testing.T) {
	Configure(false, false, true)
	defer Configure(false, false, false)

	err := Errorf("bad thing %d", 7)
	if err == nil {
		t.Fatal("Errorf returned no error")
	}
	if !strings.Contains(err.Error(), "bad thing 7") {
		t.Errorf("err = %q, want the formatted message", err)
	}
}

// %w must survive, or a caller cannot unwrap to find the real cause.
func TestErrorfSupportsWrappedErrors(t *testing.T) {
	Configure(false, false, true)
	defer Configure(false, false, false)

	sentinel := errNotFound{}
	err := Errorf("reading %s: %w", "go.mod", sentinel)
	if !strings.Contains(err.Error(), "go.mod") {
		t.Errorf("err = %q, want the context included", err)
	}
	if !isSentinel(err) {
		t.Errorf("err = %#v, want the wrapped cause preserved", err)
	}
}

func TestNoColorIsRecorded(t *testing.T) {
	Configure(false, false, true)
	if !NoColorEnabled() {
		t.Error("NoColorEnabled should be true")
	}
	Configure(false, false, false)
	if NoColorEnabled() {
		t.Error("NoColorEnabled should be false")
	}
}

// With --no-color the output must contain no escape sequences at all.
func TestNoColorProducesNoEscapeSequences(t *testing.T) {
	_, errOut := capture(t, false, false, func() {
		Info("plain text")
		Success("ok")
		Warning("hmm")
		Error("no")
		Step("step")
	})
	if strings.Contains(errOut, "\x1b[") {
		t.Errorf("--no-color output still contains escape sequences: %q", errOut)
	}
}

type errNotFound struct{}

func (errNotFound) Error() string { return "not found" }

func isSentinel(err error) bool {
	for e := err; e != nil; {
		if _, ok := e.(errNotFound); ok {
			return true
		}
		u, ok := e.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		e = u.Unwrap()
	}
	return false
}
