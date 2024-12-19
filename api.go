package main

type APIserver struct {
	listenAddr string
}

func NewAPIserver (listenAddr string) *APIserver {
	return &APIserver{
		listenAddr: listenAddr,
	}
}