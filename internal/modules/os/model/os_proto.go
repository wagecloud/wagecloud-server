package osmodel

import osv1 "github.com/wagecloud/wagecloud-server/gen/pb/os/v1"

func OSModelToProto(os OS) *osv1.OS {
	return &osv1.OS{
		Id:        os.ID,
		Name:      os.Name,
		CreatedAt: os.CreatedAt,
	}
}

func OSProtoToModel(proto *osv1.OS) OS {
	return OS{
		ID:        proto.Id,
		Name:      proto.Name,
		CreatedAt: proto.CreatedAt,
	}
}

func ArchModelToProto(arch Arch) *osv1.Arch {
	return &osv1.Arch{
		Id:        arch.ID,
		Name:      arch.Name,
		CreatedAt: arch.CreatedAt,
	}
}
