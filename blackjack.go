package main

import (
	"fmt"
	"math/rand"
	"time"
)

//Structs we be using

type Card struct {
	num   string
	suit  string
	value int
}

func (c Card) ToString() string {
	if c.value != 1 {
		return fmt.Sprintf("%v of %v (%v)", c.num, c.suit, c.value)
	}
	return fmt.Sprintf("Ace of %v (1 or 11)", c.suit)
}

type Player struct {
	hand     map[uint8]Card
	hasStood bool
}

func (p Player) GetValue() (v int) {
	aceCount := 0
	for _, c := range p.hand {
		v += c.value
		if c.value == 1 {
			aceCount++
		}
	}
	if aceCount != 0 {
		for i := 0; i < aceCount; i++ {
			if v+10 <= 21 {
				v += 10
			}
		}
	}
	return
}

func (p Player) Hand2String() (s string) {
	forLoopIters := uint8(len(p.hand) - 1)
	for i := uint8(0); i < forLoopIters; i++ {
		s = fmt.Sprintf("%v[%v], ", s, p.hand[i].ToString())
	}
	s = fmt.Sprintf("%v[%v]", s, p.hand[forLoopIters].ToString())
	return
}

//returns a new index for the deck
func (p *Player) HitMe(deck [52]Card, i uint8) (nI uint8) {
	p.hand[uint8(len(p.hand))] = deck[i]
	nI = i + 1
	return
}

type sRand struct {
	r rand.Rand
}

//global vars

var nums = [13]string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen",
	"King"}
var suits = [4]string{"Spades", "Clubs", "Hearts", "Diamonds"}

//useful functions

//the reason GMASD() is O(n)
func RandomInRange(r sRand, start int, end int) int {
	diff := end - start
	rNum := r.r.Intn(diff)
	return rNum + start
}

func GiveMeAShuffledDeck() [52]Card {
	ret := [52]Card{}
	src := time.Now().UnixNano()
	rg := rand.New(rand.NewSource(src))
	r := sRand{*rg}

	//make me an ordered deck, doesn't deserve it's own function
	for i, _ := range ret {
		ret[i] = Card{
			nums[i%13],
			suits[i/13],
			i%13 + 1, //gotta love semicolon auto insertion
		}
	}

	//shuffle it!
	for c := 0; c < 51; c++ {
		rir := RandomInRange(r, c+1, 52)
		temp := ret[rir]
		ret[rir] = ret[c]
		ret[c] = temp
	}

	return ret
}

func NewGame() (d [52]Card, dI uint8, user, dealer Player) {
	d = GiveMeAShuffledDeck()
	playerHand, dealerHand := make(map[uint8]Card), make(map[uint8]Card)
	playerHand[0] = d[0]
	playerHand[1] = d[1]
	dealerHand[0] = d[2]
	dealerHand[1] = d[3]
	user, dealer = Player{playerHand, false}, Player{dealerHand, false}
	dI = 4
	return
}

func RunDealerLogic(deck [52]Card, deckI uint8, u, d Player) (nD Player) {
	nD = d
	dV := nD.GetValue()
	uV := u.GetValue()
	for (dV < uV) || (dV == uV && dV < 17) {
		deckI = nD.HitMe(deck, deckI)
		dV = nD.GetValue()
	}
	return
}

func main() {

	deck, deckI, user, dealer := NewGame()
	for !user.hasStood && user.GetValue() < 22 {
		fmt.Printf("Dealer:\n[%v], [ hidden ]\nYour hand:\n%v (%v)\nhit (h) or stand (s)?\n",
			dealer.hand[0].ToString(), user.Hand2String(), user.GetValue())
		var input string
		_, err := fmt.Scanf("%s", &input)
		if err != nil {
			fmt.Println(err)
		} else if input == "h" {
			deckI = user.HitMe(deck, deckI)
		} else if input == "s" {
			user.hasStood = true
		}
	}
	finalUserValue := user.GetValue()
	finalDealerValue := dealer.GetValue()
	if finalUserValue > 21 {
		fmt.Println("You busted!")
	} else {
		dealer = RunDealerLogic(deck, deckI, user, dealer)
		finalDealerValue = dealer.GetValue()
		if finalDealerValue > 21 {
			fmt.Println("Dealer busted!")
		} else if finalUserValue > finalDealerValue {
			fmt.Println("You won!")
		} else if finalDealerValue > finalUserValue {
			fmt.Println("Dealer won! the house always wins >:^)")
		} else {
			fmt.Println("tie")
		}
	}
	fmt.Printf("Dealer's hand:\n%v (%v)\nYour hand:\n%v (%v)\n",
		dealer.Hand2String(), finalDealerValue, user.Hand2String(), finalUserValue)
}
