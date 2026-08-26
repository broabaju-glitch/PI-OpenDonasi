package utils

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

// SetupTemplates sets the template functions and loads the HTML templates
func SetupTemplates(r *gin.Engine) {
	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"formatRupiah": func(amount float64) string {
			intAmt := int64(amount)
			str := fmt.Sprintf("%d", intAmt)
			n := len(str)
			if n <= 3 {
				return str
			}
			var result []byte
			for i, c := range str {
				if i > 0 && (n-i)%3 == 0 {
					result = append(result, '.')
				}
				result = append(result, byte(c))
			}
			return string(result)
		},
	})

	r.LoadHTMLGlob("templates/admin/*")
}
