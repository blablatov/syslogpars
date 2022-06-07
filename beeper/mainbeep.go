package mainbeep

import (
	"beeper"
	"sync"
)

func MainBeep(wg sync.WaitGroup) {
	defer wg.Done()
	// beep once
	//beeper.Beep()

	// beep three times
	//beeper.Beep(3)

	beeper.Melody("**--**--**--**")
	// beep, beep, pause, pause, beep, beep
}
