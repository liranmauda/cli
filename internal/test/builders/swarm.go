package builders

import (
	"time"

	"github.com/moby/moby/api/types/swarm"
)

// Swarm creates a swarm with default values.
// Any number of swarm function builder can be pass to augment it.
func Swarm(swarmBuilders ...func(*swarm.Swarm)) *swarm.Swarm {
	t1 := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	s := &swarm.Swarm{
		ClusterInfo: swarm.ClusterInfo{
			ID: "swarm",
			Meta: swarm.Meta{
				CreatedAt: t1,
			},
			Spec: swarm.Spec{},
		},
		JoinTokens: swarm.JoinTokens{
			Worker:  "worker-join-token",
			Manager: "manager-join-token",
		},
	}

	for _, builder := range swarmBuilders {
		builder(s)
	}

	return s
}

// Autolock set the swarm into autolock mode
func Autolock() func(*swarm.Swarm) {
	return func(s *swarm.Swarm) {
		s.Spec.EncryptionConfig.AutoLockManagers = true
	}
}
