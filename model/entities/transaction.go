package entities

import "time"

type Transaction struct {
	AccountID       uint      `gorm:"column:account_id"`
	OperationTypeID uint      `gorm:"column:operation_type_id"`
	Amount          float64   `gorm:"column:amount"`
	EventDate       time.Time `gorm:"column:event_date;default:CURRENT_TIMESTAMP()"`
}

func (t Transaction) TableName() string {
	return "transactions"
}
