package main

import (
	"fmt"
	"github.com/AhmettCelik/blog-aggregator/internal/config"
)

func startGator() {
	var cfg config.Config = config.Read()
	cfg.SetUser("Xangetsu")
	cfg = config.Read()
	fmt.Println(cfg)
}
