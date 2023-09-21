# webdav

simple webdav server with auth

## Usage

config.yml

```yml
port: 80
accounts:
  - user: 'demo'
    password: '123456789'
    dir: ./data
```

run with binary

```
webdav -c config.yml
```

run with docker

```
docker run tmaize/webdav -p 80:80 -v config.yml:/app/config.yml
```
