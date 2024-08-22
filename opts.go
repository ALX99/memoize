package memoize

type (
	// MemoizerOption is a function that sets an option on the memoizer
	MemoizerOption func(*memoizer[any])
)

// WithCache sets the [Cache] to use
func WithCache(cache Cache) MemoizerOption {
	return func(o *memoizer[any]) {
		o.cache = cache
	}
}
