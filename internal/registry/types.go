package registry

type Registry[T any] struct {
	Name string
	items map[string]T
}
