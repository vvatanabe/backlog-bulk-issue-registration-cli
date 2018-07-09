package internal

type Callback func(*CallbackOptions)

func NewDefaultCallbackOptions() *CallbackOptions {
	return &CallbackOptions{
		Each: func() {},
		Before: func() {},
		After: func() {},
	}
}

type CallbackOptions struct {
	Each   func()
	Before func()
	After  func()
}

func Before(f func()) Callback {
	return func(c *CallbackOptions) {
		c.Before = f
	}
}

func Each(f func()) Callback {
	return func(c *CallbackOptions) {
		c.Each = f
	}
}

func After(f func()) Callback {
	return func(c *CallbackOptions) {
		c.After = f
	}
}