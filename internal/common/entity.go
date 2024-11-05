package common

type Request struct {
	UserID    string
}

type SearchRequest struct {
	Limit  uint16
	Page   uint16
	Search string
}

type AllResponse struct {
	List  []any
	Total uint32
}
