package usecases

import "go/clean/entities"

type OrderRepository interface {
	Save(order entities.Order) error
}
