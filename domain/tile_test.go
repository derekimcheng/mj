package domain

import (
	"reflect"
	"testing"
)

func Test_NewTile(t *testing.T) {
	suit := NewSuit("Bamboo", SuitTypeSimple, 10)
	type args struct {
		suit    *Suit
		ordinal int
	}
	tests := []struct {
		name    string
		args    args
		want    *Tile
		wantErr bool
	}{
		{
			"nil suit",
			args{nil, 0},
			nil,
			true,
		},
		{
			"negative ordinal",
			args{suit, -1},
			nil,
			true,
		},
		{
			"ordinal too big",
			args{suit, 10},
			nil,
			true,
		},
		{
			"happy path 1",
			args{suit, 0},
			&Tile{suit, 0},
			false,
		},
		{
			"happy path 2",
			args{suit, 5},
			&Tile{suit, 5},
			false,
		},
		{
			"happy path 3",
			args{suit, 9},
			&Tile{suit, 9},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTile(tt.args.suit, tt.args.ordinal)
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
