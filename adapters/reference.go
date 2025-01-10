package adapters

type ReferenceReturn struct {
	Tipo  string
	Value string
}

func NewAdapterReferenceReturn(tipo, value string) *ReferenceReturn {
	return &ReferenceReturn{
		Tipo:  tipo,
		Value: value,
	}

}
