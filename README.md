# Video_Analytics

## Tech stack

- **Language** : Go 1.20+
- **Database** : MariaDB 11
- **Cache/Queue** : Redis Streams
- **Mailer** : SMTP

## External Package

- [Gin](https://github.com/gin-gonic/gin) – HTTP web framework
- [GORM](https://gorm.io) – ORM
- [Redis](github.com/redis/go-redis/v9) – Redis client
- [uuid](github.com/satori/go.uuid) – Unique IDs
- [Prometheus](github.com/prometheus/client_golang/prometheus) – Metrics integration

### clone URL
```sh
https://github.com/SagarVijaya/Video_Analytics.git
```
### **1 Install Dependencies**

```sh
go mod tidy
```

### **2 Run the app**
```sh
go run main.go

```

### **3 Run with Docker**
```sh
docker build -t video-ads-backend .
docker run -p 3011:3011 -p 9090:9090 video-ads-backend
```

### **4 Run with Docker Compose**

```sh
docker compose up -d
```

## sample Request

```sh
curl hcurl --location 'http://localhost:3011/ads/analytics?from=2025-06-18T00%3A00%3A00Z&to=2025-06-19T00%3A00%3A00Z'

curl curl --location --request GET 'http://localhost:3011/ads/analytics?minutes=30' \'

curl curl --location 'http://localhost:3011/ads'

curl curl --location 'http://localhost:3011/ads/click' \
--header 'Content-Type: application/json' \
--data '{
  "ad_id": 1,
  "clicked_at": "2025-06-22T12:45:00Z",
  "ip": "192.168.1.1110",
  "playback_second": 15
}'

```

## sample Response 1
- Endpoint :  /ads

```json
{
  "status": "S",
  "errMsg": "",
  "adsList": [
    {
      "ID": 1,
      "image_url": "https://cdn.example.com/ad1.jpg",
      "target_url": "https://example.com/product/123",
      "created_at": "2025-06-21T10:00:00Z"
    }
  ]
}

```

- Endpoint :  /ads/click
### sample Request 2
```json
{
  "ad_id": 1,
  "ip": "192.168.1.100",
  "playback_second": 5
}
```

### sample Response 2
```json
{"data":{"end":"2025-05-06","overall_details":[{"quantity":3,"Category":"Clothing"},{"quantity":3,"Category":"Shoes"},{"quantity":4,"Category":"Electronics"}],"start":"2016-01-01"},"message":"Top Product For Category","status":"S"}

```

### sample Response 3
- Endpoint :  /BackUpData

```json
{
  "status": "S",
  "Msg": "Data Inserted Successfully"
}


```

### sample Response 4
- Endpoint :  /ads/analytics?minutes=30

```json
{
    "adsAnalytics": [
        {
            "ad_id": 0,
            "total_clicks": 1,
            "recent_clicks": 0
        }
    ],
    "errMsg": "",
    "status": "S"
}
```


### sample Response 5
- Endpoint :  /ads/analytics?from=2025-06-18T00:00:00Z&to=2025-06-19T00:00:00Z

```json
{
    "adsAnalytics": [
        {
            "ad_id": 1,
            "clicks": 0,
            "impressions": 0,
            "ctr": 0
        },
        {
            "ad_id": 2,
            "clicks": 0,
            "impressions": 0,
            "ctr": 0
        }
    ],
    "errMsg": "",
    "status": "S"
}
```




