package main

import (
	"html/template"
	"log"
	"regexp"
	"time"

	"ireul.com/orm"

	_ "ireul.com/mysql"
	"ireul.com/redis"
	"ireul.com/web"
	"ireul.com/web/binding"
)

// NameRegexp regexp for record name
var NameRegexp = regexp.MustCompile(`\A[A-Za-z0-9._-]+\z`)

// CreateForm form of create
type CreateForm struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Token string `json:"token"`
}

func redisKey(s string) string {
	return "link." + s
}

func main() {
	log.SetPrefix("[linkage] ")

	// config
	c, err := ParseConfigFile("config.yaml")
	if err != nil {
		panic(err)
	}

	// db
	db, err := NewDB(c.DatabaseURL, c.Env != web.PROD)
	if err != nil {
		panic(err)
	}

	// redis
	r, err := redis.Open(c.RedisURL)
	if err != nil {
		panic(err)
	}

	// web
	w := web.Classic()
	w.SetEnv(c.Env)
	w.Use(web.Renderer())
	w.Use(Renderer())
	w.Map(c)
	w.Map(db)
	w.Map(r)

	w.Use(func(ctx *web.Context, c Config) {
		ctx.Data["Title"] = c.Title
	})

	w.Get(
		"/",
		func(r Render) {
			r.HTML(200, "index")
		},
	)

	w.Post(
		"/create",
		binding.Bind(CreateForm{}),
		func(ctx *web.Context, r Render, db *DB, rd *redis.Client, f CreateForm, c Config) {
			if f.Token != c.Token {
				r.Error(400, "token mismatch")
				return
			}
			if !NameRegexp.MatchString(f.Name) {
				r.Error(400, "name contains invalid characters")
				return
			}
			if f.URL == "" {
				r.Error(400, "url is empty")
				return
			}
			d := Record{Name: f.Name, URL: f.URL}
			if err := db.Create(&d).Error; err != nil {
				r.Error(400, err.Error())
				return
			}
			rd.Set(redisKey(d.Name), d.URL, time.Minute*5)
			r.JSON(200, d)
		},
	)
	w.Get("/:name", func(ctx *web.Context, r Render, db *DB, rd *redis.Client) {
		name := ctx.Params(":name")
		r.Data("Name", name)

		key := redisKey(name)
		url, err := rd.Get(key).Result()
		if err == nil {
			rd.Expire(key, time.Minute*5)
			r.Data("URL", url)
			r.Data("JSURL", template.JSEscapeString(url))
			r.HTML(200, "link")
			return
		}
		if err != nil && err != redis.Nil {
			log.Println("internal error, name =", name, ", err=", err.Error())
			r.Error(500, "internal error")
			return
		}
		d := Record{}
		if err := db.Where("name = ?", name).First(&d).Error; err != nil {
			if err == orm.ErrRecordNotFound {
				r.HTML(404, "not_found")
			} else {
				log.Println("internal error, name =", name, ", err=", err.Error())
				r.Error(500, "internal error")
			}
			return
		}
		rd.Set(key, d.URL, time.Minute*5)
		r.Data("URL", d.URL)
		r.Data("JSURL", template.JSEscapeString(d.URL))
		r.HTML(200, "link")
	})

	w.Run("127.0.0.1", c.Port)
}
