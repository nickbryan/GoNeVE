package engine

import "github.com/google/uuid"

// Entity represents a single entity/world object managed by the engine. An Entity is identified by its unique ID.
//
// Due to this engines implementation of an Entity Component System, the Entity does not know about the components
// that are associated with it. Instead, an Entity should be registered with a component manager and it should be
// up to that component manager to keep track of the data relevant to the Entity within that system.
type Entity struct {
	ID string // The unique identifier for this entity.
}

// EntityManager is responsible for the creation and destruction of Entities within the engine.
type EntityManager struct {
	entities map[string]*Entity
}

// NewEntityManager will create and initialise a new EntityManager.
func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities: make(map[string]*Entity),
	}
}

// Create a new Entity within the EntityManager. Each Entity will be given a V4 UUID.
func (em *EntityManager) Create() *Entity {
	e := &Entity{ID: uuid.New().String()}

	em.entities[e.ID] = e

	return e
}

// CreateWith allows an Entity to be created with a predefined identifier. This is useful in networking
// when the Server has passed an Entity ID to the Client and the Client need to register it within their
// EntityManager.
//
// If an Entity with the given identifier already exists then that Entity will be returned. This is to prevent collisions
// or overwrites within the EntityManager.
func (em *EntityManager) CreateWith(id string) *Entity {
	if e, ok := em.entities[id]; ok {
		return e
	}

	e := &Entity{ID: id}
	em.entities[e.ID] = e
	return e
}

// Get will return the Entity associated with the given ID or nil of one does not exist.
func (em *EntityManager) Get(id string) *Entity {
	return em.entities[id]
}

// Alive will check to see if the Entity exists within the EntityManagers internal Entity list.
func (em *EntityManager) Alive(e *Entity) bool {
	_, ok := em.entities[e.ID]
	return ok
}

// Destroy will remove the given Entity from the EntityManager.
func (em *EntityManager) Destroy(e *Entity) {
	delete(em.entities, e.ID)
}
