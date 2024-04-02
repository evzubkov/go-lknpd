# go-lknpd
Unofficial client of the HTTP-API site [Мой налог для самозанятых](https://lknpd.nalog.ru/).
The login method is used via a mobile phone. To generate device id and refresh-token use *first-login*:
```
    go run first-login/main.go
```
Enter the phone number (in the format 79XXXXXXXXXX), the confirmation code from SMS.

# Usage
```
    // It is important to indicate the required time zone so that the time in the receipt is formed correctly
    client := lknpd.NewClient(
		"Asia/Barnaul",
		os.Getenv("DEVICE_ID"),
		os.Getenv("REFRESH_TOKEN"))

    // Creating a check
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

    // Deleting a check
	if err = client.CancelSale(lknpd.CancelSaleRequest{
		CancelType:  lknpd.Cancel,
		ReceiptUUID: saleId,
	}); err != nil {
		log.Panic(err)
	}
```


### Sources and other implementations
[Автоматизация для самозанятых: как интегрировать налог с IT проектом](https://habr.com/ru/post/436656/)

JS lib [alexstep/moy-nalog](https://github.com/alexstep/moy-nalog)

PHP lib [shoman4eg/moy-nalog](https://github.com/shoman4eg/moy-nalog)

Go lib (with authorization by login/password) [shoman4eg/go-moy-nalog](https://github.com/shoman4eg/go-moy-nalog)