package cart

import (
	"context"
	"log"

	"github.com/cart-overflow/common/pkg/core"
	"github.com/cart-overflow/common/pkg/metadata"
)

type service struct {
	cache *cache
	log   *log.Logger
}

func newService(
	cache *cache,
	log *log.Logger,
) *service {
	return &service{cache, log}
}

// MARK: - Set Cart Item

type setCartItemRequest struct {
	productId string
	amount    int64
	md        *metadata.Metadata
}

type setCartItemResponse struct {
	productId string
}

func (s *service) SetCartItem(
	ctx context.Context,
	req *setCartItemRequest,
) (*setCartItemResponse, error) {
	if req.amount <= 0 {
		err := s.cache.deleteCartItem(ctx, req.md.UserId, req.productId)
		if err != nil {
			s.log.Printf("set cart deletion failure: %v", req.amount)
			return nil, core.NewErr(
				core.ErrInvalidArgument,
				"DELETE_CART_ITEM_FAILURE",
				"failed to remove cart item",
				nil,
			)
		}
		return &setCartItemResponse{productId: req.productId}, nil
	}

	err := s.cache.setCartItem(ctx, req.md.UserId, req.productId, req.amount)
	if err != nil {
		s.log.Printf("failed to update cart item: %v", err)
		return nil, core.NewErr(
			core.ErrInternal,
			"UPDATE_CART_ITEM_FAILURE",
			"failed to update cart item",
			nil,
		)
	}

	return &setCartItemResponse{productId: req.productId}, nil
}

// MARK: - Get Cart Item

type getCartItemRequest struct {
	productId string
	md        *metadata.Metadata
}

type getCartItemResponse struct {
	item *CartItem
}

func (s *service) GetCartItem(
	ctx context.Context,
	req *getCartItemRequest,
) (*getCartItemResponse, error) {
	cartItem, err := s.cache.getCartItem(ctx, req.md.UserId, req.productId)
	if err != nil {
		s.log.Printf("failed to get cart item: %v", err)
		return nil, core.NewErr(
			core.ErrInternal,
			"GET_CART_ITEM_FAILURE",
			"failed to get cart item",
			nil,
		)
	}

	return &getCartItemResponse{item: cartItem}, nil
}

// MARK: - Get Cart

type getCartRequest struct {
	md *metadata.Metadata
}

type getCartResponse struct {
	items []*CartItem
}

func (s *service) GetCart(
	ctx context.Context,
	req *getCartRequest,
) (*getCartResponse, error) {
	cart, err := s.cache.getCart(ctx, req.md.UserId)
	if err != nil {
		s.log.Printf("failed to get cart: %v", err)
		return nil, core.NewErr(
			core.ErrInternal,
			"GET_CART_FAILURE",
			"failed to get cart",
			nil,
		)
	}

	return &getCartResponse{cart}, nil
}

// MARK: - Clear Cart

type clearCartRequest struct {
	md *metadata.Metadata
}

func (s *service) ClearCart(
	ctx context.Context,
	req *clearCartRequest,
) error {
	err := s.cache.clearCart(ctx, req.md.UserId)
	if err != nil {
		s.log.Printf("failed to clear cart: %v", err)
		return core.NewErr(
			core.ErrInternal,
			"DELETE_CART_FAILURE",
			"failed to delete cart",
			nil,
		)
	}
	return nil
}
