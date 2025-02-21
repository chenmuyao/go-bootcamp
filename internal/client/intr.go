package client

import (
	"context"
	"math/rand/v2"
	"sync/atomic"

	intrv1 "github.com/chenmuyao/go-bootcamp/api/proto/gen/intr/v1"
	"google.golang.org/grpc"
)

type InteractiveClient struct {
	remote intrv1.InteractiveServiceClient
	local  intrv1.InteractiveServiceClient

	threshold atomic.Int32
}

// CancelCollect implements intrv1.InteractiveServiceClient.
func (i *InteractiveClient) CancelCollect(
	ctx context.Context,
	in *intrv1.CancelCollectRequest,
	opts ...grpc.CallOption,
) (*intrv1.CancelCollectResponse, error) {
	return i.selectClient().CancelCollect(ctx, in, opts...)
}

// CancelLike implements intrv1.InteractiveServiceClient.
func (i *InteractiveClient) CancelLike(
	ctx context.Context,
	in *intrv1.CancelLikeRequest,
	opts ...grpc.CallOption,
) (*intrv1.CancelLikeResponse, error) {
	return i.selectClient().CancelLike(ctx, in, opts...)
}

// Collect implements intrv1.InteractiveServiceClient.
func (i *InteractiveClient) Collect(
	ctx context.Context,
	in *intrv1.CollectRequest,
	opts ...grpc.CallOption,
) (*intrv1.CollectResponse, error) {
	return i.selectClient().Collect(ctx, in, opts...)
}

// Get implements intrv1.InteractiveServiceClient.
func (i *InteractiveClient) Get(
	ctx context.Context,
	in *intrv1.GetRequest,
	opts ...grpc.CallOption,
) (*intrv1.GetResponse, error) {
	return i.selectClient().Get(ctx, in, opts...)
}

// GetByIDs implements intrv1.InteractiveServiceClient.
func (i *InteractiveClient) GetByIDs(
	ctx context.Context,
	in *intrv1.GetByIDsRequest,
	opts ...grpc.CallOption,
) (*intrv1.GetByIDsResponse, error) {
	return i.selectClient().GetByIDs(ctx, in, opts...)
}

// GetTopLike implements intrv1.InteractiveServiceClient.
func (i *InteractiveClient) GetTopLike(
	ctx context.Context,
	in *intrv1.GetTopLikeRequest,
	opts ...grpc.CallOption,
) (*intrv1.GetTopLikeResponse, error) {
	return i.selectClient().GetTopLike(ctx, in, opts...)
}

// IncrReadCnt implements intrv1.InteractiveServiceClient.
func (i *InteractiveClient) IncrReadCnt(
	ctx context.Context,
	in *intrv1.IncrReadCntRequest,
	opts ...grpc.CallOption,
) (*intrv1.IncrReadCntResponse, error) {
	return i.selectClient().IncrReadCnt(ctx, in, opts...)
}

// Like implements intrv1.InteractiveServiceClient.
func (i *InteractiveClient) Like(
	ctx context.Context,
	in *intrv1.LikeRequest,
	opts ...grpc.CallOption,
) (*intrv1.LikeResponse, error) {
	return i.selectClient().Like(ctx, in, opts...)
}

// MustBatchGet implements intrv1.InteractiveServiceClient.
func (i *InteractiveClient) MustBatchGet(
	ctx context.Context,
	in *intrv1.MustBatchGetRequest,
	opts ...grpc.CallOption,
) (*intrv1.MustBatchGetResponse, error) {
	return i.selectClient().MustBatchGet(ctx, in, opts...)
}

func (i *InteractiveClient) selectClient() intrv1.InteractiveServiceClient {
	// [0, 100)
	num := rand.Int32N(100)
	if num < i.threshold.Load() {
		return i.remote
	}
	return i.local
}

func (i *InteractiveClient) UpdateThreshold(val int32) {
	i.threshold.Store(val)
}

func NewInteractiveClient(
	remote intrv1.InteractiveServiceClient,
	local intrv1.InteractiveServiceClient,
) *InteractiveClient {
	return &InteractiveClient{
		remote:    remote,
		local:     local,
		threshold: atomic.Int32{},
	}
}

var _ intrv1.InteractiveServiceClient = &InteractiveClient{}
