package main

import (
	"fmt"
	"github.com/Miguelm-del/api-auth-jwt-go/configs"
)

func main() {
	config, _ := configs.Load(".")
	fmt.Println(config.DBDriver)
}
