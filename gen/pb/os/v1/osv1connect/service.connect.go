// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: os/v1/service.proto

package osv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/wagecloud/wagecloud-server/gen/pb/os/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// OSServiceName is the fully-qualified name of the OSService service.
	OSServiceName = "os.v1.OSService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// OSServiceGetOSProcedure is the fully-qualified name of the OSService's GetOS RPC.
	OSServiceGetOSProcedure = "/os.v1.OSService/GetOS"
	// OSServiceListOSsProcedure is the fully-qualified name of the OSService's ListOSs RPC.
	OSServiceListOSsProcedure = "/os.v1.OSService/ListOSs"
	// OSServiceCreateOSProcedure is the fully-qualified name of the OSService's CreateOS RPC.
	OSServiceCreateOSProcedure = "/os.v1.OSService/CreateOS"
	// OSServiceUpdateOSProcedure is the fully-qualified name of the OSService's UpdateOS RPC.
	OSServiceUpdateOSProcedure = "/os.v1.OSService/UpdateOS"
	// OSServiceDeleteOSProcedure is the fully-qualified name of the OSService's DeleteOS RPC.
	OSServiceDeleteOSProcedure = "/os.v1.OSService/DeleteOS"
	// OSServiceGetArchProcedure is the fully-qualified name of the OSService's GetArch RPC.
	OSServiceGetArchProcedure = "/os.v1.OSService/GetArch"
	// OSServiceListArchsProcedure is the fully-qualified name of the OSService's ListArchs RPC.
	OSServiceListArchsProcedure = "/os.v1.OSService/ListArchs"
	// OSServiceCreateArchProcedure is the fully-qualified name of the OSService's CreateArch RPC.
	OSServiceCreateArchProcedure = "/os.v1.OSService/CreateArch"
	// OSServiceUpdateArchProcedure is the fully-qualified name of the OSService's UpdateArch RPC.
	OSServiceUpdateArchProcedure = "/os.v1.OSService/UpdateArch"
	// OSServiceDeleteArchProcedure is the fully-qualified name of the OSService's DeleteArch RPC.
	OSServiceDeleteArchProcedure = "/os.v1.OSService/DeleteArch"
)

// OSServiceClient is a client for the os.v1.OSService service.
type OSServiceClient interface {
	// Get OS by ID
	GetOS(context.Context, *connect.Request[v1.GetOSRequest]) (*connect.Response[v1.GetOSResponse], error)
	// List OSs
	ListOSs(context.Context, *connect.Request[v1.ListOSsRequest]) (*connect.Response[v1.ListOSsResponse], error)
	// Create OS
	CreateOS(context.Context, *connect.Request[v1.CreateOSRequest]) (*connect.Response[v1.CreateOSResponse], error)
	// Update OS
	UpdateOS(context.Context, *connect.Request[v1.UpdateOSRequest]) (*connect.Response[v1.UpdateOSResponse], error)
	// Delete OS
	DeleteOS(context.Context, *connect.Request[v1.DeleteOSRequest]) (*connect.Response[v1.DeleteOSResponse], error)
	// Get Arch by ID
	GetArch(context.Context, *connect.Request[v1.GetArchRequest]) (*connect.Response[v1.GetArchResponse], error)
	// List Archs
	ListArchs(context.Context, *connect.Request[v1.ListArchsRequest]) (*connect.Response[v1.ListArchsResponse], error)
	// Create Arch
	CreateArch(context.Context, *connect.Request[v1.CreateArchRequest]) (*connect.Response[v1.CreateArchResponse], error)
	// Update Arch
	UpdateArch(context.Context, *connect.Request[v1.UpdateArchRequest]) (*connect.Response[v1.UpdateArchResponse], error)
	// Delete Arch
	DeleteArch(context.Context, *connect.Request[v1.DeleteArchRequest]) (*connect.Response[v1.DeleteArchResponse], error)
}

// NewOSServiceClient constructs a client for the os.v1.OSService service. By default, it uses the
// Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewOSServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) OSServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	oSServiceMethods := v1.File_os_v1_service_proto.Services().ByName("OSService").Methods()
	return &oSServiceClient{
		getOS: connect.NewClient[v1.GetOSRequest, v1.GetOSResponse](
			httpClient,
			baseURL+OSServiceGetOSProcedure,
			connect.WithSchema(oSServiceMethods.ByName("GetOS")),
			connect.WithClientOptions(opts...),
		),
		listOSs: connect.NewClient[v1.ListOSsRequest, v1.ListOSsResponse](
			httpClient,
			baseURL+OSServiceListOSsProcedure,
			connect.WithSchema(oSServiceMethods.ByName("ListOSs")),
			connect.WithClientOptions(opts...),
		),
		createOS: connect.NewClient[v1.CreateOSRequest, v1.CreateOSResponse](
			httpClient,
			baseURL+OSServiceCreateOSProcedure,
			connect.WithSchema(oSServiceMethods.ByName("CreateOS")),
			connect.WithClientOptions(opts...),
		),
		updateOS: connect.NewClient[v1.UpdateOSRequest, v1.UpdateOSResponse](
			httpClient,
			baseURL+OSServiceUpdateOSProcedure,
			connect.WithSchema(oSServiceMethods.ByName("UpdateOS")),
			connect.WithClientOptions(opts...),
		),
		deleteOS: connect.NewClient[v1.DeleteOSRequest, v1.DeleteOSResponse](
			httpClient,
			baseURL+OSServiceDeleteOSProcedure,
			connect.WithSchema(oSServiceMethods.ByName("DeleteOS")),
			connect.WithClientOptions(opts...),
		),
		getArch: connect.NewClient[v1.GetArchRequest, v1.GetArchResponse](
			httpClient,
			baseURL+OSServiceGetArchProcedure,
			connect.WithSchema(oSServiceMethods.ByName("GetArch")),
			connect.WithClientOptions(opts...),
		),
		listArchs: connect.NewClient[v1.ListArchsRequest, v1.ListArchsResponse](
			httpClient,
			baseURL+OSServiceListArchsProcedure,
			connect.WithSchema(oSServiceMethods.ByName("ListArchs")),
			connect.WithClientOptions(opts...),
		),
		createArch: connect.NewClient[v1.CreateArchRequest, v1.CreateArchResponse](
			httpClient,
			baseURL+OSServiceCreateArchProcedure,
			connect.WithSchema(oSServiceMethods.ByName("CreateArch")),
			connect.WithClientOptions(opts...),
		),
		updateArch: connect.NewClient[v1.UpdateArchRequest, v1.UpdateArchResponse](
			httpClient,
			baseURL+OSServiceUpdateArchProcedure,
			connect.WithSchema(oSServiceMethods.ByName("UpdateArch")),
			connect.WithClientOptions(opts...),
		),
		deleteArch: connect.NewClient[v1.DeleteArchRequest, v1.DeleteArchResponse](
			httpClient,
			baseURL+OSServiceDeleteArchProcedure,
			connect.WithSchema(oSServiceMethods.ByName("DeleteArch")),
			connect.WithClientOptions(opts...),
		),
	}
}

// oSServiceClient implements OSServiceClient.
type oSServiceClient struct {
	getOS      *connect.Client[v1.GetOSRequest, v1.GetOSResponse]
	listOSs    *connect.Client[v1.ListOSsRequest, v1.ListOSsResponse]
	createOS   *connect.Client[v1.CreateOSRequest, v1.CreateOSResponse]
	updateOS   *connect.Client[v1.UpdateOSRequest, v1.UpdateOSResponse]
	deleteOS   *connect.Client[v1.DeleteOSRequest, v1.DeleteOSResponse]
	getArch    *connect.Client[v1.GetArchRequest, v1.GetArchResponse]
	listArchs  *connect.Client[v1.ListArchsRequest, v1.ListArchsResponse]
	createArch *connect.Client[v1.CreateArchRequest, v1.CreateArchResponse]
	updateArch *connect.Client[v1.UpdateArchRequest, v1.UpdateArchResponse]
	deleteArch *connect.Client[v1.DeleteArchRequest, v1.DeleteArchResponse]
}

// GetOS calls os.v1.OSService.GetOS.
func (c *oSServiceClient) GetOS(ctx context.Context, req *connect.Request[v1.GetOSRequest]) (*connect.Response[v1.GetOSResponse], error) {
	return c.getOS.CallUnary(ctx, req)
}

// ListOSs calls os.v1.OSService.ListOSs.
func (c *oSServiceClient) ListOSs(ctx context.Context, req *connect.Request[v1.ListOSsRequest]) (*connect.Response[v1.ListOSsResponse], error) {
	return c.listOSs.CallUnary(ctx, req)
}

// CreateOS calls os.v1.OSService.CreateOS.
func (c *oSServiceClient) CreateOS(ctx context.Context, req *connect.Request[v1.CreateOSRequest]) (*connect.Response[v1.CreateOSResponse], error) {
	return c.createOS.CallUnary(ctx, req)
}

// UpdateOS calls os.v1.OSService.UpdateOS.
func (c *oSServiceClient) UpdateOS(ctx context.Context, req *connect.Request[v1.UpdateOSRequest]) (*connect.Response[v1.UpdateOSResponse], error) {
	return c.updateOS.CallUnary(ctx, req)
}

// DeleteOS calls os.v1.OSService.DeleteOS.
func (c *oSServiceClient) DeleteOS(ctx context.Context, req *connect.Request[v1.DeleteOSRequest]) (*connect.Response[v1.DeleteOSResponse], error) {
	return c.deleteOS.CallUnary(ctx, req)
}

// GetArch calls os.v1.OSService.GetArch.
func (c *oSServiceClient) GetArch(ctx context.Context, req *connect.Request[v1.GetArchRequest]) (*connect.Response[v1.GetArchResponse], error) {
	return c.getArch.CallUnary(ctx, req)
}

// ListArchs calls os.v1.OSService.ListArchs.
func (c *oSServiceClient) ListArchs(ctx context.Context, req *connect.Request[v1.ListArchsRequest]) (*connect.Response[v1.ListArchsResponse], error) {
	return c.listArchs.CallUnary(ctx, req)
}

// CreateArch calls os.v1.OSService.CreateArch.
func (c *oSServiceClient) CreateArch(ctx context.Context, req *connect.Request[v1.CreateArchRequest]) (*connect.Response[v1.CreateArchResponse], error) {
	return c.createArch.CallUnary(ctx, req)
}

// UpdateArch calls os.v1.OSService.UpdateArch.
func (c *oSServiceClient) UpdateArch(ctx context.Context, req *connect.Request[v1.UpdateArchRequest]) (*connect.Response[v1.UpdateArchResponse], error) {
	return c.updateArch.CallUnary(ctx, req)
}

// DeleteArch calls os.v1.OSService.DeleteArch.
func (c *oSServiceClient) DeleteArch(ctx context.Context, req *connect.Request[v1.DeleteArchRequest]) (*connect.Response[v1.DeleteArchResponse], error) {
	return c.deleteArch.CallUnary(ctx, req)
}

// OSServiceHandler is an implementation of the os.v1.OSService service.
type OSServiceHandler interface {
	// Get OS by ID
	GetOS(context.Context, *connect.Request[v1.GetOSRequest]) (*connect.Response[v1.GetOSResponse], error)
	// List OSs
	ListOSs(context.Context, *connect.Request[v1.ListOSsRequest]) (*connect.Response[v1.ListOSsResponse], error)
	// Create OS
	CreateOS(context.Context, *connect.Request[v1.CreateOSRequest]) (*connect.Response[v1.CreateOSResponse], error)
	// Update OS
	UpdateOS(context.Context, *connect.Request[v1.UpdateOSRequest]) (*connect.Response[v1.UpdateOSResponse], error)
	// Delete OS
	DeleteOS(context.Context, *connect.Request[v1.DeleteOSRequest]) (*connect.Response[v1.DeleteOSResponse], error)
	// Get Arch by ID
	GetArch(context.Context, *connect.Request[v1.GetArchRequest]) (*connect.Response[v1.GetArchResponse], error)
	// List Archs
	ListArchs(context.Context, *connect.Request[v1.ListArchsRequest]) (*connect.Response[v1.ListArchsResponse], error)
	// Create Arch
	CreateArch(context.Context, *connect.Request[v1.CreateArchRequest]) (*connect.Response[v1.CreateArchResponse], error)
	// Update Arch
	UpdateArch(context.Context, *connect.Request[v1.UpdateArchRequest]) (*connect.Response[v1.UpdateArchResponse], error)
	// Delete Arch
	DeleteArch(context.Context, *connect.Request[v1.DeleteArchRequest]) (*connect.Response[v1.DeleteArchResponse], error)
}

// NewOSServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewOSServiceHandler(svc OSServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	oSServiceMethods := v1.File_os_v1_service_proto.Services().ByName("OSService").Methods()
	oSServiceGetOSHandler := connect.NewUnaryHandler(
		OSServiceGetOSProcedure,
		svc.GetOS,
		connect.WithSchema(oSServiceMethods.ByName("GetOS")),
		connect.WithHandlerOptions(opts...),
	)
	oSServiceListOSsHandler := connect.NewUnaryHandler(
		OSServiceListOSsProcedure,
		svc.ListOSs,
		connect.WithSchema(oSServiceMethods.ByName("ListOSs")),
		connect.WithHandlerOptions(opts...),
	)
	oSServiceCreateOSHandler := connect.NewUnaryHandler(
		OSServiceCreateOSProcedure,
		svc.CreateOS,
		connect.WithSchema(oSServiceMethods.ByName("CreateOS")),
		connect.WithHandlerOptions(opts...),
	)
	oSServiceUpdateOSHandler := connect.NewUnaryHandler(
		OSServiceUpdateOSProcedure,
		svc.UpdateOS,
		connect.WithSchema(oSServiceMethods.ByName("UpdateOS")),
		connect.WithHandlerOptions(opts...),
	)
	oSServiceDeleteOSHandler := connect.NewUnaryHandler(
		OSServiceDeleteOSProcedure,
		svc.DeleteOS,
		connect.WithSchema(oSServiceMethods.ByName("DeleteOS")),
		connect.WithHandlerOptions(opts...),
	)
	oSServiceGetArchHandler := connect.NewUnaryHandler(
		OSServiceGetArchProcedure,
		svc.GetArch,
		connect.WithSchema(oSServiceMethods.ByName("GetArch")),
		connect.WithHandlerOptions(opts...),
	)
	oSServiceListArchsHandler := connect.NewUnaryHandler(
		OSServiceListArchsProcedure,
		svc.ListArchs,
		connect.WithSchema(oSServiceMethods.ByName("ListArchs")),
		connect.WithHandlerOptions(opts...),
	)
	oSServiceCreateArchHandler := connect.NewUnaryHandler(
		OSServiceCreateArchProcedure,
		svc.CreateArch,
		connect.WithSchema(oSServiceMethods.ByName("CreateArch")),
		connect.WithHandlerOptions(opts...),
	)
	oSServiceUpdateArchHandler := connect.NewUnaryHandler(
		OSServiceUpdateArchProcedure,
		svc.UpdateArch,
		connect.WithSchema(oSServiceMethods.ByName("UpdateArch")),
		connect.WithHandlerOptions(opts...),
	)
	oSServiceDeleteArchHandler := connect.NewUnaryHandler(
		OSServiceDeleteArchProcedure,
		svc.DeleteArch,
		connect.WithSchema(oSServiceMethods.ByName("DeleteArch")),
		connect.WithHandlerOptions(opts...),
	)
	return "/os.v1.OSService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case OSServiceGetOSProcedure:
			oSServiceGetOSHandler.ServeHTTP(w, r)
		case OSServiceListOSsProcedure:
			oSServiceListOSsHandler.ServeHTTP(w, r)
		case OSServiceCreateOSProcedure:
			oSServiceCreateOSHandler.ServeHTTP(w, r)
		case OSServiceUpdateOSProcedure:
			oSServiceUpdateOSHandler.ServeHTTP(w, r)
		case OSServiceDeleteOSProcedure:
			oSServiceDeleteOSHandler.ServeHTTP(w, r)
		case OSServiceGetArchProcedure:
			oSServiceGetArchHandler.ServeHTTP(w, r)
		case OSServiceListArchsProcedure:
			oSServiceListArchsHandler.ServeHTTP(w, r)
		case OSServiceCreateArchProcedure:
			oSServiceCreateArchHandler.ServeHTTP(w, r)
		case OSServiceUpdateArchProcedure:
			oSServiceUpdateArchHandler.ServeHTTP(w, r)
		case OSServiceDeleteArchProcedure:
			oSServiceDeleteArchHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedOSServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedOSServiceHandler struct{}

func (UnimplementedOSServiceHandler) GetOS(context.Context, *connect.Request[v1.GetOSRequest]) (*connect.Response[v1.GetOSResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("os.v1.OSService.GetOS is not implemented"))
}

func (UnimplementedOSServiceHandler) ListOSs(context.Context, *connect.Request[v1.ListOSsRequest]) (*connect.Response[v1.ListOSsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("os.v1.OSService.ListOSs is not implemented"))
}

func (UnimplementedOSServiceHandler) CreateOS(context.Context, *connect.Request[v1.CreateOSRequest]) (*connect.Response[v1.CreateOSResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("os.v1.OSService.CreateOS is not implemented"))
}

func (UnimplementedOSServiceHandler) UpdateOS(context.Context, *connect.Request[v1.UpdateOSRequest]) (*connect.Response[v1.UpdateOSResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("os.v1.OSService.UpdateOS is not implemented"))
}

func (UnimplementedOSServiceHandler) DeleteOS(context.Context, *connect.Request[v1.DeleteOSRequest]) (*connect.Response[v1.DeleteOSResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("os.v1.OSService.DeleteOS is not implemented"))
}

func (UnimplementedOSServiceHandler) GetArch(context.Context, *connect.Request[v1.GetArchRequest]) (*connect.Response[v1.GetArchResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("os.v1.OSService.GetArch is not implemented"))
}

func (UnimplementedOSServiceHandler) ListArchs(context.Context, *connect.Request[v1.ListArchsRequest]) (*connect.Response[v1.ListArchsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("os.v1.OSService.ListArchs is not implemented"))
}

func (UnimplementedOSServiceHandler) CreateArch(context.Context, *connect.Request[v1.CreateArchRequest]) (*connect.Response[v1.CreateArchResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("os.v1.OSService.CreateArch is not implemented"))
}

func (UnimplementedOSServiceHandler) UpdateArch(context.Context, *connect.Request[v1.UpdateArchRequest]) (*connect.Response[v1.UpdateArchResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("os.v1.OSService.UpdateArch is not implemented"))
}

func (UnimplementedOSServiceHandler) DeleteArch(context.Context, *connect.Request[v1.DeleteArchRequest]) (*connect.Response[v1.DeleteArchResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("os.v1.OSService.DeleteArch is not implemented"))
}
