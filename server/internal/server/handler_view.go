package server

import (
	"math"

	"github.com/ByteTheCookies/cookieserver/internal/config"
	"github.com/ByteTheCookies/cookieserver/internal/logger"
	"github.com/ByteTheCookies/cookieserver/internal/models"
	"github.com/ByteTheCookies/cookieserver/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) HandleIndexPage(c *fiber.Ctx) error {
	if err := CookieAuthMiddleware(c); err != nil {
		return err
	}

	limit := c.QueryInt("limit", config.LIMIT)
	if limit <= 0 {
		limit = config.LIMIT
	}

	logger.Log.Info().Int("Limit", limit).Msg("Index page request")
	data := models.ViewParamsDashboard{
		Limit: limit,
	}
	return c.Render("pages/dashboard", data, "layouts/main")
}

func (s *FiberServer) HandleLoginPage(c *fiber.Ctx) error {
	return c.Render("pages/login", map[string]any{}, "layouts/main")
}

func (s *FiberServer) HandlePartialsPagination(c *fiber.Ctx) error {
	if err := CookieAuthMiddleware(c); err != nil {
		return err
	}

	limit, err := c.ParamsInt("limit", config.LIMIT)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Parametro limit non valido")
	}
	logger.Log.Debug().Int("limit", limit).Msg("Paginated flags request")

	totalFlags, err := s.db.FlagsNumber(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Errore nel recupero dei dati")
	}

	offset := c.QueryInt("offset", config.OFFSET)

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

func (s *FiberServer) HandlePartialsFlags(c *fiber.Ctx) error {
	if err := CookieAuthMiddleware(c); err != nil {
		return err
	}

	limit, err := c.ParamsInt("limit", config.LIMIT)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Parametro limit non valido")
	}

	offset := c.QueryInt("offset", config.OFFSET)
	logger.Log.Debug().Int("offset", offset).Int("limit", limit).Msg("Paginated flags request")

	flags, err := s.db.GetPagedFlags(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Errore nel recupero dei dati")
	}

	logger.Log.Debug().Int("n_flags", len(flags)).Msg("Paginated flags response")

	data := models.ViewParamsFlags{
		Flags: flags,
	}

	return c.Render("partials/flags_rows", data)
}
