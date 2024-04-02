# go-lknpd
Неофициальный клиент HTTP-API сайта [Мой налог для самозанятых](https://lknpd.nalog.ru/).
Используется метод логина через мобильный телефон. Для генерации id девайса и refresh-token используйте *first-login*::
```
    go run first-login/main.go
```
Вводим телефон (в формате 79XXXXXXXXX), код подтверждения из СМС.

# Использование
```
    // Важно указать нужный часовой пояс, чтобы правильно формировалось время в чеке
    client := lknpd.NewClient(
		"Asia/Barnaul",
		os.Getenv("DEVICE_ID"),
		os.Getenv("REFRESH_TOKEN"))

    // Создание чека
	saleId, err := client.CreateSale(lknpd.CreateSaleRequest{
		PaymentType: lknpd.Cash,
		Services: []*lknpd.Service{
			{
				Label:    "Информационная услуга",
				Amount:   decimal.NewFromFloat(50.5),
				Quantity: 1,
			},
		},
		Client: &lknpd.Customer{IncomeType: lknpd.Individual}})
	if err != nil {
		log.Panic(err)
	}

    // Удаление чека
	if err = client.CancelSale(lknpd.CancelSaleRequest{
		CancelType:  lknpd.Cancel,
		ReceiptUUID: saleId,
	}); err != nil {
		log.Panic(err)
	}
```


### Источники и другие реализации
[Автоматизация для самозанятых: как интегрировать налог с IT проектом](https://habr.com/ru/post/436656/)

JS lib [alexstep/moy-nalog](https://github.com/alexstep/moy-nalog)

PHP lib [shoman4eg/moy-nalog](https://github.com/shoman4eg/moy-nalog)

Go lib (с авторизация по логин/пароль) [shoman4eg/go-moy-nalog](https://github.com/shoman4eg/go-moy-nalog)