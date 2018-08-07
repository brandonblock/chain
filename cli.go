package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// CLI contains the blockchain for reference via command line
type CLI struct {
	bc *Blockchain
}

const (
	addblock   = "addblock"
	printchain = "printchain"
)

// Run starts the CLI listener
func (cli *CLI) Run() {
	cli.validateArgs()

	// create commands and subcommands
	addBlockCmd := cli.newCmd(addblock)
	addBlockData := addBlockCmd.String("data", "", "Block Data")

	printChainCmd := cli.newCmd(printchain)

	// handle commands
	switch os.Args[1] {
	case addblock:
		if err := addBlockCmd.Parse(os.Args[2:]); err != nil {
			cli.printUsage(os.Args[1])
			os.Exit(1)
		}
	case printchain:
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			cli.printUsage(os.Args[1])
			os.Exit(1)
		}
	default:
		cli.printUsage(os.Args[1])
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block, err := bci.Next()
		if err != nil {
			return
		}

		fmt.Printf("Previous Has: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		p := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(p.Validate()))
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("block added")
}

func (cli *CLI) newCmd(cmd string) *flag.FlagSet {
	return flag.NewFlagSet(cmd, flag.ExitOnError)
}

func (cli *CLI) printUsage(input string) {
	if input != "" {
		fmt.Println(fmt.Sprintf("\"%s\" not understood", input))
	}
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage("")
		os.Exit(1)
	}
}
