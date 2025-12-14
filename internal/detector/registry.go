package detector

import (
	"strings"

	reg "github.com/m-mdy-m/psx/internal/registry"
)

type DetectorRegistry struct {
	registry	*reg.Registry[Detector]
	order		[]string
}

func NewDetectorRegistry() *DetectorRegistry {
	return &DetectorRegistry{
		registry: reg.New[Detector]("detector"),
		order:     []string{},
	}
}
func (r *DetectorRegistry) Register(detector Detector) {
	sig := detector.GetSignature()
	key := strings.ToLower(sig.Name)

	if _,ok := r.registry.Get(key);!ok{
		r.order = append(r.order,key)
	}
	r.registry.Add(key,detector)
}
func (r *DetectorRegistry) Get(name string) (Detector, bool) {
	return r.registry.Get(strings.ToLower(name))
}
func (r *DetectorRegistry) All() []Detector {
	items := r.registry.All()
	result := make([]Detector, 0, len(items))

	for _, key := range r.order {
		if d, ok := items[key]; ok {
			result = append(result, d)
		}
	}

	return result
}
func (r *DetectorRegistry) GetNames() []string {
	return append([]string{}, r.order...)
}

