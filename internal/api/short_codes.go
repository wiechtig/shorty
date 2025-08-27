package api

import (
	"context"
	"errors"
	"log/slog"
)

func (s Server) CreateShortCode(ctx context.Context, request CreateShortCodeRequestObject) (CreateShortCodeResponseObject, error) {
	//TODO implement me
	slog.Info("CreateShortCode called")

	return nil, errors.New("CreateShortCode not implemented")
}

func (s Server) ListShortCodes(ctx context.Context, request ListShortCodesRequestObject) (ListShortCodesResponseObject, error) {
	//TODO implement me
	slog.Info("ListShortCodes called")

	return nil, errors.New("ListShortCodes not implemented")
}
