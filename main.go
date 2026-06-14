package main

import (
	"arena-vip/app"
	"os"
	"strings"

	_ "arena-vip/pages/contato"
	_ "arena-vip/pages/depoimentos"
	_ "arena-vip/pages/faq"
	_ "arena-vip/pages/horarios"
	_ "arena-vip/pages/landing"
	_ "arena-vip/pages/loja"
	_ "arena-vip/pages/privacidade"
	_ "arena-vip/pages/professores"
	_ "arena-vip/pages/resultados"
)

func main() {
	loadDotEnv()
	r := app.GetInstance().Router
	r.Static("/static", "./static")
	r.Run(":8080")
}

func loadDotEnv() {
	data, err := os.ReadFile(".env")
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		os.Setenv(strings.TrimSpace(key), strings.TrimSpace(val))
	}
}
