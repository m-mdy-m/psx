package detector

// ProjectType represents detected project information
type ProjectType struct {
	Primary        string            // Main language (nodejs, go, rust, python, etc.)
	Secondary      []string          // Additional languages detected
	Framework      string            // Framework name (react, django, gin, etc.)
	PackageManager string            // Package manager (npm, cargo, pip, etc.)
	BuildTool      string            // Build tool (webpack, maven, gradle, etc.)
	Structure      string            // Project structure (application, library, workspace)
	Version        string            // Language/framework version if detectable
	Features       map[string]bool   // Additional features (typescript, esm, etc.)
	Confidence     float64           // Detection confidence (0.0 to 1.0)
}

// DetectionResult contains complete detection information
type DetectionResult struct {
	Type        ProjectType  // Detected project type
	Files       []string     // Key files that led to detection
	Description string       // Human-readable description
	Suggestions []string     // Suggestions for improvement
}

// LanguageSignature defines how to detect a language
type LanguageSignature struct {
	Name           string   // Language identifier (nodejs, go, rust, etc.)
	DisplayName    string   // Human-readable name (Node.js, Go, Rust, etc.)
	PrimaryFiles   []string // Files that strongly indicate this language
	SecondaryFiles []string // Files that weakly indicate this language
	Extensions     []string // File extensions (.js, .go, .rs, etc.)
	MinConfidence  float64  // Minimum confidence threshold
}

// Detector interface for language detection
type Detector interface {
	// GetSignature returns the language signature
	GetSignature() LanguageSignature

	// Detect attempts to detect this language in the project
	Detect(projectPath string) (*DetectionResult, error)

	// CanDetect quickly checks if detection is possible
	CanDetect(projectPath string) bool
}

// DetectorRegistry manages all language detectors
type DetectorRegistry struct {
	detectors map[string]Detector
	order     []string // Detection order
}

// NewDetectorRegistry creates a new registry
func NewDetectorRegistry() *DetectorRegistry {
	return &DetectorRegistry{
		detectors: make(map[string]Detector),
		order:     []string{},
	}
}

// Register adds a detector to the registry
func (r *DetectorRegistry) Register(detector Detector) {
	sig := detector.GetSignature()
	r.detectors[sig.Name] = detector
	r.order = append(r.order, sig.Name)
}

// Get retrieves a detector by name
func (r *DetectorRegistry) Get(name string) (Detector, bool) {
	detector, exists := r.detectors[name]
	return detector, exists
}

// All returns all registered detectors in order
func (r *DetectorRegistry) All() []Detector {
	result := make([]Detector, 0, len(r.order))
	for _, name := range r.order {
		if detector, exists := r.detectors[name]; exists {
			result = append(result, detector)
		}
	}
	return result
}

// GetNames returns all registered detector names
func (r *DetectorRegistry) GetNames() []string {
	return append([]string{}, r.order...) // Return a copy
}

