package main

import (
	"log"

	"github.com/faiface/mainthread"

	"github.com/nickbryan/GoNeVE/pkg/render"

	"github.com/nickbryan/GoNeVE/pkg/event"
	"github.com/nickbryan/GoNeVE/pkg/input"

	"github.com/nickbryan/GoNeVE/pkg/engine"
)

func main() {
	//c := network.NewClient(":7777")
	//if err := c.Dial(); err != nil {
	//	log.Fatalf("unable to dial server: %v", err)
	//}
	//
	//c.Handle(&impl.Player{EntityManager: engine.NewEntityManager()})

	// Need to figure out best way to add window, should this be a renderer? how do we split clear screen poll events etc?

	mainthread.Run(func() {
		em := event.NewManager()

		inputSys := input.CreateSystem(em)
		inputSys.AddKeyCommands(input.KeyW, input.Pressed, input.KeyCommandExecutorFunc(func(_ float64) {
			log.Println("press")
		}))

		renderer, err := render.CreateDefaultSystem(em)
		if err != nil {
			log.Fatalf("unable to create default render system: %v", err)
		}

		e := engine.New()

		e.AddPreSimulator(inputSys)

		e.AddPreRenderer(renderer)
		e.AddRenderer(renderer)
		e.AddPostRenderer(renderer)

		e.Run()
	})
}
