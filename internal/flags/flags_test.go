package flags

import (
	"strings"
	"testing"
	"time"
)

func TestDefaultsAreValid(t *testing.T) {
	o := Defaults()
	if err := o.Validate(); err != nil {
		t.Fatalf("the default option set does not validate: %v", err)
	}
	if o.Check.OutputFormat != FormatTable {
		t.Errorf("default output format = %q, want %q", o.Check.OutputFormat, FormatTable)
	}
	if o.Check.FailOn != "error" {
		t.Errorf("default --fail-on = %q, want error", o.Check.FailOn)
	}
	if o.Check.Level != "all" {
		t.Errorf("default --level = %q, want all", o.Check.Level)
	}
	if o.Fix.Interactive {
		t.Error("fix must not prompt unless asked; a CI run must never block on stdin")
	}
	if o.Fix.Answers == nil {
		t.Error("Answers must be initialised, or a fix that reads an answer writes to a nil map")
	}
}

func TestVerboseIsSuppressedByQuiet(t *testing.T) {
	o := Defaults()
	o.Global.Verbose = true
	if !o.Verbose() {
		t.Error("--verbose alone should enable verbose output")
	}

	o.Global.Quiet = true
	if o.Verbose() {
		t.Error("--quiet must win over --verbose; a quiet run is quiet")
	}

	o.Global.Verbose = false
	if o.Verbose() {
		t.Error("--quiet alone should not enable verbose output")
	}
}

func TestValidateRejectsUnknownEnums(t *testing.T) {
	tests := []struct {
		name  string
		apply func(*Options)
		want  string
	}{
		{"output format", func(o *Options) { o.Check.OutputFormat = "yaml" }, "unsupported output format"},
		{"level", func(o *Options) { o.Check.Level = "fatal" }, "unsupported severity level"},
		{"fail-on", func(o *Options) { o.Check.FailOn = "always" }, "unsupported --fail-on"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			o := Defaults()
			tc.apply(&o)
			err := o.Validate()
			if err == nil {
				t.Fatalf("Validate accepted %s", tc.name)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("err = %q, want it to mention %q", err, tc.want)
			}
		})
	}
}

func TestValidateAcceptsEveryDocumentedFormat(t *testing.T) {
	for _, f := range []string{
		FormatTable, FormatCompact, FormatJSON, FormatNDJSON,
		FormatSARIF, FormatGitHub, FormatJUnit, FormatMarkdown,
	} {
		o := Defaults()
		o.Check.OutputFormat = f
		if err := o.Validate(); err != nil {
			t.Errorf("Validate rejected the documented format %q: %v", f, err)
		}
	}
}

func TestValidateAcceptsEveryDocumentedLevelAndThreshold(t *testing.T) {
	for _, l := range []string{"all", "error", "warning", "info"} {
		o := Defaults()
		o.Check.Level = l
		if err := o.Validate(); err != nil {
			t.Errorf("Validate rejected level %q: %v", l, err)
		}
	}
	for _, f := range []string{"error", "warning", "none"} {
		o := Defaults()
		o.Check.FailOn = f
		if err := o.Validate(); err != nil {
			t.Errorf("Validate rejected --fail-on %q: %v", f, err)
		}
	}
}

// --dry-run promises nothing is written, so --force, which is about overwriting,
// cannot both be in effect.
func TestValidateRejectsDryRunWithForce(t *testing.T) {
	o := Defaults()
	o.Fix.DryRun = true
	o.Fix.Force = true

	err := o.Validate()
	if err == nil {
		t.Fatal("--dry-run and --force should not combine")
	}
	if !strings.Contains(err.Error(), "--dry-run") || !strings.Contains(err.Error(), "--force") {
		t.Errorf("err = %q, want it to name both flags", err)
	}
}

func TestValidateAllowsDryRunAloneAndForceAlone(t *testing.T) {
	o := Defaults()
	o.Fix.DryRun = true
	if err := o.Validate(); err != nil {
		t.Errorf("--dry-run alone should be valid: %v", err)
	}

	o = Defaults()
	o.Fix.Force = true
	if err := o.Validate(); err != nil {
		t.Errorf("--force alone should be valid: %v", err)
	}
}

func TestMachineFormatsAreIdentified(t *testing.T) {
	machine := map[string]bool{
		FormatJSON:   true,
		FormatNDJSON: true,
		FormatSARIF:  true,
		FormatJUnit:  true,
		FormatTable:  false,
		FormatGitHub: false,
	}
	for format, want := range machine {
		o := Defaults()
		o.Check.OutputFormat = format
		if got := o.machineFormat(); got != want {
			t.Errorf("machineFormat(%q) = %v, want %v", format, got, want)
		}
	}
}

func TestInteractiveIsOptIn(t *testing.T) {
	o := Defaults()
	if o.Fix.InteractiveWanted() {
		t.Error("fix must not want interaction by default")
	}
	o.Fix.Interactive = true
	if !o.Fix.InteractiveWanted() {
		t.Error("-i should request interaction")
	}
}

// The interval floor keeps a mistyped --interval from pinning a CPU core.
func TestWatchTimingDefaultsAreSane(t *testing.T) {
	d := Defaults()
	if d.Watch.Interval < 100*time.Millisecond {
		t.Errorf("default interval %v is below the 100ms floor the command enforces", d.Watch.Interval)
	}
	if d.Watch.Debounce <= 0 {
		t.Errorf("default debounce %v must be positive", d.Watch.Debounce)
	}
	if d.Watch.Debounce >= d.Watch.Interval {
		t.Errorf("debounce %v must be shorter than interval %v", d.Watch.Debounce, d.Watch.Interval)
	}
}
