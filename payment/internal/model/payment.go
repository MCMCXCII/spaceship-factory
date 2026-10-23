package model

// PaymentMethod — способ оплаты заказа
type PaymentMethod int

const (
	// PaymentMethodUnspecified — способ оплаты не указан
	PaymentMethodUnspecified PaymentMethod = iota
	// PaymentMethodCard — банковская карта
	PaymentMethodCard
	// PaymentMethodSBP — система быстрых платежей
	PaymentMethodSBP
	// PaymentMethodCreditCard — кредитная карта
	PaymentMethodCreditCard
	// PaymentMethodInvestorMoney — деньги инвестора
	PaymentMethodInvestorMoney
)

// String возвращает имя способа оплаты для логов
func (m PaymentMethod) String() string {
	switch m {
	case PaymentMethodCard:
		return "CARD"
	case PaymentMethodSBP:
		return "SBP"
	case PaymentMethodCreditCard:
		return "CREDIT_CARD"
	case PaymentMethodInvestorMoney:
		return "INVESTOR_MONEY"
	default:
		return "UNSPECIFIED"
	}
}
