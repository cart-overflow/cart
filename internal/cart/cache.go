package cart

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/valkey-io/valkey-go"
)

type cache struct {
	cl  valkey.Client
	now func() time.Time
	log *log.Logger
}

func newCache(cl valkey.Client, now func() time.Time, log *log.Logger) *cache {
	return &cache{cl, now, log}
}

type cachedCartItem struct {
	Amount  int64  `json:"amount"`
	AddedAt *int64 `json:"added_at"`
}

func (c *cache) setCartItem(
	ctx context.Context,
	userId string,
	productId string,
	amount int64,
) error {
	cl, cancel := c.cl.Dedicate()
	defer cancel()

	resp := cl.Do(ctx, cl.B().Watch().Key(cartKey(userId)).Build())
	err := resp.Error()
	if err != nil {
		return err
	}

	item := cachedCartItem{}
	resp = cl.Do(ctx, cl.B().Hget().Key(cartKey(userId)).Field(productId).Build())
	if resp.Error() == valkey.Nil {
		addedAt := c.now().UnixMicro()
		item.AddedAt = &addedAt
	} else {
		err := resp.DecodeJSON(&item)
		if err != nil {
			return err
		}
	}
	item.Amount = amount

	jsonItem, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal cart item: %v", err)
	}

	resps := cl.DoMulti(
		ctx,
		cl.B().Multi().Build(),
		cl.B().
			Hset().
			Key(cartKey(userId)).
			FieldValue().
			FieldValue(productId, valkey.BinaryString(jsonItem)).
			Build(),
		cl.B().Exec().Build(),
	)

	for _, r := range resps {
		err = r.Error()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cache) getCartItem(
	ctx context.Context,
	userId string,
	productId string,
) (*CartItem, error) {
	resp := c.cl.Do(ctx, c.cl.B().Hget().Key(cartKey(userId)).Field(productId).Build())
	err := resp.Error()
	if err == valkey.Nil {
		return &CartItem{
			amount:    0,
			productId: productId,
		}, nil
	}

	var item cachedCartItem
	err = resp.DecodeJSON(&item)
	if err != nil {
		return nil, err
	}

	return &CartItem{
		productId: productId,
		amount:    item.Amount,
		added_at:  item.AddedAt,
	}, nil
}

func (c *cache) deleteCartItem(ctx context.Context, userId string, productId string) error {
	return c.cl.Do(ctx, c.cl.B().Hdel().Key(cartKey(userId)).Field(productId).Build()).Error()
}

func (c *cache) clearCart(ctx context.Context, userId string) error {
	return c.cl.Do(ctx, c.cl.B().Del().Key(cartKey(userId)).Build()).Error()
}

func (c *cache) getCart(ctx context.Context, userId string) ([]*CartItem, error) {
	resp := c.cl.Do(ctx, c.cl.B().Hgetall().Key(cartKey(userId)).Build())
	err := resp.Error()
	if err == valkey.Nil {
		return []*CartItem{}, nil
	}

	rawCart, err := resp.ToMap()
	if err != nil {
		return nil, err
	}

	cart := make([]*CartItem, 0, len(rawCart))
	for key, msg := range rawCart {
		jsonItem, err := msg.AsBytes()
		if err != nil {
			return nil, err
		}

		var item cachedCartItem
		err = json.Unmarshal(jsonItem, &item)
		if err != nil {
			return nil, err
		}

		cart = append(cart, &CartItem{
			productId: key,
			amount:    item.Amount,
			added_at:  item.AddedAt,
		})
	}

	return cart, nil
}

// MARK: Helpres

func cartKey(userId string) string {
	return "cart:" + userId
}
