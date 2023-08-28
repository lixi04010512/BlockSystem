package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"lxblockchain/wallet"
	"math/big"
	"net/http"
	"path"
	"strconv"
)

const tempDir = "walletServer/htmltemplate"

type WalletServer struct {
	port    uint16
	gateway string //区块链的节点地址
}

type Transaction struct {
	senderAddress  string
	receiveAddress string
	value          *big.Int
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(tempDir, "index_bootstrap.html"))
		t.Execute(w, "")
	default:
		log.Printf("ERROR: 非法的HTTP请求方式")
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	//设置允许的方法
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		m, _ := myWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: 非法的HTTP请求方式")
	}
}

func (ws *WalletServer) walletByPrivatekey(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:

		w.Header().Add("Content-Type", "application/json")
		privatekey := req.FormValue("privatekey")
		myWallet := wallet.LoadWallet(privatekey)
		m, _ := myWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) History(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:

		// 将数据转换为JSON字符串
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 设置响应头部为JSON类型
		w.Header().Set("Content-Type", "application/json")

		// 将JSON字符串写入响应体
		w.Write(jsonData)
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}

}

var data []map[string]interface{}

func (ws *WalletServer) handleTraData(w http.ResponseWriter, req *http.Request) {
	// 解析请求体中的JSON数据

	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 打印解析后的数据
	fmt.Println("Received data:", data)

	// 返回响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data received successfully"))

}

func (ws *WalletServer) Run() {

	fs := http.FileServer(http.Dir("walletServer/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//网络通信
	http.HandleFunc("/tra-data", ws.handleTraData)

	http.HandleFunc("/", ws.Index)

	//生成新钱包
	http.HandleFunc("/wallet", ws.Wallet)

	//加载钱包
	http.HandleFunc("/walletByPrivatekey", ws.walletByPrivatekey)

	//显示历史交易
	http.HandleFunc("/history", ws.History)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))

}
