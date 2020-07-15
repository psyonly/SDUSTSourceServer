//Package for convenient to use JavaScript/H5 or other tools in dev
package gotools

import "html/template"

func JsAlert(str string) template.HTML {
	return template.HTML("<script>alert(' " + str + " ')</script>")
}
