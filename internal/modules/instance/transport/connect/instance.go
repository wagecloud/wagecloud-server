package instanceconnect

// type ImplementedInstanceServiceHandler struct {
// 	instancev1connect.UnimplementedInstanceServiceHandler
// 	service instancesvc.Service
// }

// func NewImplementedAccountServiceHandler(service instancesvc.Service) instancev1connect.InstanceServiceHandler {
// 	return &ImplementedInstanceServiceHandler{
// 		service: service,
// 	}
// }

// // Get instance by ID
// func (t *ImplementedInstanceServiceHandler) GetInstance(ctx context.Context, req *connect.Request[instancev1.GetInstanceRequest]) (*connect.Response[instancev1.GetInstanceResponse], error) {
// 	result, err := t.service.GetInstance(ctx, instancesvc.GetInstanceParams{
// 		ID: req.Msg.Id,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&instancev1.GetInstanceResponse{
// 		Instance: instancemodel.InstanceModelToProto(result),
// 	}), nil
// }

// // List instances
// func (t *ImplementedInstanceServiceHandler) ListInstances(ctx context.Context, req *connect.Request[instancev1.ListInstancesRequest]) (*connect.Response[instancev1.ListInstancesResponse], error) {
// 	result, err := t.service.ListInstances(ctx, instancesvc.ListInstancesParams{
// 		PaginationParams: commonmodel.PaginationParamsProtoToModel(req.Msg.Pagination),
// 		Account:          accountmodel.AuthenticatedAccountProtoToModel(req.Msg.Account),
// 		OsID:             req.Msg.OsId,
// 		ArchID:           req.Msg.ArchId,
// 		Name:             req.Msg.Name,
// 		CpuFrom:          req.Msg.CpuFrom,
// 		CpuTo:            req.Msg.CpuTo,
// 		RamFrom:          req.Msg.RamFrom,
// 		RamTo:            req.Msg.RamTo,
// 		StorageFrom:      req.Msg.StorageFrom,
// 		StorageTo:        req.Msg.StorageTo,
// 		CreatedAtFrom:    req.Msg.CreatedAtFrom,
// 		CreatedAtTo:      req.Msg.CreatedAtTo,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&instancev1.ListInstancesResponse{
// 		Instances:  slice.Map(result.Data, instancemodel.InstanceModelToProto),
// 		Pagination: commonmodel.PaginateResultModelToProto(result),
// 	}), nil
// }

// // Create instance
// func (t *ImplementedInstanceServiceHandler) CreateInstance(ctx context.Context, req *connect.Request[instancev1.CreateInstanceRequest]) (*connect.Response[instancev1.CreateInstanceResponse], error) {
// 	result, err := t.service.CreateInstance(ctx, instancesvc.CreateInstanceParams{
// 		Account:           accountmodel.AuthenticatedAccountProtoToModel(req.Msg.Account),
// 		Name:              req.Msg.Name,
// 		SSHAuthorizedKeys: req.Msg.SshAuthorizedKeys,
// 		Password:          req.Msg.Password,
// 		LocalHostname:     req.Msg.LocalHostname,
// 		OsID:              req.Msg.OsId,
// 		ArchID:            req.Msg.ArchId,
// 		Memory:            req.Msg.Memory,
// 		Cpu:               req.Msg.Cpu,
// 		Storage:           req.Msg.Storage,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&instancev1.CreateInstanceResponse{
// 		Instance: instancemodel.InstanceModelToProto(result),
// 	}), nil
// }

// // Update instance
// func (t *ImplementedInstanceServiceHandler) UpdateInstance(ctx context.Context, req *connect.Request[instancev1.UpdateInstanceRequest]) (*connect.Response[instancev1.UpdateInstanceResponse], error) {
// 	result, err := t.service.UpdateInstance(ctx, instancesvc.UpdateInstanceParams{
// 		ID:        req.Msg.Id,
// 		Name:      req.Msg.Name,
// 		NetworkID: req.Msg.NetworkId,
// 		OsID:      req.Msg.OsId,
// 		ArchID:    req.Msg.ArchId,
// 		Cpu:       req.Msg.Cpu,
// 		Ram:       req.Msg.Ram,
// 		Storage:   req.Msg.Storage,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&instancev1.UpdateInstanceResponse{
// 		Instance: instancemodel.InstanceModelToProto(result),
// 	}), nil
// }

// // Delete instance
// func (t *ImplementedInstanceServiceHandler) DeleteInstance(ctx context.Context, req *connect.Request[instancev1.DeleteInstanceRequest]) (*connect.Response[instancev1.DeleteInstanceResponse], error) {
// 	err := t.service.DeleteInstance(ctx, instancesvc.DeleteInstanceParams{
// 		ID: req.Msg.Id,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&instancev1.DeleteInstanceResponse{}), nil
// }

// // Start instance
// func (t *ImplementedInstanceServiceHandler) StartInstance(ctx context.Context, req *connect.Request[instancev1.StartInstanceRequest]) (*connect.Response[instancev1.StartInstanceResponse], error) {
// 	err := t.service.StartInstance(ctx, instancesvc.StartInstanceParams{
// 		ID: req.Msg.Id,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&instancev1.StartInstanceResponse{}), nil
// }

// // Stop instance
// func (t *ImplementedInstanceServiceHandler) StopInstance(ctx context.Context, req *connect.Request[instancev1.StopInstanceRequest]) (*connect.Response[instancev1.StopInstanceResponse], error) {
// 	err := t.service.StopInstance(ctx, instancesvc.StopInstanceParams{
// 		ID: req.Msg.Id,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return connect.NewResponse(&instancev1.StopInstanceResponse{}), nil
// }
