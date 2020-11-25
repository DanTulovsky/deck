package deck

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/nfnt/resize"
	gim "github.com/ozankasikci/go-image-merge"

	ppb "github.com/DanTulovsky/pepper-poker-v2/proto"
)

var (
	// Divider is the divider between the hole and board in the terminal
	Divider = flag.String("card_divider", "../images/deck/deck1/blank.png", "divider between the hole and board")
)

// Deck is a deck of cards.
type Deck struct {
	cards []*Card
}

// NewDeck returns a new, unshuffled deck of cards.
func NewDeck() *Deck {

	suits := []ppb.CardSuit{
		ppb.CardSuit_Club,
		ppb.CardSuit_Diamond,
		ppb.CardSuit_Spade,
		ppb.CardSuit_Heart,
	}

	ranks := []ppb.CardRank{
		ppb.CardRank_Ace,
		ppb.CardRank_King,
		ppb.CardRank_Queen,
		ppb.CardRank_Jack,
		ppb.CardRank_Ten,
		ppb.CardRank_Nine,
		ppb.CardRank_Eight,
		ppb.CardRank_Seven,
		ppb.CardRank_Six,
		ppb.CardRank_Five,
		ppb.CardRank_Four,
		ppb.CardRank_Three,
		ppb.CardRank_Two,
	}

	nSuits := len(suits)
	nRanks := len(ranks)

	var c = make([]*Card, nSuits*nRanks)
	for i := 0; i < nSuits; i++ {
		for j := 0; j < nRanks; j++ {
			c[i*nRanks+j] = NewCard(suits[i], ranks[j])
			ppb.CardSuit_Club.Descriptor()
		}
	}
	return &Deck{cards: c[:]}
}

// NewShuffledDeck returns a shuffled deck
func NewShuffledDeck() *Deck {
	d := NewDeck()
	d.Shuffle()
	return d
}

// Shuffle shuffles the deck.
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.cards), func(x, y int) { d.cards[x], d.cards[y] = d.cards[y], d.cards[x] })
}

// IsEmpty returns whether the deck has any remaining cards.
func (d *Deck) IsEmpty() bool {
	return len(d.cards) == 0
}

// Next returns the next card from the deck, or an error if there are no remaining cards.
func (d *Deck) Next() (*Card, error) {
	if d.IsEmpty() {
		return nil, errors.New("deck is empty")
	}
	c := d.cards[0]
	d.cards = d.cards[1:]
	return c, nil
}

// Remove remove the given card from the deck
func (d *Deck) Remove(card *Card) error {
	if d.IsEmpty() {
		return errors.New("deck is empty")
	}

	for i, c := range d.cards {
		if c.GetRank() == card.GetRank() && c.GetSuit() == card.GetSuit() {
			// Remove the element at index i from a.
			copy(d.cards[i:], d.cards[i+1:])   // Shift a[i+1:] left one index.
			d.cards[len(d.cards)-1] = nil      // Erase last element (write zero value).
			d.cards = d.cards[:len(d.cards)-1] // Truncate slice.
			return nil
		}
	}
	return nil
}

// Return returns the card to the bottom of the deck
func (d *Deck) Return(card *Card) error {
	if CardInList(card, d.cards) {
		return fmt.Errorf("cannot add duplicate card %v", card)
	}

	d.cards = append(d.cards, card)
	return nil
}

// RandomCard returns a random card (not from a deck)
func RandomCard() *Card {
	return &Card{
		Card: &ppb.Card{
			Suite: RandomSuit(),
			Rank:  RandomRank(),
		},
	}
}

// RandomSuit returns a suit at random
func RandomSuit() ppb.CardSuit {
	return ppb.CardSuit(rand.Int31n(4))
}

// RandomSuitNotIn returns a suit at random not in suits
func RandomSuitNotIn(suits ...ppb.CardSuit) ppb.CardSuit {
	if len(suits) == 4 {
		log.Fatal("there are only 4 suits total")
	}

	var suit ppb.CardSuit

	for {
		suit = RandomSuit()
		for _, s := range suits {
			if suit == s {
				continue
			}
			return suit
		}
	}
}

// RandomRank returns a random rank
func RandomRank() ppb.CardRank {
	return ppb.CardRank(rand.Int31n(13))
}

// RandomRankNotIn returns a random rank that is not in ranks
func RandomRankNotIn(ranks ...ppb.CardRank) ppb.CardRank {
	if len(ranks) == 13 {
		log.Fatal("there are only 13 ranks total")
	}

OUTER:
	for {
		rank := ppb.CardRank(rand.Int31n(13))
		for _, r := range ranks {
			if rank == r {
				continue OUTER
			}
		}
		return rank
	}
}

// RandomRankAbove returns a random rank that is > r
func RandomRankAbove(r ppb.CardRank) ppb.CardRank {
	var rank ppb.CardRank

	for rank <= r {
		rank = ppb.CardRank(rand.Int31n(13))
	}

	return rank
}

// CardsImage combines all the cards into a single horizontal image and returns it
func CardsImage(cards []*Card, divider bool) (image.Image, error) {
	grids := []*gim.Grid{}

	length := len(cards)
	if divider {
		length++
	}

	for i := 0; i < len(cards); i++ {
		// first two cards are the hole, add an empty space after
		grids = append(grids,
			&gim.Grid{
				ImageFilePath: cards[i].ImageFileName(),
			},
		)

		// Append divider
		if divider && i == 1 {
			grids = append(grids,
				&gim.Grid{
					ImageFilePath: *Divider,
				},
			)
		}
	}

	for _, c := range cards {
		grids = append(grids,
			&gim.Grid{
				ImageFilePath: c.ImageFileName(),
			},
		)
	}

	var width = uint(180 * length)
	// height := int(float64(width) / 1.5)
	img, err := gim.New(grids, length, 1).Merge()
	if err != nil {
		return nil, err
	}

	return resize.Resize(width, 0, img, resize.NearestNeighbor), nil

}
