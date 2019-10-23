package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SupplyChaincode struct {
}

type ValInfo struct {
	TransactionId string `json:"transactionId"`
	ProductId     string `json:"productId"`
	From          string `json:"from"`
	To            string `json:"to"`
	Amount        int    `json:"amount,string"`
	Description   string `json:"description"`
	Detail        string `json:"detail"`
	subTid        string `json:"subTid"`
}

func (t *SupplyChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Supply-ChainCode Init!!")
	return shim.Success(nil)
}

func (t *SupplyChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("supply Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "makeTransaction" {
		// Make transaction Id and data
		return t.makeTransaction(stub, args)
	} else if function == "queryAll" {
		// the old "Query" is now implemtned in invoke
		return t.queryAll(stub, args)
	} else if function == "moveProduct" {
		// Move val djdfrom A to B or B to A
		return t.moveProduct(stub, args)
	} else if function == "queryById" {
		// query info by transaction Id
		return t.queryById(stub, args)
	}

	return shim.Error("Invalid invoke function name.")
}

// make product with this function
func (t *SupplyChaincode) makeTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	transactionId := args[0]
	productId := args[1]
	from := args[2]
	to := args[3]
	amount, _ := strconv.Atoi(args[4])
	description := args[5]
	detail := args[6]
	subTid := args[7]

	fmt.Println("log>Input Transaction value : " + transactionId)

	valInfo := &ValInfo{transactionId, productId, from, to, amount, description, detail, subTid}
	valInfoBytes, err := json.Marshal(valInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(transactionId, valInfoBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("#func makeTransaction >> putState complete")
	return shim.Success(nil)

}

func (t *SupplyChaincode) moveProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("CALL >> moveProduct[Chaincode]")
	from := args[0]
	to := args[1]
	pid := args[2]
	fromTid := from + "_" + pid
	toTid := to + "_" + pid
	X, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}

	FromBytes, err := stub.GetState(fromTid)
	if err != nil {
		return shim.Error("Failed to get from transactionId:" + err.Error())
	} else if FromBytes == nil {
		return shim.Error("Data does not exist")
	}

	toBytes, err := stub.GetState(toTid)
	if err != nil {
		return shim.Error("Failed to get to transactionId:" + err.Error())
	} else if toBytes == nil {
		return shim.Error("Data does not exist")
	}

	fromVal := ValInfo{}
	err = json.Unmarshal(FromBytes, &fromVal)
	if err != nil {
		return shim.Error(err.Error())
	}

	toVal := ValInfo{}
	err = json.Unmarshal(toBytes, &toVal)
	if err != nil {
		return shim.Error(err.Error())
	}

	fromVal.Amount = fromVal.Amount - X
	toVal.Amount = toVal.Amount + X

	fromValBytes, _ := json.Marshal(fromVal)
	err = stub.PutState(fromTid, fromValBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	toValBytes, _ := json.Marshal(toVal)
	err = stub.PutState(toTid, toValBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

// query all
func (t *SupplyChaincode) queryAll(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("call query method")
	queryString := "{\"selector\":{}}"
	fmt.Println("queryString" + queryString)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults

	fmt.Println("call constructQueryResponseFromIterator")
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		fmt.Println("what")
		fmt.Println(string(queryResponse.Value))
		fmt.Println("what")

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

func (t *SupplyChaincode) queryById(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("queryById 함수호출")
	id := args[0]
	queryString := fmt.Sprintf("{\"selector\":{\"id\":\"%s\"}}", id)
	queryResults, err := getQueryResultForQueryString(stub, queryString)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)

}

func main() {
	err := shim.Start(new(SupplyChaincode))
	if err != nil {
		fmt.Printf("Error starting Supply chaincode: %s", err)
	}
}
