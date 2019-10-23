package main

import (
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
	OrderId                string
	ProductId              string
	ConsumerId             string
	ManufactureId          string
	WholesalerId           string
	RetailerId             string
	LogisticsId            string
	Status                 string
	MaterialProcessDate     string
	ManufactureProcessDate string
	WholesaleProcessDate   string
	ShippingProcessDate    string
	RetailProcessDate      string
	OrderPrice             int
	ShippingPrice          int
	DeliveryDate           string
}

func (t *SupplyContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return setupProductSupplyChainOrder(stub)
}

func (t *SupplyContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createMaterial" {
		return t.createMaterial(stub, args)
	} else if function == "manufactureProcessing" {

		return t.manufactureProcessing(stub, args)
	} else if function == "wholesalerDistribute" {

		return t.wholesalerDistribute(stub, args)
	} else if function == "initiateShipment" {

		return t.initiateShipment(stub, args)
	} else if function == "deliverToRetail" {

		return t.deliverToRetail(stub, args)
	} else if function == "completeOrder" {

		return t.completeOrder(stub, args)
	} else if function == "query" {

		return t.query(stub, args)
	}

	return shim.Error("Invalid function name")
}

func setupProductSupplyChainOrder(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	orderId := args[0]
	consumerId := args[1]
	orderPrice, _ := strconv.Atoi(args[2])
	shippingPrice, _ := strconv.Atoi(args[3])
	SupplyContract := product{
		OrderId:       orderId,
		ConsumerId:    consumerId,
		OrderPrice:    orderPrice,
		ShippingPrice: shippingPrice,
		Status:        "order initiated"}

	productBytes, _ := json.Marshal(SupplyContract)
	stub.PutState(SupplyContract.OrderId, productBytes)

	return shim.Success(nil)
}

func (f *SupplyContract) createMaterial(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	productBytes, _ := stub.GetState(orderId)
	fd := product{}
	json.Unmarshal(productBytes, &fd)

	if fd.Status == "order initiated" {
		fd.ProductId = "PRODUCT_1"
		currentts := time.Now()
		fd.MaterialProcessDate = currentts.Format("2019-01-02 15:04:05")
		fd.Status = "raw product created"
	} else {
		fmt.Printf("Order not initiated yet")
	}

	productBytes, _ = json.Marshal(fd)
	stub.PutState(orderId, productBytes)

	return shim.Success(nil)
}

func (f *SupplyContract) manufactureProcessing(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	productBytes, err := stub.GetState(orderId)
	fd := product{}
	err = json.Unmarshal(productBytes, &fd)
	if err != nil {
		return shim.Error(err.Error())
	}

	if fd.Status == "raw product created" {
		fd.ManufactureId = "Manufacture_1"
		currentts := time.Now()
		fd.ManufactureProcessDate = currentts.Format("2006-01-02 15:04:05")
		fd.Status = "manufacture Process"
	} else {
		fd.Status = "Error"
		fmt.Printf("Raw product not initiated yet")
	}

	productBytes0, _ := json.Marshal(fd)
	err = stub.PutState(orderId, productBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
func (f *SupplyContract) wholesalerDistribute(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	productBytes, err := stub.GetState(orderId)
	fd := product{}
	err = json.Unmarshal(productBytes, &fd)
	if err != nil {
		return shim.Error(err.Error())
	}

	if fd.Status == "manufacture Process" {
		fd.WholesalerId = "Wholesaler_1"
		currentts := time.Now()
		fd.WholesaleProcessDate = currentts.Format("2019-01-02 15:04:05")
		fd.Status = "wholesaler distribute"
	} else {
		fd.Status = "Error"
		fmt.Printf("Manufacture not initiated yet")
	}

	productBytes0, _ := json.Marshal(fd)
	err = stub.PutState(orderId, productBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (f *SupplyContract) initiateShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	productBytes, err := stub.GetState(orderId)
	fd := product{}
	err = json.Unmarshal(productBytes, &fd)
	if err != nil {
		return shim.Error(err.Error())
	}

	if fd.Status == "wholesaler distribute" {
		fd.LogisticsId = "LogisticsId_1"
		currentts := time.Now()
		fd.ShippingProcessDate = currentts.Format("2006-01-02 15:04:05")
		fd.Status = "initiated shipment"
	} else {
		fmt.Printf("Wholesaler not initiated yet")
	}

	productBytes0, _ := json.Marshal(fd)
	err = stub.PutState(orderId, productBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (f *SupplyContract) deliverToRetail(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	productBytes, err := stub.GetState(orderId)
	fd := product{}
	err = json.Unmarshal(productBytes, &fd)
	if err != nil {
		return shim.Error(err.Error())
	}

	if fd.Status == "initiated shipment" {
		fd.RetailerId = "Retailer_1"
		currentts := time.Now()
		fd.RetailProcessDate = currentts.Format("2019-01-02 15:04:05")
		fd.Status = "Retailer started"

	} else {
		fmt.Printf("Shipment not initiated yet")
	}

	productBytes0, _ := json.Marshal(fd)
	err = stub.PutState(orderId, productBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (f *SupplyContract) completeOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	orderId := args[0]
	productBytes, err := stub.GetState(orderId)
	fd := product{}
	err = json.Unmarshal(productBytes, &fd)
	if err != nil {
		return shim.Error(err.Error())
	}

	if fd.Status == "Retailer started" {
		currentts := time.Now()
		fd.DeliveryDate = currentts.Format("2019-01-02 15:04:05")
		fd.Status = "Consumer received order"
	} else {
		fmt.Printf("Retailer not initiated yet")
	}

	productBytes0, _ := json.Marshal(fd)
	err = stub.PutState(orderId, productBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (f *SupplyContract) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ENIITY string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expected ENIITY Name")
	}

	ENIITY = args[0]
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

func main() {

	err := shim.Start(new(SupplyContract))
	if err != nil {
		fmt.Printf("Error creating new product Contract: %s", err)
	}
}
