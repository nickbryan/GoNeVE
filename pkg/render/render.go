package render

import (
	"fmt"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/faiface/mainthread"
	"github.com/nickbryan/GoNeVE/pkg/engine"
)

type System struct {
	WindowManager WindowManager
	win           Window
	resizeOnce    sync.Once
}

func CreateDefaultSystem(pub engine.Publisher) (*System, error) {

	s := &System{
		WindowManager: &glfwWindowManager{},
	}
	err := s.WindowManager.Initialise(pub)
	if err != nil {
		return nil, fmt.Errorf("unable to initialise WindowManager: %v", err)
	}

	// TODO: config options
	win, err := s.WindowManager.CreateWindow(1080, 720, "Test")
	if err != nil {
		return nil, fmt.Errorf("unable to create window: %v", err)
	}

	s.win = win
	// TODO: config option
	s.WindowManager.SetSwapInterval(0)

	err = mainthread.CallErr(func() error {
		err = gl.Init()
		if err != nil {
			return fmt.Errorf("failed to initialise OpenGL: %v", err)
		}

		version := gl.GoStr(gl.GetString(gl.VERSION))
		// TODO: shift this?
		fmt.Println("OpenGL version: ", version)

		gl.Viewport(0, 0, 1080, 720)

		gl.Enable(gl.DEPTH_TEST)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *System) PreRender(_ float64) {
	s.WindowManager.PollEvents()

	mainthread.Call(func() {
		// TODO: config option for color
		gl.ClearColor(0.57, 0.71, 0.77, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	})

	// TODO: replace with glfw3.3 fix when released
	s.resizeOnce.Do(func() {
		x, y := s.win.GetPos()
		s.win.SetPos(x+1, y)
	})
}

func (s *System) Render(ipl float64) {

}

func (s *System) PostRender(_ float64) {
	s.win.SwapBuffers()
}
