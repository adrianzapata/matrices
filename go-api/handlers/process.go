package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/interseguro/matrices-go-api/client"
	"github.com/interseguro/matrices-go-api/matrix"
	"github.com/interseguro/matrices-go-api/qr"
)

type processRequest struct {
	Matrix matrix.Matrix `json:"matrix"`
}

type processResponse struct {
	Original matrix.Matrix         `json:"original"`
	Rotated  matrix.Matrix         `json:"rotated"`
	Q        matrix.Matrix         `json:"q"`
	R        matrix.Matrix         `json:"r"`
	Stats    *client.StatsResponse `json:"stats"`
}

// ProcessMatrix wires together the whole pipeline described in the
// challenge:
//  1. Validate the incoming matrix.
//  2. Rotate it 90 degrees clockwise (the transform named explicitly in the
//     "Arquitectura de la solución" section).
//  3. Compute its QR factorization (the transform named explicitly in the
//     "Funcionalidad requerida" section — see README for how the two
//     descriptions were reconciled).
//  4. Forward the resulting matrices (rotated, Q, R) to the Node.js API,
//     which computes aggregate statistics over them.
//  5. Return everything to the caller in one response.
func ProcessMatrix(nodeClient *client.NodeClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req processRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
		}

		if err := matrix.Validate(req.Matrix); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		rotated := matrix.RotateClockwise90(req.Matrix)

		q, r, err := qr.Decompose(req.Matrix)
		if err != nil {
			if errors.Is(err, matrix.ErrEmptyMatrix) || errors.Is(err, matrix.ErrNotRectangular) {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to compute QR factorization"})
		}

		stats, err := nodeClient.FetchStats(c.Context(), map[string]matrix.Matrix{
			"rotated": rotated,
			"q":       q,
			"r":       r,
		})
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "failed to reach statistics API: " + err.Error()})
		}

		return c.JSON(processResponse{
			Original: req.Matrix,
			Rotated:  rotated,
			Q:        q,
			R:        r,
			Stats:    stats,
		})
	}
}
