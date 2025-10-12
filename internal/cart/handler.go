package cart

import (
	"context"

	"github.com/cart-overflow/cart-api/pkg/pb"
	"github.com/cart-overflow/common/pkg/metadata"
	"github.com/cart-overflow/common/pkg/rpc"
)

type Handler struct {
	s *service
	pb.UnimplementedCartServiceServer
}

func newHandler(s *service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) SetCartItem(
	ctx context.Context,
	req *pb.SetCartItemRequest,
) (*pb.SetCartItemResponse, error) {
	res, err := h.s.SetCartItem(ctx, &setCartItemRequest{
		productId: req.ProductId,
		amount:    req.Amount,
		md:        metadata.FromRpcCtx(ctx),
	})
	if err != nil {
		return nil, rpc.MapServiceErr(err)
	}

	return &pb.SetCartItemResponse{ProductId: res.productId}, nil
}

func (h *Handler) GetCartItem(
	ctx context.Context,
	req *pb.GetCartItemRequest,
) (*pb.GetCartItemResponse, error) {
	res, err := h.s.GetCartItem(ctx, &getCartItemRequest{
		productId: req.ProductId,
		md:        metadata.FromRpcCtx(ctx),
	})
	if err != nil {
		return nil, rpc.MapServiceErr(err)
	}

	return &pb.GetCartItemResponse{Item: &pb.CartItem{
		ProductId: res.item.productId,
		Amount:    res.item.amount,
		AddedAt:   res.item.added_at,
	}}, nil
}

func (h *Handler) GetCart(
	ctx context.Context,
	req *pb.GetCartRequest,
) (*pb.GetCartResponse, error) {
	res, err := h.s.GetCart(ctx, &getCartRequest{md: metadata.FromRpcCtx(ctx)})
	if err != nil {
		return nil, rpc.MapServiceErr(err)
	}

	items := make([]*pb.CartItem, len(res.items))
	for i, cartItem := range res.items {
		items[i] = &pb.CartItem{
			ProductId: cartItem.productId,
			Amount:    cartItem.amount,
			AddedAt:   cartItem.added_at,
		}
	}

	return &pb.GetCartResponse{Items: items}, nil
}

func (h *Handler) ClearCart(
	ctx context.Context,
	req *pb.ClearCartRequest,
) (*pb.ClearCartResponse, error) {
	err := h.s.ClearCart(ctx, &clearCartRequest{
		md: metadata.FromRpcCtx(ctx),
	})
	if err != nil {
		return nil, rpc.MapServiceErr(err)
	}
	return &pb.ClearCartResponse{}, nil
}
