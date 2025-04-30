# Online Library of Congress

## Endpoints

### Get Books (book search)
GET /books

Method: GET
Query Parameters: 
* title (wildcard case insensitive search)
* genre (wildcard case insensitive search)
* publisher (wildcard case insensitive search)
* isbn (exact match)
* author (wildcard case insensitive search)

##### Sample Query:

```
curl --location 'http://localhost:8080/book?title=The&isbn=978-0-439-02348-1&publisher=har&genre=DRA&author=gri'
```

### Get Book (book retrieval)
GET /book/{id}

Method: GET
Query Parameters: None
URL Parameters: 
* bookID (required)

##### Sample Query:

```
curl --location 'http://localhost:8080/book/1'
```


### Checkout a Book 
POST /book/checkout

Method: POST
Query Parameters: None
URL Parameters: None
Body:
    {
        bookID string
        userID string
    }

##### Sample Query:

```
curl --location 'http://localhost:8080/book/checkout' \
--header 'Content-Type: application/json' \
--data '{
    "book_id": 5,
    "user_id": 2,
    "created_by": "david"
}'
```

