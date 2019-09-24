package engine

import (
	"reflect"
	"testing"
)

func TestNewEntityManager(t *testing.T) {
	t.Run("returns a pointer to a new EntityManager", func(t *testing.T) {
		got := reflect.TypeOf(NewEntityManager())
		want := reflect.TypeOf(&EntityManager{})
		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})
}

func TestEntityManager(t *testing.T) {
	t.Run("can Create a new Entity to be used within the ECS", func(t *testing.T) {
		em := NewEntityManager()
		got := reflect.TypeOf(em.Create())
		want := reflect.TypeOf(&Entity{})
		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})

	t.Run("can Destroy an existing entity", func(t *testing.T) {
		em := NewEntityManager()
		e := em.Create()
		em.Destroy(e)
		if em.Alive(e) {
			t.Errorf("entity was not removed from the EntityManager")
		}
	})

	t.Run("does not error when Destroy is called on a destroyed Entity", func(t *testing.T) {
		em := NewEntityManager()
		e := em.Create()
		em.Destroy(e)
		em.Destroy(e)
	})

	t.Run("can lookup an Entity by it's ID", func(t *testing.T) {
		em := NewEntityManager()
		e := em.Create()
		got := em.Get(e.ID)
		if e != got {
			t.Errorf("could not Get Entity via ID: got %v expected %v", got, e)
		}

		if e := em.Get("some_id"); e != nil {
			t.Errorf("expected Get to return nil: got %v", e)
		}
	})

	t.Run("can create an Entity from an existing id", func(t *testing.T) {
		id := "d0800ab1-0ce1-4c0c-8d19-32ef5d4074c7"
		em := NewEntityManager()
		e := em.CreateWith(id)
		if reflect.TypeOf(e) != reflect.TypeOf(&Entity{}) {
			t.Errorf("returned value is not an Entity: got %T expected %T", e, &Entity{})
		}

		if e.ID != id {
			t.Errorf("entity does not have the correct id: got %s expected %s", e.ID, id)
		}

		got := em.Get(id)
		if got != e {
			t.Errorf("returned Entity from CreateWith does not match lookup with same ID via Get: got %v expected %v", got, e)
		}

		got = em.CreateWith(e.ID)
		if got != e {
			t.Errorf("expected existing Entity to be returned from CreateWith when passing the ID of an existing Entity")
		}
	})
}

func TestEntityManager_Alive(t *testing.T) {
	t.Run("returns false if the Entity does not exist within the EntityManager", func(t *testing.T) {
		em := NewEntityManager()
		if em.Alive(&Entity{}) {
			t.Errorf("expected Alive to return false")
		}
	})

	t.Run("returns true if the Entity  exists within the EntityManager", func(t *testing.T) {
		em := NewEntityManager()
		if !em.Alive(em.Create()) {
			t.Errorf("expected Alive to return true")
		}
	})
}
