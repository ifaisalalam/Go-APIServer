package shortener

// CreateShortURLInput is the input for creating a new short URL.
type CreateShortURLInput struct {
	LongURL  string
	ShortURL string
}

// CreateShortURLOutput contains the CreateShortURLInput.ShortURL if it was created successfully.
type CreateShortURLOutput struct {
	ShortURL string
}

// GetTargetURLInput contains the ShortURL to be used for retrieving the GetTargetURLOutput.
type GetTargetURLInput struct {
	ShortURL string
}

// GetTargetURLOutput is the result containing the LongURL for the GetTargetURLInput.ShortURL.
type GetTargetURLOutput struct {
	LongURL string
}
