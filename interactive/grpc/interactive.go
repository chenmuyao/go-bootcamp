package grpc

import (
	"context"

	"github.com/chenmuyao/generique/gslice"
	intrv1 "github.com/chenmuyao/go-bootcamp/api/proto/gen/intr/v1"
	"github.com/chenmuyao/go-bootcamp/interactive/domain"
	"github.com/chenmuyao/go-bootcamp/interactive/service"
)

type InteractiveServiceServer struct {
	intrv1.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}

// CancelCollect implements intrv1.InteractiveServiceServer.
func (i *InteractiveServiceServer) CancelCollect(
	ctx context.Context,
	request *intrv1.CancelCollectRequest,
) (*intrv1.CancelCollectResponse, error) {
	err := i.svc.CancelCollect(
		ctx,
		request.GetBiz(),
		request.GetId(),
		request.GetCid(),
		request.GetUid(),
	)
	return &intrv1.CancelCollectResponse{}, err
}

// CancelLike implements intrv1.InteractiveServiceServer.
func (i *InteractiveServiceServer) CancelLike(
	ctx context.Context,
	request *intrv1.CancelLikeRequest,
) (*intrv1.CancelLikeResponse, error) {
	err := i.svc.CancelLike(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &intrv1.CancelLikeResponse{}, err
}

// Collect implements intrv1.InteractiveServiceServer.
func (i *InteractiveServiceServer) Collect(
	ctx context.Context,
	request *intrv1.CollectRequest,
) (*intrv1.CollectResponse, error) {
	err := i.svc.Collect(ctx, request.GetBiz(), request.GetId(), request.GetCid(), request.GetUid())
	return &intrv1.CollectResponse{}, err
}

// Get implements intrv1.InteractiveServiceServer.
func (i *InteractiveServiceServer) Get(
	ctx context.Context,
	request *intrv1.GetRequest,
) (*intrv1.GetResponse, error) {
	intr, err := i.svc.Get(ctx, request.GetBiz(), request.GetId(), request.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResponse{
		Intr: i.toDTO(intr),
	}, nil
}

// GetByIDs implements intrv1.InteractiveServiceServer.
func (i *InteractiveServiceServer) GetByIDs(
	ctx context.Context,
	request *intrv1.GetByIDsRequest,
) (*intrv1.GetByIDsResponse, error) {
	intrs, err := i.svc.GetByIDs(ctx, request.GetBiz(), request.GetIds())
	if err != nil {
		return nil, err
	}
	res := make(map[int64]*intrv1.Interactive, len(intrs))
	for key, val := range intrs {
		res[key] = i.toDTO(val)
	}
	return &intrv1.GetByIDsResponse{
		Intrs: res,
	}, nil
}

// GetTopLike implements intrv1.InteractiveServiceServer.
func (i *InteractiveServiceServer) GetTopLike(
	ctx context.Context,
	request *intrv1.GetTopLikeRequest,
) (*intrv1.GetTopLikeResponse, error) {
	likes, err := i.svc.GetTopLike(ctx, request.GetBiz(), int(request.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &intrv1.GetTopLikeResponse{Ids: likes}, nil
}

// IncrReadCnt implements intrv1.InteractiveServiceServer.
func (i *InteractiveServiceServer) IncrReadCnt(
	ctx context.Context,
	request *intrv1.IncrReadCntRequest,
) (*intrv1.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, request.GetBiz(), request.GetBizId())
	return &intrv1.IncrReadCntResponse{}, err
}

// Like implements intrv1.InteractiveServiceServer.
func (i *InteractiveServiceServer) Like(
	ctx context.Context,
	request *intrv1.LikeRequest,
) (*intrv1.LikeResponse, error) {
	err := i.svc.Like(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &intrv1.LikeResponse{}, err
}

// MustBatchGet implements intrv1.InteractiveServiceServer.
func (i *InteractiveServiceServer) MustBatchGet(
	ctx context.Context,
	request *intrv1.MustBatchGetRequest,
) (*intrv1.MustBatchGetResponse, error) {
	intrs, err := i.svc.MustBatchGet(ctx, request.GetBiz(), request.GetIds())
	if err != nil {
		return nil, err
	}
	res := gslice.Map(intrs, func(id int, src domain.Interactive) *intrv1.Interactive {
		return i.toDTO(src)
	})

	return &intrv1.MustBatchGetResponse{Intrs: res}, nil
}

func (i *InteractiveServiceServer) toDTO(intr domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizID,
		ReadCnt:    intr.ReadCnt,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Liked:      intr.Liked,
		Collected:  intr.Collected,
	}
}

func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{svc: svc}
}
