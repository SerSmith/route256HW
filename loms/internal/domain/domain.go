package domain

type OrderStatus string

const (
	NewStatus             OrderStatus = "new"
	AwaitingPaymentStatus OrderStatus = "awaiting payment"
	FailedStatus          OrderStatus = "failed"
	PayedStatus           OrderStatus = "payed"
	CanceledStatus        OrderStatus = "cancelled"
)

type ItemOrder struct {
	SKU   uint32
	Count uint16
}

type Order struct {
	User   int64
	Items  []*ItemOrder
	Status OrderStatus
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}

type Model struct {
}

func New() *Model {
	return &Model{}
}
