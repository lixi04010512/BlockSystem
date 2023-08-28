package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"lxblockchain/utils"
	"math/big"
	"strings"

	"github.com/fatih/color"
)

type Transaction struct {
	SenderAddress  string
	ReceiveAddress string
	Value          *big.Int
}

func NewTransaction(sender string, receive string, value *big.Int) *Transaction {
	t := Transaction{sender, receive, value}
	return &t
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value *big.Int, senderPublicKey *ecdsa.PublicKey, i *utils.Signature) bool {
	n := NewTransaction(sender, recipient, value)

	if sender == MINING_ACCOUNT_ADDRESS {
		bc.TransactionPool = append(bc.TransactionPool, n)
		bc.LastBlock().TxSize = uint16(len(bc.TransactionPool))
		fmt.Println("pool", bc.TransactionPool)
		return true
	}

	if bc.GetBalance(sender).Cmp(value) == -1 {
		log.Printf("ERROR: %s,没有足够的钱", sender)
		return false
	}

	if bc.VerifyTransactionSignature(senderPublicKey, i, n) {
		bc.TransactionPool = append(bc.TransactionPool, n)
		bc.LastBlock().TxSize = uint16(len(bc.TransactionPool))
		fmt.Println("pool1", bc.TransactionPool)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
	}
	return true

}

func (bc *Blockchain) VerifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey, u *utils.Signature, t *Transaction) bool {
	i, _ := json.Marshal(t)
	b := sha256.Sum256([]byte(i))
	return ecdsa.Verify(senderPublicKey, b[:], u.R, u.S)
}

var Tra []Transaction

func (t *Transaction) Print() {
	color.Red("%s\n", strings.Repeat("~", 30))
	color.Cyan("发送地址             %s\n", t.SenderAddress)
	color.Cyan("接受地址             %s\n", t.ReceiveAddress)
	color.Cyan("金额                 %d\n", t.Value)

	tr := Transaction{
		SenderAddress:  t.SenderAddress,
		ReceiveAddress: t.ReceiveAddress,
		Value:          t.Value,
	}

	Tra = append(Tra, tr)

}

func TradeHistory() []map[string]interface{} {
	fmt.Println("ii", Tra)
	// 创建一个切片用于存储每个数据项的map
	data := make([]map[string]interface{}, len(Tra))

	// 遍历每个Transaction结构体实例，将数据存储为key-value对的形式
	for i, t := range Tra {
		item := make(map[string]interface{})
		item["senderAddress"] = t.SenderAddress
		item["receiveAddress"] = t.ReceiveAddress
		item["value"] = t.Value
		data[i] = item
	}
	return data
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string   `json:"sender_blockchain_address"`
		Recipient string   `json:"recipient_blockchain_address"`
		Value     *big.Int `json:"value"`
	}{
		Sender:    t.SenderAddress,
		Recipient: t.ReceiveAddress,
		Value:     t.Value,
	})
}

func (bc *Blockchain) GetTransactionByHash(hash []byte) {
	var array32 [32]byte
	copy(array32[:], hash)

	for i, block := range bc.Block {
		if block.Hash == array32 {
			fmt.Println("=====区块号为：", i, " 的交易为：=========")
			for _, j := range block.Transactions {
				j.Print()
			}
		}
	}

}
