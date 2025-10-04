
package blockchain

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "time"
)

// Property represents a real property asset
type Property struct {
    ID          string `json:"id"`          // Unique property identifier (e.g., parcel number)
    Owner       string `json:"owner"`       // Current owner (e.g., wallet address or name)
    Description string `json:"description"` // Description of the property
    Location    string `json:"location"`    // Location details
    Value       int64  `json:"value"`       // Estimated value in some unit (e.g., USD)
}

// Transaction represents a transfer or update of a property
type Transaction struct {
    From      string   `json:"from"`      // Sender (previous owner)
    To        string   `json:"to"`        // Receiver (new owner)
    Property  Property `json:"property"`  // The property being transferred/updated
    Timestamp int64    `json:"timestamp"` // Unix timestamp
}

// Block represents a block in the blockchain
type Block struct {
    Index        int           `json:"index"`
    Timestamp    int64         `json:"timestamp"`
    Transactions []Transaction `json:"transactions"`
    PrevHash     string        `json:"prev_hash"`
    Hash         string        `json:"hash"`
    Nonce        int           `json:"nonce"` // For proof of work
}

// Blockchain represents the chain of blocks and property registry
type Blockchain struct {
    Chain     []Block              `json:"chain"`
    PendingTx []Transaction        `json:"pending_tx"` // Pending transactions
    Registry  map[string]Property  `json:"registry"`   // Property ID to Property mapping for quick lookup
}

// NewBlockchain creates a new blockchain with a genesis block
func NewBlockchain() *Blockchain {
    genesisBlock := createGenesisBlock()
    return &Blockchain{
        Chain:     []Block{genesisBlock},
        PendingTx: []Transaction{},
        Registry:  make(map[string]Property),
    }
}

// createGenesisBlock creates the first block
func createGenesisBlock() Block {
    genesisTx := Transaction{
        From:      "genesis",
        To:        "genesis",
        Property:  Property{ID: "genesis_property", Owner: "genesis", Description: "Genesis Property", Location: "N/A", Value: 0},
        Timestamp: time.Now().Unix(),
    }
    block := Block{
        Index:        0,
        Timestamp:    time.Now().Unix(),
        Transactions: []Transaction{genesisTx},
        PrevHash:     "0",
        Hash:         "",
        Nonce:        0,
    }
    block.Hash = block.calculateHash()
    return block
}

// calculateHash computes the hash of the block
func (b *Block) calculateHash() string {
    data := fmt.Sprintf("%d%d%s%s%d", b.Index, b.Timestamp, b.PrevHash, b.transactionsToString(), b.Nonce)
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}

// transactionsToString converts transactions to a string for hashing
func (b *Block) transactionsToString() string {
    txBytes, _ := json.Marshal(b.Transactions)
    return string(txBytes)
}

// MineBlock performs simple proof of work
func (b *Block) MineBlock(difficulty int) {
    target := string(make([]byte, difficulty))
    for i := range target {
        target = target[:i] + "0" + target[i+1:]
    }
    for b.Hash[:difficulty] != target {
        b.Nonce++
        b.Hash = b.calculateHash()
    }
}

// AddTransaction adds a transaction to pending list
func (bc *Blockchain) AddTransaction(tx Transaction) error {
    // Validate transaction
    if tx.Property.ID == "" {
        return errors.New("invalid property ID")
    }
    if tx.From == "" || tx.To == "" {
        return errors.New("invalid from/to addresses")
    }
    // Check if property exists and ownership
    if prop, exists := bc.Registry[tx.Property.ID]; exists {
        if prop.Owner != tx.From {
            return errors.New("sender does not own the property")
        }
    } else {
        // New property, ensure from is empty or creator
        if tx.From != "" {
            return errors.New("cannot transfer non-existent property")
        }
        tx.From = "creator" // For new properties
    }
    bc.PendingTx = append(bc.PendingTx, tx)
    return nil
}

// MinePendingTransactions mines a new block with pending tx
func (bc *Blockchain) MinePendingTransactions(miner string, difficulty int) {
    if len(bc.PendingTx) == 0 {
        return // Nothing to mine
    }
    lastBlock := bc.Chain[len(bc.Chain)-1]
    newBlock := Block{
        Index:        lastBlock.Index + 1,
        Timestamp:    time.Now().Unix(),
        Transactions: bc.PendingTx,
        PrevHash:     lastBlock.Hash,
    }
    newBlock.MineBlock(difficulty)

    // Update registry
    for _, tx := range bc.PendingTx {
        prop := tx.Property
        prop.Owner = tx.To
        bc.Registry[tx.Property.ID] = prop
    }

    bc.Chain = append(bc.Chain, newBlock)
    bc.PendingTx = []Transaction{}
}

// GetProperty retrieves a property from registry
func (bc *Blockchain) GetProperty(id string) (Property, error) {
    if prop, exists := bc.Registry[id]; exists {
        return prop, nil
    }
    return Property{}, errors.New("property not found")
}

// ValidateChain checks the integrity of the blockchain
func (bc *Blockchain) ValidateChain() bool {
    for i := 1; i < len(bc.Chain); i++ {
        curr := bc.Chain[i]
        prev := bc.Chain[i-1]
        if curr.PrevHash != prev.Hash {
            return false
        }
        if curr.Hash != curr.calculateHash() {
            return false
        }
    }
    return true
}