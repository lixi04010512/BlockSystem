package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"lxblockchain/utils"
	"math/big"

	"github.com/btcsuite/btcd/btcutil/base58"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

// 新建账户
func NewWallet() *Wallet {

	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey
	i := sha256.New()
	i.Write(w.publicKey.X.Bytes())
	i.Write(w.publicKey.Y.Bytes())
	digest := i.Sum(nil)
	address := base58.Encode(digest)
	w.blockchainAddress = address

	return w
}

// 为什么要写以下返回私钥和公钥的方法
func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

func (t *Transaction) Getvalue() *big.Int {
	return t.value
}

// 通过已有私钥加载钱包
func LoadWallet(privkey string) *Wallet {
	privateKey := privkey
	privateKeyInt := new(big.Int)
	privateKeyInt.SetString(privateKey, 16)
	// 曲线
	curve := elliptic.P256()
	// 获取公钥
	x, y := curve.ScalarBaseMult(privateKeyInt.Bytes())
	publicKey := ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
	n := new(Wallet)
	w := sha256.New()
	w.Write(publicKey.X.Bytes())
	w.Write(publicKey.Y.Bytes())
	digest := w.Sum(nil)
	address := base58.Encode(digest)
	n.blockchainAddress = address
	n.publicKey = &publicKey
	n.privateKey = &ecdsa.PrivateKey{
		PublicKey: publicKey,
		D:         privateKeyInt,
	}
	return n
}

type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      *big.Int
}

func (wallet *Wallet) Transfer(recipientBlockchainAddress string, value *big.Int) *Transaction {
	return &Transaction{wallet.privateKey, wallet.publicKey, wallet.blockchainAddress, recipientBlockchainAddress, value}
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"private_key"`
		PublicKey         string `json:"public_key"`
		BlockchainAddress string `json:"blockchain_address"`
	}{
		PrivateKey:        w.PrivateKeyStr(),
		PublicKey:         w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),
	})
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	i, _ := json.Marshal(t)
	a := sha256.Sum256([]byte(i))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, a[:])
	return &utils.Signature{R: r, S: s}
}
