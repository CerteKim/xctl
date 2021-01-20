package xctl

import (
	"context"
	"fmt"
	"log"

	logcmd "github.com/xtls/xray-core/app/log/command"
	proxymancmd "github.com/xtls/xray-core/app/proxyman/command"
	statscmd "github.com/xtls/xray-core/app/stats/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/common/uuid"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/vmess"
	"google.golang.org/grpc"
)

//GenerateUUID ...
//	Create a random UUID in string
//
func GenerateUUID() string {
	u := uuid.New()
	return u.String()
}

//ServiceClient ...
//  Struct holding the rpc client.
type ServiceClient struct {
	APIAddress  string
	APIPort     uint32
	statClient  statscmd.StatsServiceClient
	proxyClient proxymancmd.HandlerServiceClient
	logClient   logcmd.LoggerServiceClient
}

// NewServiceClient ...
// Create ServiceClient obj
func NewServiceClient(addr string, port uint32) *ServiceClient {
	cmdConn, err := grpc.Dial(fmt.Sprintf("%s:%d", addr, port), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return nil
	}

	svr := ServiceClient{APIAddress: addr, APIPort: port,
		statClient:  statscmd.NewStatsServiceClient(cmdConn),
		proxyClient: proxymancmd.NewHandlerServiceClient(cmdConn),
		logClient:   logcmd.NewLoggerServiceClient(cmdConn),
	}
	return &svr
}

//QueryStats ...
// List all stats from v2ray. users didn't generate any traffic won't included.
func (h *ServiceClient) QueryStats(pattern string, reset bool) map[string]int64 {
	sresp, err := h.statClient.QueryStats(context.Background(), &statscmd.QueryStatsRequest{
		Pattern: pattern,
		Reset_:  reset,
	})

	result := make(map[string]int64)
	if err != nil {
		log.Printf("failed to call grpc command: %v", err)
	} else {
		// log.Printf("%v", sresp)
		for _, stat := range sresp.Stat {
			result[stat.Name] = stat.Value
		}
	}

	return result
}

//GetStats ...
//  Get stats data of a sepecied item, key from QueryStats query.
func (h *ServiceClient) GetStats(name string, reset bool) (string, int64) {
	sresp, err := h.statClient.GetStats(context.Background(), &statscmd.GetStatsRequest{
		Name:   name,
		Reset_: reset,
	})

	if err != nil {
		log.Printf("%v", err)
		return "", 0
	}

	return sresp.Stat.Name, sresp.Stat.Value
}

// AddUser ...
//   Add a user to an inbound on the fly. The effect is not permentnent.
func (h *ServiceClient) AddUser(inboundTag string, email string, level uint32, uuid string, alterID uint32) {
	resp, err := h.proxyClient.AlterInbound(context.Background(), &proxymancmd.AlterInboundRequest{
		Tag: inboundTag,
		Operation: serial.ToTypedMessage(&proxymancmd.AddUserOperation{
			User: &protocol.User{
				Level: level,
				Email: email,
				Account: serial.ToTypedMessage(&vmess.Account{
					Id:               uuid,
					AlterId:          alterID,
					SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
				}),
			},
		}),
	})

	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("%v", resp)
	}
}

// RemoveUser ...
//   Remove a user from an Inbound on the fly. The effect is not permentnent.
func (h *ServiceClient) RemoveUser(inboundTag string, email string) {
	resp, err := h.proxyClient.AlterInbound(context.Background(), &proxymancmd.AlterInboundRequest{
		Tag: inboundTag,
		Operation: serial.ToTypedMessage(&proxymancmd.RemoveUserOperation{
			Email: email,
		}),
	})

	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("%v", resp)
	}
}

// RestartLogger
// IDK if it will work :D
func (h *ServiceClient) RestartLogger() {
	resp, err := h.logClient.RestartLogger(context.Background(), &logcmd.RestartLoggerRequest{})
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("%v", resp)
	}
}

func (h *ServiceClient) adi(in *core.InboundHandlerConfig) {
	resp, err := h.proxyClient.AddInbound(context.Background(), &proxymancmd.AddInboundRequest{
		Inbound: in,
	})
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("%v", resp)
	}
}

func (h *ServiceClient) ado(out *core.OutboundHandlerConfig) {
	resp, err := h.proxyClient.AddOutbound(context.Background(), &proxymancmd.AddOutboundRequest{
		Outbound: out,
	})
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("%v", resp)
	}
}

func (h *ServiceClient) rmi(tag string) {
	resp, err := h.proxyClient.RemoveInbound(context.Background(), &proxymancmd.RemoveInboundRequest{
		Tag: tag,
	})
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("%v", resp)
	}
}

func (h *ServiceClient) rmp(tag string) {
	resp, err := h.proxyClient.RemoveOutbound(context.Background(), &proxymancmd.RemoveOutboundRequest{
		Tag: tag,
	})
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("%v", resp)
	}
}
