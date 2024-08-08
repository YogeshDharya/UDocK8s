package main

import (
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"path/filepath"
	"github.com/spf13/viper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)
type ApiResponse struct{
     Pagination struct{
		Limit int `json:"limit"`
		Offset int `json:"offset"`
		Count int `json:"count"`
		Total int `json:"total"`
	 }`json: "pagination"`
	 Data []Article `json:"data"`
}
type Article struct {
        Author        string `json:"author"`
        Title           string `json:"title"`
        Description string `json:"description"`
        URL             string `json:"url"`
		Source  string `json:"source"` 
        Image       string `json:"image"`
        Category    string `json:"category"`
        Language    string `json:"language"`
		Country  string  `json:"country"`
        PublishedAt   string `json:"published_at"`
}
func initConfig(){
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("../../")
	viper.AutomaticEnv()
	viper.BindEnv("NEWS_SERVICE_PORT")
	viper.BindEnv("API_TOKEN")
	err := viper.ReadInConfig()
	if err != nil{
		panic("Error Reading Config env file : " + err.Error())
	}
}
func main(){
	//Load Environment Variable 
  //no need to have absolute file path anymore ?? 
	initConfig()
	
	publicDir, err := filepath.Abs("../../public")
	if err != nil {
		log.Fatal("Can't resolve html template directory !")
	}
	engine := html.New(publicDir, ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	apiKey := viper.GetString("API_TOKEN")
	if apiKey == "" {
		log.Fatal("API key not Set !")
	}
	port := viper.GetString("NEWS_SERVICE_PORT")
	if port == ""{
		log.Fatal("Couldn't fetch port details !")
	}
	log.Println("News Blog Works Fine Now !")
	app.Get("/news",func(c *fiber.Ctx) error {
		news, err := fetchNews(apiKey)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
        }
	return c.Render("index",news)
	})
	app.Static("/","./public")
	log.Fatal(app.Listen(fmt.Sprintf(":%s",port)))
}
func fetchNews(APIkey string) (*ApiResponse, error){
	url := "https://api.currentsapi.services/v1/latest-news?apiKey=" + APIkey 
	res, err := http.Get(url);
	if err != nil {
		return nil , err 
	} 
	defer res.Body.Close() 
	var newsResponse ApiResponse
	if err := json.NewDecoder(res.Body).Decode(&newsResponse); err != nil {
		return nil , err
	}
	return &newsResponse,nil 
}