package main
import(
	"encoding/json"
	"log"
	"net/http"
	"github.com/spf13/viper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)
func initConfig(){
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.BindEnv("NEWS_SERVICE_PORT")
	viper.BindEnv("USER_SERVICE_PORT")
} 
func main(){
	app := fiber.New()
	newsPort := viper.GetString("NEWS_SERVICE_PORT")
	userPort := viper.GetString("USER_SERVICE_PORT")
	app.Use("/",func(c *fiber.Ctx) error {
		if c.Path() == "/news" {
				return proxy.Do(c,"http://localhost"+newsPort)
		}
        return proxy.Do(c,"http://localhost:"+ userPort)
	})
	app.Static("/","./public")
	log.Fatal(app.Listen(":3003"))
}
