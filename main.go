package main

import "os"

func main() {
	bc, err := NewBlockchain()
	if err != nil {
		os.Exit(1)
	}
	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()
}
