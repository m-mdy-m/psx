package resources

import "sort"

// ManagersFor returns the package managers declared for a language, ordered by
// priority, filtered to those whose manifest is actually present.
//
// Detecting from the filesystem rather than assuming the first entry avoids
// generating a pnpm lockfile hint for an npm project.
func ManagersFor(lang string, present map[string]bool) []PackageManager {
	profile, ok := languages.Languages[normalizeKey(lang)]
	if !ok {
		return nil
	}

	all := append([]PackageManager(nil), profile.PackageManagers...)
	sort.SliceStable(all, func(i, j int) bool { return all[i].Priority < all[j].Priority })

	if len(present) == 0 {
		return all
	}

	var matched []PackageManager
	for _, m := range all {
		if m.File != "" && present[m.File] {
			matched = append(matched, m)
		}
	}
	if len(matched) > 0 {
		return matched
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
