package mpesa

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Thorium234/afritechonline/backend/internal/models"
	"github.com/Thorium234/afritechonline/backend/internal/payments"
)

// Config holds M-Pesa API credentials.
type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	Shortcode      string
	Passkey        string
	Environment    string // "sandbox" or "production"
	CallbackURL    string
}

// Client handles M-Pesa API operations.
type Client struct {
	config       *Config
	httpClient   *http.Client
	accessToken  string
	tokenExpires time.Time
}

// NewClient creates a new M-Pesa API client.
func NewClient(config *Config) *Client {
	return &Client{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// STKPushRequest represents an M-Pesa STK Push request.
type STKPushRequest struct {
	Amount           float64 `json:"amount"`
	PhoneNumber      string  `json:"phone_number"`
	AccountReference string  `json:"account_reference"`
	TransactionDesc  string  `json:"transaction_desc"`
}

// STKPushResponse represents an M-Pesa STK Push response.
type STKPushResponse struct {
	MerchantRequestID   string `json:"merchant_request_id"`
	CheckoutRequestID   string `json:"checkout_request_id"`
	ResponseCode        string `json:"response_code"`
	ResponseDescription string `json:"response_description"`
	CustomerMessage     string `json:"customer_message"`
}

// CallbackPayload represents an M-Pesa callback payload.
type CallbackPayload struct {
	Body struct {
		stkCallback struct {
			MerchantRequestID string `json:"merchant_request_id"`
			CheckoutRequestID string `json:"checkout_request_id"`
			ResultCode        int    `json:"result_code"`
			ResultDesc        string `json:"result_desc"`
			CallbackMetadata  struct {
				Item []struct {
					Name  string `json:"name"`
					Value interface{} `json:"value"`
				} `json:"Item"`
			} `json:"callback_metadata"`
		} `json:"stkCallback"`
	} `json:"Body"`
}

// InitiateSTKPush starts an M-Pesa STK Push payment.
func (c *Client) InitiateSTKPush(ctx context.Context, req STKPushRequest) (*STKPushResponse, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	password, timestamp := generateSTKPassword(c.config.Shortcode, c.config.Passkey)

	payload := map[string]interface{}{
		"BusinessShortCode": c.config.Shortcode,
		"Password":         password,
		"Timestamp":        timestamp,
		"TransactionType":  "CustomerPayBillOnline",
		"Amount":           req.Amount,
		"PartyA":           req.PhoneNumber,
		"PartyB":           c.config.Shortcode,
		"PhoneNumber":      req.PhoneNumber,
		"CallBackURL":      c.config.CallbackURL,
		"AccountReference": req.AccountReference,
		"TransactionDesc":  req.TransactionDesc,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s/mpesa/stkpush/v1/processrequest", c.getBaseURL())
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var stkResp STKPushResponse
	if err := json.Unmarshal(respBody, &stkResp); err != nil {
		return nil, err
	}

	return &stkResp, nil
}

// ValidateCallback validates the M-Pesa callback signature.
func (c *Client) ValidateCallback(body []byte, signature string) bool {
	h := hmac.New(md5.New, []byte(c.config.ConsumerSecret))
	h.Write(body)
	expected := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

// ParseCallback extracts payment details from a callback payload.
func (c *Client) ParseCallback(body []byte) (*models.Payment, error) {
	var payload CallbackPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	resultCode := payload.Body.stkCallback.ResultCode
	status := payments.StatusCompleted
	if resultCode != 0 {
		status = payments.StatusFailed
	}

	// Extract metadata
	var mpesaReceipt string
	for _, item := range payload.Body.stkCallback.CallbackMetadata.Item {
		switch item.Name {
		case "PhoneNumber":
			_ = fmt.Sprintf("%v", item.Value)
		case "MpesaReceiptNumber":
			mpesaReceipt = fmt.Sprintf("%v", item.Value)
		}
	}

	return &models.Payment{
		Method:    payments.PaymentMethodMPesa,
		Reference: mpesaReceipt,
		Status:    status,
	}, nil
}

func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	if c.accessToken != "" && time.Now().Before(c.tokenExpires) {
		return c.accessToken, nil
	}

	url := fmt.Sprintf("%s/oauth/v1/generate?grant_type=client_credentials&client_id=%s&client_secret=%s",
		c.getBaseURL(), c.config.ConsumerKey, c.config.ConsumerSecret)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenExpires = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)
	return c.accessToken, nil
}

func (c *Client) getBaseURL() string {
	if c.config.Environment == "production" {
		return "https://api.safaricom.co.ke"
	}
	return "https://sandbox.safaricom.co.ke"
}

func generateSTKPassword(shortcode, passkey string) (string, string) {
	timestamp := time.Now().Format("20060102150405")
	data := shortcode + passkey + timestamp
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:]), timestamp
}
