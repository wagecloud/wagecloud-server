package instancemodel

import instancev1 "github.com/wagecloud/wagecloud-server/gen/pb/instance/v1"

func InstanceModelToProto(instance Instance) *instancev1.Instance {
	return &instancev1.Instance{
		Id:        instance.ID,
		Name:      instance.Name,
		OsId:      instance.OSID,
		ArchId:    instance.ArchID,
		Cpu:       instance.CPU,
		Ram:       instance.RAM,
		Storage:   instance.Storage,
		CreatedAt: instance.CreatedAt,
		UpdatedAt: instance.UpdatedAt,
	}
}

func InstanceProtoToModel(instance *instancev1.Instance) Instance {
	return Instance{
		ID:        instance.Id,
		Name:      instance.Name,
		OSID:      instance.OsId,
		ArchID:    instance.ArchId,
		CPU:       instance.Cpu,
		RAM:       instance.Ram,
		Storage:   instance.Storage,
		CreatedAt: instance.CreatedAt,
		UpdatedAt: instance.UpdatedAt,
	}
}

func NetworkModelToProto(network Network) *instancev1.Network {
	return &instancev1.Network{
		Id:        network.ID,
		CreatedAt: network.CreatedAt,
		PrivateIp: network.PrivateIP,
	}
}

func NetworkProtoToModel(network *instancev1.Network) Network {
	return Network{
		ID:        network.Id,
		CreatedAt: network.CreatedAt,
		PrivateIP: network.PrivateIp,
	}
}
