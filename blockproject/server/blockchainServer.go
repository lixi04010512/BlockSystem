package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"log"
	"lxblockchain/block"
	"lxblockchain/wallet"
	"math/big"
	"net/http"
	"strconv"
)

var cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)
var Miners *wallet.Wallet
var Blockchains *block.Blockchain

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		// NewBlockchain与以前的方法不一样,增加了地址和端口2个参数,是为了区别不同的节点
		bc = block.NewBlockchain(minersWallet.BlockchainAddress())
		cache["blockchain"] = bc
		color.Magenta("===矿工帐号信息====\n")
		color.Magenta("矿工private_key\n %v\n", minersWallet.PrivateKeyStr())
		color.Magenta("矿工publick_key\n %v\n", minersWallet.PublicKeyStr())
		color.Magenta("矿工blockchain_address\n %s\n", minersWallet.BlockchainAddress())
		color.Magenta("===============\n")
	}
	return bc
}

func (bcs *BlockchainServer) Chain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		miner := wallet.LoadWallet(block.MINING_ACCOUNT_ADDRESS)
		Miners = miner
		minerAddress := miner.BlockchainAddress()
		blockchain := block.NewBlockchain(minerAddress)
		Blockchains = blockchain
		Blockchains.Mining()
		w.Write([]byte("success！"))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Mine(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		Blockchains.Mining()
		w.Write([]byte("success！"))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) GetBlockByNumber(w http.ResponseWriter, req *http.Request) {
	blockNumber := req.URL.Query().Get("blockNumber")

	num, err := strconv.ParseUint(blockNumber, 10, 64)
	if err != nil {
		fmt.Print(err)
	}
	switch req.Method {
	case http.MethodGet:
		color.Cyan("####################GetBlockByNumber###############")
		Blockchains.GetBlockByNumber(num)
		w.Write([]byte("success！"))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) GetBlockByHash(w http.ResponseWriter, req *http.Request) {
	blockHash := req.URL.Query().Get("blockHash")

	switch req.Method {
	case http.MethodGet:
		color.Cyan("####################GetBlockByHash###############")
		for i := 0; i < 20; i++ {
			bc := Blockchains.OneBlock(i).Hash1()
			hashString := hex.EncodeToString(bc[:])
			if hashString == blockHash {
				Blockchains.GetBlockByHash([]byte(bc[:]))
				break
			}
		}
		w.Write([]byte("success！"))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	account1Address := req.URL.Query().Get("account1Address")
	account2Address := req.URL.Query().Get("account2Address")
	money := req.URL.Query().Get("money")
	account1 := wallet.LoadWallet(account1Address)
	switch req.Method {
	case http.MethodGet:
		{
			log.Printf("\n\n\n")
			log.Println("接受到wallet发送的交易")
			var bigs, _ = new(big.Int).SetString(money, 10)
			trade := account1.Transfer(account2Address, bigs)
			isAdded := Blockchains.AddTransaction(account1Address, account2Address, trade.Getvalue(), account1.PublicKey(), trade.GenerateSignature())
			color.HiGreen("这笔交易验证通过吗? %v\n", isAdded)
			Blockchains.Mining()
			w.Write([]byte("success！"))
		}
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Save(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		{
			// 保存区块链数据
			err := block.SaveBlockchain(Blockchains)
			if err != nil {
				log.Fatal(err)
			}
		}
		w.Write([]byte("success！"))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (bcs *BlockchainServer) Load(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		miner := wallet.LoadWallet(block.MINING_ACCOUNT_ADDRESS)
		Miners = miner

		blockchain := &block.Blockchain{}

		// 加载保存的区块链数据
		err := block.LoadBlockchain(blockchain)

		Blockchains = blockchain
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte("success！"))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (bcs *BlockchainServer) GetBalance(w http.ResponseWriter, req *http.Request) {
	accountAddress := req.URL.Query().Get("accountAddress")
	switch req.Method {
	case http.MethodGet:
		big := big.NewInt(0)
		big = Blockchains.GetBalance(accountAddress)
		fmt.Println(accountAddress, "的余额是：", big)
		message := accountAddress + "的余额是：" + big.String()
		w.Write([]byte(message))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (bcs *BlockchainServer) GetTransactionByHash(w http.ResponseWriter, req *http.Request) {
	Hash := req.URL.Query().Get("Hash")
	switch req.Method {
	case http.MethodGet:
		for i := 0; i < 20; i++ {
			bc := Blockchains.OneBlock(i).Hash1()
			hashString := hex.EncodeToString(bc[:])
			if hashString == Hash {
				Blockchains.GetTransactionByHash([]byte(bc[:]))
				break
			}
		}
		w.Write([]byte("success！"))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (bcs *BlockchainServer) Prints(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		Blockchains.Print()
		w.Write([]byte("success！"))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (bcs *BlockchainServer) sendTraDataToWalletServer(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		// 构建请求体
		fmt.Println("block.Tra", block.Tra)
		a := block.TradeHistory()

		// 将数据转换为JSON字符串
		jsonData, err := json.Marshal(a)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("jsonData", jsonData)
		// 创建HTTP请求
		req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/tra-data", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}

		// 设置请求头部为JSON类型
		req.Header.Set("Content-Type", "application/json")

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// 检查响应状态码
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Request failed with status: %s", resp.Status)
		}

		log.Println("Data sent successfully")
		w.Write([]byte("success！"))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (bcs *BlockchainServer) Run() {

	//区块链持久化--加载数据库
	http.HandleFunc("/load", bcs.Load)

	//新建一条链
	http.HandleFunc("/chain", bcs.Chain)

	//挖矿，区块添加
	http.HandleFunc("/mine", bcs.Mine)

	//交易
	http.HandleFunc("/transactions", bcs.Transactions)

	//根据number查block
	http.HandleFunc("/GetBlockByNumber", bcs.GetBlockByNumber)

	//根据hash查block
	http.HandleFunc("/GetBlockByHash", bcs.GetBlockByHash)

	//查询账户余额
	http.HandleFunc("/GetBalance", bcs.GetBalance)

	//根据hash查交易
	http.HandleFunc("/GetTransactionByHash", bcs.GetTransactionByHash)

	//打印block和交易信息
	http.HandleFunc("/prints", bcs.Prints)

	//区块链同步，网络通信
	http.HandleFunc("/sendTraDataToWalletServer", bcs.sendTraDataToWalletServer)

	//区块链持久化--将链上的信息保存到数据库
	http.HandleFunc("/save", bcs.Save)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(bcs.Port())), nil))

}
