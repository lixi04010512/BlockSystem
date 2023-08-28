package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"lxblockchain/block"
	"lxblockchain/wallet"
	"math/big"
)

func init() {

	color.Green("     ██╗ ██████╗ ██╗  ██╗███╗   ██╗██╗  ██╗ █████╗ ██╗")
	color.Green("     ██║██╔═══██╗██║  ██║████╗  ██║██║  ██║██╔══██╗██║")
	color.Green("     ██║██║   ██║███████║██╔██╗ ██║███████║███████║██║")
	color.Green("██   ██║██║   ██║██╔══██║██║╚██╗██║██╔══██║██╔══██║██║")
	color.Green("╚█████╔╝╚██████╔╝██║  ██║██║ ╚████║██║  ██║██║  ██║██║")
	color.Green("╚════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝")

	color.Red("██████╗ ██╗      ██████╗  ██████╗██╗  ██╗ ██████╗██╗  ██╗ █████╗ ██╗███╗   ██╗")
	color.Red("██╔══██╗██║     ██╔═══██╗██╔════╝██║ ██╔╝██╔════╝██║  ██║██╔══██╗██║████╗  ██║")
	color.Red("██████╔╝██║     ██║   ██║██║     █████╔╝ ██║     ███████║███████║██║██╔██╗ ██║")
	color.Red("██╔══██╗██║     ██║   ██║██║     ██╔═██╗ ██║     ██╔══██║██╔══██║██║██║╚██╗██║")
	color.Red("██████╔╝███████╗╚██████╔╝╚██████╗██║  ██╗╚██████╗██║  ██║██║  ██║██║██║ ╚████║")
	color.Red("╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝")

	log.SetPrefix("Blockchain: ")
}

func main() {
	// 1、根据矿工私钥加载钱包,输出矿工区块链地址 (LoadWallet)
	miner := wallet.LoadWallet(block.MINING_ACCOUNT_ADDRESS)
	minerAddress := miner.BlockchainAddress()
	fmt.Println("minerAddress", minerAddress)

	// 2、生成2个账户account1、account2的钱包,输出account1、account2区块链地址 (NewWallet)
	account1 := wallet.NewWallet()
	account2 := wallet.NewWallet()
	account1Address := account1.BlockchainAddress()
	account2Address := account2.BlockchainAddress()
	fmt.Println("account1Address", account1Address)
	fmt.Println("account2Address", account2Address)

	// 3、新建一条链
	blockchain := block.NewBlockchain(minerAddress)
	blockchain.Mining()

	// 4、转账交易 矿工->account1 数量2e+19 (200000000000000000000)
	var bigs, _ = new(big.Int).SetString("20000000000000000000", 10)
	trade := miner.Transfer(account1Address, bigs)
	isAdded := blockchain.AddTransaction(minerAddress, account1Address, trade.Getvalue(), miner.PublicKey(), trade.GenerateSignature())
	color.HiGreen("这笔交易验证通过吗? %v\n", isAdded)
	blockchain.Mining()

	// 5、转账交易 account1->account2 数量2000
	trade1 := account1.Transfer(account2Address, big.NewInt(2000))
	isAdded = blockchain.AddTransaction(account1Address, account2Address, trade1.Getvalue(), account1.PublicKey(), trade1.GenerateSignature())
	color.HiGreen("这笔交易验证通过吗? %v\n", isAdded)

	// 6、打包区块 Mining
	blockchain.Mining()

	// 7、查询区块 GetBlockByNumber
	blockchain.GetBlockByNumber(1)

	// 8、查询区块 GetBlockByHash
	bc := blockchain.LastBlock().Hash()
	blockchain.GetBlockByHash([]byte(bc[:]))

	// 9、查询交易GetTransactionByHash
	blockchain.GetTransactionByHash([]byte(bc[:]))

	// 10、输出区块信息（区块头和区块交易) Print
	blockchain.Print()

}
