package deck

import (
	"testing"
)

func TestNewDeck(t *testing.T) {
	got := NewDeck()

	if len(got.cards) != 52 {
		t.Errorf("NewDeck() returned deck of length %v; expected: 52", len(got.cards))
	}
}

func TestDeck_Remove(t *testing.T) {
	tests := []struct {
		name    string
		card    Card
		wantErr bool
	}{
		{
			card:    NewRandomCard(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDeck()
			if err := d.Remove(tt.card); (err != nil) != tt.wantErr {
				t.Errorf("Deck.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if CardInList(tt.card, d.cards) {
					t.Errorf("Deck.Remove(%v) failed to remove card", tt.card)
				}
			}
		})
	}
}
