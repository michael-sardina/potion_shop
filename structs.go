package main

type Store struct {
	Gold          int
	Inventory     map[string]int
	SelectedColor string
	Prices        map[string]int
}

type View struct {
	Store        Store
	TotalPotions int
	CurrentPrice int
}
