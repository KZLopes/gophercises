package main

import (
	"fmt"
	"main/internal/secret"
)

func main() {
	v := secret.NewMemoryVault("mocked-encoding-key")
	err := v.Set("demo_key", "some value")
	if err != nil {
		panic(err)
	}

	plain, err := v.Get("demo_key")
	if err != nil {
		panic(err)
	}

	fmt.Println("Plain:", plain)
}
