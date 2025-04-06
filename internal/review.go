package internal

type Review struct {
	Body        string `json:"body"`
	SubmittedAt string `json:"submitted_at"`
	State       string `json:"state"`
	User        struct {
		Login string `json:"login"`
	} `json:"user"`
}
