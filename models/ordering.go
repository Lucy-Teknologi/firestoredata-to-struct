package models

import (
	"time"
)

type OrderGroup struct {
	UUID       string `firestore:"uuid"`
	BranchUUID string `firestore:"branch_uuid"`

	CreatedBy string `firestore:"created_by"`

	Space *SpaceInfo `firestore:"space"`
	// Delivery *DeliveryInfo `firestore:"delivery"`
	// Queue    *QueueInfo    `firestore:"queue"`

	Orders          map[string]Order         `firestore:"orders"`
	CancelledOrders []CancelledOrder         `firestore:"cancelled_orders"`
	Taxes           map[string]OrderTax      `firestore:"taxes"`
	Discounts       map[string]OrderDiscount `firestore:"discounts"`
	Coupons         []OrderCoupon            `firestore:"coupons"`

	CancelledAt  *time.Time `firestore:"cancelled_at"`
	CancelledBy  string     `firestore:"cancelled_by,omitempty"`
	CancelReason string     `firestore:"cancel_reason,omitempty"`

	CreatedAt time.Time  `firestore:"created_at"`
	DeletedAt *time.Time `firestore:"deleted_at"`
}

type CancelledOrder struct {
	// This is the cancel request uuid not the order uuid
	UUID           string `firestore:"uuid"`
	BranchUUID     string `firestore:"branch_uuid"`
	OrderGroupUUID string `firestore:"order_group_uuid"`
	Order          Order  `firestore:"order"`

	CancelledAt  time.Time `firestore:"cancelled_at"`
	CancelledBy  string    `firestore:"cancelled_by"`
	CancelReason string    `firestore:"cancel_reason"`

	Type             []string        `firestore:"type"`
	NewQuantity      int             `firestore:"new_quantity,omitempty"`
	NewVariant       *OrderVariant   `firestore:"new_variant,omitempty"`
	AddedModifiers   []OrderModifier `firestore:"added_modifiers,omitempty"`
	RemovedModifiers []OrderModifier `firestore:"removed_modifiers,omitempty"`

	Indexes  []string  `firestore:"indexes"`
	Approval *Approval `firestore:"approval"`
}

type Order struct {
	UUID      string `firestore:"uuid"`
	CreatedBy string `firestore:"created_by"`

	Item      OrderItem                `firestore:"item"`
	Modifiers map[string]OrderModifier `firestore:"modifiers"`
	Discounts map[string]OrderDiscount `firestore:"discounts"`

	Custom           bool         `firestore:"custom"`
	Takeaway         bool         `firestore:"takeaway"`
	Quantity         int          `firestore:"quantity"`
	RefundedQuantity int          `firestore:"refunded_quantity"`
	Note             string       `firestore:"note"`
	Waiter           *OrderWaiter `firestore:"waiter"`

	CreatedAt time.Time  `firestore:"created_at"`
	UpdatedAt time.Time  `firestore:"updated_at"`
	DeletedAt *time.Time `firestore:"deleted_at"`
}

type OrderItem struct {
	UUID         string        `firestore:"uuid"`
	Name         string        `firestore:"name"`
	CategoryName string        `firestore:"category_name"`
	CategoryUUID string        `firestore:"category_uuid"`
	Label        string        `firestore:"label"`
	Description  string        `firestore:"description"`
	ImagePath    *string       `firestore:"image_path"`
	Price        float64       `firestore:"price"`
	Variant      *OrderVariant `firestore:"variant"`

	Materials []OrderMaterial `firestore:"materials,omitempty"`
}

type OrderVariant struct {
	UUID        string  `firestore:"uuid"`
	Label       string  `firestore:"label"`
	ImagePath   *string `firestore:"image_path"`
	Description string  `firestore:"description"`
	Price       float64 `firestore:"price"`

	Materials []OrderMaterial `firestore:"materials,omitempty"`
}

type OrderModifier struct {
	UUID             string  `firestore:"uuid"`
	Name             string  `firestore:"name"`
	Quantity         int     `firestore:"quantity"`
	RefundedQuantity int     `firestore:"refunded_quantity"`
	Price            float64 `firestore:"price"`

	Materials []OrderMaterial `firestore:"materials,omitempty"`
}

type OrderWaiter struct {
	UUID string `firestore:"uuid"`
	Name string `firestore:"name"`
}

type SpaceInfo struct {
	ZoneUUID    string `firestore:"zone_uuid"`
	SpaceUUID   string `firestore:"space_uuid"`
	SpaceNumber int    `firestore:"space_number"`
	SpaceAlias  string `firestore:"space_alias"`
	ZoneName    string `firestore:"zone_name"`
	GroupCode   string `firestore:"group_code"`
}

// type QueueInfo struct {
// 	UUID     string        `firestore:"uuid"`
// 	Number   int           `firestore:"number"`
// 	Customer *CustomerInfo `firestore:"customer"`

// 	ScheduledAt *time.Time `firestore:"scheduled_at"`
// }

// type OldDeliveryInfo struct {
// 	Number   int           `firestore:"number"`
// 	Driver   string        `firestore:"driver"`
// 	Partner  string        `firestore:"partner"`
// 	Customer *CustomerInfo `firestore:"customer"`
// }

// type DeliveryInfo struct {
// 	Number   string        `firestore:"number"`
// 	Driver   string        `firestore:"driver"`
// 	Partner  string        `firestore:"partner"`
// 	Customer *CustomerInfo `firestore:"customer"`
// }

type OrderTax struct {
	UUID  string  `firestore:"uuid"`
	Name  string  `firestore:"name"`
	Value float32 `firestore:"value"`
}

type OrderDiscount struct {
	UUID    string  `firestore:"uuid"`
	Name    string  `firestore:"name"`
	Fixed   float64 `firestore:"fixed"`
	Percent float32 `firestore:"percent"`
}

type OrderCoupon struct {
	UUID       string `firestore:"uuid"`
	CouponUUID string `firestore:"coupon_uuid"`
	Name       string `firestore:"name"`
	Fixed      int64  `firestore:"fixed"`
}

type OrderMaterial struct {
	UUID              string `firestore:"uuid"`
	Name              string `firestore:"name"`
	MeasurementUUID   string `firestore:"measurement_uuid"`
	MeasurementAmount string `firestore:"measurement_amount"`
}

type Approval struct {
	By      string    `json:"by" firestore:"by"`
	At      time.Time `json:"at" firestore:"at"`
	Comment string    `json:"comment" firestore:"comment"`
}
