package Elevator

import (
	."./../Driver"
	."./../Udp"
	"encoding/json"
	//."fmt"
)



//Lager en liste for antall heiser, men adresse, posisjon og retning.
//Om en heis faller ut (vi ikke får meldinger etter viss tid), settes floor og direction til -2.
//Om en bestilling ikke har blitt tatt, regner gjenværende heiser ut ny cost seg i mellom.



var GlobalUp [N_FLOORS]bool
var GlobalDown [N_FLOORS]bool

// My inside orders
var Inside [N_FLOORS]bool


// My up orders
var Up [N_FLOORS]bool


// My down orders
var Down [N_FLOORS]bool


var Receive_ch = make(chan Udp_message)
var Send_ch = make(chan Udp_message)

var MyFloor = -1
var LastFloor = 0
var MyDirection = -1	// 1 = UP, 0 = DOWN, -1 = stands still
var MyAddress string
var Defect bool

var DoorOpen = false
var OpenDoor = make(chan int)
var GotMessage = make(chan string)
var Alive = make(chan string)
var NotOnNet = make(chan string)
//var orderHandled = make(chan int)


type Order struct {

	//My posistion:
	MyFloor int
	MyDirection int
	//Order:
	Floor int				
	Direction int 		// 1 = UP, 0 = DOWN, -1 = stands still
	OrderHandled bool 
	NewOrder bool
	DoorOpen bool
	//My orders:
	Up [N_FLOORS]bool
	Down [N_FLOORS]bool
	Inside [N_FLOORS]bool
		
}



/////////////////////////////////////////////////////////////////////////////////////////////



func UpdateMyOrders(receivedOrder Order) {

	if receivedOrder.OrderHandled {
	
		Inside[receivedOrder.Floor] = false
		if receivedOrder.Direction == 1 {
			Up[receivedOrder.Floor] = false
		} else if receivedOrder.Direction == 0 {
			Down[receivedOrder.Floor] = false
		} else if receivedOrder.Direction == -1 {
			if MyDirection == 1 || receivedOrder.Floor == 0 {
				Up[receivedOrder.Floor] = false
			} 
			if MyDirection == 0 || receivedOrder.Floor == N_FLOORS -1  {
				Down[receivedOrder.Floor] = false
			}
		}
		
		
	} else if receivedOrder.NewOrder {
	
		if receivedOrder.Direction == 0 {
			Down[receivedOrder.Floor] = true
			
		} else if receivedOrder.Direction == 1 {
			Up[receivedOrder.Floor] = true
			
		} else if receivedOrder.Direction == -1 {
			Inside[receivedOrder.Floor] = true
			Elev_set_button_lamp(BUTTON_COMMAND, receivedOrder.Floor, 1)
			
		} else {
			println("Unvalid floor or direction")
		}	
		
	} else {
		println("Error in UpdateMyOrders")
	}

}



/////////////////////////////////////////////////////////////////////////////////////////////



func UpdateGlobalOrders(order Order){

	if order.NewOrder {
		if order.Direction == 1 {
			GlobalUp[order.Floor] = true
		} else if order.Direction == 0 {
			GlobalDown[order.Floor] = true
		}
	} else if order.OrderHandled {
		if order.Direction == 1 {
			GlobalUp[order.Floor] = false
		} else if order.Direction == 0 {
			GlobalDown[order.Floor] = false
		}
	} else {
		println("Error in UpdateGlobalOrders")
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



func SetButtonLight(order Order, IP string) {
	
	if order.NewOrder && order.Direction == 0 {
		Elev_set_button_lamp(BUTTON_CALL_DOWN, order.Floor, 1)
		
	} else if order.NewOrder && order.Direction == 1 {
		Elev_set_button_lamp(BUTTON_CALL_UP, order.Floor, 1)
		
	} else if order.OrderHandled {
		
		if IP == MyAddress {
			Elev_set_button_lamp(BUTTON_COMMAND, order.Floor, 0)
		}
		
		if (order.Direction == 1 && order.Floor != N_FLOORS-1) || (order.Floor == 0) {
			Elev_set_button_lamp(BUTTON_CALL_UP, order.Floor, 0)
			
		} else if (order.Direction == 0 && order.Floor != 0) || (order.Floor == N_FLOORS-1){
			Elev_set_button_lamp(BUTTON_CALL_DOWN, order.Floor, 0)
			
		} else if order.Direction == -1 {
			if MyDirection == 1 {
				Elev_set_button_lamp(BUTTON_CALL_UP, order.Floor, 0)
			} else if MyDirection == 0 {
				Elev_set_button_lamp(BUTTON_CALL_DOWN, order.Floor, 0)
			}
		}
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



// Returns the direction if the elevator should take an order from "floor". 
// Returns two if no orders.
func GetOrder(direction int, floor int) int {
	
	if Inside[floor] {
		return -1
	}
	if Up[floor] && Down[floor] {
		return direction
	}
	if Up[floor] && (direction == 1 || direction == -1 || floor == 0 || !CheckOrdersUnderFloor(floor)) {
		return 1
	}
	if Down[floor] && (direction == 0 || direction == -1 || floor == N_FLOORS-1 || !CheckOrdersAboveFloor(floor)) {
		return 0
	}
	
	return 2
}



/////////////////////////////////////////////////////////////////////////////////////////////



func CheckOrdersUnderFloor(floor int) bool {
	for i:=0; i<floor; i++ {
		if (Up[i] || Down[i] || Inside[i]) {
			return true
		}
	}
	return false
}



/////////////////////////////////////////////////////////////////////////////////////////////



func CheckOrdersAboveFloor(floor int) bool {
	for i:=floor+1; i<N_FLOORS; i++ {
		if (Up[i] || Down[i] || Inside[i]) {
			return true
		}
	}
	return false
}



/////////////////////////////////////////////////////////////////////////////////////////////



func EmptyQueue() bool {
	for i:=0; i<N_FLOORS; i++ {
		if (Up[i] || Down[i] || Inside[i]) {
			return false
		}
	}
	return true
}



/////////////////////////////////////////////////////////////////////////////////////////////



//Makes struct into byte array and sends it on channel
func SendOrder(order Order) {
	b, err := json.Marshal(order)
	
	if (err != nil) {
		println("Send Order Error: ", err)
	}
	
	var message Udp_message
	message.Raddr = "broadcast"
	message.Data = b
	message.Length = 1024
	
	Send_ch <- message
}



/////////////////////////////////////////////////////////////////////////////////////////////



func SetDirectionToOrder(orderDir int) {
	
	if (EmptyQueue()) {
		MyDirection = -1
		
	} else if GetOrder(orderDir, LastFloor) == orderDir {
		if orderDir == 0 && !CheckOrdersAboveFloor(LastFloor) {
			MyDirection = 0
			OpenDoor <- 0
			
		} else if orderDir == 1 && !CheckOrdersUnderFloor(LastFloor) {
			MyDirection = 1
			OpenDoor <- 1
		}
	} else {
		if (MyDirection == 0) && !(CheckOrdersUnderFloor(LastFloor)) {
			MyDirection = 1
		} else if (MyDirection == 1) && !(CheckOrdersAboveFloor(LastFloor)) {
			MyDirection = 0
		} else if MyDirection == -1 {
			if CheckOrdersAboveFloor(LastFloor) {
				MyDirection = 1
			} else if CheckOrdersUnderFloor(LastFloor) {
				MyDirection = 0
			}
		}
	}
}








