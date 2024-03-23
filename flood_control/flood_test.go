package flood

import (
	"context"
	"testing"
)

var (
	ID1 int64 = 1234
	ID2 int64 = 34
	ID3 int64 = 223338
)

func TestFloodP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	fm, err := InitDBFlood(ctx)
	if err != nil {
		t.Errorf("%s unexpected error: %#v", "DB err", err)
	}
	defer cancel()

	{
		ch, err := fm.Check(ctx, ID1)
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}
		if ch == true {
			t.Errorf("%s wrong result, expected %#v, got %#v", "Access test (token)", false, ch)
		}
		fm.Check(ctx, ID1)
		ch, err = fm.Check(ctx, ID2)
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}
		if ch == true {
			t.Errorf("%s wrong result, expected %#v, got %#v", "Access test (token)", false, ch)
		}
		fm.Check(ctx, ID1)
		fm.Check(ctx, ID1)
		fm.Check(ctx, ID1)
		fm.Check(ctx, ID1)
		ch, err = fm.Check(ctx, ID1)
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}
		if ch != true {
			t.Errorf("wrong result, expected %#v, got %#v", true, ch)
		}

	}

}

func TestFloodN(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	_, err := InitDBFlood(ctx)
	if err != nil {
		t.Errorf("%s unexpected error: %#v", "DB err", err)
	}
	defer cancel()

}
