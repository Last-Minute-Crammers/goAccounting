package thirdparty

import (
	aiService "goAccounting/internal/service/thirdparty/ai"
)

type Group struct {
	AIServiceGroup *aiService.ServiceGroup
}

var GroupApp = &Group{
	AIServiceGroup: aiService.NewServiceGroup(),
}
