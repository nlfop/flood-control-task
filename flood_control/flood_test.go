package flood

import (
	"context"
	"sync"
	"testing"
	"time"
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
	fm.N = 4
	fm.K = 3
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

	time.Sleep(5 * time.Second)

	{
		ch, err := fm.Check(ctx, ID1)
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}
		if ch == true {
			t.Errorf("%s wrong result, expected %#v, got %#v", "Access test (token)", false, ch)
		}

	}

	{
		wg := &sync.WaitGroup{}

		for i := 0; i < 5; i++ {

			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				fm.Check(ctx, ID3)
				wg.Done()
			}(wg)

		}
		wg.Wait()
		ch, err := fm.Check(ctx, ID3)
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}
		if ch != true {
			t.Errorf("wrong result, expected %#v, got %#v", true, ch)
		}
	}
}
