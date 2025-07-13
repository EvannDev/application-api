package server

import (
	"application-api/internal/argocd"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	v1 := s.Group("/api/v1")
	v1.Get("/apps/:namespace/:name", s.GetApplicationDetailsHandler)
	v1.Get("/apps/:namespace", s.GetAllApplicationsHandler)
}

func (s *FiberServer) GetAllApplicationsHandler(c *fiber.Ctx) error {
	log.Printf("Received request for all application namespace: %s", c.Params("namespace"))

	applications, err := argocd.GetAllApplications(c.Params("namespace"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp := fiber.Map{
		"applications": []fiber.Map{},
	}

	for _, app := range applications {
		appData := fiber.Map{
			"name":       app.Name,
			"version":    app.Status.Sync.Revision,
			"syncStatus": app.Status.Sync.Status,
			"health":     app.Status.Health.Status,
		}
		resp["applications"] = append(resp["applications"].([]fiber.Map), appData)
	}

	return c.JSON(resp)
}

func (s *FiberServer) GetApplicationDetailsHandler(c *fiber.Ctx) error {
	log.Printf("Received request for application: %s in namespace: %s", c.Params("name"), c.Params("namespace"))

	applications, err := argocd.GetApplicationsByName(c.Params("name"), c.Params("namespace"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp := fiber.Map{
		"name":       applications.Name,
		"namespace":  applications.Namespace,
		"status":     applications.Status.Sync.Status,
		"health":     applications.Status.Health.Status,
		"conditions": applications.Status.Conditions,
		"resources":  applications.Status.Resources,
		"source": fiber.Map{
			"repoURL":        applications.Spec.Source.RepoURL,
			"path":           applications.Spec.Source.Path,
			"targetRevision": applications.Spec.Source.TargetRevision,
		},
		"destination": fiber.Map{
			"server":    applications.Spec.Destination.Server,
			"namespace": applications.Spec.Destination.Namespace,
		},
		"project":   applications.Spec.Project,
		"createdAt": applications.ObjectMeta.CreationTimestamp.String(),
	}

	return c.JSON(resp)
}
