// main.go
package main

import (
    "blockchain"
    "fmt"
    "time"
)

func main() {
    bc := blockchain.NewBlockchain()

    // Example: Add a new property
    tx1 := blockchain.Transaction{
        From:      "", // New property
        To:        "owner1",
        Property:  blockchain.Property{ID: "prop1", Description: "House in City A", Location: "123 Street", Value: 500000},
        Timestamp: time.Now().Unix(),
    }
    bc.AddTransaction(tx1)

    // Mine
    bc.MinePendingTransactions("miner1", 4)

    // Transfer property
    tx2 := blockchain.Transaction{
        From:      "owner1",
        To:        "owner2",
        Property:  blockchain.Property{ID: "prop1"}, // Value and desc can be updated if needed
        Timestamp: time.Now().Unix(),
    }
    bc.AddTransaction(tx2)

    // Mine again
    bc.MinePendingTransactions("miner2", 4)

    // Print chain
    fmt.Println("Blockchain:")
    for _, block := range bc.Chain {
        fmt.Printf("Block %d: Hash %s, Prev %s, Tx Count %d\n", block.Index, block.Hash, block.PrevHash, len(block.Transactions))
    }

    // Validate
    fmt.Printf("Chain valid: %v\n", bc.ValidateChain())

    // Get property
    prop, _ := bc.GetProperty("prop1")
    fmt.Printf("Property prop1 owner: %s\n", prop.Owner)

    // Start API in a goroutine or separately
    /bc.StartAPI(8080) // Uncomment to run the API
}