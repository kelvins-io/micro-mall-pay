package startup

import (
	"context"
	"net/http"

	"gitee.com/cristiane/micro-mall-pay/http_server"
	"gitee.com/cristiane/micro-mall-pay/proto/micro_mall_sku_proto/sku_business"
	"gitee.com/cristiane/micro-mall-pay/server"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// RegisterGRPCServer 此处注册pb的Server
func RegisterGRPCServer(grpcServer *grpc.Server) error {
	sku_business.RegisterSkuBusinessServiceServer(grpcServer, server.NewSkuBusinessServer())
	return nil
}

// RegisterGateway 此处注册pb的Gateway
func RegisterGateway(ctx context.Context, gateway *runtime.ServeMux, endPoint string, dopts []grpc.DialOption) error {
	if err := sku_business.RegisterSkuBusinessServiceHandlerFromEndpoint(ctx, gateway, endPoint, dopts); err != nil {
		return err
	}
	return nil
}

// RegisterHttpRoute 此处注册http接口
func RegisterHttpRoute(serverMux *http.ServeMux) error {
	serverMux.HandleFunc("/swagger/", http_server.SwaggerHandler)
	return nil
}
