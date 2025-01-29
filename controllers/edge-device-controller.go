package handlers

type Databse interface {
	Init() error
	Create() error
	Replace() error
	Update() error
	Delete() error
}
