package main

import(
	."./Elevator"
	."./Timer"
	."./Driver"
)



func main() {
	
	Init()
	

	go CheckButtonCallUp()
	go CheckButtonCallDown()
	go CheckButtonCommand()
	go RunElevator()
	go UpdateFloor()
	go DoorControl()	
	go ReceiveMessage()
	go SendUpdateMessage()
	
	
	
	s := make(chan int)
	go Stop(s)
	
	select {
	case <- s:
		Elev_set_motor_direction(0)
		break
	}
	
	
}







