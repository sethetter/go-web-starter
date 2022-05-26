# go-web-starter

An opinionated web app starting point.

## Tools

- Go server-rendered pages, backed by Postgres
- Web framework: [Gin](https://gin-gonic.com)
- Validations: [govalidator](https://github.com/asaskevich/govalidator)
- Query building: [sqlx](http://jmoiron.github.io/sqlx/) and [squirrel](https://github.com/Masterminds/squirrel)
- Minimal CSS, no frameworks
- Frontend interactivity with [htmx](http://htmx.org)

## Security

- [Cross-site request forgery (CSRF)](https://developer.mozilla.org/en-US/docs/Glossary/CSRF)
- [Cross-origin resource sharing (CORS)](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- A bunch of other security headers
  - [HTTP strict transport security (HSTS)](https://developer.mozilla.org/en-US/docs/Glossary/HSTS)
  - [Content security policy (CSP)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy)
  - [Mime type sniffing](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Content-Type-Options)
  - [Frame blocking](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options)
  - [Cross origin resource policy](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cross-Origin-Resource-Policy)
  - [Cross origin opener policy](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cross-Origin-Opener-Policy)

## Features

- Session handling
- Basic email service
- Passwordless login
