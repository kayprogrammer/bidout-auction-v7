package routes
import 	"github.com/gofiber/fiber/v2"
	
type HealthCheckSchema struct {
	Success		string		`json:"success" example:"pong"`
}

// @Summary HealthCheck
// @Description This endpoint checks the health of our application.
// @Tags HealthCheck
// @Success 200 {object} HealthCheckSchema
// @Router /healthcheck [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"success": "pong"})
}