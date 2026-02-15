package helpers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"

	"backend-commerce/config"
	"backend-commerce/models"

)

func GetSnapToken(order models.Order, user models.User) (string, string, error) {
	serverKey := config.GetEnv("MIDTRANS_SERVER_KEY", "")
	isProd := config.GetEnv("MIDTRANS_IS_PRODUCTION", "false") == "true"

	var s = snap.Client{}
	env := midtrans.Sandbox
	if isProd {
		env = midtrans.Production
	}
	s.New(serverKey, env)

	var calculatedGrossAmount int64 = 0
	var items []midtrans.ItemDetails
	for _, item := range order.Items {
		name := item.Product.Name
		if len(name) > 50 {
			name = name[:50]
		}
		items = append(items, midtrans.ItemDetails{
			ID:    fmt.Sprintf("PROD-%d", item.ProductId),
			Name:  name,
			Price: int64(item.Price),
			Qty:   int32(item.Quantity),
		})
		calculatedGrossAmount += int64(item.Price) * int64(item.Quantity)
	}

	if order.ShippingCost > 0 {
		items = append(items, midtrans.ItemDetails{
			ID:    "SHIPPING",
			Name:  fmt.Sprintf("Shipping (%s - %s)", order.Courier, order.Service),
			Price: int64(order.ShippingCost),
			Qty:   1,
		})
		calculatedGrossAmount += int64(order.ShippingCost)
	}

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.Id,
			GrossAmt: calculatedGrossAmount,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
		Items: &items,
	}
	snapResp, err := s.CreateTransaction(req)
	if err != nil {
		return "", "", err
	}

	return snapResp.Token, snapResp.RedirectURL, nil
}

func VerifySignature(orderId, statusCode, grossAmount, signatureKey string) bool {
	serverKey := config.GetEnv("MIDTRANS_SERVER_KEY", "")
	signatureString := orderId + statusCode + grossAmount + serverKey

	hasher := sha512.New()
	hasher.Write([]byte(signatureString))
	expectedSignature := hex.EncodeToString(hasher.Sum(nil))

	return expectedSignature == signatureKey
}
