package user


type User struct{
	ID int64	`gorm:"primaryKey"`
	Email string
	Password string
}
