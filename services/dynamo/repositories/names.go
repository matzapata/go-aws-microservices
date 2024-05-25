package repositories

type NamesRepository interface {
	CreateName(name string) (string, error)
}
