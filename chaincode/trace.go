package main

import (
	"fmt"
	// "reflect"
	"strconv"
	"time"
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
	Operator    string
	OperateTime int64
}

type feedInfo struct {
	FeedName string
}

type waterQuality struct {
	WhetherQualified bool
	CheckAgent       string
	AnimalDensity    int32 //蟹苗密度
}

type transfer struct {
	From string
	To   string
}

type store struct {
	Temperature int
	Wetness     int
}

type crab struct {
	Id                      int64
	PoolId                  int64
	InitInformation         baseInfo
	FeedInformation         []feedInfo
	WaterQualityInformation []waterQuality
	TransferInformation     []transfer
	StoreInformation        []store
}

type craber interface {
	changeFeed(_feedName string) bool
	changeWaterQuality(_whetherQualified bool, _checkAgent string, _animalDensity int32) bool
	changeTransfer(from string, to string) bool
	changeStore(temperature int, wetness int) bool
}

//
func newCrab(id int64, poolId int64, Operator string) *crab {
	var result = new(crab)
	result.Id = id
	result.PoolId = poolId
	result.InitInformation = baseInfo{Operator, time.Now().Unix()}
	return result
}

func (cra *crab) changeFeed(_feedName string) bool {
	cra.FeedInformation = append(cra.FeedInformation, feedInfo{_feedName})
	return true
}

func (cra *crab) changeWaterQuality(_whetherQualified bool, _checkAgent string, _animalDensity int32) bool {
	cra.WaterQualityInformation = append(cra.WaterQualityInformation, waterQuality{_whetherQualified, _checkAgent, _animalDensity})
	return true
}
func (cra *crab) changeTransfer(from string, to string) bool {
	cra.TransferInformation = append(cra.TransferInformation, transfer{from, to})
	return true
}
func (cra *crab) changeStore(temperature int, wetness int) bool {
	cra.StoreInformation = append(cra.StoreInformation, store{temperature, wetness})
	return true
}

type trace struct {
	Crabs             map[int64]*crab
	FeedTrace         map[int64][]baseInfo
	WaterQualityTrace map[int64][]baseInfo
	TransferTrace     map[int64][]baseInfo
	StoreTrace        map[int64][]baseInfo
	// exist map[int64]bool
	// var map_variable map[key_data_type]value_data_type
}

//

func newTrace() trace {
	result := trace{}
	result.Crabs = make(map[int64]*crab)
	result.FeedTrace = make(map[int64][]baseInfo)
	result.WaterQualityTrace = make(map[int64][]baseInfo)
	result.TransferTrace = make(map[int64][]baseInfo)
	result.StoreTrace = make(map[int64][]baseInfo)
	return result
}

type tracer interface {
	addcrab(_id int64, _poolId int64, _Operator string) bool
	isExist(_id int64) bool
	pushFeed(_id int64, _feedName string, _opratorName string) bool
	pushWaterQuality(_id int64, _whetherQualified bool, _checkAgent string, _animalDensity int32, _opratorName string) bool
	pushTransfer(_id int64, from string, to string, _opratorName string) bool
	pushStore(_id int64, temperature int, wetness int, _opratorName string) bool
}

func (tra *trace) isExist(_id int64) bool {
	_, result := tra.Crabs[_id]
	return result

}
func (tra *trace) addcrab(_id int64, _poolId int64, _operator string) bool {
	if tra.isExist(_id) {
		return false
	} else {
		tra.Crabs[_id] = newCrab(_id, _poolId, _operator)
		return true
	}
}

func (tra *trace) pushFeed(_id int64, _feedName string, _operatorName string) bool {
	if !tra.isExist(_id) {
		return false
	} else {
		tra.Crabs[_id].changeFeed(_feedName)
		tra.FeedTrace[_id] = append(tra.FeedTrace[_id], baseInfo{_operatorName, time.Now().Unix()})
		return true
	}
}
func (tra *trace) pushWaterQuality(_id int64, _whetherQualified bool, _checkAgent string, _animalDensity int32, _opratorName string) bool {
	if !tra.isExist(_id) {
		return false
	} else {
		tra.Crabs[_id].changeWaterQuality(_whetherQualified, _checkAgent, _animalDensity)
		tra.WaterQualityTrace[_id] = append(tra.WaterQualityTrace[_id], baseInfo{_opratorName, time.Now().Unix()})
		return true
	}
}
func (tra *trace) pushTransfer(_id int64, from string, to string, _opratorName string) bool {
	if !tra.isExist(_id) {
		return false
	} else {
		tra.Crabs[_id].changeTransfer(from, to)
		tra.TransferTrace[_id] = append(tra.TransferTrace[_id], baseInfo{_opratorName, time.Now().Unix()})
		return true
	}
}
func (tra *trace) pushStore(_id int64, temperature int, wetness int, _opratorName string) bool {
	if !tra.isExist(_id) {
		return false
	} else {
		tra.Crabs[_id].changeStore(temperature, wetness)
		tra.StoreTrace[_id] = append(tra.StoreTrace[_id], baseInfo{_opratorName, time.Now().Unix()})
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
	// var getResult string
	//先创建对象
	var tra trace
	var result bool
	existCrab := true
	fmt.Println(result)
	tra = newTrace()
	traceTest, _ := stub.GetState("trace")
	json.Unmarshal(traceTest, &tra)

	fmt.Println("调用前获取到是：")

	jsonss, _ := json.Marshal(tra)
	fmt.Println(string(jsonss))
	if fn == "getTraceInfo" { // assume 'get' even if fn is nil
		return shim.Success(jsonss)
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
	} else if fn == "addcrab" {
		if len(args) != 3 {
			err = argsNumError(len(args))
		} else {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			poolId, _ := strconv.ParseInt(args[1], 10, 64)
			operator := args[2]
			existCrab = tra.addcrab(id, poolId, operator)
			fmt.Println("添加完成,*tra是")
			fmt.Println(tra)
		}
	} else if fn == "pushFeed" {
		if len(args) != 3 {
			err = argsNumError(len(args))
		} else {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			feedName := args[1]
			operator := args[2]
			existCrab = tra.pushFeed(id, feedName, operator)

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
			existCrab = tra.pushWaterQuality(id, whetherQualified, checkAgent, int32(animalDensity), operator)

		}
	} else if fn == "pushTransfer" {
		if len(args) != 4 {
			err = argsNumError(len(args))
		} else {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			from := args[1]
			to := args[2]
			operator := args[3]
			existCrab = tra.pushTransfer(id, from, to, operator)

		}
	} else if fn == "pushStore" {
		if len(args) != 4 {
			err = argsNumError(len(args))
		} else {
			id, _ := strconv.ParseInt(args[0], 10, 64)
			temperature, _ := strconv.ParseInt(args[1], 10, 32)
			wetness, _ := strconv.ParseInt(args[2], 10, 32)
			operator := args[3]
			existCrab = tra.pushStore(id, int(temperature), int(wetness), operator)

		}
	}

	//判断是否因为不存在id出错
	if existCrab == false {
		if fn == "addcrab" {
			err = fmt.Errorf("id已存在")
			fmt.Println("id已存在")
		} else {
			err = fmt.Errorf("id不存在")
			fmt.Println("id不存在")
		}
	}

	//判断调用是否出错
	if err != nil {
		fmt.Println("参数错误")
		return shim.Error(err.Error())
	}
	invokeResult, jsonErr := json.Marshal(tra)
	if jsonErr != nil {
		err = fmt.Errorf("trace对象转json失败")
		return shim.Error(err.Error())
	} else {
		fmt.Println("要存储了，json是")
		fmt.Println(invokeResult)
		err = stub.PutState("trace", invokeResult)
		if err != nil {
			err = fmt.Errorf("更改trace失败 %s", err)
			return shim.Error(err.Error())
		}
	}
	// Return the result as success payload
	return shim.Success(invokeResult)
}

func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}

}
