// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Signup() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container mx-auto h-screen flex flex-col justify-center items-center\"><form method=\"POST\" hx-post=\"/signup\" hx-trigger=\"submit\" hx-swap=\"beforeend\" class=\"w-full lg:w-1/2 shadow-sm shadow-stone-300 p-10 rounded-lg\"><p class=\"font-bold text-3xl p-6\">Sign up</p><div class=\"flex flex-col gap-2\"><div class=\"flex flex-col md:flex-row items-center space-x-6\"><span class=\"container flex flex-col\"><label class=\"font-semibold text-lg\">First Name</label> <input class=\"border border-black p-2 rounded-md\" id=\"firstname\" name=\"firstname\" type=\"text\" placeholder=\"doe\" required></span> <span class=\"container flex flex-col\"><label class=\"font-semibold text-lg\">Last Name</label> <input class=\"border border-black p-2 rounded-md\" id=\"lastname\" name=\"lastname\" type=\"text\" placeholder=\"doe\" required></span></div></div><div class=\"flex flex-col gap-2\"><label class=\"font-semibold text-lg\">Email</label> <input class=\"border border-stone-500 p-2 rounded-md\" id=\"email\" name=\"email\" type=\"text\" placeholder=\"johndoe@gmail.com\" required></div><div class=\"flex flex-col gap-2\"><label class=\"font-semibold text-lg\">Password</label> <input class=\"border border-stone-500 p-2 rounded-md\" id=\"signup-password\" name=\"password\" type=\"text\" minlength=\"3\" maxlength=\"20\" placeholder=\"••••••••\" required></div><div class=\"container mt-5 w-full flex items-center gap-5\"><button class=\"flex justify-center space-x-2 font-semibold items-center p-2 px-3 rounded-md bg-zinc-900 text-white\" type=\"submit\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"lucide lucide-log-in\"><path d=\"M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4\"></path> <polyline points=\"10 17 15 12 10 7\"></polyline> <line x1=\"15\" x2=\"3\" y1=\"12\" y2=\"12\"></line></svg><p>Sign up</p></button> <a href=\"/login\" class=\"mt-3 text-center font-semibold flex items-center space-x-2 underline underline-offset-4\"><p>Already have an account</p></a></div></form></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
