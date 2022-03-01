# product-review service overview
Product review engine

A product has the following attributes: 

● ID (required) 

● Name (required) 

● Description (required) 

● Thumbnail Image URL (not required) 

● List of Reviews. Each review has the following attributes 
    
    ○ Reviewer Name (required)
    ○ Written Review (not required) 
    ○ Rating (required. integer 0 to 5) 


### **This service provides following API Endpoints**:
### *Get Products Endpoint*
```
This endpoint returns the list of products. 
Each product returned will have the ID, Name, Description, Thumbnail Image URL, and Overall Rating (the average of all of the product’s reviews). The list of each product’s reviews will not be returned by this endpoint.
``` 
### *Add Review Endpoint*
```
This endpoint provides functionality to add a review to the product. This endpoint will receive the Product ID, Reviewer Name, Written Review and Rating. 
``` 



### **API Spec**:
**Endpoint**: GET /api/products

**Supported query params**:
- limit [integer]: to limit number of resources in a response
- skip [integer]: to skip first n number of resources
- name [string]: name pattern for filtering resources based on name, supports regex

**Response sample**:
```
Status code: 200
Payload:
[
    {
        "id": 3350861610,
        "name": "AtWYa",
        "description": "4t47iO0GIemHfRs1rgV8",
        "thumbnail_img_url": "",
        "overall_rating": 4
    },
    {
        "id": 1597659264,
        "name": "VomjZ",
        "description": "Vuj1cJiWp5swZzqLT7Su",
        "thumbnail_img_url": "",
        "overall_rating": 0
    }
]
```

**Endpoint**: POST /api/products/:id/reviews

**Request payload/body**:
```
{
    "reviewer_name": "Suraj Bhosale",
    "rating": 4,
    "written_review": "This is a good product!"
}
```

**Response sample**:
```
Status code: 201
Payload:
{
    "reviewer_name": "Suraj Bhosale",
    "rating": 4,
    "written_review": "This is a good product!"
}
```

**All the API endpoints are authenticated, and are currently using basic auth (username, passowrd) for authenticating the requests.**

ENV variables 
`PR_USERNAME` and `PR_PASSWORD` needs to be exported with custom values, currently server identifies these creds only, and later client http requests need to pass this values as basic auth header for authenticating the requests.



# Build, Test, Deploy service
Project has a Makefile which has recipes to build, test, deploy the service locally. There are a few other useful recipes for development too.

There are a few configuration parameters which are read from environment variables. 

```
PR_USERNAME=test # username for basic auth
PR_PASSWORD=test # password for basic auth
MONGODB_URL=mongodb://localhost:27017
MONGODB_NAME=product-review
```
Project has a `.env` file containing sample/default config values which can be used for testing locally.

### *Building the service*:

`make build` can be used to compile the binary for this service


### *Executing the unit tests*:
`make test` has recipe for executed unit tests


### *Deploy service locally and perform e2e usecase testing*:
`make up` will bring up the docker-compose stack for other dependent services, currently product-review depends on mongodb service for database

`make seed-db` can be used to seed the database with some initial product data for test purpose.

`make run` will run the dependent services in docker-compose stack and will start the `product-review` service locally.


To perform e2e API tests, sample postman client requests can be found under `tools/product-review.postman_collection.json`



## Code quality  pre-commit/preps:

Makefile also contains a few commands for 
`lint, vet, fmt` code

