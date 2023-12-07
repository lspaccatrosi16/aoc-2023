package challenge

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Card int

const (
	ACE Card = iota
	KING
	QUEEN
	JOKER
	TEN
	NINE
	EIGHT
	SEVEN
	SIX
	FIVE
	FOUR
	THREE
	TWO
)

func (card Card) String() string {
	switch card {
	case ACE:
		return "ACE"
	case KING:
		return "KING"

	case QUEEN:
		return "QUEEN"

	case JOKER:
		return "JOKER"

	case TEN:
		return "TEN"

	case NINE:
		return "NINE"

	case EIGHT:
		return "EIGHT"

	case SEVEN:
		return "SEVEN"

	case SIX:
		return "SIX"

	case FIVE:
		return "FIVE"

	case FOUR:
		return "FOUR"

	case THREE:
		return "THREE"

	case TWO:
		return "TWO"

	default:
		return "INVALID"
	}
}
func (c Card) Strength() int {
	if c == JOKER {
		return 1

	}
	return 14 - int(c)
}

func (car Card) FromRune(r rune) Card {
	var c Card

	switch r {
	case 'A':
		c = ACE
	case 'K':
		c = KING
	case 'Q':
		c = QUEEN
	case 'J':
		c = JOKER
	case 'T':
		c = TEN
	case '9':
		c = NINE
	case '8':
		c = EIGHT
	case '7':
		c = SEVEN
	case '6':
		c = SIX
	case '5':
		c = FIVE
	case '4':
		c = FOUR
	case '3':
		c = THREE
	case '2':
		c = TWO
	}

	return c
}

type HandType int

const (
	FIVEK HandType = iota
	FOURK
	FULLH
	THREEK
	TWOP
	ONEP
	HIGHC
)

func (h HandType) String() string {
	switch h {
	case FIVEK:
		return "FIVE OF A KIND"
	case FOURK:
		return "FOUR OF A KIND"
	case FULLH:
		return "FULL HOUSE"
	case THREEK:
		return "THREE OF KIND"
	case TWOP:
		return "TWO PAIRS"
	case ONEP:
		return "ONE PAIR"
	case HIGHC:
		return "HIGH CARD"
	}

	return "INVALID"
}

func (h HandType) Strength() int {
	return 7 - int(h)
}

type Hand struct {
	Cards [5]Card
	Bet   int
	Rank  int
	Idx   int
	Type  HandType
}

func (h *Hand) Winnings() int {
	return h.Bet * h.Rank
}

func (h *Hand) GenType() {
	seenTypes := map[Card]int{}

	numJokers := 0

	for _, c := range h.Cards {
		if c == JOKER {
			numJokers++
		} else {
			seenTypes[c] += 1
		}
	}

	maxFound := 0

	for _, i := range seenTypes {
		if i > maxFound {
			maxFound = i
		}
	}

	maxFound += numJokers

	if len(seenTypes) == 1 {
		h.Type = FIVEK
	} else if len(seenTypes) == 2 {
		if maxFound == 4 {
			h.Type = FOURK
		} else {
			h.Type = FULLH
		}
	} else if len(seenTypes) == 3 {
		if maxFound == 3 {
			h.Type = THREEK
		} else {
			h.Type = TWOP
		}
	} else if len(seenTypes) == 4 {
		h.Type = ONEP
	} else {
		h.Type = HIGHC
	}
}

type HandList []Hand

func (h HandList) Len() int {
	return len(h)
}

func (h HandList) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h HandList) Less(i, j int) bool {
	if h[i].Type != h[j].Type {
		iStr := h[i].Type.Strength()
		jStr := h[j].Type.Strength()
		return iStr < jStr
	}
	for k := 0; k < 5; k++ {
		if h[i].Cards[k].String() != h[j].Cards[k].String() {
			return h[i].Cards[k].Strength() < h[j].Cards[k].Strength()
		}
	}

	fmt.Println("Exactly identical score")
	return false
}

func Challenge(input string) string {
	lines := strings.Split(input, "\n")

	hands := HandList{}

	for j, l := range lines {
		hasJokers := false

		components := strings.Split(l, " ")
		cards := [5]Card{}
		for i := 0; i < 5; i++ {
			card := new(Card).FromRune(rune(components[0][i]))

			if card == JOKER {
				hasJokers = true
			}

			cards[i] = card
		}

		bet, err := strconv.Atoi(components[1])
		if err != nil {
			panic("int conv error")
		}

		hand := Hand{
			Cards: cards,
			Bet:   bet,
			Idx:   j,
		}

		if hasJokers {
			searchedCards := [5]Card{}
			for i := 0; i < 5; i++ {
				if cards[i] == JOKER {
					maxS := 0
					maxC := Card(0)
					for k := 0; k < 13; k++ {
						newCards := [5]Card{}
						for m, c := range cards {
							if c == JOKER {
								newCards[m] = Card(k)
							} else {
								newCards[m] = c
							}
						}

						provHand := Hand{Cards: [5]Card(newCards)}
						provHand.GenType()
						provStrength := provHand.Type.Strength()

						if provStrength > maxS {
							maxS = provStrength
							maxC = Card(k)
						}
					}
					searchedCards[i] = maxC
				} else {
					searchedCards[i] = cards[i]
				}
			}
			hand.Cards = searchedCards
			hand.GenType()

			hand.Cards = cards
		} else {
			hand.GenType()
		}

		hands = append(hands, hand)
	}

	sort.Sort(hands)
	tally := 0

	for i := 0; i < len(hands); i++ {
		hands[i].Rank = i + 1
		w := hands[i].Winnings()

		tally += w
	}

	return strconv.Itoa(tally)
}
