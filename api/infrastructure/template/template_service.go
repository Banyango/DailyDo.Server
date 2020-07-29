package template

import (
	. "github.com/Banyango/gifoody_server/api/infrastructure/os"
	"html/template"
	"io"
)

type TemplateInterface interface {
	RenderTemplate(buffer io.Writer, templateName string, data interface{}) error
}

type TemplateService struct {
	osService OSInterface
}

func NewTemplateService(service OSInterface) *TemplateService {
	m := new(TemplateService)
	m.osService = service
	return m
}

func (t *TemplateService) RenderTemplate(buffer io.Writer, templateName string, data interface{}) error {
	templateDir := t.osService.GetEnv("TEMPLATE_DIR")

	file, err := t.osService.GetFile(templateDir, templateName)
	if err != nil {
		return err
	}

	parsedTemplate, err := template.ParseFiles(file.Name())
	if err != nil {
		return err
	}

	if err := parsedTemplate.Execute(buffer, data); err != nil {
		return err
	}

	return nil
}
