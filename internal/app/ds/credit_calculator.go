package ds

import "math"

// CreditCalculator содержит формулы для расчета параметров кредита
// ВНИМАНИЕ: Эти функции пока не используются в приложении, но могут быть подключены позже

// CalculateMonthlyPayment рассчитывает ежемесячный аннуитетный платеж по кредиту
// Формула: M = S * (r * (1 + r)^n) / ((1 + r)^n - 1)
// где:
//   - S (amount) - сумма кредита (основной долг)
//   - r (monthlyRate) - месячная процентная ставка (годовая ставка / 12 / 100)
//   - n (months) - срок кредита в месяцах
//
// Пример использования:
//
//	monthlyPayment := CalculateMonthlyPayment(1000000, 17.9, 36)
//	// для кредита 1 млн рублей под 17.9% годовых на 36 месяцев
func CalculateMonthlyPayment(amount float64, annualRate float64, months int) float64 {
	// Преобразуем годовую ставку в месячную (в десятичной форме)
	monthlyRate := annualRate / 12 / 100

	// Если ставка = 0, то просто делим сумму на количество месяцев
	if monthlyRate == 0 {
		return amount / float64(months)
	}

	// Формула аннуитетного платежа
	// M = S * (r * (1 + r)^n) / ((1 + r)^n - 1)
	coefficient := math.Pow(1+monthlyRate, float64(months))
	monthlyPayment := amount * (monthlyRate * coefficient) / (coefficient - 1)

	return monthlyPayment
}

// CalculateTotalPayment рассчитывает общую сумму выплат по кредиту
func CalculateTotalPayment(monthlyPayment float64, months int) float64 {
	return monthlyPayment * float64(months)
}

// CalculateOverpayment рассчитывает переплату по кредиту
func CalculateOverpayment(totalPayment float64, amount float64) float64 {
	return totalPayment - amount
}

// CalculateTotalMonthlyPaymentForApplication рассчитывает суммарный ежемесячный платеж
// по всем кредитам в заявке
// Это основная формула для предметной области - показывает клиенту общую ежемесячную нагрузку
func CalculateTotalMonthlyPaymentForApplication(credits []struct {
	Amount     float64
	AnnualRate float64
	Months     int
}) float64 {
	totalMonthly := 0.0

	for _, credit := range credits {
		monthlyPayment := CalculateMonthlyPayment(credit.Amount, credit.AnnualRate, credit.Months)
		totalMonthly += monthlyPayment
	}

	return totalMonthly
}
