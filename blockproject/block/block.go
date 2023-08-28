package block

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/fatih/color"

	_ "github.com/go-sql-driver/mysql"
)

var MINING_DIFFICULT = 0x3333

const MINING_ACCOUNT_ADDRESS = "lixi"

var MINING_REWARD, _ = new(big.Int).SetString("50000000000000000000", 10)

type Block struct {
	Nonce        *big.Int
	Timestamp    uint64
	Number       *big.Int
	Difficulty   *big.Int
	PreviousHash [32]byte
	Hash         [32]byte
	TxSize       uint16
	Transactions []*Transaction
}

func NewBlock(nonce *big.Int, number *big.Int, previousHash [32]byte, txs []*Transaction, difficulty *big.Int) *Block {
	b := new(Block)
	b.Timestamp = uint64(time.Now().UnixNano())
	b.Nonce = nonce
	b.Number = number
	b.Difficulty = difficulty
	b.PreviousHash = previousHash
	b.Hash = b.Hash1()
	b.Transactions = txs
	b.TxSize = uint16(len(b.Transactions))
	return b
}

func (b *Block) Print() {
	log.Printf("%-15v:%30d\n", "nonce", b.Nonce)
	log.Printf("%-15v:%30d\n", "number", b.Number)
	log.Printf("%-15v:%30d\n", "txSize", b.TxSize)
	log.Printf("%-15v:%30x\n", "previous_hash", b.PreviousHash)
	log.Printf("%-15v:%30x\n", "difficulty", b.Difficulty)
	log.Printf("%-15v:%30x\n", "hash", b.Hash)
	for _, i := range b.Transactions {
		i.Print()
	}
}

type Blockchain struct {
	Block           []*Block       //区块
	Coinbase        string         //区块奖励地址
	TransactionPool []*Transaction //将新的交易加入交易池、从交易池中删除已经被确认的交易等
}

// 新建一条链的第一个区块
// NewBlockchain(blockchainAddress string) *Blockchain
// 函数定义了一个创建区块链的方法，它接收一个字符串类型的参数 blockchainAddress，
// 它返回一个区块链类型的指针。在函数内部，它创建一个区块链对象并为其设置地址，
// 然后创建一个创世块并将其添加到区块链中，最后返回区块链对象。
func NewBlockchain(blockchainAddress string) *Blockchain {
	i := new(Blockchain)
	b := &Block{}
	i.CreateBlock(big.NewInt(0), big.NewInt(0), b.Hash1(), big.NewInt(0)) //创世纪块
	i.Coinbase = blockchainAddress
	return i
}

// (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block
//  函数是在区块链上创建新的区块，它接收两个参数：一个int类型的nonce和一个字节数组类型的 previousHash，
//  返回一个区块类型的指针。在函数内部，它使用传入的参数来创建一个新的区块，
//  然后将该区块添加到区块链的链上，并清空交易池。

func (bc *Blockchain) CreateBlock(nonce *big.Int, number *big.Int, previousHash [32]byte, difficulty *big.Int) *Block {
	b := NewBlock(nonce, number, previousHash, bc.TransactionPool, difficulty)
	bc.Block = append(bc.Block, b)
	bc.TransactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) Print() {
	for i, b := range bc.Block {
		color.Green("%s BLOCK %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		b.Print()
	}
	color.Yellow("%s\n\n\n", strings.Repeat("*", 50))
}

func (b *Block) Hash1() [32]byte {
	i, _ := json.Marshal(b)
	return sha256.Sum256([]byte(i))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp  uint64   `json:"timestamp"`
		Nonce      *big.Int `json:"nonce"`
		Number     *big.Int `json:"number"`
		Difficulty *big.Int `json:"difficulty"`
	}{
		Timestamp:  b.Timestamp,
		Nonce:      b.Nonce,
		Number:     b.Number,
		Difficulty: b.Difficulty,
	})
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.Block[len(bc.Block)-1]
}

func (bc *Blockchain) OneBlock(a int) *Block {
	return bc.Block[a]
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.TransactionPool {
		transactions = append(transactions,
			NewTransaction(t.SenderAddress,
				t.ReceiveAddress,
				t.Value))
	}
	return transactions
}

func byteToBigInt(b [32]byte) *big.Int {
	bytes := b[:]
	result := new(big.Int).SetBytes(bytes)
	return result
}

func (bc *Blockchain) ValidProof(nonce *big.Int,
	previousHash [32]byte,
	transactions []*Transaction,
	difficulty int,
) bool {
	bigint_difficulty := big.NewInt(int64(difficulty))
	target := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	target = new(big.Int).Div(target, bigint_difficulty)
	tmpBlock := Block{Nonce: nonce, PreviousHash: previousHash, Transactions: transactions, Timestamp: 0}
	result := byteToBigInt(tmpBlock.Hash1())

	return target.Cmp(result) > 0
}

func (bc *Blockchain) GetSpendTime(blocknum int) uint64 {
	if blocknum == 0 {
		return 0
	}
	return uint64(bc.Block[blocknum].Timestamp - bc.Block[blocknum-1].Timestamp)
}

func (bc *Blockchain) ProofOfWork() (*big.Int, int) {
	transaction := bc.CopyTransactionPool() //选择交易？控制交易数量？
	previousHash := bc.LastBlock().Hash1()
	nonce := big.NewInt(0)
	t := time.Now()
	if bc.GetSpendTime(len(bc.Block)-1) < 2e+8 {
		MINING_DIFFICULT += 5000
	} else {
		if MINING_DIFFICULT >= 100000 {
			MINING_DIFFICULT -= 5000
		}
	}
	for !bc.ValidProof(nonce, previousHash, transaction, MINING_DIFFICULT) {
		nonce.Add(nonce, big.NewInt(1))
	}
	end := time.Now()

	log.Printf("POW spend Time:%s DIFF%d", end.Sub(t), MINING_DIFFICULT)

	return nonce, MINING_DIFFICULT
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_ACCOUNT_ADDRESS, bc.Coinbase, MINING_REWARD, nil, nil)
	nonce, difficulty := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash
	number := big.NewInt(0)
	number.Add(bc.LastBlock().Number, big.NewInt(1))
	bc.CreateBlock(nonce, number, previousHash, big.NewInt(int64(difficulty)))
	color.Red("action=mining, status=success")
	return true
}

func (bc *Blockchain) GetBalance(accountAddress string) *big.Int {
	big := big.NewInt(0)
	for _, bcs := range bc.Block {
		for _, tra := range bcs.Transactions {
			if accountAddress == tra.ReceiveAddress {
				big = big.Add(big, tra.Value)
			}
			if accountAddress == tra.SenderAddress {
				big = big.Sub(big, tra.Value)
			}
		}
	}
	return big
}

func (blockchain *Blockchain) GetBlockByNumber(blockid uint64) {
	for i, block := range blockchain.Block {
		if big.NewInt(int64(i)).Cmp(big.NewInt(int64(blockid))) == 0 {
			fmt.Println("=====根据number查block==========")
			log.Printf("%-15v:%30d\n", "nonce", block.Nonce)
			log.Printf("%-15v:%30d\n", "timestamp", block.Timestamp)
			log.Printf("%-15v:%30d\n", "number", block.Number)
			log.Printf("%-15v:%30d\n", "difficulty", block.Difficulty)
			log.Printf("%-15v:%30x\n", "previousHash", block.PreviousHash)
			log.Printf("%-15v:%30x\n", "hash", block.Hash)
			log.Printf("%-15v:%30d\n", "txSize", block.TxSize)
		}
	}

}

func (blockchain *Blockchain) GetBlockByHash(hash []byte) {
	var array32 [32]byte
	copy(array32[:], hash)
	for _, block := range blockchain.Block {
		if block.Hash == array32 {
			fmt.Println("=====根据hash查block==========")
			log.Printf("%-15v:%30d\n", "nonce", block.Nonce)
			log.Printf("%-15v:%30d\n", "timestamp", block.Timestamp)
			log.Printf("%-15v:%30d\n", "number", block.Number)
			log.Printf("%-15v:%30d\n", "difficulty", block.Difficulty)
			log.Printf("%-15v:%30x\n", "previousHash", block.PreviousHash)
			log.Printf("%-15v:%30x\n", "hash", block.Hash)
			log.Printf("%-15v:%30d\n", "txSize", block.TxSize)
		}
	}

}

// 加载保存的区块链数据
func LoadBlockchain(blockchain *Blockchain) error {

	db, err := sql.Open("mysql", "root:root123456@tcp(localhost:3306)/block")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建一个新的 Block 对象
	block := &Block{}

	// 查询数据库中的 Blockchain 数据
	rows, err := db.Query("SELECT * FROM blockchain")
	if err != nil {
		return err
	}
	defer rows.Close()

	// 声明变量用于存储从数据库中扫描的值
	var id int64
	var nonceInt64 int64
	var timestamp uint64
	var numberInt64 int64
	var difficultyInt64 int64
	var previousHash []byte
	var hash []byte
	var txSize uint16
	var coinbase string

	// 遍历查询结果
	for rows.Next() {

		// 扫描数据库中的值到相应的变量
		err := rows.Scan(
			&id,
			&nonceInt64,
			&timestamp,
			&numberInt64,
			&difficultyInt64,
			&previousHash,
			&hash,
			&txSize,
			&coinbase,
		)
		if err != nil {
			return err
		}

		// 将扫描的值进行类型转换
		nonce := big.NewInt(nonceInt64)
		number := big.NewInt(numberInt64)
		difficulty := big.NewInt(difficultyInt64)

		// 将转换后的值分配给 Block 结构体中的字段
		block.Nonce = nonce
		block.Timestamp = timestamp
		block.Number = number
		block.Difficulty = difficulty
		block.PreviousHash = convertToByteArray(previousHash)
		block.Hash = convertToByteArray(hash)
		block.TxSize = txSize
		blockchain.Coinbase = coinbase

		// 将加载的 Block 添加到 Blockchain 中
		blockchain.Block = append(blockchain.Block, block)
	}

	// 查询数据库中的 transactions 数据
	rows1, err := db.Query("SELECT * FROM transactions where block_id=?", id)
	if err != nil {
		return err
	}
	defer rows1.Close()

	// 声明变量用于存储从数据库中扫描的值
	var id1 int64
	var block_id int64
	var sender string
	var recipient string
	var amount string

	for rows1.Next() {
		err := rows1.Scan(
			&id1,
			&block_id,
			&sender,
			&recipient,
			&amount,
		)
		if err != nil {
			return err
		}

		bigInt := new(big.Int)
		_, success := bigInt.SetString(amount, 10)
		if !success {
			fmt.Println("Failed to parse string as big integer")
			return nil
		}
		transaction := &Transaction{}
		transaction.SenderAddress = sender
		transaction.ReceiveAddress = recipient
		transaction.Value = bigInt
		block.Transactions = append(block.Transactions, transaction)
	}

	return nil
}

func convertToByteArray(bytes []byte) [32]byte {
	var byteArray [32]byte

	if len(bytes) >= len(byteArray) {
		copy(byteArray[:], bytes[:32])
	} else {
		copy(byteArray[:], bytes)
	}

	return byteArray
}

func (b *Blockchain) UnmarshalJSON(data []byte) error {
	// 在这里实现自定义的 JSON 反序列化逻辑

	var customBlockchain struct {
		Block           []*Block       `json:"Block"`
		Coinbase        string         `json:"Coinbase"`
		TransactionPool []*Transaction `json:"TransactionPool"`
	}

	jsonData := string(data)
	err := json.Unmarshal([]byte(jsonData), &customBlockchain)
	if err != nil {
		return err
	}

	// 将解析得到的值赋给 MyStruct 的字段
	b.Block = customBlockchain.Block
	b.Coinbase = customBlockchain.Coinbase
	b.TransactionPool = customBlockchain.TransactionPool
	fmt.Println("bbb", b.Block)

	return nil
}

// 保存区块链数据
func SaveBlockchain(blockchain *Blockchain) error {
	db, err := sql.Open("mysql", "root:root123456@tcp(localhost:3306)/block")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 插入语句
	query := `
        INSERT INTO blockchain (nonce, timestamp, number, difficulty, previous_hash, hash, tx_size, coinbase)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `

	// 循环插入每个区块
	for _, block := range blockchain.Block {
		result, err := db.Exec(
			query,
			block.Nonce.String(),
			block.Timestamp,
			block.Number.String(),
			block.Difficulty.String(),
			block.PreviousHash[:],
			block.Hash[:],
			block.TxSize,
			blockchain.Coinbase,
		)
		lastInsertID, err := result.LastInsertId()
		for _, tx := range block.Transactions {
			_, err := db.Exec("INSERT INTO transactions (block_id,sender, recipient, amount) VALUES ( ?,?, ?, ?)",
				lastInsertID, tx.SenderAddress, tx.ReceiveAddress, tx.Value.String())
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}

	return nil

}
