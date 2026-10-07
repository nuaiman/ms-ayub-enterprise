package models

import (
	"database/sql"
)

type Models struct {
	User                   UserModel
	Log                    LogModel
	Salary                 SalaryModel
	Broker                 BrokerModel
	Majhi                  MajhiModel
	Godown                 GodownModel
	Customer               CustomerModel
	Lot                    LotModel
	Store                  StoreModel
	StoreAdjustment        StoreAdjustmentModel
	StoreTransfer          StoreTransferModel
	LotTransfer            LotTransferModel
	GodownBill             GodownBillModel
	MajhiBill              MajhiBillModel
	CustomerStoreBill      CustomerStoreBillModel
	CustomerDeliveryBill   CustomerDeliveryBillModel
	CustomerAdditionalBill CustomerAdditionalBillModel
	BillPayment            BillPaymentModel
	Expense                ExpenseModel
	Income                 IncomeModel
	Delivery               DeliveryModel
	Damage                 DamageModel
	Invoice                InvoiceModel
	InvoiceItem            InvoiceItemModel
	Discount               DiscountModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		User:                   UserModel{DB: db},
		Log:                    LogModel{DB: db},
		Salary:                 SalaryModel{DB: db},
		Broker:                 BrokerModel{DB: db},
		Majhi:                  MajhiModel{DB: db},
		Godown:                 GodownModel{DB: db},
		Customer:               CustomerModel{DB: db},
		Lot:                    LotModel{DB: db},
		Store:                  StoreModel{DB: db},
		StoreAdjustment:        StoreAdjustmentModel{DB: db},
		StoreTransfer:          StoreTransferModel{DB: db},
		LotTransfer:            LotTransferModel{DB: db},
		GodownBill:             GodownBillModel{DB: db},
		MajhiBill:              MajhiBillModel{DB: db},
		CustomerStoreBill:      CustomerStoreBillModel{DB: db},
		CustomerDeliveryBill:   CustomerDeliveryBillModel{DB: db},
		CustomerAdditionalBill: CustomerAdditionalBillModel{DB: db},
		BillPayment:            BillPaymentModel{DB: db},
		Expense:                ExpenseModel{DB: db},
		Income:                 IncomeModel{DB: db},
		Delivery:               DeliveryModel{DB: db},
		Damage:                 DamageModel{DB: db},
		Invoice:                InvoiceModel{DB: db},
		InvoiceItem:            InvoiceItemModel{DB: db},
		Discount:               DiscountModel{DB: db},
	}
}
