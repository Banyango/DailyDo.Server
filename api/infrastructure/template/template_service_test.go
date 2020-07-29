package template

import (
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	//file, _ := ioutil.TempFile("./", "123.html")
	//file.WriteString("Hi")
	//defer os.Remove(file.Name())
	//
	//envInterface := new(mocks.OSInterface)
	//envInterface.On("GetEnv", "TEMPLATE_DIR").Return("app/")
	//envInterface.On("GetFile", "app/","test.html").Return(file, nil)
	//
	//templateService := NewTemplateService(envInterface)
	//
	//var b bytes.Buffer
	//err := templateService.RenderTemplate(&b, "test.html", nil)
	//
	//assert.Nil(t, err)
	//assert.Equal(t, "Hi",string(b.Bytes()))
}
