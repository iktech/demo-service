package data

import "github.com/iktech/demo-service/core"

type Customers map[int]*core.Customer

var CustomersList Customers = map[int]*core.Customer{
	1: {
		ID:        1,
		FirstName: "Sarah",
		LastName:  "Brennan",
		Address:   "Croom, Co. Limerick",
	},
	2: {
		ID:        2,
		FirstName: "Eva",
		LastName:  "Olson",
		Address:   "3, Patrick st, Limerick",
	},
	4: {
		ID:        4,
		FirstName: "James",
		LastName:  "Brennan",
		Address:   "Croom, Co. Limerick",
	},
	16: {
		ID:        16,
		FirstName: "Dermot",
		LastName:  "Finnegan",
		Address:   "Anacotty, Co. Limerick",
	},
}
