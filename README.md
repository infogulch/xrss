# About

This is a web-based rss reader built on top of
[caddy-xtemplate](https://github.com/infogulch/caddy-xtemplate), a php-like
hypertext preprocessor with Go html/template syntax, integrated into Caddy a
world-class web server.

Notable Features

* UI built with tailwindcss.js
* Adaptive UI suitable for small screens
* Add feeds
* Show items from all feeds or just selected feed (sometimes)
* Uses [gofeed](github.com/mmcdole/gofeed) feed reader library

Future

* Maybe eventually update feeds
* Read status, filter by read status
* Sort options
* Scrape item links; display full article content inline
* Multiple users

![screenshot](screenshot.png)

## Developing

```
GOFLAGS='-tags="sqlite_json"' CGO_ENABLED=1 xcaddy build --with github.com/infogulch/caddy-xtemplate --with github.com/greenpau/caddy-security --with github.com/infogulch/xrss

GOFLAGS='-tags="sqlite_json"' CGO_ENABLED=1 XCADDY_DEBUG=1 xcaddy build --with github.com/infogulch/caddy-xtemplate=../xtemplate --with github.com/infogulch/xrss=. --with github.com/greenpau/caddy-security
```

## Resources

https://www.hyperui.dev/

https://hurl.dev/docs/hurl-file.html

https://github.com/go-echarts/go-echarts

https://nene.github.io/prettier-sql-playground/

https://github.com/cornelk/goscrape

https://github.com/robfig/cron

Small view should use: https://play.tailwindcss.com/bQJZaR1c6O + https://tailwindcss.com/docs/hover-focus-and-other-states#focus-within

7 Restful Routes: https://gist.github.com/alexpchin/09939db6f81d654af06b

https://github.com/caddyserver/xcaddy/pull/62
