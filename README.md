# la-warehouse

## 1. Инструкцию по запуску сервиса

    make postgres
    make migrate
    make run


## 2. Инструкцию по запуску тестов

    make test

## 3. Описание API методов с запросом и ответом

### Резервирование товара на складе для доставки:

#### Request:

    curl --location 'localhost:8080/rpc' \
    --header 'Content-Type: application/json' \
    --data '{
    "method": "Warehouse.Reservation",
    "params": [{
        "warehouseId": 2,
        "products":
        [
           {
                "id": 1,
                "amount": 102
           },
           {
               "id": 2,
                "amount": 214
           },
           {
               "id": 2,
                "amount": 4123
           },
           {
               "id": 2,
                "amount": 101
           },
           {
               "id": 3,
                "amount": 632
           },
           {
               "id": 4,
                "amount": 748
           },
           {
               "id": 7,
                "amount": 4211
           },
           {
               "id": 11,
                "amount": 1211
           },
           {
               "id": 1,
                "amount": -123
           }
       ]
    }], "id": 1
    }'

#### Response:

    {
    "result": {
        "products": [
            {
                "id": 1,
                "status": "ok"
            },
            {
                "id": 2,
                "status": "ok"
            },
            {
                "id": 2,
                "status": "ok"
            },
            {
                "id": 2,
                "status": "ok"
            },
            {
                "id": 3,
                "status": "ok"
            },
            {
                "id": 4,
                "status": "ok"
            },
            {
                "id": 7,
                "status": "ok"
            },
            {
                "id": 11,
                "status": "products not exist"
            },
            {
                "id": 1,
                "status": "amount must be greater 0"
            }
        ]
    },
    "error": null,
    "id": 1
    }

### Освобождение резерва товаров на складе:

#### Request:

    curl --location 'localhost:8080/rpc' \
    --header 'Content-Type: application/json' \
    --data '{
    "method": "Warehouse.ReservationFree",
    "params": 
    [
        {
            "warehouseId": 2,
            "products": [1, 4, 6, 7, 11]
        }
    ], "id": 1
    }'

#### Response:

    {
    "result": {
        "products": [
            {
                "id": 1,
                "status": "ok"
            },
            {
                "id": 4,
                "status": "ok"
            },
            {
                "id": 6,
                "status": "warehouse_product not exist"
            },
            {
                "id": 7,
                "status": "ok"
            },
            {
                "id": 11,
                "status": "products not exist"
            }
        ]
    },
    "error": null,
    "id": 1
    }

### Получение кол-ва оставшихся товаров на складе:

#### Request 1:

    curl --location 'localhost:8080/rpc' \
    --header 'Content-Type: application/json' \
    --data '{
    "method": "Warehouse.GetAmountProducts",
    "params": 
    [
        {
            "warehouseId": 2
        }
    ], "id": 1
    }'

#### Response 1:

    {
    "result": {
        "amount": 5070
    },
    "error": null,
    "id": 1
    }

----
#### Request 2:

    curl --location 'localhost:8080/rpc' \
    --header 'Content-Type: application/json' \
    --data '{
    "method": "Warehouse.GetAmountProducts",
    "params": 
    [
        {
            "warehouseId": 7
        }
    ], "id": 1
    }'

#### Response 2:

    {
    "result": {
        "amount": 0
    },
    "error": null,
    "id": 1
    }