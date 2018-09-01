package main

import (
	"fmt"
	"strconv"
	"time"
	// "reflect"
	//shim是中间层
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
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
	animalDensity    int32 //蟹苗密度
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
	id                      int64
	poolId                  int64
	initInformation         baseInfo
	feedInformation         []feedInfo
	waterQualityInformation []waterQuality
	transferInformation     []transfer
	storeInformation        []store
}

type craber interface {
	changeFeed(_feedName string) bool
	changeWaterQuality(_whetherQualified bool, _checkAgent string, _animalDensity int32) bool
	changeTransfer(from string, to string) bool
	changeStore(temperature int, wetness int) bool
}

//
func newCrab(id int64, poolId int64, operator string) *crab {
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

func (cra *crab) changeWaterQuality(_whetherQualified bool, _checkAgent string, _animalDensity int32) bool {
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
	crabs             map[int64]*crab
	feedTrace         map[int64][]baseInfo
	waterQualityTrace map[int64][]baseInfo
	transferTrace     map[int64][]baseInfo
	storeTrace        map[int64][]baseInfo
	// exist map[int64]bool
	// var map_variable map[key_data_type]value_data_type
}

//

func newTrace() *trace {
	result := &trace{}
	result.crabs = make(map[int64]*crab)
	result.feedTrace = make(map[int64][]baseInfo)
	result.waterQualityTrace = make(map[int64][]baseInfo)
	result.transferTrace = make(map[int64][]baseInfo)
	result.storeTrace = make(map[int64][]baseInfo)
	return result
}

type tracer interface {
	addcrab(_id int64, _poolId int64, _operator string) bool
	isExist(_id int64) bool
	pushFeed(_id int64, _feedName string, _opratorName string) bool
	pushWaterQuality(_id int64, _whetherQualified bool, _checkAgent string, _animalDensity int32, _opratorName string) bool
	pushTransfer(_id int64, from string, to string, _opratorName string) bool
	pushStore(_id int64, temperature int, wetness int, _opratorName string) bool
}

func (tra *trace) isExist(_id int64) bool {
	_, result := tra.crabs[_id]
	return result

}
func (tra *trace) addcrab(_id int64, _poolId int64, _operator string) bool {
	if tra.isExist(_id) {
		return false
	} else {
		tra.crabs[_id] = newCrab(_id, _poolId, _operator)
		return true
	}
}

func (tra *trace) pushFeed(_id int64, _feedName string, _operatorName string) bool {
	if !tra.isExist(_id) {
		return false
	} else {
		tra.crabs[_id].changeFeed(_feedName)
		tra.feedTrace[_id] = append(tra.feedTrace[_id], baseInfo{_operatorName, time.Now().Unix()})
		return true
	}
}
func (tra *trace) pushWaterQuality(_id int64, _whetherQualified bool, _checkAgent string, _animalDensity int32, _opratorName string) bool {
	if !tra.isExist(_id) {
		return false
	} else {
		tra.crabs[_id].changeWaterQuality(_whetherQualified, _checkAgent, _animalDensity)
		tra.waterQualityTrace[_id] = append(tra.waterQualityTrace[_id], baseInfo{_opratorName, time.Now().Unix()})
		return true
	}
}
func (tra *trace) pushTransfer(_id int64, from string, to string, _opratorName string) bool {
	if !tra.isExist(_id) {
		return false
	} else {
		tra.crabs[_id].changeTransfer(from, to)
		tra.transferTrace[_id] = append(tra.transferTrace[_id], baseInfo{_opratorName, time.Now().Unix()})
		return true
	}
}
func (tra *trace) pushStore(_id int64, temperature int, wetness int, _opratorName string) bool {
	if !tra.isExist(_id) {
		return false
	} else {
		tra.crabs[_id].changeStore(temperature, wetness)
		tra.storeTrace[_id] = append(tra.storeTrace[_id], baseInfo{_opratorName, time.Now().Unix()})
		return true
	}
}

// chaincode初始化或者升级的时候调用
//将SimpleAsset与Init绑定,new一个SimapleAsset的时候其将拥有init函数
//stub shim.ChaincodeStubInterface是函数参数，类型是shim.ChaincodeStubInterface，参数名字叫stub
//peer.Response是返回值
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// ：号用来声明一个未被声明的变量并赋值
	// args := stub.GetStringArgs()
	// if len(args) != 2 {
	// 	return shim.Error("Incorrect arguments. Expecting a key and a value")
	// }

	// // Set up any variables or assets here by calling stub.PutState()

	// // We store the key and the value on the ledger
	// err := stub.PutState(args[0], []byte(args[1]))
	// if err != nil {
	// 	return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	// }
	// return shim.Success(nil)
	traceInfo := newTrace()
	result, jsonErr := json.Marshal(traceInfo)
	if jsonErr != nil {
		return shim.Error(fmt.Errorf("json转化错误 %s", jsonErr).Error())
	} else {
		putErr := stub.PutState("trace", result)
		if putErr != nil {

			return shim.Error(fmt.Errorf("新建trace遇到错误 %s", putErr).Error())
		} else {
			return shim.Success(nil)
		}
	}

}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.ts
func argsNumError(argsLen int) error {
	return fmt.Errorf("参数个数应该是" + string(argsLen))
}

//不管是get还是set都将调用该函数
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	// _ := args

	fmt.Println("参数是")
	fmt.Println(args)

	var err error
	var getResult string
	//先创建对象
	var tra *trace
	var result bool
	fmt.Println(result)
	tra = newTrace()
	traceTest, _ := stub.GetState("trace")
	json.Unmarshal(traceTest, tra)

	fmt.Println("调用前获取到是：")
	fmt.Println(*tra)
	if fn == "getTraceInfo" { // assume 'get' even if fn is nil
		getResult, err = getTraceInfo(stub)
		fmt.Println("get函数获取到的是")
		fmt.Println(getResult)
		return shim.Success([]byte(getResult))
	} else if fn == "isExist" {
		//需要装成
		if len(args) != 1 {
			err = argsNumError(len(args))
		} else {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			test := strconv.FormatBool(tra.isExist(id))
			fmt.Println("判断函数返回的是" + test)
			return shim.Success([]byte(test))
		}
		// addcrab(_id int64, _poolId int64, _operator string) bool
		// 	isExist(_id int64) bool
		// 	pushFeed(_id int64, _feedName string, _opratorName string) bool
		// 	pushWaterQuality(_id int64, _whetherQualified bool, _checkAgent string, _animalDensity int32, _opratorName string) bool
		// 	pushTransfer(_id int64, from string, to string, _opratorName string) bool
		// 	pushStore(_id int64, temperature int, wetness int, _opratorName string) bool
	} else if fn == "addcrab" {
		if len(args) != 3 {
			err = argsNumError(len(args))
		} else {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			poolId, _ := strconv.ParseInt(args[1], 10, 64)
			operator := args[2]
			result = tra.addcrab(id, poolId, operator)
			fmt.Println("添加完成,*tra是")
			fmt.Println(*tra)
		}
	} else if fn == "pushFeed" {
		if len(args) != 3 {
			err = argsNumError(len(args))
		} else {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			feedName := args[1]
			operator := args[2]
			result = tra.pushFeed(id, feedName, operator)
			fmt.Println("pushfeed完成,*tra是")
			fmt.Println(*tra)
		}
	} else if fn == "pushWaterQuality" {
		if len(args) != 5 {
			err = argsNumError(len(args))
		} else {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			whetherQualified, _ := strconv.ParseBool(args[1])
			checkAgent := args[2]
			animalDensity, _ := strconv.ParseInt(args[3], 10, 32)
			operator := args[4]
			result = tra.pushWaterQuality(id, whetherQualified, checkAgent, int32(animalDensity), operator)
			fmt.Println("pushWater完成,*tra是")
			fmt.Println(*tra)
		}
	} else if fn == "pushTransfer" {
		if len(args) != 4 {
			err = argsNumError(len(args))
		} else {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			from := args[1]
			to := args[2]
			operator := args[3]
			result = tra.pushTransfer(id, from, to, operator)
			fmt.Println("pushTransfer完成,*tra是")
			fmt.Println(*tra)
		}
	} else if fn == "pushStore" {
		if len(args) != 4 {
			err = argsNumError(len(args))
		} else {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			temperature, _ := strconv.ParseInt(args[1], 10, 32)
			wetness, _ := strconv.ParseInt(args[2], 10, 32)
			operator := args[3]
			result = tra.pushStore(id, int(temperature), int(wetness), operator)
			fmt.Println("pushStore完成,*tra是")
			fmt.Println(json.Marshal(*tra))
		}
	}

	invokeResult, jsonErr := json.Marshal(*tra)
	if jsonErr != nil {
		err = fmt.Errorf("trace对象转json失败")
	} else {
		fmt.Println("要存储了，json是")
		fmt.Println(invokeResult)
		err = stub.PutState("trace", invokeResult)
		if err != nil {
			err = fmt.Errorf("更改trace失败 %s", err)
		} else {
			fmt.Println("调用后获取到是：")

			//可能还没确认
			var tra1 *trace
			tra1 = newTrace()
			traceAfter, _ := stub.GetState("trace")
			json.Unmarshal(traceAfter, tra1)
			fmt.Println(*tra1)
		}
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success(invokeResult)
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one

func getTraceInfo(stub shim.ChaincodeStubInterface) (string, error) {
	data, err := stub.GetState("trace")
	if err != nil {
		return "", fmt.Errorf("遇到错误：%s", err)
	}
	//账本是空的
	if data != nil {
		return "账本是空的", nil
	}
	// fmt.Println(data)
	fmt.Println("getTrace里面数据是")
	fmt.Println(string(data))
	//账本存在且不为空
	return string(data), nil
}

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

func main() {
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
	// if err := shim.Start(new(SimpleAsset)); err != nil {
	// 	fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	// }

}
