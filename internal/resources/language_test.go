package resources

import (
	"testing"
)

// The generic profile declares `package_managers: []`. Slicing that to report a
// "preferred" manager panicked, which crashed `psx detect` on every project psx
// could not recognise.
func TestManagersForALanguageWithNoManagers(t *testing.T) {
	for _, lang := range []string{"generic", "Generic", ""} {
		got := ManagersFor(lang, map[string]bool{"go.mod": true})
		if len(got) != 0 {
			t.Errorf("ManagersFor(%q) = %v, want none", lang, got)
		}
	}
}

func TestManagersForAnUnknownLanguage(t *testing.T) {
	got := ManagersFor("cobol", map[string]bool{"main.go": true})
	if len(got) != 0 {
		t.Errorf("ManagersFor for an undeclared language = %v, want none", got)
	}
}

// With no manifest present the preferred manager is still reported, so generated
// instructions name a real tool rather than nothing.
func TestManagersForFallsBackToThePreferredManager(t *testing.T) {
	got := ManagersFor("go", map[string]bool{"unrelated.txt": true})
	if len(got) != 1 {
		t.Fatalf("got %d managers, want 1: %+v", len(got), got)
	}
	if got[0].File != "go.mod" {
		t.Errorf("preferred manager file = %q, want go.mod", got[0].File)
	}
}

// Detection from the filesystem avoids suggesting pnpm for an npm project.
func TestManagersAreFilteredByTheManifestActuallyPresent(t *testing.T) {
	got := ManagersFor("nodejs", map[string]bool{"package-lock.json": true, "package.json": true})
	if len(got) != 1 {
		t.Fatalf("got %d managers, want 1: %+v", len(got), got)
	}
	if got[0].Name != "npm" {
		t.Errorf("manager = %q, want npm for a package-lock.json project", got[0].Name)
	}
}

func TestManagersAreOrderedByPriority(t *testing.T) {
	got := ManagersFor("nodejs", nil)
	if len(got) < 2 {
		t.Skipf("nodejs declares %d managers, not enough to order", len(got))
	}
	for i := 1; i < len(got); i++ {
		if got[i-1].Priority > got[i].Priority {
			t.Errorf("managers are not ordered by priority: %+v", got)
			break
		}
	}
}

func TestSupportedLanguagesExcludesGenericAndIsSorted(t *testing.T) {
	got := SupportedLanguages()
	if len(got) == 0 {
		t.Fatal("SupportedLanguages returned nothing")
	}
	for i, l := range got {
		if l == "generic" {
			t.Error("generic must not be offered as a language to detect")
		}
		if i > 0 && got[i-1] > l {
			t.Errorf("SupportedLanguages is not sorted: %v", got)
			break
		}
	}
}

// Every language profile that detection can produce must survive a call with a
// populated and an empty manifest set.
func TestManagersForEveryDeclaredLanguage(t *testing.T) {
	for _, lang := range SupportedLanguages() {
		t.Run(lang, func(t *testing.T) {
			if _, ok := LanguageProfile(lang); !ok {
				t.Fatalf("SupportedLanguages lists %q but LanguageProfile does not know it", lang)
			}
			// Neither call may panic, whatever the manifest set.
			_ = ManagersFor(lang, nil)
			_ = ManagersFor(lang, map[string]bool{})
			_ = ManagersFor(lang, map[string]bool{"package.json": true})
		})
	}
}
