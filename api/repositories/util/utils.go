package util

type StoreResult struct {
	Data  interface{}
	Total int
	Err   error
}

type StoreChannel chan StoreResult
