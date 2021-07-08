**WARNING: This is not production ready!**

# Docs

* .env file in the root folder of the executable is used to load configuration.
Tokens:
  - `MONGO_CONN_URL` = mongodb connection string, ex. format: mongodb+srv://username:password@serverhostname.com
  - `SERVER_IP` = what hostname should the server work on? Default: `localhost:24080` (golang formatting is supported, `:24080` is a valid value, opens the server for global network)


# API

*Available Endpoints:* 
**Get Average BTC price with a five minute interval for a given date.**
* Request type: GET\
`/api/getDay?date=DD.MM.YYYY`

**Response**
- A JSON object that either contains data or a single error. (Even though error object is provided, use http status codes to determine whether your request was fulfilled or not)

Example of data:
```json
{
	"date": "08.07.2021",

	"stamps": [

		{
			"time": "20:25",
			"price": "32900.24125271"
		},

		{
			"time": "20:30",
			"price": "32901.30868121"
		},
		{
			"time": "20:35",
			"price": "32903.48645331"
		}
	]
}
```

Example of error: 

```json
{
  "error": "<ERROR>"
}
```

**Get Current Average BTC price in the last five minutes.**
* Request type: GET\
`/api/getCurrent

**Response**
- A JSON object that either contains data or a single error.

Example of data:
```json
{
	"time": "20:58",
	"price": "32807.96841525"
}
```

Example of error: 

```json
{
  "error": "<ERROR>"
}
```

