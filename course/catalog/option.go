package catalog

import (
	"encoding/base64"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ListOptions struct {
	Limit   uint64
	Page    uint64
	Preload bool
}

func (f ListOptions) GetOffset() uint64 {
	return f.Page * f.Limit
}

type ListOption func(*ListOptions)

func WithMaxResults(limit uint64) ListOption {
	return func(o *ListOptions) {
		if limit > 0 {
			o.Limit = limit
		}
	}
}

func WithPreload() ListOption {
	return func(o *ListOptions) {
		o.Preload = true
	}
}

func WithNextPage(nextPage string) ListOption {
	return func(o *ListOptions) {
		if nextPage != "" {
			pt, _ := decode(nextPage)
			o.Page = pt.page
		}
	}
}

type pageToken struct {
	page uint64
}

func (p pageToken) encode() string {
	token := fmt.Sprintf("%d", p.page)
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func decode(token string) (pageToken, error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return pageToken{}, err
	}
	var page uint64
	if _, err := fmt.Sscanf(string(decoded), "%d", &page); err != nil {
		return pageToken{}, err
	}
	return pageToken{page}, nil
}

type FindOptions struct {
	Tx *sqlx.Tx
}

type FindOption func(*FindOptions)

func WithFindTx(tx *sqlx.Tx) FindOption {
	return func(o *FindOptions) {
		o.Tx = tx
	}
}

type UpdateOptions struct {
	Tx *sqlx.Tx
}

type UpdateOption func(*UpdateOptions)

func WithUpdateTx(tx *sqlx.Tx) UpdateOption {
	return func(o *UpdateOptions) {
		o.Tx = tx
	}
}
