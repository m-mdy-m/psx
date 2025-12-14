package detector

type ProjectType struct {
	Primary        string
	Version        string
	Features       map[string]bool
}

type DetectionResult struct {
	Type        ProjectType
	Files       []string
}

type LanguageSignature struct {
	Name           string
	PrimaryFiles   []string
}

type Detector interface {
	GetSignature() LanguageSignature
	Detect(projectPath string) (*DetectionResult, error)
	CanDetect(projectPath string) bool
}
