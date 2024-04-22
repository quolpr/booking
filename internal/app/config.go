package app

type config struct {
	port string
}

func newConfig() *config {
	return &config{
		port: "8080",
	}
}
