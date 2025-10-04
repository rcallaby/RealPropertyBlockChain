// api.go
package blockchain

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"
)

// StartAPI starts the HTTP server for the blockchain API
func (bc *Blockchain) StartAPI(port int) {
    http.HandleFunc("/chain", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(bc.Chain)
    })

    http.HandleFunc("/property/", func(w http.ResponseWriter, r *http.Request) {
        parts := strings.Split(r.URL.Path, "/")
        if len(parts) != 3 {
            http.Error(w, "Invalid path", http.StatusBadRequest)
            return
        }
        id := parts[2]
        prop, err := bc.GetProperty(id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        json.NewEncoder(w).Encode(prop)
    })

    http.HandleFunc("/transaction", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        var tx Transaction
        if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        if err := bc.AddTransaction(tx); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        w.WriteHeader(http.StatusCreated)
        fmt.Fprintln(w, "Transaction added")
    })

    http.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {
        miner := r.URL.Query().Get("miner")
        if miner == "" {
            miner = "anonymous"
        }
        difficultyStr := r.URL.Query().Get("difficulty")
        difficulty, _ := strconv.Atoi(difficultyStr)
        if difficulty == 0 {
            difficulty = 4 // Default difficulty
        }
        bc.MinePendingTransactions(miner, difficulty)
        fmt.Fprintln(w, "Block mined")
    })

    http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
        valid := bc.ValidateChain()
        json.NewEncoder(w).Encode(map[string]bool{"valid": valid})
    })

    fmt.Printf("Starting API on port %d\n", port)
    http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}