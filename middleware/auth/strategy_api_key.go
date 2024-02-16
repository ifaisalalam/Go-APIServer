package auth

import (
	"net/http"
	"os"
)

const (
	ApiKeyHttpHeaderName = "X-API-KEY"
)

const (
	StrategyApiKeyInHeader StrategyName = "__/strategy-api-key-in-header"
)

var (
	ApiKeyValue = os.Getenv("DEFAULT_API_KEY")
)

func NewApiKeyNoOpStrategy() Strategy {
	return NewApiKeyStrategyWithValidator(noOpValidator)
}

func NewApiKeyStrategyWithValidator(provider Validator) Strategy {
	return &apiKeyStrategy{validator: provider}
}

type apiKeyStrategy struct {
	validator Validator
}

func (s *apiKeyStrategy) Authenticate(req *http.Request) (bool, error) {
	if token := req.Header.Get(ApiKeyHttpHeaderName); token != "" {
		return s.validator.Validate(token)
	}
	return false, nil
}

func (s *apiKeyStrategy) Is(name StrategyName) bool {
	return name == StrategyApiKeyInHeader
}

type Validator interface {
	Validate(apiKey string) (bool, error)
}

type validatorFunc func(apiKey string) (bool, error)

func (f validatorFunc) Validate(apiKey string) (bool, error) {
	return f(apiKey)
}

var noOpValidator = validatorFunc(func(apiKey string) (bool, error) {
	return true, nil
})
