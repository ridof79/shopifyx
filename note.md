{
    "id" : "uuid", // generate uuid from db(postgres)
	"name": "", // not null, minLength 5, maxLength 60
	"price": 10000, // not null, min 0
	"imageUrl" : "", // not null, url=true
	"stock" : 10, // not null, min 0
	"condition": "new | second", // not null, must only accept enum
	"tags": [""], // string not null, minItems 0
	"isPurchaseable": true, // not null
    "userId" : "user_id"
}


{
	"name": "", 
	"price": "",
	"imageUrl" : "", 
	"stock" : "", 
	"condition": "",
	"tags": [""], 
	"isPurchaseable":  ""
}


{
	"id" : "uuid", // generate uuid from db(postgres)
	"bankName":"name", // not null, minLength 5, maxLength 15
	"bankAccountName":"accName", // not null, minLength 5, maxLength 15
	"bankAccountNumber": "0981", // not null, minLength 5, maxLength 15
	"userId" : "user_id"
}

{
	"bankName":"name", 
	"bankAccountName":"accName",
	"bankAccountNumber": "0981", 
}

{
	"id":"uuid"
	"productId":"", // not null, must be a correct product id
	"bankAccountId":"", // not null, must be a correct bank account id
	"paymentProofImageUrl":"", // not null, must be a correct url
	"quantity":10 // not null, min 1
	"userId": "user_id" 
}

migrate -database "postgres://postgres:admin@localhost:5433/shopifyx_data?sslmode=disable" -path db/migrations up