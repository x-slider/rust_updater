package main

import (
	"fmt"
	"os"

	"rust_updater/cmd/updater"
)

func main() {
	fmt.Println("Rust Server Updater - Starting...")

	if err := updater.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Update completed successfully.")
}
