package freedom

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AlexChe360/procash/internal/config"
)

type PaymentResponse struct {
	Status      string `xml:"pg_status"`
	PaymentId   string `xml:"pg_payment_id"`
	RedirectURL string `xml:"pg_redirect_url"`
	Sig         string `xml:"pg_sig"`
}

func GenerateURL(cfg config.Config, amount int, description string) (map[string]string, error) {
	merchantID := cfg.MerchantID
	secretKey := cfg.PaymentSecretKey

	userId := cfg.PaymentUserId
	paymentURL := cfg.PaymentURL

	if merchantID == "" || secretKey == "" {
		return nil, errors.New("missing merchant_id or secret_key")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderId := fmt.Sprintf("%05d", r.Intn(100000))
	salt := fmt.Sprintf("%x", rand.Int63())

	params := url.Values{
		"pg_order_id":      {orderId},
		"pg_merchant_id":   {merchantID},
		"pg_amount":        {fmt.Sprintf("%d", amount)},
		"pg_description":   {description},
		"pg_salt":          {salt},
		"pg_payment_route": {"frame"},
		"pg_user_id":       {userId},
	}

	sig := generateSignature(params, secretKey)
	params.Set("pg_sig", sig)

	log.Println("[FreedomPay] Sending request:")
	for k := range params {
		log.Printf("%s = %s\n", k, params.Get(k))
	}

	resp, err := http.PostForm(paymentURL, params)
	if err != nil {
		return nil, fmt.Errorf("send_request error: %w", err)
	}
	defer resp.Body.Close()

	return parseResponse(resp)
}

func generateSignature(params url.Values, secret string) string {
	order := []string{
		"pg_amount",
		"pg_description",
		"pg_merchant_id",
		"pg_order_id",
		"pg_payment_route",
		"pg_salt",
		"pg_user_id",
	}

	values := []string{"init_payment.php"}
	for _, k := range order {
		values = append(values, params.Get(k))
	}
	values = append(values, secret)

	data := strings.Join(values, ";")
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func parseResponse(resp *http.Response) (map[string]string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println("[FreedomPay] Response body:", string(body))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response error: %s", resp.Status)
	}

	var r PaymentResponse
	if err := xml.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("xml parse error: %s", err)
	}

	return map[string]string{
		"status":       r.Status,
		"payment_id":   r.PaymentId,
		"redirect_url": r.RedirectURL,
		"sig":          r.Sig,
	}, nil
}
