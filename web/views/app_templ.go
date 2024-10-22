// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	ff_entity "ff/internal/feature_flag/entity"
	"ff/web/components"
)

func AppPage() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>HTMX Feature Flags Demo</title><!-- Include HTMX from CDN --><script src=\"https://unpkg.com/htmx.org@1.9.4\"></script><!-- Include Tailwind CSS from CDN --><link href=\"https://cdn.jsdelivr.net/npm/tailwindcss@2.1.2/dist/tailwind.min.css\" rel=\"stylesheet\"><!--  Font awesome --><script src=\"https://kit.fontawesome.com/934cef5fae.js\" crossorigin=\"anonymous\"></script><!-- Include Hyperscript from CDN --><script src=\"https://unpkg.com/hyperscript.org@0.9.13\"></script><style>\n\t\tbody {\n\t\t\tfont-family: Arial, sans-serif;\n\t\t}\n\n\t\tbutton {\n\t\t\tpadding: 10px 20px;\n\t\t\tborder: none;\n\t\t\tborder-radius: 5px;\n\t\t\tcursor: pointer;\n\t\t}\n\t</style><script type=\"text/javascript\">\n\t\tdocument.addEventListener(\"DOMContentLoaded\", (event) => {\n\t\t\tdocument.body.addEventListener('htmx:beforeSwap', function (evt) {\n\t\t\t\tif ([400, 404, 409, 401, 403].includes(evt.detail.xhr.status)) {\n\t\t\t\t\tconsole.log(\"setting status to paint\");\n\t\t\t\t\t// allow 400 errors to swap as we are using this as a signal that\n\t\t\t\t\t// a form was submitted with bad data and want to rerender with the\n\t\t\t\t\t// errors\n\t\t\t\t\t//\n\t\t\t\t\t// set isError to false to avoid error logging in console\n\t\t\t\t\tevt.detail.shouldSwap = true;\n\t\t\t\t\tevt.detail.isError = false;\n\t\t\t\t}\n\t\t\t});\n\t\t});\n\t</script></head><body class=\"pt-8 pl-8 pb-10 pr-6\"><script>\n\t\t// Enable HTMX logging\n\t\thtmx.logAll();\n\t</script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Header().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<main id=\"container\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</main>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Modal(false, ff_entity.FeatureFlagResponse{}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Message(false, "", false).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
