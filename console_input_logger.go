package coninlogger

import (
	"github.com/jeet-parekh/winapi"
)

type (
	// ConsoleInputLogger provides an interface to receive ConsoleInputEvents.
	ConsoleInputLogger struct {
		handle uintptr
		messages chan *ConsoleInputEvent
		oldMode uintptr
		captureMode uintptr
		enabled bool
	}

	// ConsoleInputEvent contains information about a single console input event.
	ConsoleInputEvent struct {
		KeyDown bool
		RepeatCount uint16
		VirtualKeyCode uint16
		VirtualScanCode uint16
		UnicodeCharacter uint16
		ControlKeyState uint32
	}
)

// NewConsoleInputLogger creates and returns a ConsoleInputLogger.
// The parameter is the buffer size of the channel through which the messages would be sent.
func NewConsoleInputLogger(bufferSize uint) *ConsoleInputLogger {
	cil := &ConsoleInputLogger {
		handle : getStdHandle(winapi.STD_INPUT_HANDLE),
		messages : make(chan *ConsoleInputEvent, bufferSize),
		captureMode : winapi.ENABLE_WINDOW_INPUT | winapi.ENABLE_MOUSE_INPUT | winapi.ENABLE_PROCESSED_INPUT,
	}
	cil.oldMode = getConsoleMode(cil.handle)
	return cil
}

// GetMessageChannel returns the channel through which the messages would be sent.
func (cil *ConsoleInputLogger) GetMessageChannel() chan *ConsoleInputEvent {
	return cil.messages
}

// Start the ConsoleInputLogger.
func (cil *ConsoleInputLogger) Start() {
	cil.enabled = true
	setConsoleMode(cil.handle, cil.captureMode)
	go func() {
		for {
			inputRecord := []winapi.INPUT_RECORD_KEY {
				winapi.INPUT_RECORD_KEY {
					EventType : uint16(winapi.KEY_EVENT),
				},
			}
			var count uintptr
			_, err := winapi.ReadConsoleInputKey(cil.handle, inputRecord, 1, &count)
			if err.Error() != _SUCCESS { panic(err.Error) }
			keyInput := inputRecord[0].Event
			coninmsg := &ConsoleInputEvent {
				RepeatCount : keyInput.WRepeatCount,
				VirtualKeyCode : keyInput.WVirtualKeyCode,
				VirtualScanCode : keyInput.WVirtualScanCode,
				UnicodeCharacter : keyInput.UChar,
				ControlKeyState : keyInput.DwControlKeyState,
			}
			if keyInput.BKeyDown == 0 {
				coninmsg.KeyDown = false
			} else {
				coninmsg.KeyDown = true
			}
			if cil.enabled {
				cil.messages <- coninmsg
			} else {
				break
			}
		}
	}()
}

// Stop the ConsoleInputLogger and close the message channel.
func (cil *ConsoleInputLogger) Stop() {
	cil.enabled = false
	setConsoleMode(cil.handle, cil.oldMode)
	close(cil.messages)
}

// getStdHandle gets and returns the console handle.
func getStdHandle(nStdHandle uintptr) uintptr {
	hStdin, err := winapi.GetStdHandle(nStdHandle)
	if err.Error() != _SUCCESS { panic(err.Error) }
	return hStdin
}

// getConsoleMode gets and returns the console mode.
func getConsoleMode(handle uintptr) uintptr {
	var mode uintptr
	_, err := winapi.GetConsoleMode(handle, &mode)
	if err.Error() != _SUCCESS { panic(err.Error()) }
	return mode
}

// setConsoleMode sets the console mode.
func setConsoleMode(handle uintptr, mode uintptr) {
	_, err := winapi.SetConsoleMode(handle, mode)
	if err.Error() != _SUCCESS { panic(err.Error()) }
}
