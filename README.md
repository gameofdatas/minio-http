# minio-rest
A http wrapper over Minio

## introduction
A http wrapper over Minio which will provide end users to upload and download files.

## api-endpoints

### /api/upload

The endpoint allows you upload files in the form of base64 encoded string to specific bucket, folder with retention period.

#### request

`POST /api/upload/`
#### request body
```
type reqBody struct {
bucket string         // name of the bukcet where file should be uploaded
file_path string      // path where file should be upload
file_name string      // name of the file to be uploaded
file_content string   // base64 encoded string of file
retention_ms int      // milliseconds after the object should be deleted from MinIO
}
```
#### sample request:
```
POST /api/upload HTTP/1.1
Host: localhost
Content-Type: application/json
Content-Length: 133

{"bucket":"soruce-bucket","file_path":"folder-name","file_name":"metadata.json","file_content":"c2FtcGxl","retention_ms":900000}
```
#### response
HTTP/1.1 201 Created
Date: Thu, 24 Feb 2011 12:36:30 GMT
Status: 201 Created
Connection: close
Location: /api/upload/

### /api/download
The endpoint allows you down files from specific bucket, folder in base64 encoded string upon authorization.
#### request body
```
type reqBody struct {
bucket string         // name of the bukcet where file should be downloaded from
file_path string      // path where file should be downloaded from
file_name string      // name of the file to be downloaded
}
```
#### sample request:
```
POST /api/download HTTP/1.1
Host: localhost
Content-Type: application/json
Content-Length: 133

{"bucket":"soruce-bucket","file_path":"folder-name","file_name":"metadata.json"}
```

### Response
```text
HTTP/1.1 201 Created
Date: Thu, 24 Feb 2011 12:36:30 GMT
Status: 201 Created
Connection: close
```
### response body
```
c2FtcGxl  # base64 encoded file content
```

### /api/health

Return the status of the application availability. It returns a json response

```json
{"alive":true,"since":"7h22m29.850902461s","version":"0.1.0","build_date":"","go_version":"go1.18.3","commit":""}
```

## development

env variable needed by minio-rest
```text
API_PORT                 # port number to expose to (string)  eg: ':8080'
MINIO_ENDPOINT           # minio server endpoint
MINIO_ACCESS_KEY         # minio server username to connect
MINIO_SECRET_KEY         # minio secret to connect
```

To run or start the minio-rest

`go run main.go`

To run unit tests 

`go test ./...`

## logging

It logs every request comes across minio-rest with the request body and respective response.

The logs fields looks like:

```text
	srvName     string  # service name
	ID          string  # adaptation ID
	RemoteIP    string  # remote ip 
	Host        string  # host machine
	Method      string  # request method
	Path        string  # api path
	Protocol    string  # http protocol
	StatusCode  int     # response status code
	Latency     float64 # total time took to address the request and return response
	Error       error   # error message if any error
	Stack       []byte  # debug stack trace for error
	RequestBody []byte  # request body
```

```text
{"level":"info","srv":"minio-rest","remote_ip":"127.0.0.1","host":"localhost:8081","method":"POST","path":"/api/download","protocol":"http","status_code":200,"latency":0.018843997,"tag":"request","request_body":"{\"bucket\":\"soruce-bucket3\",\"file_path\":\"x-adapatation-id\",\"file_name\":\"metadata.json\"}","time":1657521249893,"message":"success"}
```
