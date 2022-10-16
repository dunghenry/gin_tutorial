# Gin + MongoDB + Docker + Authentication

### Run auto hot reload:

```js
air run main.go
```

### Run server

```js
go run main.go
```

### Define template

```js
{{ template "header.tmpl" . }}
```

### Using docker file

```js
docker build -t username/image_name:version .
```

```js
docker run -it --name container_name --rm -p 3000:3000 username/image_name:version
```

### Using docker-compose

```js
docker-compose up
```

### Run and build using docker-compose

```js
docker-compose up --build
```
