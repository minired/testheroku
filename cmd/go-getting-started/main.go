package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v1"
	"database/sql"
	_"github.com/lib/pq"
)

var (
	Dbm *gorp.DbMap
)
func main() {
	InitDB()
	port := os.Getenv("PORT")
	
	if port == "" {
		port = "4747"
		log.Println("[-] No PORT environment variable detected. Setting to ", port)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("../../templates/*.tmpl.html")
	router.Static("../../static", "static")

	router.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	
	
	router.GET("/someJSON", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
    })

    router.GET("/moreJSON/:id", func(c *gin.Context) {
        // You also can use a struct
        var msg struct {
            Name    string `json:"user"`
            Message string
            Number  int
        }
		id := c.Param("id")
		name, _ := Dbm.SelectStr("SELECT \"Name\" FROM \"Test\" where \"ID\" =$1", id)
		
        msg.Name = name
        msg.Message = "hey"
        msg.Number = 123
        // Note that msg.Name becomes "user" in the JSON
        // Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}
        c.JSON(http.StatusOK, msg)
    })

    router.GET("/someXML", func(c *gin.Context) {
        c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
    })

	router.Run(":" + port)
}


///// DB



func InitDB() {
    // connect to db using standard Go database/sql API
    // use whatever database/sql driver you wish
	databseurl := "postgres://euvolywvhwcsfq:DXqvao2sMmQbgHGlJO877oxfgj@ec2-107-22-187-89.compute-1.amazonaws.com:5432/d553fh1pbahm1t"
    db, err := sql.Open("postgres", databseurl)
    checkErr(err, "postgres.sql.Open failed")

    // construct a gorp DbMap
    Dbm = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
	//dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")

    // create the table. in a production system you'd generally
    // use a migration tool, or create the tables via scripts
    //err = dbmap.CreateTablesIfNotExists()
    //checkErr(err, "Create tables failed")

    //return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
