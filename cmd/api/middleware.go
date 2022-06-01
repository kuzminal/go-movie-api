package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)
	// отдельная функция очистки списка клиентов сервиса
	// выполняется отдельно от потока лимитера
	// чтобы избежать ошибок устанавливаем блокировку
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// если лимитер включен, то выполняем проверки ограничений
		if app.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr) // получаем IP-адрес пользователя из запроса
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			mu.Lock()

			if _, found := clients[ip]; !found {
				// инициализация нового лимитера с максимальным кол-вом запросов "за один раз" = 4? в секунду = 2,
				//если по IP не нашли в мапе уже инициализированный лимитер
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}
			clients[ip].lastSeen = time.Now()
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}
