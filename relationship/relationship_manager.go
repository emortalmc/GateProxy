package relationship

import "go.minekube.com/gate/pkg/util/uuid"

var (
	LastMessageMap = make(map[uuid.UUID]uuid.UUID)
)
