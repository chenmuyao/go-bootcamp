package client

import (
	"context"

	"github.com/chenmuyao/generique/gslice"
	intrv1 "github.com/chenmuyao/go-bootcamp/api/proto/gen/intr/v1"
	"github.com/chenmuyao/go-bootcamp/interactive/domain"
	"github.com/chenmuyao/go-bootcamp/interactive/service"
	"google.golang.org/grpc"
)

type LocalInteractiveAdapter struct {
	svc service.InteractiveService
}

// CancelCollect implements intrv1.InteractiveServiceClient.
func (l *LocalInteractiveAdapter) CancelCollect(
	ctx context.Context,
	in *intrv1.CancelCollectRequest,
	opts ...grpc.CallOption,
) (*intrv1.CancelCollectResponse, error) {
	err := l.svc.CancelCollect(ctx, in.GetBiz(), in.GetId(), in.GetCid(), in.GetUid())
	return &intrv1.CancelCollectResponse{}, err
}

// CancelLike implements intrv1.InteractiveServiceClient.
func (l *LocalInteractiveAdapter) CancelLike(
	ctx context.Context,
	in *intrv1.CancelLikeRequest,
	opts ...grpc.CallOption,
) (*intrv1.CancelLikeResponse, error) {
	err := l.svc.CancelLike(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.CancelLikeResponse{}, err
}

// Collect implements intrv1.InteractiveServiceClient.
func (l *LocalInteractiveAdapter) Collect(
	ctx context.Context,
	in *intrv1.CollectRequest,
	opts ...grpc.CallOption,
) (*intrv1.CollectResponse, error) {
	err := l.svc.Collect(ctx, in.GetBiz(), in.GetId(), in.GetCid(), in.GetUid())
	return &intrv1.CollectResponse{}, err
}

// Get implements intrv1.InteractiveServiceClient.
func (l *LocalInteractiveAdapter) Get(
	ctx context.Context,
	in *intrv1.GetRequest,
	opts ...grpc.CallOption,
) (*intrv1.GetResponse, error) {
	res, err := l.svc.Get(ctx, in.GetBiz(), in.GetId(), in.GetUid())
	return &intrv1.GetResponse{Intr: l.toDTO(res)}, err
}

// GetByIDs implements intrv1.InteractiveServiceClient.
func (l *LocalInteractiveAdapter) GetByIDs(
	ctx context.Context,
	in *intrv1.GetByIDsRequest,
	opts ...grpc.CallOption,
) (*intrv1.GetByIDsResponse, error) {
	intrs, err := l.svc.GetByIDs(ctx, in.GetBiz(), in.GetIds())
	if err != nil {
		return nil, err
	}
	res := make(map[int64]*intrv1.Interactive, len(intrs))
	for key, val := range intrs {
		res[key] = l.toDTO(val)
	}
	return &intrv1.GetByIDsResponse{
		Intrs: res,
	}, nil
}

// GetTopLike implements intrv1.InteractiveServiceClient.
func (l *LocalInteractiveAdapter) GetTopLike(
	ctx context.Context,
	in *intrv1.GetTopLikeRequest,
	opts ...grpc.CallOption,
) (*intrv1.GetTopLikeResponse, error) {
	likes, err := l.svc.GetTopLike(ctx, in.GetBiz(), int(in.GetLimit()))
	if err != nil {
		return nil, err
	}
	return &intrv1.GetTopLikeResponse{Ids: likes}, nil
}

// IncrReadCnt implements intrv1.InteractiveServiceClient.
func (l *LocalInteractiveAdapter) IncrReadCnt(
	ctx context.Context,
	in *intrv1.IncrReadCntRequest,
	opts ...grpc.CallOption,
) (*intrv1.IncrReadCntResponse, error) {
	err := l.svc.IncrReadCnt(ctx, in.GetBiz(), in.GetBizId())
	return &intrv1.IncrReadCntResponse{}, err
}

// Like implements intrv1.InteractiveServiceClient.
func (l *LocalInteractiveAdapter) Like(
	ctx context.Context,
	in *intrv1.LikeRequest,
	opts ...grpc.CallOption,
) (*intrv1.LikeResponse, error) {
	err := l.svc.Like(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.LikeResponse{}, err
}

// MustBatchGet implements intrv1.InteractiveServiceClient.
func (l *LocalInteractiveAdapter) MustBatchGet(
	ctx context.Context,
	in *intrv1.MustBatchGetRequest,
	opts ...grpc.CallOption,
) (*intrv1.MustBatchGetResponse, error) {
	intrs, err := l.svc.MustBatchGet(ctx, in.GetBiz(), in.GetIds())
	if err != nil {
		return nil, err
	}
	res := gslice.Map(intrs, func(id int, src domain.Interactive) *intrv1.Interactive {
		return l.toDTO(src)
	})

	return &intrv1.MustBatchGetResponse{Intrs: res}, nil
}

func (i *LocalInteractiveAdapter) toDTO(intr domain.Interactive) *intrv1.Interactive {
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

func NewLocalInteractiveAdapter(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	return &LocalInteractiveAdapter{
		svc: svc,
	}
}

var _ intrv1.InteractiveServiceClient = &LocalInteractiveAdapter{}
