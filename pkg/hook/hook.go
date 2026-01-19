package hook

type Hook func() error

type Registry struct {
	hooks []Hook
}

var (
	PrePrompt  Registry
	PostCommand Registry
)

func (r *Registry) Register(h Hook) {
	r.hooks = append(r.hooks, h)
}

func (r Registry) Run() error {
	for _, h := range r.hooks {
		if err := h(); err != nil {
			return err
		}
	}

	return nil
}
