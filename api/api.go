package api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/cafebazaar/booker-reservation/proto"
)

type reservationService struct{}

func (srv *reservationService) Get(c context.Context, s *proto.ReservationGetRequest) (*proto.ReservationGetReply, error) {
	p := &proto.ReservationGetReply{
		ReplyProperties: proto.ReplyPropertiesTemplate(),
	}

	return p, nil
}

func RegisterServer(grpcServer *grpc.Server) {
	proto.RegisterReservationServer(grpcServer, new(reservationService))
}
