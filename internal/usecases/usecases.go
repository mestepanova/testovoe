package usecases

type Usecases struct {
	storage Storage
}

func NewUsecases(storage Storage) *Usecases {
	return &Usecases{
		storage: storage,
	}
}
