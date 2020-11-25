package deck

import (
	"log"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"

	ppb "github.com/DanTulovsky/pepper-poker-v2/proto"
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	os.Exit(m.Run())
}

func TestCard_Sorting(t *testing.T) {
	tests := []struct {
		name  string
		cards []*Card
		want  []*Card
	}{
		{
			cards: []*Card{
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Two),
				NewCard(ppb.CardSuit_Club, ppb.CardRank_King),
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Ten),
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Jack),
			},
			want: []*Card{
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Two),
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Ten),
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Jack),
				NewCard(ppb.CardSuit_Club, ppb.CardRank_King),
			},
		},
		{
			cards: []*Card{
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Two),
				NewCard(ppb.CardSuit_Heart, ppb.CardRank_King),
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Ten),
				NewCard(ppb.CardSuit_Spade, ppb.CardRank_Jack),
			},
			want: []*Card{
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Two),
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Ten),
				NewCard(ppb.CardSuit_Spade, ppb.CardRank_Jack),
				NewCard(ppb.CardSuit_Heart, ppb.CardRank_King),
			},
		},
		{
			cards: []*Card{
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Two),
				NewCard(ppb.CardSuit_Spade, ppb.CardRank_Jack),
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Ten),
				NewCard(ppb.CardSuit_Heart, ppb.CardRank_Two),
			},
			want: []*Card{
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Two),
				NewCard(ppb.CardSuit_Heart, ppb.CardRank_Two),
				NewCard(ppb.CardSuit_Club, ppb.CardRank_Ten),
				NewCard(ppb.CardSuit_Spade, ppb.CardRank_Jack),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(SortByCards(tt.cards))
			if !CardsEqual(tt.cards, tt.want) {
				t.Errorf("lists not equal:\n[have]> %v\n[want]> %v", tt.cards, tt.want)
			}
		})
	}
}
