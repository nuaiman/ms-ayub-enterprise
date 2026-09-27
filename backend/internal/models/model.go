package models

import (
	"database/sql"
)

type Models struct {
	User             UserModel
	Log              LogModel
	Salary           SalaryModel
	Broker           BrokerModel
	Majhi            MajhiModel
	Godown           GodownModel
	Rent             RentModel
	Customer         CustomerModel
	Lot              LotModel
	Store            StoreModel
	Damage           DamageModel
	Delivery         DeliveryModel
	DeliveryItem     DeliveryItemModel
	Transport        TransportModel
	Vehicle          VehicleModel
	Expense          ExpenseModel
	AdditionalCharge AdditionalChargeModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		User:             UserModel{DB: db},
		Log:              LogModel{DB: db},
		Salary:           SalaryModel{DB: db},
		Broker:           BrokerModel{DB: db},
		Majhi:            MajhiModel{DB: db},
		Godown:           GodownModel{DB: db},
		Rent:             RentModel{DB: db},
		Customer:         CustomerModel{DB: db},
		Lot:              LotModel{DB: db},
		Store:            StoreModel{DB: db},
		Damage:           DamageModel{DB: db},
		Delivery:         DeliveryModel{DB: db},
		DeliveryItem:     DeliveryItemModel{DB: db},
		Transport:        TransportModel{DB: db},
		Vehicle:          VehicleModel{DB: db},
		Expense:          ExpenseModel{DB: db},
		AdditionalCharge: AdditionalChargeModel{DB: db},
	}
}
