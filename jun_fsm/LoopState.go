package jun_fsm

import (
	"context"
	"time"
)

type LoopState struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	*time.Ticker
	customData *CustomData
}
