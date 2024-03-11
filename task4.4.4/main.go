package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"runtime"
	"time"
)

func monitorGoroutines(prevGoroutines int) {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		currentGoroutines := runtime.NumGoroutine()
		fmt.Printf("Текущее количество горутин: %d\n", currentGoroutines)

		if float64(currentGoroutines) > float64(prevGoroutines)*1.2 {
			fmt.Println("⚠️ Предупреждение: Количество горутин увеличилось более чем на 20%!")
		} else if float64(currentGoroutines) < float64(prevGoroutines)*0.8 {
			fmt.Println("⚠️ Предупреждение: Количество горутин уменьшилось более чем на 20%!")
		}

		prevGoroutines = currentGoroutines
	}
}

func main() {
	g, _ := errgroup.WithContext(context.Background())

	// Мониторинг горутин
	go func() {
		monitorGoroutines(runtime.NumGoroutine())
	}()

	// Имитация активной работы приложения с созданием горутин
	for i := 0; i < 64; i++ {
		g.Go(func() error {
			time.Sleep(5 * time.Second)
			return nil
		})
		time.Sleep(80 * time.Millisecond)
	}

	// Ожидание завершения всех горутин
	if err := g.Wait(); err != nil {
		fmt.Println("Ошибка:", err)
	}
}
