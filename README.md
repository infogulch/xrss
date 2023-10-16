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
# build with xcaddy:
GOFLAGS='-tags="sqlite_json"' CGO_ENABLED=1 xcaddy build --with github.com/infogulch/xtemplate/caddy --with github.com/greenpau/caddy-security --with github.com/infogulch/xrss

# build for local development using xtemplate checked out next to xrss
GOFLAGS='-tags="sqlite_json"' CGO_ENABLED=1 XCADDY_DEBUG=1 xcaddy build --with github.com/infogulch/xtemplate/caddy=../xtemplate/caddy  --with github.com/infogulch/xtemplate=../xtemplate --with github.com/infogulch/xrss=. --with github.com/greenpau/caddy-security

tailwindcss -i static/main.css -o static/site.css --content './templates/**/*' -w
```

## Resources

https://hurl.dev/docs/hurl-file.html

https://github.com/go-echarts/go-echarts

https://nene.github.io/prettier-sql-playground/

https://github.com/cornelk/goscrape

https://github.com/robfig/cron

7 Restful Routes: https://gist.github.com/alexpchin/09939db6f81d654af06b

https://github.com/caddyserver/xcaddy/pull/62

https://github.com/infogulch/xrss/commits/master.atom

### CSS

* https://css-tricks.com/
* https://moderncss.dev/
* https://smolcss.dev/
* https://modernfontstacks.com/

* https://classless.de/classless.css
* https://purifycss.online/
* https://tailwindcss.com/docs/
* https://www.hyperui.dev/
* Responsive layout demo: https://play.tailwindcss.com/KJTJ5n574r
* https://github.com/argyleink/open-props/
* https://yqnn.github.io/svg-path-editor/
