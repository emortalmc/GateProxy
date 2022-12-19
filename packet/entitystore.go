package packet

import (
	"math"
	"sync"

	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/util/uuid"
)

var EntityStore entityStore

type entityStore struct{ sync.Map }

func (s *entityStore) Subscribe(mgr event.Manager) {
	event.Subscribe(mgr, math.MaxInt, func(e *proxy.ServerConnectedEvent) {
		s.Store(e.Player().ID(), e.EntityID())
	})
	event.Subscribe(mgr, math.MaxInt, func(e *proxy.DisconnectEvent) {
		s.Delete(e.Player().ID())
	})
}

func (s *entityStore) EntityID(playerID uuid.UUID) int {
	if entityID, ok := s.Load(playerID); ok {
		return entityID.(int)
	}
	return 0
}
