/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

         http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// CounterChaincode example simple Chaincode implementation
type CounterChaincode struct {
}

// カウンター情報
type Counter struct {
	Name string `json:"name"`
	Counts int `json:"counts"`
}

const numOfCounters int = 3

// カウンター情報の初期値を設定
func (cc *CounterChaincode) Init(stub *shim.ChaincodeStub, function string, args []string)   ([]byte, error) {
	var counters [numOfCounters]Counter
	var countersBytes [numOfCounters][]byte
	
	// カウンター情報を生成
	counters[0] = Counter{Name: "Office Worker", Counts: 0}
	counters[1] = Counter{Name: "Home Worker", Counts:0}
	counters[2] = Counter{Name: "Student", Counts:0}
	
	// カウンター情報をワールドステートメントに追加
	for i := 0; i < len(counters); i++ {
		// JSON形式に変換
		countersBytes[i], _ = json.Marshal(counters[i])
		// ワールドステートに追加
		stub.PutState(strconv.Itoa(i), countersBytes[i])
	}
	
	return nil, nil
}

// カウンター情報を更新
func (cc *CounterChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	// function名でハンドリング
	if funtion == "countUp" {
		// カウントアップを実行
		return cc.countUp(stub, args)
	}
	
	return nil, errors.New("Received unknown function")
}


// カウンター情報を参照
func (cc *CounterChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	// function名でハンドリング
	if function == "refresh" {
		// カウンター情報を取得
		return cc.getCounters(stub, args)
	}
	
	return nil, errors.New("Received unknown function")
}

// カウントアップを実行
func (cc *CounterChaincode) countUp(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	// ワールドステートから選択されたカウンター情報を取得
	counterId := args[0]
	counterJson, _ := stub.GetState(counterId)
	
	// 取得したJSON形式の情報をCounterに変換
	counter := Counter{} // Counter構造体の初期化
	json.Unmarshal(counterJson, &counter)
	
	// カウントアップ
	counter.Counts++
	
	// ワールドステートに更新後の値を追加
	counterJson, _ = json.Marshal(counter)
	stub.PutState(counterId, counterJson)
	
	return nil, nil
}

// Validationg Peerに接続し、チェーンコードを実行
func main() {
	err := shim.Start(new(CounterChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
