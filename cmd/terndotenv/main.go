package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os/exec"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations",
		"./internal/store/pgstore/migrations",
		"--config",
		"./internal/store/pgstore/migrations/tern.conf",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Command execution failed with %s\n", err)
		fmt.Printf("Output: ", string(output))
		panic(err)
	}

	fmt.Printf("Command executed successfuly ", string(output))
}
