package functions

import (
	"github.com/bitver/zigbee-steward/coordinator"
	"github.com/bitver/zigbee-steward/logger"
	"github.com/dyrkin/zcl-go"
)

var log = logger.MustGetLogger("functions")

type Functions struct {
	generic *GenericFunctions
	cluster *ClusterFunctions
}

func New(coordinator *coordinator.Coordinator, zcl *zcl.Zcl) *Functions {
	return &Functions{
		generic: &GenericFunctions{coordinator: coordinator},
		cluster: &ClusterFunctions{
			global: &GlobalClusterFunctions{
				coordinator: coordinator,
				zcl:         zcl,
			},
			local: NewLocalClusterFunctions(coordinator, zcl),
		},
	}
}

func (f *Functions) Generic() *GenericFunctions {
	return f.generic
}

func (f *Functions) Cluster() *ClusterFunctions {
	return f.cluster
}
