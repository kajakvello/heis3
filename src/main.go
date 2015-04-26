package main

import(
	."./Elevator"
	."./Timer"
	."./Driver"
)


// HAR FIKSET:
// - Nappe ut og inn nettverkskabler funker, woho :D
// - Knappelys er fikset. Ordre må slettes når døren åpnes, hvis ikke blir alt kaos. Brukte 1000 timer på å prøve å fikse det.
// - AliveTimer merker at både seg selv og andre heiser mister strøm og blir defekte.
// - Order-structen ligger i OrderRegister nå. Blir ryddigere og funker helt fint med importering.
// - Lage egen break-funksjon
// - Tror jeg har fikset cost - om alle heiser hadde samme Lastfloor, men var på vei til en annen bestilling ville begge få cost = 0. Fikset!
// - Har laget globalorderTimer, men har ikke testet. Må finne ut hvor timer skal settes og oppdateres
// - Har ikke hatt problem med init


//TODO:
// - AliveTimers funker nesten, men heisene krangler om cost når disse kjører. Må fikse cost.
// - Test kode i RunElevator() for å få Defekt heis tilbake i systemet
// - GlobalUp, GlobalDown og UpdateGlobalOrders ligger i orderRegister, men brukes ikke til noe ennå. Bør sjekkes med jevne mellomrom slik at ALLE ordre håndteres
// - Init() kjøres ikke alltid ordentlig før programmet termineres. Er dette Stop() sin feil?
// - Heis som blir koblet fra nettverket: Andre heiser tar bestillingene, men heisen kommer seg ikke tilbake på nettet, og fungerer ikke normalt lenger



//Begynn på så lite nytt som mulig. Fikser vi bugsene vi har nå blir heisen supersmud!!




func main() {
	
	//Initialiser heis
	Init()
	
	

	go CheckButtonCallUp()
	go CheckButtonCallDown()
	go CheckButtonCommand()
	go RunElevator()
	go UpdateFloor()
	
	
	go DoorControl()

	
	go ReceiveMessage()
	go SendUpdateMessage()
	
	go PrintStatus()
	
	
	

	s := make(chan int)
	go Stop(s)
	
	select {			// Programmet sluttet å fucke Init() da denne ble byttet ut med tom select{}
						//Kanksje heisen ikke funker om Stop allerede er aktivert av tidligere folk som har kjørt heisen??
	case <- s:
		Elev_set_motor_direction(0)
		break
	} 
	
	
}






