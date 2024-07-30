package maps

type MyMap struct {
	data map[string]string
}

func New() *MyMap {
	return &MyMap{
		data: make(map[string]string),
	}
}

func (m *MyMap) Add(key string, value string) {
	m.data[key] = value
}

func (m *MyMap) Get(key string) (string, bool) {
	value, exists := m.data[key]
	return value, exists
}

func (m *MyMap) Delete(key string) {
	delete(m.data, key)
}
