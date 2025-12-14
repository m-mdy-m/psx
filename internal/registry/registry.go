package registry


func New[T any](name string)* Registry[T] {
	return  &Registry[T]{
		Name:name,
		items: make(map[string]T),
	}
}
func (r *Registry[T]) Add(name string,value T){
	r.items[name] = value
}
func (r * Registry[T]) Get(name string) (T,bool){
	v,ok :=r.items[name]
	return v,ok
}
func (r *Registry[T])All()map[string]T{
	clone :=make(map[string]T,len(r.items))
	for k,v := range r.items{
		clone[k]=v
	}
	return clone
}
func (r*Registry[T])Delete(name string){
	delete(r.items,name)
}

