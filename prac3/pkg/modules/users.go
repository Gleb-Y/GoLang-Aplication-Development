package modules

type User struct {
	ID    int    `db:"id"    json:"id"`
	Name  string `db:"name"  json:"name"`
	Email string `db:"email" json:"email"`
	Age   int    `db:"age"   json:"age"`
	Phone string `db:"phone" json:"phone"`
}
