package render

import (
	"runtime"

	"github.com/nickbryan/GoNeVE/pkg/engine"

	"github.com/nickbryan/GoNeVE/pkg/input"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/faiface/mainthread"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// Window encapsulates both a top-level win and an OpenGL context.
type Window interface {
	SwapBuffers()
	ShouldClose() bool
	SetPos(x, y int)
	GetPos() (x, y int)
}

// glfwWindow wraps the GLFWwindow functionality to satisfy the Window interface.
type glfwWindow struct {
	win *glfw.Window
}

// SwapBuffers swaps the front and back buffers of the specified win when rendering with OpenGL.
// If the swap interval is greater than zero, the GPU driver waits the specified number of screen
// simulations before swapping the buffers.
func (w *glfwWindow) SwapBuffers() {
	mainthread.Call(func() {
		w.win.SwapBuffers()
	})
}

// ShouldClose returns the value of the close flag of the specified win.
func (w *glfwWindow) ShouldClose() bool {
	var shouldClose bool

	mainthread.Call(func() {
		shouldClose = w.win.ShouldClose()
	})

	return shouldClose
}

// SetPos sets the position of the Window.
func (w *glfwWindow) SetPos(x, y int) {
	mainthread.Call(func() {
		w.win.SetPos(x, y)
	})
}

// GetPos returns the current position of the Window.
func (w *glfwWindow) GetPos() (x, y int) {
	mainthread.Call(func() {
		x, y = w.win.GetPos()
	})

	return x, y
}

// WindowManager encapsulates all shared win functionality.
type WindowManager interface {
	Initialise(publisher engine.Publisher) error
	Teardown()
	CreateWindow(width, height int, title string) (Window, error)
	PollEvents()
	SetSwapInterval(interval int)
}

// glfwWindowManager wraps the shared glfwWindow functionality.
type glfwWindowManager struct {
	publisher engine.Publisher
}

// Initialise initialises GLFW and sets appropriate win hints.
func (wm *glfwWindowManager) Initialise(publisher engine.Publisher) error {
	wm.publisher = publisher

	return mainthread.CallErr(func() error {
		if err := glfw.Init(); err != nil {
			return err
		}

		glfw.WindowHint(glfw.ContextVersionMajor, 4)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.Samples, 8)

		// Required for OSX.
		if runtime.GOOS == "darwin" {
			glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		}

		return nil
	})
}

// Teardown will destroy any remaining win, monitor and cursor objects,
// restore any modified gamma ramps, re-enable the screensaver if it had
// been disabled and free any other resources allocated by GLFW.
func (wm *glfwWindowManager) Teardown() {
	mainthread.Call(func() {
		glfw.Terminate()
	})
}

// CreateWindow creates a win, its associated OpenGL context and initialises the GLFW callbacks.
func (wm *glfwWindowManager) CreateWindow(width, height int, title string) (Window, error) {
	var (
		err     error
		glfwWin *glfw.Window
		win     *glfwWindow
	)

	err = mainthread.CallErr(func() error {
		glfwWin, err = glfw.CreateWindow(width, height, title, nil, nil)
		if err != nil {
			return err
		}

		win = &glfwWindow{
			win: glfwWin,
		}

		win.win.MakeContextCurrent()

		win.win.SetFramebufferSizeCallback(func(win *glfw.Window, width int, height int) {
			// TODO: move this to event to remove dependency
			gl.Viewport(0, 0, int32(width), int32(height))
		})

		win.win.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, mod glfw.ModifierKey) {
			if action == glfw.Press {
				wm.publisher.Publish(
					input.KeyEventMessage{
						Action:   input.KeyPressed,
						Key:      input.Key(key),
						Modifier: input.ModifierKey(mod),
					},
					input.KeyPressedEvent,
				)
			}

			if action == glfw.Release {
				wm.publisher.Publish(
					input.KeyEventMessage{
						Action:   input.KeyReleased,
						Key:      input.Key(key),
						Modifier: input.ModifierKey(mod),
					},
					input.KeyReleasedEvent,
				)
			}
		})

		return nil
	})

	return win, err
}

// PollEvents processes only those events that are already in the event
// queue and then returns immediately. Processing events will cause the
// win and input callbacks associated with those events to be called.
func (wm *glfwWindowManager) PollEvents() {
	mainthread.Call(func() {
		glfw.PollEvents()
	})
}

// SetSwapInterval sets the swap interval for the current OpenGL context
// i.e. the number of screen simulations to wait from the time glfwSwapBuffers
// was called before swapping the buffers and returning.
func (wm *glfwWindowManager) SetSwapInterval(interval int) {
	mainthread.Call(func() {
		glfw.SwapInterval(interval)
	})
}
