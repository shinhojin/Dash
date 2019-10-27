package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SupplyContract struct {
}

type product struct {
	OrderId   string
	Id        string
	ProductId string
	Amount    int
	priceEach int
	Detail    string
	Date      string
	SubId     string
}

func (t *SupplyContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Supply-ChainCode Init!!")
	return shim.Success(nil)
}

func (t *SupplyContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "setProduct" {
		return t.setProduct(stub, args)
	} else if function == "moveProduct" {
		return t.moveProduct(stub, args)
	} else if function == "useProduct" {
		return t.useProduct(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	} else if function == "history" {
		return t.history(stub, args)
	}

	return shim.Error("Invalid function name")
}

func (t *SupplyContract) setProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("call setProduct")
	orderId := args[0]
	id := args[1]
	productId := args[2]
	amount, _ := strconv.Atoi(args[3])
	price, _ := strconv.Atoi(args[4])
	detail := args[5]
	sub := args[6]
	date := time.Now().Format("2006-01-02 15:04:05")
	SupplyContract := product{
		OrderId:   orderId,
		Id:        id,
		ProductId: productId,
		Amount:    amount,
		priceEach: price,
		Detail:    detail,
		Date:      date,
		SubId:     sub}

	productBytes, _ := json.Marshal(SupplyContract)
	stub.PutState(orderId, productBytes)

	return shim.Success(nil)
}

func (t *SupplyContract) moveProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("call moveProduct")
	from := args[0]
	to := args[1]
	amount, _ := strconv.Atoi(args[2])
	date := time.Now().Format("2006-01-02 15:04:05")

	FromBytes, err := stub.GetState(from)
	if err != nil {
		return shim.Error("Failed to get from orderId:" + err.Error())
	} else if FromBytes == nil {
		return shim.Error("Geterr: Data does not exist")
	}
	toBytes, err := stub.GetState(to)
	if err != nil {
		return shim.Error("Failed to get to orderId:" + err.Error())
	} else if toBytes == nil {
		return shim.Error("Geterr: Data does not exist")
	}

	fromData := product{}
	err = json.Unmarshal(FromBytes, &fromData)
	if err != nil {
		return shim.Error(err.Error())
	}
	toData := product{}
	err = json.Unmarshal(toBytes, &toData)
	if err != nil {
		return shim.Error(err.Error())
	}
	if fromData.Amount >= amount {
		fromData.Amount = fromData.Amount - amount
		toData.Amount = toData.Amount + amount
		fromData.Date = date
		toData.Date = date
	} else {
		return shim.Error("not enough products")
	}

	fromDataBytes, _ := json.Marshal(fromData)
	err = stub.PutState(from, fromDataBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	toDataBytes, _ := json.Marshal(toData)
	err = stub.PutState(to, toDataBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SupplyContract) useProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("call useProduct")
	from := args[0]
	amount, _ := strconv.Atoi(args[1])
	date := time.Now().Format("2006-01-02 15:04:05")

	FromBytes, err := stub.GetState(from)
	if err != nil {
		return shim.Error("Failed to get from orderId:" + err.Error())
	} else if FromBytes == nil {
		return shim.Error("Geterr: Data does not exist")
	}

	fromData := product{}
	err = json.Unmarshal(FromBytes, &fromData)
	if err != nil {
		return shim.Error(err.Error())
	}

	if fromData.Amount >= amount {
		fromData.Amount = fromData.Amount - amount
		fromData.Date = date
	} else {
		return shim.Error("not enough products")
	}

	fromDataBytes, _ := json.Marshal(fromData)
	err = stub.PutState(from, fromDataBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (f *SupplyContract) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expected ENIITY Name")
	}

	ENIITY := args[0]
	Avalbytes, err := stub.GetState(ENIITY)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + ENIITY + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil order for " + ENIITY + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(Avalbytes)
}

func (f *SupplyContract) history(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expected KEY Name")
	}
	key := args[0]
	history, _ := stub.GetHistoryForKey(key)
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for history.HasNext() {
		modification, err := history.Next()
		if err != nil {
			fmt.Println(err.Error())
			return shim.Error(err.Error())
		}
		buffer.WriteString("#")
		buffer.WriteString(string(modification.Value))
		buffer.WriteString("#")
	}
	buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

func main() {

	err := shim.Start(new(SupplyContract))
	if err != nil {
		fmt.Printf("Error creating new product Contract: %s", err)
	}
}
