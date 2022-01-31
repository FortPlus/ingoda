Service to keep lists, ban list, white list etc.
Current API version is - "v1"

### Methods

#### Get all items: GET: [/api/{version}/{LIST_NAME}]()
Response codes:
- 200 - Ok, banned list in response body
- 204 - item doesn't exists in banned list
- 500 - internal error

#### Check if item exists in banned list: GET: [/api/{version}/{LIST_NAME}/exist]()
Response codes:
- 200 - item is exists in list
- 204 - item doesn't exists in list
- 500 - internal error

#### Add record to banned list: POST: [/api/{version}/{LIST_NAME}]()
Response codes:
- 201 - created
- 400 - record data has invalid format
- 500 - internal error

#### Delete record from list: DELETE: [/api/{version}/{LIST_NAME}/{id}]()
Response codes:
- 200 - record deleted
- 401 - record not found
- 500 - internal error
