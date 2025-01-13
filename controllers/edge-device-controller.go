package handlers

type Databse interface {
	Init()
	Create()
	Replace()
	Update()
	Delete()
}

func init() error {

}
