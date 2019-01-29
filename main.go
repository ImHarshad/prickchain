package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// =========================== Blockchain code from here ===========================

// Block represents each item in the blockchain
type Block struct {
	Index int
	Timestamp string
	Data string
	Hash string
	PrevHash string
}

// Blockchain is a series of validated blocks - array of Blocks will form chain of blocks
var Blockchain []Block

var mutex = &sync.Mutex{}

// ====================== MAIN ======================

func main() {
	err := godotenv.Load("prickchain.env")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{0, t.String(), "Harsh's PRIvate bloCKCHAIN", calculateHash(genesisBlock), ""}
		spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())

}


// SHA256 hashing
func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.Data + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, Data string) Block {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Data = Data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
}

// Make sure block is valid by checking index, and cross-verifying hash of previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index + 1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// Compare the lengths of two chains from two different nodes and select the longest and updated one - TRY THIS SITUATION
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

// =========================== Blockchain code till here ===========================

// =========================== Server code from here ===========================

// Web server
func run() error {
	mux := makeMuxRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port:", httpPort)
	s := &http.Server {
		Addr: ":" + httpPort,
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create handlers - Handels GET and POST methods for your blockchain
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlockchain).Methods("POST")
	return muxRouter
}

// Write blockchain when we receive a http request
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

// Message takes incoming JSON payload for writting Data
type Message struct {
	Data string
}

// Takes JSON payload as an input for Data
func handleWriteBlockchain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var msg Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}

	defer r.Body.Close()

	mutex.Lock()
	prevBlock := Blockchain[len(Blockchain)-1]
	newBlock := generateBlock(prevBlock, msg.Data)

	if isBlockValid(newBlock, prevBlock) {
		Blockchain = append(Blockchain, newBlock)
		spew.Dump(Blockchain)
	}

	mutex.Unlock()

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}


func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

// =========================== Server code till here ===========================
