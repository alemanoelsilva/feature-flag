package main

import (
	"ff/web/app/services"
	"ff/web/app/utils"
	"html/template"
	"io"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const INDEX_PAGE = "index.html"
const PAGINATION_DEFAULT = 500

type Template struct {
	tmpl *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("web/views/**/*.html")),
	}
}

type Header struct {
	Name        string
	IsEditing   bool
	IsAssigning bool
}

type CreateFeatureFlagForm struct {
	ID               int
	Name             string
	ErrorName        bool
	Description      string
	ErrorDescription bool
	IsActive         bool
	ExpirationDate   string
	ErrorPage        string
	ShowErrorMessage bool
}

type CreateAssignmentForm struct {
	ID                      int
	FeatureFlagID           int
	FeatureFlagName         string
	IsGlobal                bool
	People                  []services.Person
	FeatureFlags            []services.FeatureFlag
	IsFeatureFlagToggleOpen bool
	Pagination              services.Pagination
	ErrorPage               string
	ShowErrorMessage        bool
}

type WebApp struct {
	Header          Header
	FeatureFlags    []services.FeatureFlag
	FeatureFlagForm CreateFeatureFlagForm
	AssignmentForm  CreateAssignmentForm
}

func initFeatureFlagForm() CreateFeatureFlagForm {
	return CreateFeatureFlagForm{
		ID:               0,
		Name:             "",
		ErrorName:        false,
		Description:      "",
		ErrorDescription: false,
		IsActive:         false,
		ExpirationDate:   "",
		ErrorPage:        "",
		ShowErrorMessage: false,
	}
}

func initAssignmentForm() CreateAssignmentForm {
	return CreateAssignmentForm{
		ID:                      0,
		FeatureFlagID:           0,
		FeatureFlagName:         "",
		IsGlobal:                false,
		People:                  []services.Person{},
		FeatureFlags:            []services.FeatureFlag{},
		IsFeatureFlagToggleOpen: false,
		Pagination: services.Pagination{
			Page:  1,
			Limit: PAGINATION_DEFAULT,
			Next:  true,
		},
		ErrorPage:        "",
		ShowErrorMessage: false,
	}
}

func initWebApp() *WebApp {
	var ff services.FeatureFlag
	return &WebApp{
		Header: Header{
			Name:        "Header",
			IsEditing:   false,
			IsAssigning: false,
		},
		FeatureFlags:    ff.GetFeatureFlag(),
		FeatureFlagForm: initFeatureFlagForm(),
		AssignmentForm:  initAssignmentForm(),
	}
}

func main() {
	e := echo.New()

	e.Renderer = newTemplate()
	e.Use(middleware.Logger())
	// e.Static("/css", "css")
	e.Static("/images", "web/images")

	webApp := initWebApp()

	/**
	* Index page
	 */
	e.GET("/", func(c echo.Context) error {
		webApp = initWebApp()
		return c.Render(200, INDEX_PAGE, webApp)
	})

	/**
	* Filter feature flag list
	 */
	e.GET("/feature-flags-filter", func(c echo.Context) error {
		ff := services.FeatureFlag{
			Name:     c.QueryParam("name"),
			IsActive: c.FormValue("isActive") == "on",
			IsGlobal: c.FormValue("isGlobal") == "on",
		}

		webApp.FeatureFlags = ff.GetFeatureFlag()
		return c.Render(200, "feature_flags", webApp.FeatureFlags)
	})

	/**
	* Form block to create feature flag
	 */
	e.GET("/new-feature-flag-form", func(c echo.Context) error {
		webApp.FeatureFlagForm = initFeatureFlagForm()
		webApp.Header = Header{
			IsEditing:   true,
			IsAssigning: false,
		}
		return c.Render(200, INDEX_PAGE, webApp)
	})

	/**
	* Form block to update feature flag
	 */
	e.GET("/feature-flag-form", func(c echo.Context) error {
		idStr := c.QueryParam("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("Error getting feature flag id: %v", err)
		}

		featureFlag := utils.FindFeatureFlagByID(id, &webApp.FeatureFlags)

		webApp.FeatureFlagForm = CreateFeatureFlagForm{
			ID:             featureFlag.ID,
			Name:           featureFlag.Name,
			Description:    featureFlag.Description,
			IsActive:       featureFlag.IsActive,
			ExpirationDate: featureFlag.ExpirationDate,
		}
		webApp.Header = Header{
			IsEditing:   true,
			IsAssigning: false,
		}
		return c.Render(200, INDEX_PAGE, webApp)
	})

	/**
	* Create new feature flag and return to the updated feature flag list
	 */
	e.POST("/feature-flag", func(c echo.Context) error {
		ff := services.FeatureFlag{
			Name:           c.FormValue("name"),
			Description:    c.FormValue("description"),
			IsActive:       c.FormValue("isActive") == "on",
			ExpirationDate: c.FormValue("expirationDate"),
		}

		apiError := ff.CreateFeatureFlag()

		if apiError.IsError {
			var errorMessages string

			if apiError.Field == "Name" {
				errorMessages += apiError.Message + "\n" // Add newline after each message
			}
			if apiError.Field == "Description" {
				errorMessages += apiError.Message + "\n"
			}
			if apiError.Field == "Page" {
				errorMessages += apiError.Message + "\n"
			}
			webApp.FeatureFlagForm = CreateFeatureFlagForm{
				Name:             ff.Name,
				ErrorName:        apiError.Field == "Name",
				Description:      ff.Description,
				ErrorDescription: apiError.Field == "Description",
				IsActive:         ff.IsActive,
				ExpirationDate:   ff.ExpirationDate,
				ErrorPage:        errorMessages,
				ShowErrorMessage: true,
			}
			webApp.Header = Header{
				IsEditing:   true,
				IsAssigning: false,
			}
		} else {
			webApp = initWebApp()
		}
		return c.Render(200, INDEX_PAGE, webApp)
	})

	/**
	* Update a feature flag and return to the updated feature flag list
	 */
	e.PUT("/feature-flag/:id", func(c echo.Context) error {
		id, err := utils.GetNumberParamFromRequest("id", c)
		if err != nil {
			return c.String(400, "Id must be an integer")
		}
		// by default, create inactive, active only after having at least one assignment (like BP)
		ff := services.FeatureFlag{
			ID:             id,
			Description:    c.FormValue("description"),
			IsActive:       c.FormValue("isActive") == "on",
			ExpirationDate: c.FormValue("expirationDate"),
		}

		apiError := ff.UpdateFeatureFlag()

		// TODO: handle 500
		if apiError.IsError {
			var errorMessages string

			if apiError.Field == "Name" {
				errorMessages += apiError.Message + "\n" // Add newline after each message
			}
			if apiError.Field == "Description" {
				errorMessages += apiError.Message + "\n"
			}
			if apiError.Field == "Page" {
				errorMessages += apiError.Message + "\n"
			}
			webApp.FeatureFlagForm = CreateFeatureFlagForm{
				ID:               id,
				Name:             webApp.FeatureFlagForm.Name,
				Description:      ff.Description,
				ErrorDescription: apiError.Field == "Description",
				IsActive:         ff.IsActive,
				ExpirationDate:   ff.ExpirationDate,
				ErrorPage:        errorMessages,
				ShowErrorMessage: true,
			}
			webApp.Header = Header{
				IsEditing:   true,
				IsAssigning: false,
			}
		} else {
			webApp = initWebApp()
		}
		return c.Render(200, INDEX_PAGE, webApp)
	})

	/**
	* Update a feature flag and return to the updated feature flag list
	 */
	e.PUT("/feature-flag/active/:id", func(c echo.Context) error {
		id, err := utils.GetNumberParamFromRequest("id", c)
		if err != nil {
			return c.String(400, "Id must be an integer")
		}

		featureFlag := utils.FindFeatureFlagByID(id, &webApp.FeatureFlags)
		ff := services.FeatureFlag{
			ID:             id,
			Description:    featureFlag.Description,
			IsActive:       !featureFlag.IsActive,
			ExpirationDate: featureFlag.ExpirationDate,
		}
		ff.UpdateFeatureFlag()

		webApp = initWebApp()
		return c.Render(200, INDEX_PAGE, webApp)
	})

	/**
	* Form block to apply assignment
	 */
	e.GET("/assignment-form", func(c echo.Context) error {
		next := c.QueryParam("next") == "true"
		idStr := c.QueryParam("id")
		var ff services.FeatureFlag
		featureFlags := ff.GetFeatureFlag()
		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Error getting feature flag id: %v", err)
			}

			ff = utils.FindFeatureFlagByID(id, &webApp.FeatureFlags)
		}

		paginationReq := services.Pagination{
			Page:  1,
			Limit: PAGINATION_DEFAULT,
		}
		if next {
			paginationReq.Page = webApp.AssignmentForm.Pagination.Page
			paginationReq.Limit = webApp.AssignmentForm.Pagination.Limit
		}

		var p services.Person
		people, pagination := p.GetPerson(ff.ID, paginationReq)
		webApp.AssignmentForm = CreateAssignmentForm{
			ID:                      0,
			FeatureFlagID:           ff.ID,
			FeatureFlagName:         ff.Name,
			IsGlobal:                ff.IsGlobal,
			People:                  people,
			FeatureFlags:            featureFlags,
			IsFeatureFlagToggleOpen: false,
			Pagination: services.Pagination{
				Page:  pagination.Page,
				Limit: pagination.Limit,
				Next:  pagination.Next,
			},
		}
		webApp.Header = Header{
			IsEditing:   false,
			IsAssigning: true,
		}
		return c.Render(200, INDEX_PAGE, webApp)
	})

	/**
	* Filter people list
	 */
	e.GET("/people-filter", func(c echo.Context) error {
		next := c.QueryParam("next") == "true"

		if webApp.AssignmentForm.FeatureFlagID == 0 {
			webApp.AssignmentForm.ErrorPage = "Select a Feature Flag before filtering the people list"
			webApp.AssignmentForm.ShowErrorMessage = true

			return c.Render(200, "people_list", webApp.AssignmentForm)
		}

		p := services.Person{
			Name:       c.QueryParam("name"),
			IsAssigned: c.FormValue("isAssigned") == "on",
		}

		paginationReq := services.Pagination{
			Page:  1,
			Limit: PAGINATION_DEFAULT,
		}
		if next {
			paginationReq.Page = webApp.AssignmentForm.Pagination.Page
			paginationReq.Limit = webApp.AssignmentForm.Pagination.Limit
		}

		people, pagination := p.GetPerson(webApp.AssignmentForm.FeatureFlagID, paginationReq)

		webApp.AssignmentForm.People = people
		webApp.AssignmentForm.Pagination = pagination
		// webApp.AssignmentForm.
		return c.Render(200, "people_list", webApp.AssignmentForm)
	})

	/**
	* Apply or remove an assignment
	 */
	e.POST("/assignment/:personId", func(c echo.Context) error {
		if webApp.AssignmentForm.FeatureFlagID == 0 {
			webApp.AssignmentForm = CreateAssignmentForm{
				ID:               0,
				FeatureFlagID:    webApp.AssignmentForm.FeatureFlagID,
				FeatureFlagName:  webApp.AssignmentForm.FeatureFlagName,
				IsGlobal:         webApp.AssignmentForm.IsGlobal,
				People:           webApp.AssignmentForm.People,
				FeatureFlags:     webApp.AssignmentForm.FeatureFlags,
				ErrorPage:        "Select a Feature Flag before making assignment to it",
				ShowErrorMessage: true,
			}

			return c.Render(200, "assignment_form", webApp.AssignmentForm)
		}

		id, err := utils.GetNumberParamFromRequest("personId", c)
		if err != nil {
			log.Fatalf("Error getting person id: %v", err)
		}
		ar := services.AssignmentRequest{
			FeatureFlagID: webApp.AssignmentForm.FeatureFlagID,
			PersonID:      id,
		}

		personToAssign := utils.FindPersonByID(id, &webApp.AssignmentForm.People)

		if !personToAssign.IsAssigned {
			ar.ApplyAssignment()
		} else {
			ar.DeleteAssignment()
		}

		var p services.Person
		people, pagination := p.GetPerson(webApp.AssignmentForm.FeatureFlagID, services.Pagination{
			Page:  1,
			Limit: PAGINATION_DEFAULT,
		})
		webApp.AssignmentForm.People = people
		webApp.AssignmentForm.Pagination = pagination
		return c.Render(200, "assignment_form", webApp.AssignmentForm)
	})

	/**
	* Apply global assignment
	 */
	e.POST("/assignment/global", func(c echo.Context) error {
		if webApp.AssignmentForm.FeatureFlagID == 0 {
			webApp.AssignmentForm = CreateAssignmentForm{
				ID:               0,
				FeatureFlagID:    webApp.AssignmentForm.FeatureFlagID,
				FeatureFlagName:  webApp.AssignmentForm.FeatureFlagName,
				IsGlobal:         webApp.AssignmentForm.IsGlobal,
				People:           webApp.AssignmentForm.People,
				FeatureFlags:     webApp.AssignmentForm.FeatureFlags,
				ErrorPage:        "Select a Feature Flag before making assignment to it",
				ShowErrorMessage: true,
			}

			return c.Render(200, INDEX_PAGE, webApp)
		}

		isGlobal := c.QueryParam("isGlobal") == "true"
		featureFlag := utils.FindFeatureFlagByID(webApp.AssignmentForm.FeatureFlagID, &webApp.AssignmentForm.FeatureFlags)

		var ff services.FeatureFlag
		ff.ID = featureFlag.ID
		ff.Name = featureFlag.Name
		ff.Description = featureFlag.Description
		ff.IsActive = featureFlag.IsActive
		ff.ExpirationDate = featureFlag.ExpirationDate
		ff.IsGlobal = isGlobal
		ff.UpdateFeatureFlag()

		var p services.Person
		people, pagination := p.GetPerson(webApp.AssignmentForm.FeatureFlagID, services.Pagination{
			Page:  1,
			Limit: PAGINATION_DEFAULT,
		})
		webApp.AssignmentForm.People = people
		webApp.AssignmentForm.Pagination = pagination
		webApp.AssignmentForm.IsGlobal = ff.IsGlobal
		// return c.Render(200, "assignment_form", webApp.AssignmentForm)
		return c.Render(200, INDEX_PAGE, webApp)
	})

	/**
	* Handle Dropdown
	 */
	e.GET("/toggle-dropdown", func(c echo.Context) error {
		var tmpl string

		if !webApp.AssignmentForm.IsFeatureFlagToggleOpen {
			webApp.AssignmentForm.IsFeatureFlagToggleOpen = true
			tmpl =
				`<div id="dropdown" class="absolute bg-white divide-y divide-gray-100 rounded-lg shadow">
					<ul class="py-2 text-sm text-gray-700">
						{{ range . }}
							<li>
								<span class="block cursor-pointer px-4 py-2 hover:bg-gray-100" hx-get="/assignment-form?id={{ .ID }}" hx-target="body" hx-swap="outerHTML swap:100ms">
									{{ .Name }}
								</span>
							</li>
						{{ end }}
					</ul>
				</div>`

		} else {
			webApp.AssignmentForm.IsFeatureFlagToggleOpen = false
			tmpl = `<div id="dropdown" hidden></div>`
		}

		t, err := template.New("dropdown").Parse(tmpl)
		if err != nil {
			return err
		}

		// Return the dropdown HTML as the response
		return t.Execute(c.Response().Writer, webApp.AssignmentForm.FeatureFlags)
	})

	e.GET("/error/dismiss", func(c echo.Context) error {

		webApp.AssignmentForm.ShowErrorMessage = false
		webApp.FeatureFlagForm.ShowErrorMessage = false
		return c.Render(200, INDEX_PAGE, webApp)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
