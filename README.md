# NASA crawler api

## ENV variables description
- **PORT** default **8080**
- **API_KEY** default **DEMO_KEY**
- **API_URL** default **https://api.nasa.gov/planetary/apod?api_key=%s&date=%s**
- **CONCURRENT_REQUESTS** default **5**

## Running docker 
- docker build --no-cache -t url-collector .
- docker run --env-file ./env.list -p 8080:8081/tcp -it url-collector

## Example
http://localhost:8080/pictures?start_date=2021-06-16&end_date=2021-06-18
