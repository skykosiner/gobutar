package localstorage

type LocalStorage struct {
	items map[string]string
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		items: make(map[string]string),
	}
}

func (L *LocalStorage) Store(key, value string) {
	L.items[key] = value
}

func (L *LocalStorage) Get(key string) string {
	return L.items[key]
}
