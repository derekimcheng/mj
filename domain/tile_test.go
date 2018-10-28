package domain

import (
	"reflect"
	"testing"
)

func Test_NewTile(t *testing.T) {
	suit := NewSuit("Bamboo", SuitTypeSimple, 10, nil)
	type args struct {
		suit    *Suit
		ordinal int
		id      int
	}
	tests := []struct {
		name    string
		args    args
		want    *Tile
		wantErr bool
	}{
		{
			"nil suit",
			args{nil, 0, 0},
			nil,
			true,
		},
		{
			"negative ordinal",
			args{suit, -1, 0},
			nil,
			true,
		},
		{
			"ordinal too big",
			args{suit, 10, 10},
			nil,
			true,
		},
		{
			"happy path 1",
			args{suit, 0, 0},
			&Tile{suit, 0, 0},
			false,
		},
		{
			"happy path 2",
			args{suit, 5, 5},
			&Tile{suit, 5, 5},
			false,
		},
		{
			"happy path 3",
			args{suit, 9, 9},
			&Tile{suit, 9, 9},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTile(tt.args.suit, tt.args.ordinal, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_StringNilFriendlyNameFunc(t *testing.T) {
	suit := NewSuit("Bamboo", SuitTypeSimple, 10, nil)
	tile, _ := NewTile(suit, 5, 0)
	if tile == nil {
		t.Fatalf("Failed to create new tile")
	}

	expected := "suit:Bamboo,ord:5,id:0"
	actual := tile.String()
	if expected != actual {
		t.Errorf("Wrong name: expected: %s, actual: %s", expected, actual)
	}
}

func Test_StringNonNilFriendlyNameFunc(t *testing.T) {
	friendlyNameFunc := func(t *Tile) string {
		return "Foo"
	}
	suit := NewSuit("Bamboo", SuitTypeSimple, 10, friendlyNameFunc)
	tile, _ := NewTile(suit, 5, 0)
	if tile == nil {
		t.Fatalf("Failed to create new tile")
	}

	expected := "Foo"
	actual := tile.String()
	if expected != actual {
		t.Errorf("Wrong name: expected: %s, actual: %s", expected, actual)
	}
}
