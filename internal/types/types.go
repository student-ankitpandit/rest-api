package types


type Student struct {
	Id int
	Name string `validate:"required"`
	Email string `validate:"required"`
	Password string `validate:"required"`
}