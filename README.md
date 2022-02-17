##Spot Hero Backend Challenge

- As per the requirements mentioned in the challenge doc.
- Application will run on port 5000 (if running on Macos Monterey, stop sharing Airplay, using 5000 port)
- Challenge implementation is in the GO language
  * Using the ``gorm`` library for object relationship mapping
  * Using the ``gorm sqlite`` library for database.
  * using the ``sqlmock`` library for test mocking.
  * using the ``mux`` library for rest endpoints routing.
  * Using the provided seeded data for rates.
- Start the Application on port 5000 through [main.go](main.go) file.
- Data is loaded into the [rates.db](rates.db), if needed to delete the file and application startup will load the data.
- Test cases are present for price, rate endpoints and model.
- Following are the sample endpoints results
  ```bash
    curl http://localhost:5000/rates
    [{"days":"mon,tues,thurs","times":"0900-2100","tz":"America/Chicago","price":1500},{"days":"fri,sat,sun","times":"0900-2100","tz":"America/Chicago","price":2000},{"days":"wed","times":"0600-1800","tz":"America/Chicago","price":1750},{"days":"mon,wed,sat","times":"0100-0500","tz":"America/Chicago","price":1000},{"days":"sun,tues","times":"0100-0700","tz":"America/Chicago","price":925}]
  
    curl http://localhost:5000/price\?start\=2015-07-01T07:00:00-05:00\&end\=2015-07-01T12:00:00-05:00
    {"price":1750}
  
    curl http://localhost:5000/price\?start\=2015-07-04T15:00:00%2B00:00\&end\=2015-07-04T20:00:00%2B00:00
    {"price":2000}
  
    curl http://localhost:5000/price\?start\=2015-07-04T07:00:00%2B05:00\&end\=2015-07-04T20:00:00%2B05:00
    "unavailable"
  ```
  

Notes
- ISO-8601 date pattern only has zone offset, hard to determine the timezon based on offset.
- Price: 2000 test case is green because not considering different time zone offset, still considering ``America/Chicago``