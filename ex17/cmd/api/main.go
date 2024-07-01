package main

import (
	"fmt"
	"main/internal/secret"
)

func main() {
	vaultM := secret.NewMemoryVault("mocked-encoding-key")
	err := vaultM.Set("demo_key", "some value")
	if err != nil {
		panic(err)
	}

	plain, err := vaultM.Get("demo_key")
	if err != nil {
		panic(err)
	}

	fmt.Println("Plain1:", plain)

	vaultF := secret.NewFileVault("my-encoding key", "./.secrets")
	err = vaultF.Set("demo_key", "some file values")
	if err != nil {
		panic(err)
	}

	plain, err = vaultF.Get("demo_key")
	if err != nil {
		panic(err)
	}

	fmt.Println("Plain2:", plain)

}
