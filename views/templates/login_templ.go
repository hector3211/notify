// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Login() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container mx-auto h-screen flex flex-col justify-center items-center\"><form method=\"POST\" hx-post=\"/login\" hx-trigger=\"submit\" hx-swap=\"beforeend\" class=\"w-full lg:w-1/2  shadow-sm shadow-stone-300 p-10 rounded-lg\"><p class=\"font-bold text-3xl p-6\">Login</p><div class=\"flex flex-col gap-2\"><label class=\"font-semibold text-lg\">Email</label> <input class=\"border border-stone-500 p-2 rounded-md\" id=\"email\" name=\"email\" type=\"text\" placeholder=\"johndoe@gmail.com\" required></div><div class=\"flex flex-col gap-2\"><label class=\"font-semibold text-lg\">Password</label> <input class=\"border border-stone-500 p-2 rounded-md\" id=\"password\" name=\"password\" type=\"text\" minlength=\"3\" maxlength=\"20\" placeholder=\"••••••••\" required></div><div class=\"container mt-5 w-full flex items-center gap-4\"><button type=\"submit\" class=\"text-center w-28 font-semibold p-2 px-3 rounded-md bg-zinc-900 text-white flex justify-center items-center space-x-2\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"lucide lucide-fingerprint\"><path d=\"M12 10a2 2 0 0 0-2 2c0 1.02-.1 2.51-.26 4\"></path> <path d=\"M14 13.12c0 2.38 0 6.38-1 8.88\"></path> <path d=\"M17.29 21.02c.12-.6.43-2.3.5-3.02\"></path> <path d=\"M2 12a10 10 0 0 1 18-6\"></path> <path d=\"M2 16h.01\"></path> <path d=\"M21.8 16c.2-2 .131-5.354 0-6\"></path> <path d=\"M5 19.5C5.5 18 6 15 6 12a6 6 0 0 1 .34-2\"></path> <path d=\"M8.65 22c.21-.66.45-1.32.57-2\"></path> <path d=\"M9 6.8a6 6 0 0 1 9 5.2v2\"></path></svg><p>Login</p></button> <a href=\"/signup\" class=\"text-center w-28 font-semibold p-2 px-3 rounded-md border border-black flex items-center space-x-2\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"20\" height=\"20\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"lucide lucide-square-arrow-out-up-right\"><path d=\"M21 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h6\"></path> <path d=\"m21 3-9 9\"></path> <path d=\"M15 3h6v6\"></path></svg><p>Sign up</p></a></div></form></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
