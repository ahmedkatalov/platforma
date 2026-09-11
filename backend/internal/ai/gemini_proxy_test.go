package ai

import (
	"bufio"
	"context"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// TestAskRoutesThroughProxy проверяет, что Ask с заданным прокси действительно
// идёт через него: для https-цели клиент шлёт CONNECT на прокси. Поднимаем
// фейковый прокси, который ловит первую строку запроса.
func TestAskRoutesThroughProxy(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	var got atomic.Value // string — первая строка, полученная прокси
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		line, _ := bufio.NewReader(conn).ReadString('\n')
		got.Store(line)
		// Отвечаем ошибкой — сам ответ Gemini нам не нужен.
		_, _ = conn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
	}()

	c := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	proxy := "http://" + ln.Addr().String()
	_, _ = c.Ask(ctx, "fake-key", "gemini-2.5-flash", proxy, "sys",
		[]Message{{Role: "user", Text: "hi"}})

	v, _ := got.Load().(string)
	if !strings.HasPrefix(v, "CONNECT generativelanguage.googleapis.com:443") {
		t.Fatalf("прокси не получил ожидаемый CONNECT, got=%q", v)
	}
}

// TestAskBadProxyErrors — неверный адрес прокси даёт понятную ошибку, а не панику.
func TestAskBadProxyErrors(t *testing.T) {
	c := NewClient()
	_, err := c.Ask(context.Background(), "k", "m", "://bad", "s",
		[]Message{{Role: "user", Text: "x"}})
	if err == nil {
		t.Fatal("ожидали ошибку на неверном прокси")
	}
}
