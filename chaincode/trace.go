package main

import (
	"fmt"
	"time"
	// "reflect"
	//shim是中间层
	// "github.com/hyperledger/fabric/core/chaincode/shim"
	// "github.com/hryperledger/fabric/protos/peer"
	// "container/list"
)

// 接受shim的函数，shim是一个接口
type SimpleAsset struct {
}

type baseInfo struct {
	operator    string
	operateTime int64
}

type feedInfo struct {
	feedName string
}

type waterQuality struct {
	whetherQualified bool
	checkAgent       string
	animalDensity    uint32 //蟹苗密度
}

type transfer struct {
	from string
	to   string
}

type store struct {
	temperature int
	wetness     int
}

type crab struct {
	id                      uint64
	poolId                  uint64
	initInformation         baseInfo
	feedInformation         []feedInfo
	waterQualityInformation []waterQuality
	transferInformation     []transfer
	storeInformation        []store
}

type craber interface {
	changeFeed(_feedName string) bool
	changeWaterQuality(_whetherQualified bool, _checkAgent string, _animalDensity uint32) bool
	changeTransfer(from string, to string) bool
	changeStore(temperature int, wetness int) bool
}

//
func newCrab(id uint64, poolId uint64, operator string) *crab {
	var result = new(crab)
	result.id = id
	result.poolId = poolId
	result.initInformation = baseInfo{operator, time.Now().Unix()}
	return result
}

func (cra *crab) changeFeed(_feedName string) bool {
	cra.feedInformation = append(cra.feedInformation, feedInfo{_feedName})
	return true
}

func (cra *crab) changeWaterQuality(_whetherQualified bool, _checkAgent string, _animalDensity uint32) bool {
	cra.waterQualityInformation = append(cra.waterQualityInformation, waterQuality{_whetherQualified, _checkAgent, _animalDensity})
	return true
}
func (cra *crab) changeTransfer(from string, to string) bool {
	cra.transferInformation = append(cra.transferInformation, transfer{from, to})
	return true
}
func (cra *crab) changeStore(temperature int, wetness int) bool {
	cra.storeInformation = append(cra.storeInformation, store{temperature, wetness})
	return true
}

type trace struct {
	crabs             map[uint64]*crab
	feedTrace         map[uint64][]baseInfo
	waterQualityTrace map[uint64][]baseInfo
	transferTrace     map[uint64][]baseInfo
	storeTrace        map[uint64][]baseInfo
	// exist map[uint64]bool
	// var map_variable map[key_data_type]value_data_type
}

//

func newTrace() *trace {
	result := &trace{}
	result.crabs = make(map[uint64]*crab)
	result.feedTrace = make(map[uint64][]baseInfo)
	result.waterQualityTrace = make(map[uint64][]baseInfo)
	result.transferTrace = make(map[uint64][]baseInfo)
	result.storeTrace = make(map[uint64][]baseInfo)
	return result
}

type tracer interface {
	addcrab(_id uint64, _poolId uint64, _operator string) bool
	isExist(_id uint64) bool
	pushFeed(_id uint64, _feedName string, _opratorName string) bool
	pushWaterQuality(_id uint64, _whetherQualified bool, _checkAgent string, _animalDensity uint32, _opratorName string) bool
	pushTransfer(_id uint64, from string, to string, _opratorName string) bool
	pushStore(_id uint64, temperature int, wetness int, _opratorName string) bool
}

func (tra *trace) isExist(_id uint64) bool {
	_, result := tra.crabs[_id]
	return result

}
func (tra *trace) addcrab(_id uint64, _poolId uint64, _operator string) bool {
	if tra.isExist(_id) {
		return false
	} else {
		tra.crabs[_id] = newCrab(_id, _poolId, _operator)
		return true
	}
}

func (tra *trace) pushFeed(_id uint64, _feedName string, _operatorName string) bool {
	if !tra.isExist(_id) {
		return false
	} else {
		tra.crabs[_id].changeFeed(_feedName)
		tra.feedTrace[_id] = append(tra.feedTrace[_id], baseInfo{_operatorName, time.Now().Unix()})
		return true
	}
}
func (tra *trace) pushWaterQuality(_id uint64, _whetherQualified bool, _checkAgent string, _animalDensity uint32, _opratorName string) bool {
	if !tra.isExist(_id) {
		return false
	} else {
		tra.crabs[_id].changeWaterQuality(_whetherQualified, _checkAgent, _animalDensity)
		tra.waterQualityTrace[_id] = append(tra.waterQualityTrace[_id], baseInfo{_opratorName, time.Now().Unix()})
		return true
	}
}
func (tra *trace) pushTransfer(_id uint64, from string, to string, _opratorName string) bool {
	if !tra.isExist(_id) {
		return false
	} else {
		tra.crabs[_id].changeTransfer(from, to)
		tra.transferTrace[_id] = append(tra.transferTrace[_id], baseInfo{_opratorName, time.Now().Unix()})
		return true
	}
}
func (tra *trace) pushStore(_id uint64, temperature int, wetness int, _opratorName string) bool {
	if !tra.isExist(_id) {
		return false
	} else {
		tra.crabs[_id].changeStore(temperature, wetness)
		tra.storeTrace[_id] = append(tra.storeTrace[_id], baseInfo{_opratorName, time.Now().Unix()})
		return true
	}
}

func main() {

	// var base baseInfo
	// base.operator = "yapie"*
	// base.operateTime = time.Now().Unix()
	// fmt.Println(reflect.TypeOf(test))
	var crabTest *crab
	// crabTest.id = 1
	// crabTest.feedInfoMation = append(crabTest.feedInfoMation, feedInfo{"test"})
	// crabTest.feedInfoMation = append(crabTest.feedInfoMation, feedInfo{"test"})
	crabTest = newCrab(1, 1233456, "yapie")
	crabTest.changeFeed("sweet feed")
	crabTest.changeTransfer("here", "there")
	crabTest.changeStore(11, 10)
	crabTest.changeWaterQuality(true, "chengdu agent", 1111)

	var traceTest *trace
	traceTest = newTrace()
	traceTest.addcrab(1, 11, "yapie")
	traceTest.pushFeed(1, "milk", "yapie")
	traceTest.pushWaterQuality(1, true, "chengdu agent", 1111, "yapie")
	traceTest.pushTransfer(1, "chengdu", "tianshui", "yapie")
	traceTest.pushStore(1, 11, 23, "yapie")
	traceTest.pushTransfer(1, "chengdu", "tianshui", "yapie")
	traceTest.pushTransfer(1, "chengdu", "tianshui", "yapie")
	// fmt.Println(*traceTest)
	fmt.Println(*traceTest.crabs[1])

	// fmt.Println((*traceTest).crabs[1].)
	fmt.Println((*traceTest).feedTrace[1])
}

// // chaincode初始化或者升级的时候调用
// //将SimpleAsset与Init绑定,new一个SimapleAsset的时候其将拥有init函数
// //stub shim.ChaincodeStubInterface是函数参数，类型是shim.ChaincodeStubInterface，参数名字叫stub
// //peer.Response是返回值
// func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
// 	// ：号用来声明一个未被声明的变量并赋值
// 	args := stub.GetStringArgs()
// 	if len(args) != 2 {
// 		return shim.Error("Incorrect arguments. Expecting a key and a value")
// 	}

// 	// Set up any variables or assets here by calling stub.PutState()

// 	// We store the key and the value on the ledger
// 	err := stub.PutState(args[0], []byte(args[1]))
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
// 	}
// 	return shim.Success(nil)
// }

// // Invoke is called per transaction on the chaincode. Each transaction is
// // either a 'get' or a 'set' on the asset created by Init function. The Set
// // method may create a new asset by specifying a new key-value pair.ts

// //不管是get还是set都将调用该函数
// func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
// 	// Extract the function and args from the transaction proposal
// 	fn, args := stub.GetFunctionAndParameters()

// 	var result string
// 	var err error

// 	if fn == "addcrab" {
// 		result, err = addcrab(stub, args)
// 	} else { // assume 'get' even if fn is nil
// 		result, err = get(stub, args)
// 	}
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	// Return the result as success payload
// 	return shim.Success([]byte(result))
// }

// // Set stores the asset (both key and value) on the ledger. If the key exists,
// // it will override the value with the new one
// func addcrab(stub shim.ChaincodeStubInterface, args []string) (string, error) {

// 	if len(args) != 3 {
// 		return "", fmt.Errorf("")
// 	}

// 	err := stub.PutState(args[0], []byte(args[1]))
// 	if err != nil {
// 		return "", fmt.Errorf("Failed to set asset: %s", args[0])
// 	}
// 	return args[1], nil
// }

// // Get returns the value of the specified asset key
// func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
// 	if len(args) != 1 {
// 		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
// 	}

// 	value, err := stub.GetState(args[0])
// 	if err != nil {
// 		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
// 	}
// 	if value == nil {
// 		return "", fmt.Errorf("Asset not found: %s", args[0])
// 	}
// 	return string(value), nil
// }

// main function starts up the chaincode in the container during instantiate
