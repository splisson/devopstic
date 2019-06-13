package entities

type User struct {
	Model
	Username  string `gorm:"type:varchar(100);unique_index"`
	FirstName string
	LastName  string
	Email string `gorm:"type:varchar(100);unique_index"`
}
