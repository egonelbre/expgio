package rig

type Registry struct {
	Editor map[string]EditorDefinition
}

func NewRegistry() *Registry {
	return &Registry{make(map[string]EditorDefinition)}
}

type EditorDefinition struct {
}
