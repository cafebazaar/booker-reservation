package api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/cafebazaar/booker-reservation/proto"
)

type reservationService struct{}

func (srv *reservationService) GetReservation(c context.Context, s *proto.GetReservationRequest) (*proto.GetReservationReply, error) {
	p := &proto.GetReservationReply{
		ReplyProperties: proto.ReplyPropertiesTemplate(),
	}

	return p, nil
}

func (srv *reservationService) PostReservation(c context.Context, s *proto.PostReservationRequest) (*proto.PostReservationReply, error) {
	p := &proto.PostReservationReply{
		ReplyProperties: proto.ReplyPropertiesTemplate(),
	}

	return p, nil
}

func RegisterServer(grpcServer *grpc.Server) {
	proto.RegisterReservationServer(grpcServer, new(reservationService))
}
