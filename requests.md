## **Price MM Requests**

#### Template

1. **Method:** ``
    - **Description:**
    - **Params:** `limit` (optional), `offset` (required).
    - **Request Example:** ``
    - **Response:** [
      {
      "id": "123",
      "name": "John Doe",
      "email": "john.doe@example.com",
      "createdAt": "2024-01-01T12:00:00Z"
      }
      ]
    - **Response codes:** 200 (OK), 500 (Internal Server Error).

### Users

1. **Method:** `GET /users`
    - **Description:** Get users list.
    - **Params:** `limit` (optional), `skip` (optional).
    - **Request Example:** `GET /users?limit=10&offset=20`
    - **Response:** [
      {
      "id": "123",
      "name": "John Doe",
      "email": "john.doe@example.com",
      "createdAt": "2024-01-01T12:00:00Z"
      },
      ...
      ]
    - **Response codes:** 200 (OK), 500 (Internal Server Error).
    -
### Vendor

1. **Method:** `GET /api/vendor/country-range-prices/`

    - **Description:**
    - **Params:** `minPrice` (required), `maxPrice` (required), `country` (required)
    - **Request Example:** `GET /api/vendor/country-prices-range?maxPrice=0.9999&minPrice=0.09&country=France`
    - **Response:** [
      {
      "id": "6686864cc871a1c294ef1baa",
      "mccmnc": "255001",
      "effectiveFrom": "2024-07-04T13:32:59.999Z",
      "effectiveTo": "2024-08-28T15:17:00Z",
      "rate": {
      "value": 0.111,
      "currency": "USD"
      },
      "productId": "01J1W68CW0DXK4XF8RVCB2WC86",
      "country": "Ukraine",
      "network": "MTS-UKR"
      }
      ]
    - **Response codes:** 200 (OK), 500 (Internal Server Error).



2. **Method** `POST /api/prices?`

    - **Description** Create/Update Vendor Prices
    - **Params** `type` (required)
    - **Request Example** `POST /api/prices?type=vendors`
      Body {
      productId: product.id,
      product: product.label,
      partnerId: job.partner.id,
      partnerName: job.partner.label,
      organizationId: orgId,
      prices: {
      mccmnc: "255001"
      priceUploadId: "66cf149e520f0fca27e20150"
      "effectiveFrom" : ISODate("2024-08-28T15:17:00.000+0000"),
      "effectiveTo" : ISODate("2099-12-31T00:00:00.000+0000"),
      "rate" : {
      "value" : 0.49345,
      "currency" : "USD",
      "type" : "NEW"
      },
      country: "Ukraine",
      status: "NEW PRICE" | "UPDATED PRICE" | "EFFECTIVE_DATE_UPDATED" | "NEW PRICE" | "CLOSED PRICE",
      productId: "01J1W68CW0DXK4XF8RVCB2WC86"
      }
      }
    - **Response codes:** 200 (OK), 500 (Internal Server Error).


### Client

#### RatePlans

1. **Method:** `GET /api/rate-plans`

    - **Description:** Get Rate Plan List
    - **Params:** `type` `skip` `limit` (required), `sort` `sortBy` `criteria` `searchBy` (optional)
    - **Request Example:** `GET /api/rate-plans/?skip=1&limit=2&criteria=productname&sortBy=name&sort=asc&type=Custom`
    - **Response:** (type = Custom) [
      {
      "id": "01J4GYBPGG9Q6GVRAZ9AA8VC3G",
      "name": "Custom Rate Plan for Inactive retestban / Inactive retestban0508",
      "planType": "Custom",
      "prices" : [
      {
      "price" : {
      "id" : "",
      "mccmnc" : "255007",
      "country" : "Ukraine",
      "productId" : "01J1W68CW0DXK4XF8RVCB2WC86",
      "priceUploadId" : "66cf149e520f0fca27e20150",
      "effectiveFrom" : ISODate("2024-08-28T15:17:00.000+0000"),
      "effectiveTo" : ISODate("2099-12-31T00:00:00.000+0000"),
      "rate" : {
      "value" : 0.49345,
      "currency" : "USD",
      "type" : "NEW"
      },
      "systemCountry" : "Ukraine",
      "systemNet" : "Kyivstar",
      "createdat" : ISODate("2024-08-28T12:19:09.110+0000"),
      "updatedAt" : ISODate("2024-08-28T12:19:09.110+0000")
      },
      "groupingId" : "01J1W68CW0DXK4XF8RVCB2WC86"
      }
      ],
      "updatedAt": "2024-08-06T07:00:12.476Z",
      "partnerId": "01J4GVKTWXH2HDCKW5DP6TDP27",
      "productId": "01J4GVMK6XM9M2ZXYSD3KNMQRV",
      "systemCountriesCount": 1
      },
      ...
      ]

    - **Response:** (type = Global) [
      {
      "id": "01HXECBSCDTCE6P53GCTKXZMBB",
      "name": "17",
      "organizationId": "01HGWKXF4ZCP89Y3F9MYCARQ5T",
      "planType": "Global",
      "prices" : [
      {
      "price" : {
      "id" : "",
      "mccmnc" : "255007",
      "country" : "Ukraine",
      "productId" : "01J1W68CW0DXK4XF8RVCB2WC86",
      "priceUploadId" : "66cf149e520f0fca27e20150",
      "effectiveFrom" : ISODate("2024-08-28T15:17:00.000+0000"),
      "effectiveTo" : ISODate("2099-12-31T00:00:00.000+0000"),
      "rate" : {
      "value" : 0.49345,
      "currency" : "USD",
      "type" : "NEW"
      },
      "systemCountry" : "Ukraine",
      "systemNet" : "Kyivstar",
      "createdat" : ISODate("2024-08-28T12:19:09.110+0000"),
      "updatedAt" : ISODate("2024-08-28T12:19:09.110+0000")
      },
      "groupingId" : "01J1W68CW0DXK4XF8RVCB2WC86"
      }
      ],
      "updatedAt": "2024-05-09T09:53:44.333Z",
      "linkedProductsCount": 0,
      "planBackwardId": "98214851-5537-80e3-920c-bba36c1c16ef"
      },
      ...
      ]
    - **Response codes:** 200 (OK), 500 (Internal Server Error).

2. **Method:** `GET /rate-plans/plan/`

    - **Description:** Get RatePlan by Product id
    - **Params:** Applied to RatePlan `productId`. Applied to RatePlan's prices: `skip` `limit` (required), `sort` `sortBy` `criteria` (optional).
    - **Request Example:** `GET /rate-plans/plan/?productId=01J1W68CW0DXK4XF8RVCB2WC86&criteria=&sortBy=&sort=asc&limit=15&skip=0`
    - **Response:** {
      "id": "01J1YPSKSGFGRTZ37C0CFA7PFB",
      "name": "Custom Rate Plan for PartnerForCountClient / testprdct123",
      "organizationId": "wc34dcti10",
      "planType": "Custom",
      "prices": [
      {
      "id": "6686864cc871a1c294ef1baa",
      "mccmnc": "255001",
      "effectiveFrom": "2024-07-04T13:32:59.999Z",
      "effectiveTo": "2024-08-28T15:17:00Z",
      "rate": {
      "value": 0.111,
      "currency": "USD",
      "type": "INCREASE"
      },
      "productId": "01J1W68CW0DXK4XF8RVCB2WC86",
      "systemCountry": "Ukraine",
      "systemNet": "MTS-UKR"
      }
      ],
      "updatedAt": "2024-08-28T12:19:09.11Z",
      "partnerId": "01HKA9SHGG2FSBFZM5556Q8JD9",
      "productId": "01J1W68CW0DXK4XF8RVCB2WC86",
      "systemCountriesCount": 1
      }
    - **Response codes:** 200 (OK), 500 (Internal Server Error).

3. **Method:** `GET /rate-plans/plan/`

    - **Description:** Get RatePlan by id
    - **Params:** Applied to RatePlan `id`. Applied to RatePlan's prices: `skip` `limit` (required), `sort` `sortBy` `criteria` (optional).
    - **Request Example:** `GET /rate-plans/plan/?id=01J1YPSKSGFGRTZ37C0CFA7PFB&criteria=&sortBy=&sort=asc&limit=15&skip=0`
    - **Response:** {
      "id": "01J1YPSKSGFGRTZ37C0CFA7PFB",
      "name": "Custom Rate Plan for PartnerForCountClient / testprdct123",
      "organizationId": "wc34dcti10",
      "planType": "Custom",
      "prices": [
      {
      "id": "6686864cc871a1c294ef1baa",
      "mccmnc": "255001",
      "effectiveFrom": "2024-07-04T13:32:59.999Z",
      "effectiveTo": "2024-08-28T15:17:00Z",
      "rate": {
      "value": 0.111,
      "currency": "USD",
      "type": "INCREASE"
      },
      "productId": "01J1W68CW0DXK4XF8RVCB2WC86",
      "systemCountry": "Ukraine",
      "systemNet": "MTS-UKR"
      }
      ],
      "updatedAt": "2024-08-28T12:19:09.11Z",
      "partnerId": "01HKA9SHGG2FSBFZM5556Q8JD9",
      "productId": "01J1W68CW0DXK4XF8RVCB2WC86",
      "systemCountriesCount": 1
      }

    - **Response codes:** 200 (OK), 500 (Internal Server Error).

4. **Method:** `POST /rate-plans/link`

    - **Description:** Link/Unlink RatePlan to/from Product
    - **Params:** `planId` (required), `ProductId` (required), `unlink`(optional)
    - **Request Example:** `POST rate-plans/link`
      Body {
      "planId": "01J1YPSKSGFGRTZ37C0CFA7PFB",
      "productId": "01J1W68CW0DXK4XF8RVCB2WC86",
      "unlink": false | true,
      }
    - **Response:**

    - **Response codes:** 204 (No content), 500 (Internal Server Error).

5. **Method** `POST /api/prices?`

    - **Description** Create/Update Custom RatePlan
    - **Params** `type` (required)
    - **Request Example** `POST /api/prices?type=clients`
      Body {
      name: `Custom Rate Plan for <Client name>/ <Product name>`,
      backwardRef: job.groupingRef,
      planType: 'Custom',
      organizationId: orgId,
      prices: {
      mccmnc: "255001"
      priceUploadId: "66cf149e520f0fca27e20150"
      "effectiveFrom" : ISODate("2024-08-28T15:17:00.000+0000"),
      "effectiveTo" : ISODate("2099-12-31T00:00:00.000+0000"),
      "rate" : {
      "value" : 0.49345,
      "currency" : "USD",
      "type" : "NEW"
      },
      country: "Ukraine",
      status: "NEW PRICE" | "UPDATED PRICE" | "EFFECTIVE_DATE_UPDATED" | "NEW PRICE" | "CLOSED PRICE",
      productId: "01J1W68CW0DXK4XF8RVCB2WC86"
      }
      }
    - **Response codes:** 200 (OK), 500 (Internal Server Error).

6. **Method** `POST /api/prices?`

    - **Description** Create/Update Global RatePlan
    - **Params** `type` (required)
    - **Request Example** `POST /api/prices?type=clients`
      Body {
      name: "Gloabal Rate plan name",
      backwardRef: job.groupingRef,
      planType: 'Global',
      organizationId: orgId,
      prices: {
      mccmnc: "255001"
      priceUploadId: "66cf149e520f0fca27e20150"
      "effectiveFrom" : ISODate("2024-08-28T15:17:00.000+0000"),
      "effectiveTo" : ISODate("2099-12-31T00:00:00.000+0000"),
      "rate" : {
      "value" : 0.49345,
      "currency" : "USD",
      "type" : "NEW"
      },
      country: "Ukraine",
      status: "NEW PRICE" | "UPDATED PRICE" | "EFFECTIVE_DATE_UPDATED" | "NEW PRICE" | "CLOSED PRICE",
      productId: "01J1W68CW0DXK4XF8RVCB2WC86"
      }
      }
    - **Response codes:** 200 (OK), 500 (Internal Server Error).

#### PriceLists

1. **Method:** `GET /price-lists/:productId`
    - **Description:** Get PriceList List
    - **Params:** `skip` `limit` (required), `sort` `sortBy` `criteria` (optional)
    - **Request Example:** `GET /price-lists/:productId/?criteria=255007&sortBy=rate&sort=desc&skip=0&limit=15`
    - **Response:** {
      "id": "01J1YPSKT07MAEEZSHWAMXDGGJ",
      "organizationId": "wc34dcti10",
      "customRatePlanId": "01J1YPSKSGFGRTZ37C0CFA7PFB",
      "globalRatePlanId": null,
      "partnerId": "01HKA9SHGG2FSBFZM5556Q8JD9",
      "productId": "01J1W68CW0DXK4XF8RVCB2WC86",
      "prices": [
      {
      "id": "6686864cc871a1c294ef1baa",
      "mccmnc": "255001",
      "effectiveFrom": "2024-07-04T13:32:59.999Z",
      "effectiveTo": "2024-08-28T15:17:00Z",
      "rate": {
      "value": 0.111,
      "currency": "USD",
      "type": "INCREASE"
      },
      "productId": "01J1W68CW0DXK4XF8RVCB2WC86",
      "systemCountry": "Ukraine",
      "systemNet": "MTS-UKR",
      "ratePlanType": "Custom"
      },
      ...
      ]
      "updatedAt": "2024-08-06T07:00:12.476Z",
      }
    - **Response codes:** 200 (OK), 500 (Internal Server Error).

#### Prices (both vendor prices and client )

1. **Method:** `GET /api/prices/timeline`

    - **Description:** Get 2 past prices, one current and one future.
    - **Params:** `productId` (required), `type` (required), `mccmnc` (required), `forward` (optional, default=1), `backward` (optional, default=2)

   ##### Client

    - **Request Example:** `GET api/prices/timeline?productId=01J5ZTH7VFY2SDVR1RN9BPYNK9&type=clients&mccmnc=450`
    - **Response:** [
      {
      "effectiveFrom": "2024-08-28T10:27:37Z",
      "effectiveTo": "2024-08-28T11:30:00Z",
      "rate": "\ufffd\ufffd0.04 EUR",
      "current": false,
      "productId": "01J5ZTH7VFY2SDVR1RN9BPYNK9"
      },
      {
      "effectiveFrom": "2024-08-28T11:30:00Z",
      "effectiveTo": "2099-12-31T00:00:00Z",
      "rate": "\ufffd\ufffd0.22 EUR",
      "current": true,
      "productId": "01J5ZTH7VFY2SDVR1RN9BPYNK9"
      }
      ]

   ##### Vendor

    - **Request Example:** `GET /api/prices/timeline?productId=01J5ZV5PA6FXBHTCY9ZJN74SZN&type=vendors&mccmnc=450005`
    - **Response:** [
      {
      "effectiveFrom": "2024-08-28T11:34:37Z",
      "effectiveTo": "2099-12-31T00:00:00Z",
      "rate": "\ufffd\ufffd0.99 EUR",
      "current": true,
      "productId": "01J5ZV5PA6FXBHTCY9ZJN74SZN"
      }
      ]

2. **Method:** `GET /api/prices/current`

    - **Description:** ----
    - **Params:** `productId` (required), `type` (required), `country` (required)

   ##### Client

    - **Request Example:** `GET /api/prices/current?productId=01J5ZTH7VFY2SDVR1RN9BPYNK9&country=South%20Korea&type=clients&skip=0&limit=20`
    - **Response:** [
      {
      "id": "66cf064fe000fccd9292f994",
      "mccmnc": "450",
      "effectiveFrom": "2024-08-28T11:30:00Z",
      "effectiveTo": "2099-12-31T00:00:00Z",
      "rate": {
      "value": 0.22,
      "currency": "EUR",
      "type": "INCREASE"
      },
      "productId": "01J5ZTH7VFY2SDVR1RN9BPYNK9",
      "systemCountry": "South Korea",
      "systemNet": "Other-KOR"
      }
      ]

   ##### Vendor

    - **Request Example:** `GET /api/prices/current?productId=01J5ZV5PA6FXBHTCY9ZJN74SZN&country=South%20Korea&type=vendors&skip=0&limit=20`
    - **Response:** [
      {
      "id": "66cf0a9e175dbaa88af813a3",
      "mccmnc": "450005",
      "effectiveFrom": "2024-08-28T11:34:37Z",
      "effectiveTo": "2099-12-31T00:00:00Z",
      "rate": {
      "value": 0.99,
      "currency": "EUR",
      "type": "NEW"
      },
      "productId": "01J5ZV5PA6FXBHTCY9ZJN74SZN",
      "systemCountry": "South Korea",
      "systemNet": "SK Telecom-KOR"
      },
      {
      "id": "66cf0a9e175dbaa88af813a4",
      "mccmnc": "450006",
      "effectiveFrom": "2024-08-28T11:34:37Z",
      "effectiveTo": "2099-12-31T00:00:00Z",
      "rate": {
      "value": 0.01,
      "currency": "EUR",
      "type": "NEW"
      },
      "productId": "01J5ZV5PA6FXBHTCY9ZJN74SZN",
      "systemCountry": "South Korea",
      "systemNet": "LG Telecom-KOR"
      }
      ]

3. **Method:** `GET /api/prices/last`

    - **Description:** ---
    - **Params:** `productId` (required), `type` (required), `country` (required)

   ##### Client

    - **Request Example:** `GET /api/prices/last?productId=01J5ZTH7VFY2SDVR1RN9BPYNK9&type=clients&skip=0&limit=20`
    - **Response:** [
      {
      "id": "66cf064fe000fccd9292f994",
      "mccmnc": "450",
      "effectiveFrom": "2024-08-28T11:30:00Z",
      "effectiveTo": "2099-12-31T00:00:00Z",
      "rate": {
      "value": 0.22,
      "currency": "EUR",
      "type": "INCREASE"
      },
      "productId": "01J5ZTH7VFY2SDVR1RN9BPYNK9",
      "systemCountry": "South Korea",
      "systemNet": "Other-KOR",
      "oldRate": {
      "value": 0.04,
      "currency": "EUR",
      "type": "NEW"
      }
      },
      {
      "id": "66cefb1ee000fccd9292f7f8",
      "mccmnc": "450006",
      "effectiveFrom": "2024-08-28T10:27:37Z",
      "effectiveTo": "2024-08-28T11:30:00Z",
      "rate": {
      "value": 0.03,
      "currency": "EUR",
      "type": "NEW"
      },
      "productId": "01J5ZTH7VFY2SDVR1RN9BPYNK9",
      "systemCountry": "South Korea",
      "systemNet": "LG Telecom-KOR",
      "oldRate": null
      },
      {
      "id": "66cefb1ee000fccd9292f7f9",
      "mccmnc": "450005",
      "effectiveFrom": "2024-08-28T10:27:37Z",
      "effectiveTo": "2024-08-28T11:30:00Z",
      "rate": {
      "value": 0.02,
      "currency": "EUR",
      "type": "NEW"
      },
      "productId": "01J5ZTH7VFY2SDVR1RN9BPYNK9",
      "systemCountry": "South Korea",
      "systemNet": "SK Telecom-KOR",
      "oldRate": null
      }
      ]

   ##### Vendor

    - **Request Example:** `GET /api/prices/last?productId=01J5ZV5PA6FXBHTCY9ZJN74SZN&type=vendors&skip=0&limit=20`
    - **Response:** [
      {
      "id": "66cf0a9e175dbaa88af813a3",
      "mccmnc": "450005",
      "effectiveFrom": "2024-08-28T11:34:37Z",
      "effectiveTo": "2099-12-31T00:00:00Z",
      "rate": {
      "value": 0.99,
      "currency": "EUR",
      "type": "NEW"
      },
      "productId": "01J5ZV5PA6FXBHTCY9ZJN74SZN",
      "systemCountry": "South Korea",
      "systemNet": "SK Telecom-KOR",
      "oldRate": null
      },
      {
      "id": "66cf0a9e175dbaa88af813a4",
      "mccmnc": "450006",
      "effectiveFrom": "2024-08-28T11:34:37Z",
      "effectiveTo": "2099-12-31T00:00:00Z",
      "rate": {
      "value": 0.01,
      "currency": "EUR",
      "type": "NEW"
      },
      "productId": "01J5ZV5PA6FXBHTCY9ZJN74SZN",
      "systemCountry": "South Korea",
      "systemNet": "LG Telecom-KOR",
      "oldRate": null
      }
      ]
