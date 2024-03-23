package main

import (
	"context"
	"fmt"
	flood "task/flood_control"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	_, err := flood.InitDBFlood(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer cancel()

}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
