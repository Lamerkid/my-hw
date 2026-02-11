package internalgrpc

import (
	context "context"

	"github.com/google/uuid"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/app"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	UnimplementedEventServiceServer
	app    app.EventService
	logger Logger
}

func NewEventService(logger Logger, app app.EventService) *Service {
	return &Service{
		app:    app,
		logger: logger,
	}
}

func (s *Service) CreateEvent(ctx context.Context, in *CreateEventRequest) (*EventResponse, error) {
	s.logger.Debug("got request CreateEvent")

	parsedUserID, err := uuid.Parse(in.UserId)
	if err != nil {
		s.logger.Error("failed to parse UserID")
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	cmd := app.CreateEventCommand{
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		StartTime:   in.GetStartTime().AsTime(),
		EndTime:     in.GetEndTime().AsTime(),
		UserID:      parsedUserID,
	}

	event, err := s.app.CreateEvent(ctx, cmd)
	if err != nil {
		s.logger.Error("failed to create event")
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	resp := &EventResponse{
		Event: &Event{
			Id:          event.ID.String(),
			Title:       event.Title,
			Description: event.Description,
			StartTime:   timestamppb.New(event.StartTime),
			EndTime:     timestamppb.New(event.EndTime),
			UserId:      event.UserID.String(),
		},
	}

	return resp, nil
}

func (s *Service) UpdateEvent(ctx context.Context, in *UpdateEventRequest) (*EventResponse, error) {
	s.logger.Debug("got request UpdateEvent on event with id: %s", in.GetId())

	parsedID, err := uuid.Parse(in.GetId())
	if err != nil {
		s.logger.Error("failed to parse UserID")
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	cmd := app.UpdateEventCommand{
		ID:          parsedID,
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		StartTime:   in.GetStartTime().AsTime(),
		EndTime:     in.GetEndTime().AsTime(),
	}

	event, err := s.app.UpdateEvent(ctx, cmd)
	if err != nil {
		s.logger.Error("failed to update event")
		return nil, status.Errorf(codes.Internal, "failed to update event: %v", err)
	}

	resp := &EventResponse{
		Event: &Event{
			Id:          event.ID.String(),
			Title:       event.Title,
			Description: event.Description,
			StartTime:   timestamppb.New(event.StartTime),
			EndTime:     timestamppb.New(event.EndTime),
			UserId:      event.UserID.String(),
		},
	}

	return resp, nil
}

func (s *Service) DeleteEvent(ctx context.Context, in *EventIdRequest) (*EventResponse, error) {
	s.logger.Debug("got request DeleteEvent on event with id: %s", in.GetId())

	parsedID, err := uuid.Parse(in.GetId())
	if err != nil {
		s.logger.Error("failed to parse UserID")
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	event, err := s.app.DeleteEvent(ctx, parsedID)
	if err != nil {
		s.logger.Error("failed to update event")
		return nil, status.Errorf(codes.Internal, "failed to update event: %v", err)
	}

	resp := &EventResponse{
		Event: &Event{
			Id:          event.ID.String(),
			Title:       event.Title,
			Description: event.Description,
			StartTime:   timestamppb.New(event.StartTime),
			EndTime:     timestamppb.New(event.EndTime),
			UserId:      event.UserID.String(),
		},
	}

	return resp, nil
}

func (s *Service) SelectEventByDay(ctx context.Context, in *SelectEventRequest) (*ArrayEventResponse, error) {
	s.logger.Debug("got request SelectEventByDay on date: %s, by week", in.Date.AsTime())

	events, err := s.app.SelectEventByDay(ctx, in.Date.AsTime())
	if err != nil {
		s.logger.Error("failed to select events")
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	pbEvents := make([]*Event, len(events))

	resp := &ArrayEventResponse{
		Events: pbEvents,
	}

	return resp, nil
}

func (s *Service) SelectEventByWeek(ctx context.Context, in *SelectEventRequest) (*ArrayEventResponse, error) {
	s.logger.Debug("got request SelectEventByWeek on date: %s, by week", in.Date.AsTime())

	events, err := s.app.SelectEventByWeek(ctx, in.Date.AsTime())
	if err != nil {
		s.logger.Error("failed to select events")
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	pbEvents := make([]*Event, len(events))

	resp := &ArrayEventResponse{
		Events: pbEvents,
	}

	return resp, nil
}

func (s *Service) SelectEventByMonth(ctx context.Context, in *SelectEventRequest) (*ArrayEventResponse, error) {
	s.logger.Debug("got request SelectEventByMonth on date: %s, by month", in.Date.AsTime())

	events, err := s.app.SelectEventByMonth(ctx, in.Date.AsTime())
	if err != nil {
		s.logger.Error("failed to select events")
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	pbEvents := make([]*Event, len(events))

	resp := &ArrayEventResponse{
		Events: pbEvents,
	}

	return resp, nil
}
