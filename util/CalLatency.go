package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type TxTimestamp struct {
	AddTxpoolTime   int64
	ApplyStartTime  int64
	ApplyFinishTime int64
}

type TxLatency struct {
	WaitApplyLatency int64
	ApplyLatency     int64
	WaitBlockLatency int64
	PowLatency       int64
}

type Hash [32]byte

type User struct {
	Name string
	Age  int8
}

type Client struct {
	TransactionHash string
	SendTime        int64
	RespTime        int64
}

type ClientTime struct {
	SendTime int64
	RespTime int64
}

var (
	waitApplyLatency    int64
	applyLatency        int64
	waitBlockLatency    int64
	powLatency          int64
	sendLatencyAvg      int64
	waitApplyLatencyAvg int64
	applyLatencyAvg     int64
	waitBlockLatencyAvg int64
	powLatencyAvg       int64
	respLatencyAvg      int64
	txNum               int64
)

func testMarshal() []byte {
	user := User{
		Name: "Tab",
		Age:  18,
	}
	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func testUnmarshal(data []byte) {
	var user User
	err := json.Unmarshal(data, &user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
}
func readClient() map[string]ClientTime {
	fp, err := os.OpenFile("Client/client.json", os.O_RDONLY, 0755)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 10000000)
	//TxTimeStamps := make(map[Hash]*TxTimestamp)
	var clients []Client
	var clientMap map[string]ClientTime
	clientMap = make(map[string]ClientTime)
	n, err := fp.Read(data)
	if err := json.Unmarshal(data[:n], &clients); err != nil {
		log.Fatal(err)
	}
	for _, value := range clients {
		clientMap[value.TransactionHash] = ClientTime{SendTime: value.SendTime, RespTime: value.RespTime}
	}
	return clientMap
}

func readBlock(index int) (int64, int64) {
	fp, err := os.OpenFile("blocktime/block"+fmt.Sprint(index)+".json", os.O_RDONLY, 0755)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 1000000)
	//TxTimeStamps := make(map[Hash]*TxTimestamp)
	var powTime []int64
	n, err := fp.Read(data)
	if err := json.Unmarshal(data[:n], &powTime); err != nil {
		log.Fatal(err)
	}

	return powTime[0], powTime[1]
}

func readBlock2(filename string) (int64, int64) {
	fp, err := os.OpenFile("/home/cshiru/Latency/timestamps/blocktime/"+filename, os.O_RDONLY, 0755)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 1000000)
	//TxTimeStamps := make(map[Hash]*TxTimestamp)
	var powTime []int64
	n, err := fp.Read(data)
	if err := json.Unmarshal(data[:n], &powTime); err != nil {
		fmt.Println(filename)
		log.Fatal(err)
	}

	return powTime[0], powTime[1]
}

func readTx(index int) map[string]TxTimestamp {
	fp, err := os.OpenFile("txtime/tx_block"+fmt.Sprint(index)+".json", os.O_RDONLY, 0755)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 1000000)
	//TxTimeStamps := make(map[Hash]*TxTimestamp)
	var TxTimeStamps map[string]TxTimestamp
	n, err := fp.Read(data)
	if err := json.Unmarshal(data[:n], &TxTimeStamps); err != nil {
		log.Fatal(err)
	}
	return TxTimeStamps
}

func readTx2(filename string) map[string]TxTimestamp {
	fp, err := os.OpenFile("/home/cshiru/Latency/timestamps/txtime/"+filename, os.O_RDONLY, 0755)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 1000000)
	//TxTimeStamps := make(map[Hash]*TxTimestamp)
	var TxTimeStamps map[string]TxTimestamp
	n, err := fp.Read(data)
	if err := json.Unmarshal(data[:n], &TxTimeStamps); err != nil {
		fmt.Println(filename)
		log.Fatal(err)
	}
	return TxTimeStamps
}

func calGlobalLatency() {
	blockNum := 171
	latencyMap := make(map[string]TxLatency)
	sendLatencyAvg = 0
	waitApplyLatencyAvg = 0
	applyLatencyAvg = 0
	waitBlockLatencyAvg = 0
	powLatencyAvg = 0
	respLatencyAvg = 0
	txNum = 0
	var clientMap map[string]ClientTime
	clientMap = readClient()
	for i := 1; i <= blockNum; i++ {
		powStart, powFinish := readBlock(i)
		var txTimestamps map[string]TxTimestamp
		txTimestamps = readTx(i)

		for hash, tx := range txTimestamps {
			sendLatency := tx.ApplyStartTime - clientMap[hash].SendTime
			waitApplyLatency := tx.ApplyStartTime - tx.AddTxpoolTime
			applyLatency := tx.ApplyFinishTime - tx.ApplyStartTime
			waitBlockLatency := powStart - tx.ApplyFinishTime
			powLatency := powFinish - powStart
			respatency := clientMap[hash].RespTime - powFinish
			sendLatencyAvg += sendLatency
			waitApplyLatencyAvg = waitApplyLatencyAvg + waitBlockLatency
			applyLatencyAvg += applyLatency
			waitBlockLatencyAvg += waitBlockLatency
			powLatencyAvg += powLatency
			respLatencyAvg += respatency
			txNum++
			latencyMap[hash] = TxLatency{WaitApplyLatency: waitApplyLatency, ApplyLatency: applyLatency, WaitBlockLatency: waitBlockLatency, PowLatency: powLatency}
		}

	}
	sendLatencyAvg /= txNum
	waitApplyLatencyAvg /= txNum
	applyLatencyAvg /= txNum
	waitBlockLatencyAvg /= txNum
	powLatencyAvg /= txNum
	respLatencyAvg /= txNum
	fmt.Println(sendLatencyAvg, waitApplyLatencyAvg, applyLatencyAvg, waitBlockLatencyAvg, powLatencyAvg, respLatencyAvg, txNum)
}

func calLocalLatency() {
	latencyMap := make(map[string]TxLatency)
	sendLatencyAvg = 0
	waitApplyLatencyAvg = 0
	applyLatencyAvg = 0
	waitBlockLatencyAvg = 0
	powLatencyAvg = 0
	respLatencyAvg = 0
	txNum = 0

	txtime, _ := ioutil.ReadDir("/home/cshiru/Latency/timestamps/txtime")
	blocktime, _ := ioutil.ReadDir("/home/cshiru/Latency/timestamps/blocktime")
	for i := 0; i < len(txtime); i++ {
		fmt.Println(txtime[i])
		powStart, powFinish := readBlock2(blocktime[i].Name())
		var txTimestamps map[string]TxTimestamp
		txTimestamps = readTx2(txtime[i].Name())
		for hash, tx := range txTimestamps {
			if tx.AddTxpoolTime == 0 {
				continue
			}
			waitApplyLatency = tx.ApplyStartTime - tx.AddTxpoolTime
			fmt.Println(tx.ApplyStartTime, tx.AddTxpoolTime, waitApplyLatency)
			applyLatency = tx.ApplyFinishTime - tx.ApplyStartTime
			//fmt.Println(applyLatency)
			waitBlockLatency = powFinish - tx.ApplyFinishTime
			fmt.Println(waitBlockLatency)
			powLatency = powFinish - powStart
			waitApplyLatencyAvg = waitApplyLatencyAvg + waitBlockLatency
			applyLatencyAvg += applyLatency
			waitBlockLatencyAvg += waitBlockLatency
			powLatencyAvg += powLatency
			txNum++
			latencyMap[hash] = TxLatency{WaitApplyLatency: waitApplyLatency, ApplyLatency: applyLatency, WaitBlockLatency: waitBlockLatency, PowLatency: powLatency}
		}
	}
	//for _, f := range txtime {
	//	fmt.Println(f.Name())
	//	powStart, powFinish := readBlock2(f.Name())
	//	var txTimestamps map[string]TxTimestamp
	//	txTimestamps = readTx2(f.Name())
	//
	//	for hash, tx := range txTimestamps {
	//		waitApplyLatency := tx.ApplyStartTime - tx.AddTxpoolTime
	//		applyLatency := tx.ApplyFinishTime - tx.ApplyStartTime
	//		waitBlockLatency := powStart - tx.ApplyFinishTime
	//		powLatency := powFinish - powStart
	//		waitApplyLatencyAvg = waitApplyLatencyAvg + waitBlockLatency
	//		applyLatencyAvg += applyLatency
	//		waitBlockLatencyAvg += waitBlockLatency
	//		powLatencyAvg += powLatency
	//		txNum++
	//		latencyMap[hash] = TxLatency{WaitApplyLatency: waitApplyLatency, ApplyLatency: applyLatency, WaitBlockLatency: waitBlockLatency, PowLatency: powLatency}
	//	}
	//}
	//for i := 1; i <= blockNum; i++ {
	//
	//
	//}
	waitApplyLatencyAvg /= txNum
	applyLatencyAvg /= txNum
	waitBlockLatencyAvg /= txNum
	powLatencyAvg /= txNum
	fmt.Println(waitApplyLatencyAvg/1000000, applyLatencyAvg/100000, waitBlockLatencyAvg/100000, powLatencyAvg/100000, txNum)
}

func main() {
	calLocalLatency()
}
