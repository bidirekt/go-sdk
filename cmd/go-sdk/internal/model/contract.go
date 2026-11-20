package model

// A location without File has no place.
type SourceLocation struct {
	File string
	Line int
}

type Diagnostic struct {
	Message  string
	Location SourceLocation
}

type Contract struct {
	Provides []Resource
	Consumes []Resource
	Schemas  []NamedSchema
}

type Resource struct {
	Provider  string
	Method    string
	Path      string
	Request   *Body
	Responses []Response
	Location  SourceLocation
}

type Response struct {
	Status int
	Body   Body
}

// A Body without Schema is written as ~.
type Body struct {
	Schema   string
	Location SourceLocation
}

// A Derived schema comes from a body, a list or a body named with bidirekt.Named: the same one declared by
// several bodies is one schema.
type NamedSchema struct {
	Name     string
	Node     Node
	Location SourceLocation
	Derived  bool
}

type Node struct {
	Shape       Shape
	Optional    bool
	Nullable    bool
	Description string
}

type Shape interface {
	shape()
}

type Primitive struct {
	Type string
}

type Object struct {
	Properties []Property
}

type Array struct {
	Items Node
}

type Map struct {
	Values Node
}

type Enum struct {
	Members []string
}

type Union struct {
	Members []UnionMember
}

type Ref struct {
	Schema string
}

type Property struct {
	Name     string
	Node     Node
	Location SourceLocation
}

type UnionMember struct {
	Key  string
	Node Node
}

func (Primitive) shape() {}
func (Object) shape()    {}
func (Array) shape()     {}
func (Map) shape()       {}
func (Enum) shape()      {}
func (Union) shape()     {}
func (Ref) shape()       {}
