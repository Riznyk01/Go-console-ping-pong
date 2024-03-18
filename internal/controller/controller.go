package controller

import "github.com/eiannone/keyboard"

func Controller(keyEventsChan chan keyboard.KeyEvent) {
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		if event.Key == keyboard.KeyEsc {
			break
		}
		keyEventsChan <- event
	}
}
