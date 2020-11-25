package deck

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path"

	// for Image() function
	_ "image/png"

	ppb "github.com/DanTulovsky/pepper-poker-v2/proto"

	imgcat "github.com/martinlindhe/imgcat/lib"
	"github.com/qeesung/image2ascii/convert"
)

var (
	suitmap = map[string]string{
		"Club":    "♣",
		"Diamond": "♦",
		"Heart":   "♥",
		"Spade":   "♠",
	}

	rankmap = map[ppb.CardRank]string{
		ppb.CardRank_Two:   "2",
		ppb.CardRank_Three: "3",
		ppb.CardRank_Four:  "4",
		ppb.CardRank_Five:  "5",
		ppb.CardRank_Six:   "6",
		ppb.CardRank_Seven: "7",
		ppb.CardRank_Eight: "8",
		ppb.CardRank_Nine:  "9",
		ppb.CardRank_Ten:   "10",
		ppb.CardRank_Jack:  "J",
		ppb.CardRank_Queen: "Q",
		ppb.CardRank_King:  "K",
		ppb.CardRank_Ace:   "A",
	}

	deckDir = flag.String("deck_dir", "../images/deck/deck1", "directory that contains the deck of cards to use")
)

// SortByCards puts the cards in order
// Low -> High: 2 -> Ace
// Low -> High: Spade, Club, Diamond, Heart
type SortByCards []*Card

func (a SortByCards) Len() int      { return len(a) }
func (a SortByCards) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByCards) Less(i, j int) bool {
	return a[i].IsLessThan(a[j])
}

// Card is a card.
type Card struct {
	Card *ppb.Card
}

// NewCard returns a new card.
func NewCard(s ppb.CardSuit, r ppb.CardRank) *Card {
	return &Card{
		&ppb.Card{
			Suite: s,
			Rank:  r,
		},
	}
}

// NewRandomCard returns a random card
func NewRandomCard() *Card {
	return &Card{
		&ppb.Card{
			Suite: RandomSuit(),
			Rank:  RandomRank(),
		},
	}
}

// ShowImage prints the image of the card to the terminal
func (c *Card) ShowImage() error {
	// enc, err := imgcat.NewEncoder(os.Stdout, imgcat.Width(imgcat.Pixels(100)), imgcat.Inline(true))
	// if err != nil {
	// 	return err
	// }

	// img, err := c.ImageFile()
	// if err := enc.Encode(img); err != nil {
	// 	return err
	// }
	// return nil

	// https://github.com/martinlindhe/imgcat
	if img, err := c.ImageFile(); err == nil {
		imgcat.Cat(img, os.Stdout)
	}

	return nil
}

// Image returns an image of the card
func (c *Card) Image() (image.Image, error) {
	f, err := c.ImageFile()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	i, _, err := image.Decode(f)
	return i, nil

}

// ImageFileName returns the name of the file containing the image
func (c *Card) ImageFileName() string {
	file := fmt.Sprintf("%s%s.png", rankmap[c.GetRank()], string(c.GetSuit().String()[0]))
	return path.Join(*deckDir, file)
}

// ImageFile returns the file containing the image
func (c *Card) ImageFile() (*os.File, error) {
	file := c.ImageFileName()

	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("unable to open %v: %v", file, err)
	}
	return f, nil
}

// AsASCII returns the card's image as ascii art
func (c *Card) AsASCII() string {

	img, err := c.Image()
	if err != nil {
		log.Fatal(err)
	}

	convertOptions := convert.DefaultOptions
	// convertOptions.FixedHeight = 40
	// convertOptions.FixedWidth = 30
	converter := convert.NewImageConverter()
	return fmt.Sprint(converter.Image2ASCIIString(img, &convertOptions))
}

// ToProto returns the card as a proto
func (c *Card) ToProto() *ppb.Card {
	return &ppb.Card{
		Suite: c.Card.Suite,
		Rank:  c.Card.Rank,
	}
}

// IsLessThan returns true if c < o
func (c *Card) IsLessThan(o *Card) bool {

	switch {
	case c.GetRank() < o.GetRank():
		return true
	case c.GetRank() > o.GetRank():
		return false
	}

	return c.GetSuit() < o.GetSuit()
}

// IsSame returns true if c is the same card as o
func (c *Card) IsSame(o *Card) bool {

	if c.GetRank() == o.GetRank() && c.GetSuit() == o.GetSuit() {
		return true
	}
	return false
}

// IsSameRank returns true if c is the same rank as o
func (c *Card) IsSameRank(o *Card) bool {

	if c.GetRank() == o.GetRank() {
		return true
	}
	return false
}

// GetSuit returns the suit of the card.
func (c *Card) GetSuit() ppb.CardSuit {
	return c.Card.GetSuite()
}

// GetRank returns the rank of the card.
func (c *Card) GetRank() ppb.CardRank {
	return c.Card.GetRank()
}

// String returns ...
func (c *Card) String() string {
	return fmt.Sprintf("%v%v", rankmap[c.Card.GetRank()], suitmap[c.Card.GetSuite().String()])
}

// CardFromProto returns a *Card from ppb.Card
func CardFromProto(cp *ppb.Card) *Card {
	return &Card{
		Card: cp,
	}
}

// CardsToProto returns cards in a proto
func CardsToProto(cards []*Card) []*ppb.Card {
	pcards := make([]*ppb.Card, len(cards))

	for i, c := range cards {
		pcards[i] = c.ToProto()
	}

	return pcards
}

// CardsFromProto returns a []*Card from []*ppb.Card
func CardsFromProto(cp []*ppb.Card) []*Card {
	cards := []*Card{}

	for _, c := range cp {
		cards = append(cards, CardFromProto(c))
	}

	return cards
}

// CardsEqual compares two slices of Cards, order matters.
func CardsEqual(a, b []*Card) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.GetRank() != b[i].GetRank() {
			return false
		}
		if v.GetSuit() != b[i].GetSuit() {
			return false
		}
	}
	return true
}

// CardInList returns true if c is in l
func CardInList(c *Card, l []*Card) bool {
	for _, card := range l {
		if c.IsSame(card) {
			return true
		}
	}
	return false
}

// RankInList returns true if r is in l
func RankInList(r ppb.CardRank, l []*Card) bool {
	for _, card := range l {
		if r == card.GetRank() {
			return true
		}
	}
	return false
}

// CountByRank returns a map of rank -> number present
func CountByRank(cards []*Card) map[ppb.CardRank]int {
	byrank := map[ppb.CardRank]int{}
	for _, c := range cards {
		if _, ok := byrank[c.GetRank()]; !ok {
			byrank[c.GetRank()] = 0
		}
		byrank[c.GetRank()]++
	}
	return byrank
}

// CardsByRank returns a map of rank -> []*Card with that rank
func CardsByRank(cards []*Card) map[ppb.CardRank][]*Card {
	byrank := make(map[ppb.CardRank][]*Card, len(cards))
	for _, c := range cards {
		if _, ok := byrank[c.GetRank()]; !ok {
			byrank[c.GetRank()] = make([]*Card, 0)
		}
		byrank[c.GetRank()] = append(byrank[c.GetRank()], c)
	}
	return byrank
}

// CountBySuit returns a map of suit -> number present
func CountBySuit(cards []*Card) map[ppb.CardSuit]int {
	// count the number of each suit
	bysuit := map[ppb.CardSuit]int{
		ppb.CardSuit_Spade:   0,
		ppb.CardSuit_Club:    0,
		ppb.CardSuit_Diamond: 0,
		ppb.CardSuit_Heart:   0,
	}

	for _, card := range cards {
		bysuit[card.Card.GetSuite()]++
	}
	return bysuit
}
