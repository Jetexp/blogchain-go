<div align="center">
  <img width="150" height="150" src="https://github.com/zikwall/blogchain/blob/master/screenshots/bc_go_300.png">
  <h1>Blog Chain</h1>
  <h5>App server written in Golang</h5>
</div>

### Deploy

- `make`

### Docker

1. `docker build -t blogchain-go-img .`
2. `docker run -d -p 3001:3001 --name blogchain-go blogchain-go-img`

### Tests

1. `go tests`

#### For teamcity

1. `go tests -json`