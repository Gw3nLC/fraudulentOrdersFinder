package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jszwec/csvutil"
)

type Order struct {
	OrderID    int
	DealID     int
	Email      string
	Street     string
	City       string
	State      string
	ZipCode    int
	CreditCard int
}

// List of fraudulent order IDs to be outputed at the end
var fraudulentOrderIds []int

// List of orders extracted from the csv input file
var orders []Order

func main() {

	/// read input file
	path := os.Args[1]
	inputFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// define csv header
	orderHeader, err := csvutil.Header(Order{}, "csv")
	if err != nil {
		log.Fatal(err)
	}

	// Build the order list from data in inputFile
	csvReader := csv.NewReader((strings.NewReader(string(inputFile[1:]))))
	dec, err := csvutil.NewDecoder(csvReader, orderHeader...)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var ord Order
		if err := dec.Decode(&ord); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, ord)
	}

	// Find fraudulent orders from the list
	findFraudulentOrders(orders)

	// print the list of fraudulent orders
	fmt.Println("List of fraudulent orders (OrderID) :")
	fmt.Println(removeDuplicateInt(fraudulentOrderIds))
}

// findFraudulentOrders loops on every order 2 by 2 and add the fraudulent ones to the fraudulentOrderIds list
func findFraudulentOrders(orders []Order) {
	for i, order1 := range orders {
		for _, order2 := range orders[i+1:] {
			if order1.checkEmail(order2) {
				// fraudulent orders, add them to the list
				fraudulentOrderIds = append(fraudulentOrderIds, order1.OrderID, order2.OrderID)
				// no need to check for adress, break
				break
			}
			if order1.CheckAdress(order2) {
				// fraudulent orders, add them to the list
				fraudulentOrderIds = append(fraudulentOrderIds, order1.OrderID, order2.OrderID)
			}
		}
	}
}

// Two orders have the same email address and deal id, but different credit card information, regardless of street address.
func (order1 Order) checkEmail(order2 Order) bool {

	// TODO handle ignored caracters + . etc with regex before compare
	if strings.EqualFold(order1.Email, order2.Email) && order1.DealID == order2.DealID && order1.CreditCard != order2.CreditCard {
		return true
	}
	return false
}

// Two orders have the same Address/City/State/Zip and deal id, but different credit card information, regardless of email address.
func (order1 Order) CheckAdress(order2 Order) bool {
	if changeStreetAbrev(order1.Street) == changeStreetAbrev(order2.Street) && order1.City == order2.City && changeStateAbrev(order1.State) == changeStateAbrev(order2.State) && order1.ZipCode == order2.ZipCode && order1.DealID == order2.DealID && order1.CreditCard != order2.CreditCard {
		return true
	}
	return false
}

// Replace location abreviation to full name
func changeStreetAbrev(adress string) string {
	streetChanged := strings.Replace(adress, "St.", "Street", -1)
	s := strings.Replace(streetChanged, "Rd.", "Road", -1)
	return s
}

// Replace abreviation of state to full name
func changeStateAbrev(state string) string {
	if state == "IL" {
		return "Illinois"
	}
	if state == "CA" {
		return "California"
	}
	if state == "NY" {
		return "New York"
	}
	return state
}

func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
