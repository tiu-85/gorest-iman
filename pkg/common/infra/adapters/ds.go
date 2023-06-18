package adapters

type DS struct {
	DB
}

func NewDS(DB DB) *DS {
	return &DS{
		DB: DB,
	}
}
