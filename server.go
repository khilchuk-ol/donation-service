package main

import (
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"donation-service/internal/data"
	"donation-service/internal/services"
)

type Server struct {
	DonationService services.DonationService
	Logger          *log.Logger
}

func NewServer(ds services.DonationService, logger *log.Logger) *Server {
	return &Server{
		DonationService: ds,
		Logger:          logger,
	}
}

func (s *Server) Run() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	s.InstallRoutes(e)

	s.Logger.Fatal(e.Start(":8080"))
}

func (s *Server) InstallRoutes(e *echo.Echo) {
	e.GET("/hello", s.helloHandler)
	e.GET("/total", s.totalHandler)

	e.GET("/donations/new", s.newDonationsHandler)

	e.GET("/donations/new/test", s.newDonationsTestHandler)

	e.GET("/donations/max", s.maxDonationHandler)
}

func (s *Server) helloHandler(c echo.Context) error {
	return c.JSON(http.StatusOK,
		struct {
			Msg string `json:"msg"`
		}{
			Msg: "hello",
		})
}

func (s *Server) totalHandler(c echo.Context) error {
	return c.JSON(http.StatusOK,
		struct {
			Total float32 `json:"total"`
		}{
			Total: s.DonationService.GetTotalSum(),
		})
}

func (s *Server) newDonationsHandler(c echo.Context) error {
	return c.JSON(http.StatusOK,
		struct {
			Donations []data.Donation `json:"donations"`
		}{
			Donations: s.DonationService.GetNewDonationsFromCache(),
		})
}

func (s *Server) newDonationsTestHandler(c echo.Context) error {
	return c.JSON(http.StatusOK,
		struct {
			Donations []data.Donation `json:"donations"`
		}{
			Donations: s.DonationService.GetNewDonationsTest(),
		})
}

func (s *Server) maxDonationHandler(c echo.Context) error {
	d, ok := s.DonationService.GetMaxDonation()

	if ok {
		return c.JSON(http.StatusOK, d)
	}

	return c.NoContent(http.StatusInternalServerError)
}
