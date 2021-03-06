package main

import (
	"github.com/xxRanger/Percome19-Crowd-Demo/back-end/appClient"
	"github.com/xxRanger/Percome19-Crowd-Demo/back-end/contract"
	"github.com/xxRanger/Percome19-Crowd-Demo/back-end/user"
	"github.com/xxRanger/Percome19-Crowd-Demo/back-end/zcrypto"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
)

const (
	AGG_CONFIG = "etc/agg.json"
	ERC20_CONFIG = "etc/erc20.json"
	LOG_FILE="etc/logfile"
)


func testPallier() {
	pub,pri,err:=zcrypto.NewPallier(128);
	cypher,err:=pub.Encrypt(big.NewInt(232131230))
	if err!=nil {
		log.Println()
		return
	}
	log.Println("N: ",pri.N," G: ",pri.G, "lambda", pri.Lambda, "mu:", pri.Mu," cipher:",cypher.C)
	m:= pri.Decrypt(cypher)
	log.Println(m)
}

func testBulletProof2() {
	rp:=zcrypto.RPProve(big.NewInt(25))
	result := zcrypto.RPVerify(rp);
	log.Println(len(rp.IPP.L),len(rp.IPP.R),len(rp.IPP.Challenges))
	if result {
		log.Println("pass test")
	} else {
		log.Println("test fail")
	}

}

func requestParserWrapper(manager *appClient.HandlerManager, agg *contract.Agg, token *contract.ERC20) func(w http.ResponseWriter, r *http.Request){
	return func (w http.ResponseWriter, r *http.Request){
		log.Println("receive a reqeust")
		var upgrader = websocket.Upgrader{}
		upgrader.CheckOrigin = func(rq *http.Request) bool { return true }
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		handler:= appHandler(c,agg,token)
		manager.Register <- handler
		handler.HandleRequest()
		manager.Unregister <- handler
		defer c.Close()
	}
}

func appHandler(c *websocket.Conn, agg *contract.Agg, token *contract.ERC20) *appClient.Handler {
	return appClient.NewHandler(c,agg,token)
}

func getNonce(agg *contract.Agg) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("get nonce")
		vars:=mux.Vars(r)
		user:= vars["user"]
		nonce,err:= agg.GetNonce(common.HexToAddress(user))
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		log.Println("user:",user)
		log.Println("nonce:",nonce)
		nonceWrapper, err:= json.Marshal(nonce)
		if err!=nil {
			log.Println(err.Error())
			log.Println(appClient.UNMARSHAL_JSON_ERROR)
			return
		}
		w.Write(nonceWrapper)
	}
}

func getChainId(agg *contract.Agg) func(w http.ResponseWriter, r* http.Request) {
	return func(w http.ResponseWriter,r *http.Request){
		log.Println("get chain id")

		chainId,err:= agg.GetChainId()
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}

		chainIdWrapper, err:= json.Marshal(chainId)
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(chainIdWrapper)
	}
}

func getEncryptedData(w http.ResponseWriter, r *http.Request) {
	log.Println("get encrypted data")
	amount:=big.NewInt(0)
	data, err:= ioutil.ReadAll(r.Body)

	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data,&amount)
	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	negative:=make([]byte,4)
	negative[0]=byte(0)
	negative[1]=byte(0)
	negative[2]=byte(0)
	if(amount.Cmp(big.NewInt(0)) == -1) {
		negative[3] = byte(1)
	} else {
		negative[3] = byte(0)
	}

	log.Println("data to encrypt:",amount)
	cipher,err := appClient.GenEncryption(amount)
	if err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	submitData := "0x"+hex.EncodeToString(negative)+cipher.Text(16)
	submitProofByte := appClient.GenBulletProof(amount)
	submitProof:= "0x"+hex.EncodeToString(submitProofByte)

	submitPayload,err := json.Marshal(&appClient.SubmitPayload{
		SubmitData:submitData,
		SubmitProof:submitProof,
	})
	if err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.Write(submitPayload)
}

func getStatistics(agg *contract.Agg) func(w http.ResponseWriter, r* http.Request) {
	return func(w http.ResponseWriter,r *http.Request){
		log.Println("get statistics")
		vars:=mux.Vars(r)
		taskId,_:= new(big.Int).SetString(vars["taskId"],10)

		countByte,err:= agg.Call(contract.FUNCTION_GET_REGISTER_NUMBER_OF_TASK,taskId)
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		count:=new(big.Int).SetBytes(countByte)
		var invalidSamples []*big.Int
		var submitValues []*big.Int
		//submitValues:=make([]*big.Int,count.Int64(),count.Int64())


		for i:=0; i< int(count.Int64()); i++ {
			submitDataByte,err:= agg.Call(contract.FUNCTION_GET_SUBMIT_DATA_OF_TASK,taskId,big.NewInt(int64(i)))
			if err!=nil {
				log.Println(err.Error())
				http.Error(w,err.Error(),http.StatusInternalServerError)
				return
			}
			submitProofByte,err:= agg.Call(contract.FUNCTION_GET_SUBMIT_PROOF_OF_TASK,taskId,big.NewInt(int64(i)))
			if err!=nil {
				log.Println(err.Error())
				http.Error(w,err.Error(),http.StatusInternalServerError)
				return
			}


			submitDataLen:= new(big.Int).SetBytes(submitDataByte[32:64])
			log.Println("static byte:",submitDataByte)
			negative:=submitDataByte[64:68]
			submitDataByte=submitDataByte[68:64+submitDataLen.Int64()]
			submitProofLen := new(big.Int).SetBytes(submitProofByte[32:64])
			submitProofByte = submitProofByte[64:64+submitProofLen.Int64()]


			rp:= new(zcrypto.RangeProof).SetBytes(submitProofByte)

			if err!=nil {
				log.Println(err.Error())
				http.Error(w,err.Error(),http.StatusInternalServerError)
				return
			}

			submitData:=new(big.Int).SetBytes(submitDataByte)
			decryptedData:= zcrypto.PriKey.Decrypt(&zcrypto.Cypher {
				C: submitData,
			})
			if !zcrypto.RPVerify(*rp) {
				if(len(invalidSamples)<5){
					if negative[3]==byte(1) {
						N,_:=big.NewInt(0).SetString(zcrypto.N,10)
						decryptedData.ModInverse(decryptedData,N)
						decryptedData.Neg(decryptedData)
					}
					invalidSamples = append(invalidSamples,decryptedData)
				}
			} else {
				submitValues= append(submitValues,decryptedData)
			}
		}

		type payload struct {
			SubmitValues []*big.Int `json:"submitValues"`
			InvalidSamples[]*big.Int `json:"invalidSamples"`
		}

		submitValuesWrapper,err :=json.Marshal(&payload{
			SubmitValues:submitValues,
			InvalidSamples:invalidSamples,
		})
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
		}
		w.Write(submitValuesWrapper)
	}
}

func requireEther(owner *user.User, agg *contract.Agg) func (http.ResponseWriter,*http.Request) {
	return func(w http.ResponseWriter,r *http.Request) {
		vars:=mux.Vars(r)
		account:=common.HexToAddress(vars["user"])
		log.Println(vars["user"],"require for ether")
		value,_:=new(big.Int).SetString(user.TRANSFER_VALUE,10)
		// check user Balance first
		// only send if user balance less than 0.001 ether
		balance,err:=agg.GetEther(account)
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		thresholdValue,_:= new(big.Int).SetString(user.THRESHOLD_VALUE,10)
		if balance.Cmp(thresholdValue) == 1 {
			errorMsg:= appClient.MSG_ALREADY_HAS_ENOUGH_ETHER
			log.Println("user: ",vars["user"],errorMsg)
			http.Error(w,errorMsg,http.StatusBadRequest)
			return
		}

		// send ether
		err =owner.Transfer(account,value)
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}

		valueWrapper,err:=json.Marshal(value)
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		log.Println(vars["user"],"got ether")
		w.Write(valueWrapper)
	}
}

func getEther(agg *contract.Agg) func (http.ResponseWriter,*http.Request) {
	return func(w http.ResponseWriter,r *http.Request) {
		vars:=mux.Vars(r)
		account:=common.HexToAddress(vars["user"])

		balance,err:=agg.GetEther(account)
		log.Println("user",vars["user"],"balance: ",balance)
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}

		balanceWrapper,err:=json.Marshal(balance)
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(balanceWrapper)
	}
}

func validate(w http.ResponseWriter,r *http.Request) {
	var password string
	data, err:= ioutil.ReadAll(r.Body)

	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	log.Println("admin password",password)
	err = json.Unmarshal(data,&password)
	if err!=nil {
		log.Println(err.Error())
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	if password==appClient.ADMIN_PASSWORD {
		accessToken:=big.NewInt(20190311)
		accessTokenWrapper,err:= json.Marshal(accessToken)
		if err!=nil {
			log.Println(err.Error())
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		w.Write(accessTokenWrapper)
	} else {
		errorMsg:=appClient.MSG_ADMIN_WRONG_PASSWORD
		log.Println(errorMsg)
		http.Error(w,errorMsg,http.StatusInternalServerError)
		return
	}
}


func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func start() {
	r := mux.NewRouter()

	manager:= appClient.NewHandlerManager()
	aggConfig,err:=contract.ReadConfig(AGG_CONFIG)
	if err!=nil {
		panic(err)
	}
	erc20Config,err:=contract.ReadConfig(ERC20_CONFIG)
	if err!=nil {
		panic(err)
	}
	agg:= contract.NewAgg(chainPort,aggConfig.Abi,common.HexToAddress(aggConfig.Address))
	token:= contract.NewERC20(chainPort,erc20Config.Abi,common.HexToAddress(erc20Config.Address))

	pk,err:= crypto.HexToECDSA(user.MASTER_KEY)
	if err!=nil {
		panic(err)
	}

	owner := user.NewUser(common.HexToAddress(user.MASTER_ADDRESS),pk,agg)
	nonce,err:=agg.GetNonce(owner.Address)
	if err!=nil {
		log.Println(err.Error())
		panic(err)
	}
	owner.Nonce = nonce

	go manager.Start()
	go manager.SubScriptContractEvent(agg)

	r.HandleFunc("/", requestParserWrapper(manager,agg,token)).Methods("GET")
	r.HandleFunc("/nonce/{user}",getNonce(agg)).Methods("GET")
	r.HandleFunc("/chainId",getChainId(agg)).Methods("GET")
	r.HandleFunc("/encryptedData",getEncryptedData).Methods("POST")
	r.HandleFunc("/statistics/{taskId}",getStatistics(agg)).Methods("GET")
	r.HandleFunc("/requireEther/{user}",requireEther(owner,agg)).Methods("GET")
	r.HandleFunc("/ether/{user}",getEther(agg)).Methods("GET")
	r.HandleFunc("/validate",validate).Methods("POST")

	fmt.Println("Running http server")
	err =http.ListenAndServe(
		"0.0.0.0:4000",
		handlers.CORS(
			handlers.AllowedMethods([]string{"get", "options", "post", "put", "head"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"Content-Type"}),
		)(r),
	)
	if err!=nil {
		log.Println(err.Error())
	}
}

var (
	chainPort string
)

func init() {
	flag.StringVar(&chainPort,"chainPort","wss://kovan.infura.io/ws/v3/ab1d2a604f814c688b4cbfbc9550fb1b","chain websocket(default wss://kovan.infura.io/ws)")
	flag.Parse()
}

func main() {
	//testPallier()
	//testBulletProof()
	//testBulletProof2()
	log.Println("chain port:",chainPort)
	f,err:= os.OpenFile(LOG_FILE,os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
	if err!=nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	start()
}
