package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Order struct {
	ID          int
	Title       string
	Icon        string
	ImageURL    string
	Rate        string
	Term        string
	Amount      string
	Feature1    string
	Feature2    string
	Feature3    string
	Description string
}

type ApplicationProduct struct {
	Order          Order
	MonthlyPayment string
}

type Application struct {
	ID          string
	FullName    string
	Income      string
	Obligations string
	Products    []ApplicationProduct
}

func (r *Repository) GetOrders() ([]Order, error) {
	orders := []Order{
		{
			ID:          1,
			Title:       "Потребительский кредит",
			Icon:        "percent",
			ImageURL:    "http://localhost:9000/images/%20credit/1.jpg",
			Rate:        "Ставка: от 17.9% годовых",
			Term:        "Срок: до 36 мес",
			Amount:      "Сумма: до 1 000 000 ₽",
			Description: "Пользовательский кредит — универсальный продукт, который можно оформить на любые личные нужды без залога и поручителей. Подходит для финансирования путешествий, покупок техники, образования или непредвиденных расходов. Заявка подаётся онлайн, решение принимается в короткие сроки. Погашение производится равными ежемесячными платежами, доступно досрочное закрытие без штрафов.",
		},
		{
			ID:          2,
			Title:       "Кредит наличными",
			Icon:        "wallet",
			ImageURL:    "http://localhost:9000/images/%20credit/2.jpg",
			Rate:        "Ставка: от 15.9% годовых",
			Term:        "Срок: до 60 мес",
			Amount:      "Сумма: до 3 000 000 ₽",
			Description: "Кредит наличными — быстрое решение для получения денежных средств на карту или наличными. Оформление занимает всего несколько минут, а решение по заявке принимается мгновенно. Средства можно использовать на любые цели без отчета перед банком. Гибкие условия погашения позволяют выбрать комфортный график платежей. Возможна доставка денег курьером на дом или в офис.",
		},
		{
			ID:          3,
			Title:       "Кредит под залог на любые цели",
			Icon:        "bank",
			ImageURL:    "http://localhost:9000/images/%20credit/3.jpg",
			Rate:        "Ставка: от 12.5% годовых",
			Term:        "Срок: до 15 лет",
			Amount:      "Сумма: до 30 000 000 ₽",
			Description: "Кредит под залог недвижимости — выгодное решение для получения крупной суммы на длительный срок. В качестве залога принимается квартира, дом или коммерческая недвижимость. Минимальный пакет документов и упрощенная процедура оформления. Низкая процентная ставка благодаря наличию обеспечения. Возможность использования средств на любые цели: ремонт, образование, развитие бизнеса или рефинансирование других кредитов.",
		},
		{
			ID:          4,
			Title:       "Кредитная карта",
			Icon:        "card",
			ImageURL:    "http://localhost:9000/images/%20credit/4.jpg",
			Rate:        "Ставка: от 19.9% годовых",
			Term:        "Льготный период: до 60 дней",
			Amount:      "Лимит: до 600 000 ₽",
			Description: "Кредитная карта — удобный инструмент для ежедневных покупок и непредвиденных расходов. Льготный период до 60 дней позволяет пользоваться средствами банка без процентов. Кэшбэк на популярные категории покупок и бонусы за использование карты. Возможность снятия наличных в любом банкомате. Управление картой и лимитом через мобильное приложение. Бесплатное обслуживание при выполнении условий банка.",
		},
		{
			ID:          5,
			Title:       "Автокредит",
			Icon:        "car",
			ImageURL:    "http://localhost:9000/images/%20credit/5.webp",
			Rate:        "Ставка: от 16.5% годовых",
			Term:        "Срок: до 60 мес",
			Amount:      "Сумма: до 5 000 000 ₽",
			Description: "Автокредит — специальное предложение для покупки нового или подержанного автомобиля. Первоначальный взнос от 10% или возможность кредитования без первого взноса. Быстрое рассмотрение заявки и оформление сделки прямо в автосалоне. Страхование КАСКО и ОСАГО по желанию клиента. Возможность досрочного погашения без штрафов и комиссий. Индивидуальные условия для зарплатных клиентов и при покупке автомобиля у партнеров банка.",
		},
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return orders, nil
}

func (r *Repository) GetOrder(id int) (Order, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
	orders, err := r.GetOrders()
	if err != nil {
		return Order{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Order{}, fmt.Errorf("заказ не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}

func (r *Repository) GetOrdersByTitle(title string) ([]Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return []Order{}, err
	}

	var result []Order
	for _, order := range orders {
		if strings.Contains(strings.ToLower(order.Title), strings.ToLower(title)) {
			result = append(result, order)
		}
	}

	return result, nil
}

func (r *Repository) GetApplications() ([]Application, error) {
	// Получаем все заказы
	orders, err := r.GetOrders()
	if err != nil {
		return nil, err
	}

	// Создаём тестовые заявки
	applications := []Application{
		{
			ID:          "A-1001",
			FullName:    "Иванов Иван Иванович",
			Income:      "120 000 ₽",
			Obligations: "25 000 ₽",
			Products: []ApplicationProduct{
				{
					Order:          orders[0], // Потребительский кредит
					MonthlyPayment: "Ежемесячный платёж ≈ 35 000 ₽",
				},
				{
					Order:          orders[3], // Кредитная карта
					MonthlyPayment: "Ежемесячный платёж ≈ 8 500 ₽",
				},
			},
		},
		{
			ID:          "A-1002",
			FullName:    "Петров Петр Петрович",
			Income:      "85 000 ₽",
			Obligations: "15 000 ₽",
			Products: []ApplicationProduct{
				{
					Order:          orders[1], // Кредит наличными
					MonthlyPayment: "Ежемесячный платёж ≈ 18 000 ₽",
				},
				{
					Order:          orders[4], // Автокредит
					MonthlyPayment: "Ежемесячный платёж ≈ 22 000 ₽",
				},
			},
		},
		{
			ID:          "A-1003",
			FullName:    "Сидорова Анна Михайловна",
			Income:      "150 000 ₽",
			Obligations: "30 000 ₽",
			Products: []ApplicationProduct{
				{
					Order:          orders[2], // Кредит под залог
					MonthlyPayment: "Ежемесячный платёж ≈ 42 000 ₽",
				},
				{
					Order:          orders[3], // Кредитная карта
					MonthlyPayment: "Ежемесячный платёж ≈ 5 000 ₽",
				},
			},
		},
	}

	return applications, nil
}

func (r *Repository) GetApplication(id string) (Application, error) {
	applications, err := r.GetApplications()
	if err != nil {
		return Application{}, err
	}

	for _, app := range applications {
		if app.ID == id {
			return app, nil
		}
	}

	return Application{}, fmt.Errorf("заявка не найдена")
}
