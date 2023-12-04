package option

import "github.com/sagernet/sing/common/json"

type TypedMessage[T any] struct {
	Message json.RawMessage
	Value   T
}

func (m *TypedMessage[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

func (m *TypedMessage[T]) UnmarshalJSON(bytes []byte) error {
	m.Message = bytes
	return json.Unmarshal(bytes, &m.Value)
}
