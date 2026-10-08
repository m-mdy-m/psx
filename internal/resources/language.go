package resources

import "sort"

// ManagersFor returns the package managers declared for a language, ordered by
// priority, filtered to those the project actually uses.
//
// Lockfiles discriminate first, because every Node manager shares the same
// package.json and they differ only by the lock they produce. Matching on the
// manifest alone offered an npm project a pnpm lockfile hint, which is the
// mistake this ordering exists to prevent.
func ManagersFor(lang string, present map[string]bool) []PackageManager {
	profile, ok := languages.Languages[normalizeKey(lang)]
	if !ok {
		return nil
	}

	all := append([]PackageManager(nil), profile.PackageManagers...)
	sort.SliceStable(all, func(i, j int) bool { return all[i].Priority < all[j].Priority })

	// A language may legitimately declare no package manager; the generic
	// profile does. Slicing it below would panic.
	if len(all) == 0 {
		return nil
	}

	if len(present) == 0 {
		return all
	}

	var byLock, byManifest []PackageManager
	for _, m := range all {
		if m.Lock != "" && present[m.Lock] {
			byLock = append(byLock, m)
		}
		if m.File != "" && present[m.File] {
			byManifest = append(byManifest, m)
		}
	}
	if len(byLock) > 0 {
		return byLock
	}
	if len(byManifest) > 0 {
		return byManifest
	}

	// No manifest present: report the preferred manager so generated
	// instructions still name a real tool.
	return all[:1]
}

func LanguageProfile(lang string) (Language, bool) {
	profile, ok := languages.Languages[normalizeKey(lang)]
	return profile, ok
}

func SupportedLanguages() []string {
	out := make([]string, 0, len(languages.Languages))
	for name := range languages.Languages {
		if name == "generic" {
			continue
		}
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func NormalizeLanguage(lang string) string {
	return NormalizeProjectType(lang)
}
