package input

const (
	KeyPressedEvent  = "input.KeyPressed"
	KeyReleasedEvent = "input.KeyReleased"
)

// KeyEventMessage encapsulates the relevant information for a keyboard event. It should be dispatched when a WindowManager
// detects a keyboard event.
//
// Upon receiving a KeyEventMessage the Manager will trigger the relevant command callbacks that the user has registered.
type KeyEventMessage struct {
	Action   Action
	Key      Key
	Modifier ModifierKey
}
