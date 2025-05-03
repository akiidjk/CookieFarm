package server

import (
	"math"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/database"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/models"
	"github.com/ByteTheCookies/cookieserver/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// HandleIndexPage renders the main dashboard page.
// It checks the cookie-based authentication and sets the default pagination limit.
func HandleIndexPage(c *fiber.Ctx) error {
	if err := CookieAuthMiddleware(c); err != nil {
		return err
	}

	limit := c.QueryInt("limit", config.DEFAULT_LIMIT)
	if limit <= 0 {
		limit = config.DEFAULT_LIMIT
	}

	logger.Log.Info().Int("Limit", limit).Msg("Index page request")
	data := models.ViewParamsDashboard{
		Limit: limit,
	}
	return c.Render("pages/dashboard", data, "layouts/main")
}

// HandleLoginPage renders the login page.
func HandleLoginPage(c *fiber.Ctx) error {
	return c.Render("pages/login", map[string]any{}, "layouts/main")
}

// HandlePartialsPagination renders only the pagination component as a partial view.
// It computes the current page and the total number of pages based on the flags count.
func HandlePartialsPagination(c *fiber.Ctx) error {
	if err := CookieAuthMiddleware(c); err != nil {
		return err
	}

	limit, err := c.ParamsInt("limit", config.DEFAULT_LIMIT)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid 'limit' parameter")
	}
	logger.Log.Debug().Int("limit", limit).Msg("Paginated flags request")

	totalFlags, err := database.FlagsNumber(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving flag count")
	}

	offset := c.QueryInt("offset", config.DEFAULT_OFFSET)

	totalPages := int(math.Ceil(float64(totalFlags) / float64(limit)))
	current := offset / limit
	pageList := utils.MakePagination(current, totalPages)

	data := models.ViewParamsPagination{
		Pagination: models.Pagination{
			Limit:    limit,
			Pages:    totalPages,
			Current:  current,
			HasPrev:  current > 0,
			HasNext:  current < totalPages-1,
			PageList: pageList,
		},
	}

	return c.Render("partials/pagination", data, "layouts/main")
}

// HandlePartialsFlags renders only the flags rows as a partial view.
// It fetches a limited and paginated list of flags from the database.
func HandlePartialsFlags(c *fiber.Ctx) error {
	if err := CookieAuthMiddleware(c); err != nil {
		return err
	}

	limit, err := c.ParamsInt("limit", config.DEFAULT_LIMIT)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid 'limit' parameter")
	}

	offset := c.QueryInt("offset", config.DEFAULT_OFFSET)
	logger.Log.Debug().Int("offset", offset).Int("limit", limit).Msg("Paginated flags request")

	flags, err := database.GetPagedFlags(uint(limit), uint(offset))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving flags")
	}

	logger.Log.Debug().Int("n_flags", len(flags)).Msg("Paginated flags response")

	data := models.ViewParamsFlags{
		Flags: flags,
	}

	return c.Render("partials/flags_rows", data)
}
