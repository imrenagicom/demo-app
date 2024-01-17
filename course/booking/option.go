package booking

import "github.com/jmoiron/sqlx"

type FindOptions struct {
	Tx           *sqlx.Tx
	DisableCache bool
}

type FindOption func(*FindOptions)

func WithFindTx(tx *sqlx.Tx) FindOption {
	return func(o *FindOptions) {
		o.Tx = tx
	}
}

func WithDisableCache() FindOption {
	return func(o *FindOptions) {
		o.DisableCache = true
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

type CreateOptions struct {
	Tx *sqlx.Tx
}

type CreateOption func(*CreateOptions)

func WithCreateTx(tx *sqlx.Tx) CreateOption {
	return func(o *CreateOptions) {
		o.Tx = tx
	}
}

type ListOptions struct {
	Tx            *sqlx.Tx
	Limit         uint64
	Page          uint64
	InvoiceNumber string
	Status        Status
}

func (f ListOptions) GetOffset() uint64 {
	return f.Page * f.Limit
}

type ListOption func(*ListOptions)

func WithFindAllTx(tx *sqlx.Tx) ListOption {
	return func(o *ListOptions) {
		o.Tx = tx
	}
}

func WithFindAllInvoiceNumber(invoiceNumber string) ListOption {
	return func(o *ListOptions) {
		o.InvoiceNumber = invoiceNumber
	}
}

func WithFindAllStatus(status Status) ListOption {
	return func(o *ListOptions) {
		o.Status = status
	}
}
