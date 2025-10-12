package cart_test

import (
	"os"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/cart-overflow/cart-api/pkg/pb"
	"github.com/cart-overflow/cart/internal/app"
	"github.com/cart-overflow/cart/internal/test"
	"github.com/cart-overflow/common/pkg/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Integration(t *testing.T) {
	setupSut(t)

	// Set Product
	t.Run("set cart item returns valid product id", setCartItem_returnsId)

	// Get Product
	t.Run("get cart item returns empty if no products", getCartItem_returnsEmptyIfNoProducts)
	t.Run("get cart item returns stored item", getCartItem_returnsStoredProduct)

	// Clear Cart
	t.Run("clear cart doesn't return error if cart is empty", clearCart_noErrorIfCartIsClear)
	t.Run("clear cart clears cart", clearCart_clearsCart)

	// Get Cart

	t.Run("get cart returns empty cart", getCart_returnsEmptyCart)
	t.Run("get cart returns not empty", getCart_returnsNotEmpty)
}

// MARK: - Set Product

func setCartItem_returnsId(t *testing.T) {
	test.CleanCache(t)

	ctx := metadata.New().WithUserId("user-id-1").RpcCtx(t.Context())
	res, err := sut.cl.SetCartItem(ctx, &pb.SetCartItemRequest{
		ProductId: "product-id-1",
		Amount:    2,
	})

	require.NoError(t, err)
	assert.Equal(t, "product-id-1", res.ProductId)
}

// MARK: - Get Product

func getCartItem_returnsEmptyIfNoProducts(t *testing.T) {
	test.CleanCache(t)

	ctx := metadata.New().WithUserId("user-id-1").RpcCtx(t.Context())
	res, err := sut.cl.GetCartItem(ctx, &pb.GetCartItemRequest{
		ProductId: "product-id-1",
	})
	require.NoError(t, err)

	assert.Equal(t, "product-id-1", res.Item.ProductId)
	assert.Equal(t, int64(0), res.Item.Amount)
	assert.Nil(t, res.Item.AddedAt)
}

func getCartItem_returnsStoredProduct(t *testing.T) {
	test.CleanCache(t)
	cleanStubTime()
	stubTime(time.UnixMilli(1))
	setCartItem(t, "user-id-1", "product-id-1", 3)
	stubTime(time.UnixMilli(2))
	setCartItem(t, "user-id-1", "product-id-1", 2)
	stubTime(time.UnixMilli(3))
	setCartItem(t, "user-id-1", "product-id-2", 1)
	stubTime(time.UnixMilli(4))
	setCartItem(t, "user-id-2", "product-id-3", 4)

	ctx := metadata.New().WithUserId("user-id-1").RpcCtx(t.Context())
	res, err := sut.cl.GetCartItem(ctx, &pb.GetCartItemRequest{
		ProductId: "product-id-1",
	})
	require.NoError(t, err)
	assert.Equal(t, "product-id-1", res.Item.ProductId)
	assert.Equal(t, int64(2), res.Item.Amount)
	assert.Equal(t, int64(1000), *res.Item.AddedAt)

	ctx = metadata.New().WithUserId("user-id-1").RpcCtx(t.Context())
	res, err = sut.cl.GetCartItem(ctx, &pb.GetCartItemRequest{
		ProductId: "product-id-2",
	})
	require.NoError(t, err)
	assert.Equal(t, "product-id-2", res.Item.ProductId)
	assert.Equal(t, int64(1), res.Item.Amount)
	assert.Equal(t, int64(3000), *res.Item.AddedAt)

	ctx = metadata.New().WithUserId("user-id-2").RpcCtx(t.Context())
	res, err = sut.cl.GetCartItem(ctx, &pb.GetCartItemRequest{
		ProductId: "product-id-3",
	})
	require.NoError(t, err)
	assert.Equal(t, "product-id-3", res.Item.ProductId)
	assert.Equal(t, int64(4), res.Item.Amount)
	assert.Equal(t, int64(4000), *res.Item.AddedAt)
}

// MARK: - Clear Cart

func clearCart_noErrorIfCartIsClear(t *testing.T) {
	test.CleanCache(t)

	ctx := metadata.New().WithUserId("user-id-1").RpcCtx(t.Context())
	_, err := sut.cl.ClearCart(ctx, &pb.ClearCartRequest{})

	require.NoError(t, err)
}

func clearCart_clearsCart(t *testing.T) {
	test.CleanCache(t)
	setCartItem(t, "user-id-1", "product-id-1", 1)
	setCartItem(t, "user-id-2", "product-id-2", 2)

	ctx := metadata.New().WithUserId("user-id-1").RpcCtx(t.Context())
	_, err := sut.cl.ClearCart(ctx, &pb.ClearCartRequest{})
	require.NoError(t, err)

	ctx = metadata.New().WithUserId("user-id-2").RpcCtx(t.Context())
	_, err = sut.cl.ClearCart(ctx, &pb.ClearCartRequest{})
	require.NoError(t, err)

	ctx = metadata.New().WithUserId("user-id-1").RpcCtx(t.Context())
	res, err := sut.cl.GetCartItem(ctx, &pb.GetCartItemRequest{
		ProductId: "product-id-1",
	})
	require.NoError(t, err)
	assert.Equal(t, "product-id-1", res.Item.ProductId)
	assert.Equal(t, int64(0), res.Item.Amount)

	ctx = metadata.New().WithUserId("user-id-2").RpcCtx(t.Context())
	res, err = sut.cl.GetCartItem(ctx, &pb.GetCartItemRequest{
		ProductId: "product-id-2",
	})
	require.NoError(t, err)
	assert.Equal(t, "product-id-2", res.Item.ProductId)
	assert.Equal(t, int64(0), res.Item.Amount)
}

// MARK: - Get Cart

func getCart_returnsEmptyCart(t *testing.T) {
	test.CleanCache(t)

	ctx := metadata.New().WithUserId("user-id-1").RpcCtx(t.Context())
	res, err := sut.cl.GetCart(ctx, &pb.GetCartRequest{})

	require.NoError(t, err)
	assert.Empty(t, res.Items)
}

func getCart_returnsNotEmpty(t *testing.T) {
	test.CleanCache(t)
	cleanStubTime()
	stubTime(time.UnixMilli(1))
	setCartItem(t, "user-id-1", "product-id-1", 3)
	stubTime(time.UnixMilli(2))
	setCartItem(t, "user-id-1", "product-id-2", 5)
	stubTime(time.UnixMilli(3))
	setCartItem(t, "user-id-1", "product-id-1", 7)
	stubTime(time.UnixMilli(4))
	setCartItem(t, "user-id-2", "product-id-1", 5)
	stubTime(time.UnixMilli(5))
	setCartItem(t, "user-id-1", "product-id-3", 10)
	stubTime(time.UnixMilli(6))
	setCartItem(t, "user-id-1", "product-id-3", 0)

	ctx := metadata.New().WithUserId("user-id-1").RpcCtx(t.Context())
	res, err := sut.cl.GetCart(ctx, &pb.GetCartRequest{})

	require.NoError(t, err)
	require.Len(t, res.Items, 2)
	AssertExists(t, res.Items, func(i *pb.CartItem) bool {
		return i.ProductId == "product-id-1" && i.Amount == 7 && *i.AddedAt == int64(1000)
	})
	AssertExists(t, res.Items, func(i *pb.CartItem) bool {
		return i.ProductId == "product-id-2" && i.Amount == 5 && *i.AddedAt == int64(2000)
	})
}

// MARK: - System Under Test

type systemUnderTest struct {
	cl pb.CartServiceClient
}

var sut *systemUnderTest
var stubTime func(t time.Time)
var cleanStubTime func()

func setupSut(t *testing.T) {
	var currentTime time.Time
	stubTime = func(t time.Time) { currentTime = t }
	cleanStubTime = func() { currentTime = time.UnixMicro(0) }

	wg := &sync.WaitGroup{}
	started := make(chan struct{})
	wg.Go(func() {
		app.Run(t.Context(),
			app.Deps{
				Getenv: func(key string) string {
					env := map[string]string{
						app.AddrKey:       "127.0.0.1:5052",
						app.ValkeyAddrKey: test.ValkeyAddr,
					}
					return env[key]
				},
				Now:     func() time.Time { return currentTime },
				Logwr:   os.Stdout,
				Started: started,
			})
	})
	t.Cleanup(func() {
		wg.Wait()
	})

	conn, err := grpc.NewClient(
		"127.0.0.1:5052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "can not connect to grpc server")
	t.Cleanup(func() {
		err := conn.Close()
		if err != nil {
			t.Logf("error on conn close: %v", err)
		}
	})

	sut = &systemUnderTest{pb.NewCartServiceClient(conn)}
	<-started
}

// MARK: Helpers

func setCartItem(t *testing.T, userId, productId string, amount int64) {
	ctx := metadata.New().WithUserId(userId).RpcCtx(t.Context())
	_, err := sut.cl.SetCartItem(ctx, &pb.SetCartItemRequest{
		ProductId: productId,
		Amount:    amount,
	})
	require.NoError(t, err)
}

func AssertExists[T any](t *testing.T, slice []T, predicate func(T) bool) {
	assert.True(t, slices.ContainsFunc(slice, predicate))
}
