package Timer


import (
	."./../Driver"
	."./../Cost"
	."./../OrderRegister"
	."time"
	//."fmt"
	
)




//Opens door for three seconds. Deletes light and order when doors open.
func DoorControl() {

	timer := NewTimer(Hour*3)
	orderDir := 2
	for {
	
		select {
		case orderDir = <- OpenDoor:
			DoorOpen = true
			Elev_set_door_open_lamp(1)
			timer.Reset(Second*3)
			if Elev_get_floor_sensor_signal() == LastFloor {
				deleteOrder := Order{LastFloor, MyDirection, LastFloor, orderDir, true, false, DoorOpen, Up, Down, Inside}
				SendOrder(deleteOrder)
			}
			
		case <- timer.C:
			Elev_set_door_open_lamp(0)
			DoorOpen = false
			
			if orderDir == 1{
				orderDir = 0
			} else if orderDir == 0 {
				orderDir = 1
			}
			go SetDirectionToOrder(orderDir)
		}
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



//Sets timer for each elevator to know that their on net. Deletes elevator from directory if no message received in 3 seconds
func MessageTimer(address string) {
	
	timer := NewTimer(3*Hour)
	for {
		select {
		
		case IP := <- GotMessage:
			if IP == address {
				timer.Reset(3*Second)
			}
		
		case <- timer.C:
			println("Elevator nr ", address, " is not on net")
			NotOnNet <- address
			
			temp := Elevators[address]
			temp.OnNet = false
			Elevators[address] = temp
			
			for i:=0; i<N_FLOORS; i++ {
				if (Elevators[address].Up)[i] {
					order := Order{LastFloor, MyDirection, i, 1, false, true, DoorOpen, Up, Down, Inside}
					go SendOrder(order)
				}
				if (Elevators[address].Down)[i] {
					order := Order{LastFloor, MyDirection, i, 0, false, true, DoorOpen, Up, Down, Inside}
					go SendOrder(order)
				}
			}
			delete (Elevators, address)
			return
			
		}	
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////


/*
func OrderTimer(order int) {

	timer := NewTimer(10*Second)
	
	if order.Direction == 1 {
		floor = order.Floor
	} else if order.Direction == 0 {
		floor = N_FLOORS-2+order.Floor
	}
	
	for {
		select {
		
		case handledFloor := <- orderHandled:
			if handledFloor == floor {
				return
			}
		}
		
		case <- timer.C:
			if order.Direction == 1 {
				Up[order.Floor] = false
			} else if order.Direction == 0 {
				Down[order.Floor] = false
			}
			newOrder := Order{floor, direction, floor, direction, false, true, DoorOpen, Up, Down, Inside}
			go SendOrder(newOrder)
			timer.Reset(10*Second)
			return
			
	}


}
*/


/////////////////////////////////////////////////////////////////////////////////////////////



//Will register if any elevator do not take its orders and set elevator to defect.
//Sends orders to other elevators if timer runs out. Deletes all outside orders and sets one order true to check if its running again
func AliveTimer(address string) {

	timer := NewTimer(3*Hour)
	oldUp := [N_FLOORS]bool{}
	oldDown := [N_FLOORS]bool{}
	oldInside := [N_FLOORS]bool{}
	
	for {	
		select {
		case IP := <- Alive: 

			if IP == address {
				temp := Elevators[IP]
				
				if Elevators[IP].Direction == -1 || Elevators[IP].DoorOpen {
					timer.Reset(10*Second)
					
				} else {
					for i:=0; i<N_FLOORS; i++ {
						if (oldUp[i] && !(Elevators[IP].Up)[i]) || (oldDown[i] && !(Elevators[IP].Down)[i]) || (oldInside[i] && !(Elevators[IP].Inside)[i]) {
							timer.Reset(10*Second)
							temp.Defect = false
						}
						oldUp[i] = temp.Up[i]
						oldDown[i] = temp.Down[i]
						oldInside[i] = temp.Inside[i]
					}
					Elevators[IP] = temp
				}
			}
		case <- timer.C:
			
			println("Elevator nr ", address, " is defect")
			temp := Elevators[address]
			temp.Defect = true
			Elevators[address] = temp
			
			for i:=0; i<N_FLOORS; i++ {
				if (Elevators[address].Up)[i] {
					order := Order{0, -1, i, 1, false, true, false, Up, Down, Inside}
					go SendOrder(order)
				}
				if (Elevators[address].Down)[i] {
					order := Order{0, -1, i, 0, false, true, false, Up, Down, Inside}
					go SendOrder(order)
				}
				oldUp[i] = false
				oldDown[i] = false
				temp.Up[i] = false
				temp.Down[i] = false
			}
			temp.Up[0] = true
			Elevators[address] = temp
		
		case IP := <- NotOnNet:
			if IP == address {
				timer.Stop()
				return
			}
		}
	}
}



/////////////////////////////////////////////////////////////////////////////////////////////



//Registers if my elevator is defect, and sends my orders to the other elevators
func SelfAliveTimer() {

	timer := NewTimer(3*Hour)
	oldUp := [N_FLOORS]bool{}
	oldDown := [N_FLOORS]bool{}
	oldInside := [N_FLOORS]bool{}
	
	for {
		select {
		case IP := <- Alive:

			if IP == MyAddress {
				
				if MyDirection == -1 || DoorOpen {
					timer.Reset(10*Second)
					
				} else {
					for i:=0; i<N_FLOORS; i++ {
						if (oldUp[i] && !Up[i]) || (oldDown[i] && !Down[i]) || (oldInside[i] && !Inside[i]) {
							timer.Reset(10*Second)
							Defect = false
						}
						oldUp[i] = Up[i]
						oldDown[i] = Down[i]
						oldInside[i] = Inside[i]
					}
				}
			}
		case <- timer.C:
		
			Defect = true
			println("IM DEFECT")
			
			for i:=0; i<N_FLOORS; i++ {
				oldUp[i] = false
				oldDown[i] = false
				Up[i] = false
				Down[i] = false
			}
			Up[0] = true
			
		}
	}
}


