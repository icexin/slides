package main

import (
	"fmt"
	"time"
)

type Account struct {
	money int
}

func (a *Account) DoPrepare() {
	time.Sleep(time.Second)
}

func (a *Account) GetGongZi(n int) {
	a.money += n
}

func (a *Account) GiveWife(n int) {
	if a.money > n {
		a.DoPrepare()
		a.money -= n
	}
}

func (a *Account) Buy(n int) {
	if a.money > n {
		a.DoPrepare()
		a.money -= n
	}
}

func (a *Account) Left() int {
	return a.money
}

func main() {
	var account Account
	account.GetGongZi(10)
	go account.GiveWife(6)
	go account.Buy(5)
	time.Sleep(2 * time.Second)
	fmt.Println(account.Left())
}
