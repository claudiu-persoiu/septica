package game

type Card struct {
	Number int
	Type   int
}

var suits = [4]string{"diamond", "hearts", "spades", "clubs"}
var ranks = [8]string{"7", "8", "9", "10", "J", "Q", "K", "A"}
