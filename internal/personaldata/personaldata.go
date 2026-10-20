package personaldata

import "fmt"

// Personal содержит персональные данные пользователя.
type Personal struct {
	Name   string  // имя пользователя.
	Weight float64 // вес пользователя в кг.
	Height float64 // рост пользователя в м.
}

// Print выводит персональные данные пользователя.
func (p Personal) Print() {
	fmt.Printf("Имя: %s\nВес: %.2f\nРост: %.2f\n", p.Name, p.Weight, p.Height)
}
