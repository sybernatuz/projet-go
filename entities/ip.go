package entities

type Ip struct {
	ID      int `gorm:"primary_key, AUTO_INCREMENT"`
	Value   string
	Attempt int
}
