package models

type IndexData struct {
	Account         *User
	IsAuthenticated bool
}

func NewUnauthenticatedIndexData() IndexData {
	return IndexData{Account: nil, IsAuthenticated: false}
}
