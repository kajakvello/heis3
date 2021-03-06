package Cost

import (
	"strconv"
	."./../Driver"
	."./../OrderRegister"
)

type ElevStatus struct{
	LastFloor int
	Direction int
	Up [N_FLOORS]bool
	Down [N_FLOORS]bool
	Inside [N_FLOORS]bool
	DoorOpen bool
	Defect bool
	OnNet bool
}

var Elevators = make(map[string]ElevStatus)



/////////////////////////////////////////////////////////////////////////////////////////////



//Calculate cost from existing orders, direction and floor. 
//The cost is multiplied by 1000 and the IP address is added to get a unike cost for each elevator.
//Returns 1 if my elevator has got the lowest cost.
func GetCost(myFloor, myDirection, orderFloor, orderDirection int, myAddress string) int {
	
	//Find my cost:
	myCost := selectCostCase(myFloor, orderFloor, myDirection, orderDirection, Up, Down, Inside)
	
	myIP, _ := strconv.Atoi(myAddress)
	myCost = (myCost*1000) + myIP
	
	//Check if other elevator got lower cost:
	for IP, elev := range Elevators {
		if (!elev.Defect && elev.OnNet) {
			elevCost := selectCostCase(elev.LastFloor, orderFloor, elev.Direction, orderDirection, elev.Up, elev.Down, elev.Inside)
		
			elevIP, _ := strconv.Atoi(IP)
			elevCost = (elevCost*1000) + elevIP
			
			if elevCost < myCost {
				return 0
			}
		}
	}

	return 1

}



/////////////////////////////////////////////////////////////////////////////////////////////



//Selects if the cost should loop over all, some or none of the exixting orders.
func selectCostCase(myFloor, orderFloor, myDirection, orderDirection int, up, down, inside [N_FLOORS]bool) int {
	cost := 0
	
	
	switch {
	
	case orderFloor > myFloor:
	
		if myDirection == -1 {
		cost += calculateCost(false, false, orderFloor, myFloor, 1, up, down, inside)
		
		} else if orderDirection == 1 {
			if myDirection == 1 {
				cost += calculateCost(false, false, orderFloor, myFloor, 1, up, down, inside)
			} else if myDirection == 0 {
				cost += calculateCost(false, true, orderFloor, myFloor, 1, up, down, inside) 
			}
		} else if orderDirection == 0 {
			if myDirection == 1 {
				cost += calculateCost(true, false, orderFloor, myFloor, 1, up, down, inside)
			} else if myDirection == 0 {
				cost += calculateCost(true, true, orderFloor, myFloor, 1, up, down, inside) 
			}
		}
	
	case orderFloor < myFloor:
	
		if myDirection == -1 {
		cost += calculateCost(false, false, myFloor, orderFloor, 0, up, down, inside)
		
		} else if orderDirection == 1 {
			if myDirection == 1 {
				cost += calculateCost(true, true, myFloor, orderFloor, 0, up, down, inside)
			} else if myDirection == 0 {
				cost += calculateCost(false, true, myFloor, orderFloor, 0, up, down, inside) 
			}
		} else if orderDirection == 0 {
			if myDirection == 1 {
				cost += calculateCost(true, false, myFloor, orderFloor, 0, up, down, inside)
			} else if myDirection == 0 {
				cost += calculateCost(false, false, myFloor, orderFloor, 0, up, down, inside) 
			}
		}
	
	case orderFloor == myFloor:
	
		if myDirection == orderDirection {
			cost += calculateCost(true, true, orderFloor, myFloor, myDirection, up, down, inside)
		} else {
			if myDirection == 0 {
				cost += calculateCost(false, true, orderFloor, myFloor, 1, up, down, inside)
			} else if myDirection == 1 {
				cost += calculateCost(true, false, orderFloor, myFloor, 0, up, down, inside)
			}
		}
	}
	return cost
}



/////////////////////////////////////////////////////////////////////////////////////////////



//Loops either over all, some or none of the existing orders and applies to cost.
//One order costs 4 points and one floor in difference costs 1 point.
func calculateCost(loopUp, loopDown bool, highestFloor, lowestFloor, pos int, up, down, inside [N_FLOORS]bool) int {
	
	cost := 0
	highestOrder := -1
	lowestOrder := N_FLOORS
	
	if highestFloor != lowestFloor {
		cost += highestFloor - lowestFloor
	
		if pos == 1 {
			for i:=lowestFloor; i<highestFloor; i++ {
				if (up[i] || inside[i]) {
					cost += 4
				}
			}
		} else if pos == 0 {
			for i:=highestFloor; i<lowestFloor; i-- {
				if (down[i] || inside[i]) {
					cost += 4
				}
			}
		}
	}
	
	if loopUp {
		for i:=highestFloor; i<N_FLOORS; i++ {
			if inside[i] {
				cost += 4
				if i > highestOrder {highestOrder = i}
			} else {
				if up[i] {
					cost += 4
					if i > highestOrder {highestOrder = i}
				}
				if (down[i] && i != highestFloor) {
					cost += 4
					if i > highestOrder {highestOrder = i}
				}
			}
		}
		if highestOrder != -1 {
			cost += (highestOrder - highestFloor)*2
		}
	}
	if loopDown {
		for i:=lowestFloor; i>=0; i-- {
			if inside[i] {
				cost += 4
				if i < lowestOrder {lowestOrder = i}
			} else {
				if (up[i] && i != lowestFloor) {
					cost += 4
					if i < lowestOrder {lowestOrder = i}
				}
				if down[i] {
					cost += 4
					if i < lowestOrder {lowestOrder = i}
				}
			}
		}
		if lowestOrder != N_FLOORS {
			cost += (lowestFloor - lowestOrder)*2
		}
	}
	return cost
}




