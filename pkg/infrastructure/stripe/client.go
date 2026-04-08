package stripe

import (
	"fmt"

	"github.com/stripe/stripe-go/v76"
	"github.com/your-org/jvairv2/configs"
)

// Client encapsula la configuración de Stripe
type Client struct {
	secretKey string
	publicKey string
}

// NewClient crea una nueva instancia del cliente de Stripe
func NewClient(cfg *configs.StripeConfig) (*Client, error) {
	if cfg.SecretKey == "" {
		return nil, fmt.Errorf("STRIPE_SECRET_KEY is required")
	}

	// Configurar la API key global del SDK
	stripe.Key = cfg.SecretKey

	return &Client{
		secretKey: cfg.SecretKey,
		publicKey: cfg.PublicKey,
	}, nil
}

// GetPublicKey retorna la clave pública de Stripe
func (c *Client) GetPublicKey() string {
	return c.publicKey
}
