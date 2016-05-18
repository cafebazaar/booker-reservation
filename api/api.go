package api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/Sirupsen/logrus"    

	"github.com/cafebazaar/booker-reservation/proto"
)

type reservationService struct{}

func (srv *reservationService) GetReservation(c context.Context, s *proto.GetReservationRequest) (*proto.GetReservationReply, error) {
	p := &proto.GetReservationReply{
		ReplyProperties: proto.ReplyPropertiesTemplate(),
	}

	reservation, err := getReservation(s.ObjectURI, s.Timestamp)
	if err != nil {
		p.ReplyProperties.StatusCode = proto.ReplyProperties_NOT_FOUND
	}
	if reservation != nil {
		p.Reservation = &proto.ReservationInstance {
			StartTimestamp: reservation.StartTimestamp,
			EndTimestamp: reservation.EndTimestamp,
			UserID: reservation.UserID,
		}
	}

	return p, nil
}

func (srv *reservationService) PostReservation(c context.Context, s *proto.PostReservationRequest) (*proto.PostReservationReply, error) {
	p := &proto.PostReservationReply{
		ReplyProperties: proto.ReplyPropertiesTemplate(),
	}

	rr := s.GetReservation()
	if rr == nil || rr.UserID == "" || rr.StartTimestamp == 0 || rr.EndTimestamp == 0 {
		p.ReplyProperties.StatusCode = proto.ReplyProperties_BAD_REQUEST
		return p, nil
	}

  	reservation, err := createReservation(s.ObjectURI, rr.StartTimestamp,
	  rr.EndTimestamp, rr.UserID)

	if err != nil {
		p.ReplyProperties.StatusCode = proto.ReplyProperties_BAD_REQUEST
		logrus.WithError(err).Debugf("Failed to createReservation(%s, %d, %d, %s)",
			s.ObjectURI, rr.StartTimestamp,
	  rr.EndTimestamp, rr.UserID)
	}
	if reservation != nil {
		p.Reservation = &proto.ReservationInstance {
			StartTimestamp: reservation.StartTimestamp,
			EndTimestamp: reservation.EndTimestamp,
			UserID: reservation.UserID,
		}
	}

	return p, nil
}

func RegisterServer(grpcServer *grpc.Server) {
	proto.RegisterReservationServer(grpcServer, new(reservationService))
}
