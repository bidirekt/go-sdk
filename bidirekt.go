package bidirekt

type Resource struct{}

type ResourceOption interface {
	resourceOption()
}

type BodyOption interface {
	bodyOption()
}

type SchemaOption interface {
	schemaOption()
}

type Name struct{}

type SchemaDeclaration struct{}

type EnumDeclaration struct{}

type bodyDeclaration struct{}

func (bodyDeclaration) resourceOption() {}

func (Name) bodyOption() {}

func (Name) schemaOption() {}

func Provides(method, path string, options ...ResourceOption) Resource {
	return Resource{}
}

func Consumes(provider, method, path string, options ...ResourceOption) Resource {
	return Resource{}
}

func Request(body any, options ...BodyOption) ResourceOption {
	return bodyDeclaration{}
}

func Response(status int, body any, options ...BodyOption) ResourceOption {
	return bodyDeclaration{}
}

func Named(name string) Name {
	return Name{}
}

func Schema(value any, options ...SchemaOption) SchemaDeclaration {
	return SchemaDeclaration{}
}

func Enum[T ~string](values ...T) EnumDeclaration {
	return EnumDeclaration{}
}
