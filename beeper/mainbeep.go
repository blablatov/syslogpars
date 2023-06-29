package mainbeep

import (
	"beeper"
	"sync"
)

func MainBeep() {
	// beep once. Подать звуковой сигнал один раз
	//beeper.Beep()

	// beep three times. Звуковой сигнал три раза
	//beeper.Beep(3)

	// beep, beep, pause, pause, beep, pause, pause, etc
	// Мелодия в цикле (*бипер, -пауза)
	beeper.Melody("**--**--**--**")
}
