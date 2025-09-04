# Online-subscriptions-aggregator
An application to aggregating data of user's online-subscriptions

To start: run Docker on your computer, then open cmd, choose directory with this app (set cd + yourpath/to/appfolder) and set command:

docker-compose up --build

Await for starting docker-compose. Soon you will get ready application!

Open browser and follow next adress:

localhost:3000

Now you can use this API for create, read, delete and update subscriptions. 

GET: Get list of subscription:   localhost:3000/subscriptions
GET: Get subscription by id:   localhost:3000/subscribe/id
POST: Add subscription:   localhost:3000/subscribe
PUT: Update subscription:   localhost:3000/subscribe/id
DELETE: Delete subscrition:   localhost:3000/subscribe/id

You also can calculate a subscriptions cost with requests:

http://localhost:3000/total-cost?start_date=2025-01&end_date=2025-02 - Find and calculate cost of all subscriptions per given period (example start date, end date)

http://localhost:3000/total-cost?start_date=2024-01&end_date=2025-11&user_id=d27618cc-cdc1-466c-aebf-c9235f5ccb9c - Calculate cost of subscriptions filtered by user_id (example start date, end date, user id)

http://localhost:3000/total-cost?start_date=2025-01&end_date=2025-11&service_name=SomeService - Calculate cost of subscriptions filtered by service_name (example start date, end date, service name)

http://localhost:3000/total-cost?start_date=2024-01&end_date=2025-12&user_id=d27618cc-cdc1-466c-aebf-c9235f5ccb9c&service_name=SomeService - Calculate cost of subscriptions filtered by service_name and user id(example start date, end date, user id, service name)

Thanks for reading!