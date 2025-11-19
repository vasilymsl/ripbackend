package ds

import (
	"fmt"
	"math"
)

// CreditCalculator содержит формулы для расчета параметров кредита и оценки кредитоспособности

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

// =======================================================================
// СКОРИНГОВАЯ МОДЕЛЬ ОЦЕНКИ КРЕДИТОСПОСОБНОСТИ ЗАЕМЩИКА
// =======================================================================

// ScoringResult содержит результат оценки кредитоспособности заемщика
type ScoringResult struct {
	CreditScore     int     // Кредитный скор (0-1000)
	Result          string  // approved / rejected / pending
	MaxCreditAmount float64 // Максимальная рекомендованная сумма кредита
	RejectionReason string  // Причина отклонения (если rejected)
	DTI             float64 // Debt-to-Income ratio (для справки)
}

// CalculateCreditworthiness выполняет оценку кредитоспособности заемщика
// на основе скоринговой модели с использованием коэффициента DTI (Debt-to-Income)
//
// Параметры:
//   - income: ежемесячный доход заемщика
//   - obligations: текущие ежемесячные обязательства заемщика
//   - requestedMonthlyPayment: запрашиваемый ежемесячный платеж по новым кредитам
//
// Скоринговая модель:
//
//	DTI = ((Обязательства + Новый платеж) / Доход) * 100%
//
//	Критерии оценки:
//	- DTI ≤ 30%  → Отличная кредитоспособность (скор 800-1000)
//	- DTI 30-40% → Хорошая кредитоспособность (скор 600-799)
//	- DTI 40-50% → Удовлетворительная (скор 400-599)
//	- DTI > 50%  → Отклонено (скор 0-399)
//
//	Максимальная сумма кредита рассчитывается как:
//	MaxCredit = (Доход - Обязательства) * 0.4 * 12 месяцев
func CalculateCreditworthiness(income, obligations, requestedMonthlyPayment float64) ScoringResult {
	result := ScoringResult{
		CreditScore:     0,
		Result:          "pending",
		MaxCreditAmount: 0,
		RejectionReason: "",
		DTI:             0,
	}

	// Валидация входных данных
	if income <= 0 {
		result.Result = "rejected"
		result.RejectionReason = "Доход должен быть больше нуля"
		return result
	}

	// Рассчитываем коэффициент долговой нагрузки (DTI - Debt-to-Income)
	totalMonthlyDebt := obligations + requestedMonthlyPayment
	dti := (totalMonthlyDebt / income) * 100
	result.DTI = dti

	// Рассчитываем максимальную рекомендованную сумму кредита
	// Формула: (Доход - Текущие обязательства) * 40% * 12 месяцев
	freeIncome := income - obligations
	if freeIncome > 0 {
		result.MaxCreditAmount = freeIncome * 0.4 * 12
	}

	// Скоринговая модель на основе DTI
	switch {
	case dti <= 30:
		// Отличная кредитоспособность
		result.CreditScore = 800 + int((30-dti)/30*200) // 800-1000
		result.Result = "approved"

	case dti > 30 && dti <= 40:
		// Хорошая кредитоспособность
		result.CreditScore = 600 + int((40-dti)/10*200) // 600-799
		result.Result = "approved"

	case dti > 40 && dti <= 50:
		// Удовлетворительная кредитоспособность
		result.CreditScore = 400 + int((50-dti)/10*200) // 400-599
		result.Result = "approved"
		result.RejectionReason = "Высокая долговая нагрузка. Рекомендуется уменьшить сумму кредита."

	case dti > 50:
		// Неудовлетворительная кредитоспособность
		result.CreditScore = int(math.Max(0, 400-((dti-50)*10))) // 0-399
		result.Result = "rejected"
		result.RejectionReason = "Превышен допустимый уровень долговой нагрузки (DTI > 50%). Текущая DTI: " +
			formatFloat(dti, 1) + "%"
	}

	// Дополнительная проверка: если свободного дохода недостаточно
	if freeIncome < requestedMonthlyPayment && result.Result == "approved" {
		result.Result = "rejected"
		result.RejectionReason = "Недостаточно свободного дохода для погашения кредита"
		result.CreditScore = int(math.Min(float64(result.CreditScore), 300))
	}

	return result
}

// formatFloat форматирует float64 с заданной точностью
func formatFloat(val float64, precision int) string {
	format := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(format, val)
}
