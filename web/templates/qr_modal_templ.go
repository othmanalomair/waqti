// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.865
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"waqti/internal/models"
)

func QRModal(creator *models.Creator, settings *models.ShopSettings, lang string, isRTL bool) templ.Component {
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
		templ_7745c5c3_Err = QRModalWithQR(creator, settings, "", lang, isRTL).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func QRModalWithQR(creator *models.Creator, settings *models.ShopSettings, qrCodeDataURL string, lang string, isRTL bool) templ.Component {
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
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!-- Modal Backdrop --><div id=\"qr-modal\" class=\"qr-modal-backdrop fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4\" x-data=\"{ show: true }\" x-show=\"show\" x-transition:enter=\"transition ease-out duration-300\" x-transition:enter-start=\"opacity-0\" x-transition:enter-end=\"opacity-100\" x-transition:leave=\"transition ease-in duration-200\" x-transition:leave-start=\"opacity-100\" x-transition:leave-end=\"opacity-0\" @click.self=\"show = false; setTimeout(() =&gt; document.getElementById(&#39;qr-modal&#39;).remove(), 200)\"><!-- Modal Content --><div class=\"qr-modal-content bg-white rounded-3xl max-w-sm w-full card-shadow transform\" x-show=\"show\" x-transition:enter=\"transition ease-out duration-300\" x-transition:enter-start=\"opacity-0 scale-95 translate-y-4\" x-transition:enter-end=\"opacity-100 scale-100 translate-y-0\" x-transition:leave=\"transition ease-in duration-200\" x-transition:leave-start=\"opacity-100 scale-100 translate-y-0\" x-transition:leave-end=\"opacity-0 scale-95 translate-y-4\"><!-- Header with Close Button --><div class=\"flex items-center justify-between p-6 pb-4\"><h2 class=\"text-xl font-bold text-slate-charcoal\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "رمز QR للمتجر")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "Store QR Code")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</h2><button @click=\"show = false; setTimeout(() =&gt; document.getElementById(&#39;qr-modal&#39;).remove(), 200)\" class=\"p-2 hover:bg-gray-100 rounded-full transition-colors\"><svg class=\"w-5 h-5 text-gray-500\" fill=\"currentColor\" viewBox=\"0 0 24 24\"><path d=\"M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z\"></path></svg></button></div><!-- Content --><div class=\"px-6 pb-6 text-center\"><!-- Creator Info --><div class=\"mb-6\"><!-- Creator Avatar -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if settings != nil && settings.LogoURL != "" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "<img src=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(settings.LogoURL)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/qr_modal.templ`, Line: 62, Col: 51}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "\" alt=\"Creator Avatar\" class=\"w-20 h-20 rounded-full object-cover mx-auto mb-4\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "<div class=\"w-20 h-20 bg-gradient-to-br from-gulf-teal to-teal-600 rounded-full flex items-center justify-center mx-auto mb-4\"><span class=\"text-white font-bold text-2xl\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if lang == "ar" {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "أ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "A")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "</span></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "<!-- Creator Name --><h3 class=\"text-lg font-bold text-slate-charcoal mb-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(creator.NameAr)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/qr_modal.templ`, Line: 78, Col: 44}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(creator.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/qr_modal.templ`, Line: 80, Col: 42}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "</h3><!-- Description -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if settings != nil {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "<p class=\"text-gray-600 text-sm\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if lang == "ar" && settings.SubHeaderAr != "" {
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(settings.SubHeaderAr)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/qr_modal.templ`, Line: 88, Col: 54}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else if settings.SubHeader != "" {
				var templ_7745c5c3_Var7 string
				templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(settings.SubHeader)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/qr_modal.templ`, Line: 90, Col: 52}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				if lang == "ar" {
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, "منشئ محتوى على Waqti.me")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, "Content Creator on Waqti.me")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 16, "</p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 17, "<p class=\"text-gray-600 text-sm\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if lang == "ar" {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 18, "منشئ محتوى على Waqti.me")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 19, "Content Creator on Waqti.me")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 20, "</p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 21, "</div><!-- QR Code Container --><div class=\"qr-code-container bg-white p-6 rounded-2xl border-2 border-gray-100 mb-6\"><!-- QR Code (Real QR Code for creator's store) --><div class=\"qr-code w-48 h-48 mx-auto bg-white rounded-lg flex items-center justify-center\" id=\"qr-code-display\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if qrCodeDataURL != "" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 22, "<!-- Server-generated QR Code --> <img src=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var8 string
			templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(qrCodeDataURL)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/qr_modal.templ`, Line: 116, Col: 52}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 23, "\" alt=\"QR Code for store\" class=\"w-48 h-48 rounded-lg\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 24, "<!-- Error message if QR generation failed --> <div class=\"text-center text-gray-500\"><svg class=\"w-12 h-12 mx-auto mb-2\" fill=\"currentColor\" viewBox=\"0 0 24 24\"><path d=\"M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z\"></path></svg><p class=\"text-sm\">Unable to generate QR code</p></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 25, "</div></div><!-- Store URL --><div class=\"bg-gray-50 rounded-xl p-4 mb-6\"><div class=\"flex items-center justify-center space-x-2\"><svg class=\"w-4 h-4 text-gulf-teal\" fill=\"currentColor\" viewBox=\"0 0 24 24\"><path d=\"M3.9 12c0-1.71 1.39-3.1 3.1-3.1h4V7H7c-2.76 0-5 2.24-5 5s2.24 5 5 5h4v-1.9H7c-1.71 0-3.1-1.39-3.1-3.1zM8 13h8v-2H8v2zm9-6h-4v1.9h4c1.71 0 3.1 1.39 3.1 3.1s-1.39 3.1-3.1 3.1h-4V17h4c2.76 0 5-2.24 5-5s-2.24-5-5-5z\"></path></svg> <span class=\"text-sm font-medium text-slate-charcoal\">waqti.me/")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var9 string
		templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(creator.Username)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/qr_modal.templ`, Line: 137, Col: 55}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 26, "</span> <button onclick=\"copyToClipboard(&#39;waqti.me/{ creator.Username }&#39;)\" class=\"p-1 hover:bg-gray-200 rounded transition-colors\" title=\"Copy link\"><svg class=\"w-4 h-4 text-gray-400\" fill=\"currentColor\" viewBox=\"0 0 24 24\"><path d=\"M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z\"></path></svg></button></div></div><!-- Download Button --><div class=\"text-center\"><button onclick=\"downloadQRCode()\" class=\"bg-gulf-teal text-white py-2 px-6 rounded-xl font-medium hover:bg-teal-600 transition-colors flex items-center justify-center space-x-2 mx-auto\" id=\"download-btn\" data-download-text=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var10 string
		templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(getDownloadText(lang))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/qr_modal.templ`, Line: 157, Col: 66}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 27, "\" data-downloaded-text=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var11 string
		templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(getDownloadedText(lang))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/qr_modal.templ`, Line: 158, Col: 70}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 28, "\"><svg class=\"w-4 h-4\" fill=\"currentColor\" viewBox=\"0 0 24 24\"><path d=\"M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z\"></path></svg> <span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 29, "تحميل رمز QR")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 30, "Download QR Code")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 31, "</span></button></div></div></div></div><style>\n        .qr-modal-backdrop {\n            animation: fadeIn 0.3s ease-out;\n        }\n\n        .qr-modal-content {\n            animation: slideUp 0.3s ease-out;\n        }\n\n        .card-shadow {\n            box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);\n        }\n\n        @keyframes fadeIn {\n            from { opacity: 0; }\n            to { opacity: 1; }\n        }\n\n        @keyframes slideUp {\n            from {\n                opacity: 0;\n                transform: scale(0.95) translateY(20px);\n            }\n            to {\n                opacity: 1;\n                transform: scale(1) translateY(0);\n            }\n        }\n    </style><script>\n        // Copy to clipboard function\n        function copyToClipboard(text) {\n            navigator.clipboard.writeText(text).then(function() {\n                // Show success feedback\n                const button = event.target.closest('button');\n                const originalTitle = button.title;\n                button.title = 'Copied!';\n                button.style.backgroundColor = '#2DD4BF';\n                button.style.color = 'white';\n\n                setTimeout(() => {\n                    button.title = originalTitle;\n                    button.style.backgroundColor = '';\n                    button.style.color = '';\n                }, 2000);\n            });\n        }\n\n        // Download QR Code function\n        function downloadQRCode() {\n            // Use local download API endpoint\n            const downloadUrl = '/api/qr/download?size=300';\n            \n            // Create a temporary link element and trigger download\n            const link = document.createElement('a');\n            link.href = downloadUrl;\n            link.download = 'qr-code-{ creator.Username }.png';\n            document.body.appendChild(link);\n            link.click();\n            document.body.removeChild(link);\n\n            // Show success feedback\n            const button = document.getElementById('download-btn');\n            const originalContent = button.innerHTML;\n            const successIcon = `\n                <svg class=\"w-4 h-4\" fill=\"currentColor\" viewBox=\"0 0 24 24\">\n                    <path d=\"M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z\"/>\n                </svg>\n            `;\n            \n            // Get localized text from data attributes\n            const downloadedText = button.getAttribute('data-downloaded-text') || 'Downloaded!';\n            button.innerHTML = successIcon + `<span>${downloadedText}</span>`;\n\n            setTimeout(() => {\n                button.innerHTML = originalContent;\n            }, 2000);\n        }\n\n    </script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func QRCodePattern() templ.Component {
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
		templ_7745c5c3_Var12 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var12 == nil {
			templ_7745c5c3_Var12 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 32, "<!-- Stylized QR Code Pattern --><svg width=\"192\" height=\"192\" viewBox=\"0 0 192 192\" class=\"qr-pattern\"><!-- Corner squares --><rect x=\"0\" y=\"0\" width=\"56\" height=\"56\" fill=\"#1E293B\" rx=\"8\"></rect> <rect x=\"8\" y=\"8\" width=\"40\" height=\"40\" fill=\"white\" rx=\"6\"></rect> <rect x=\"16\" y=\"16\" width=\"24\" height=\"24\" fill=\"#2DD4BF\" rx=\"4\"></rect> <rect x=\"136\" y=\"0\" width=\"56\" height=\"56\" fill=\"#1E293B\" rx=\"8\"></rect> <rect x=\"144\" y=\"8\" width=\"40\" height=\"40\" fill=\"white\" rx=\"6\"></rect> <rect x=\"152\" y=\"16\" width=\"24\" height=\"24\" fill=\"#2DD4BF\" rx=\"4\"></rect> <rect x=\"0\" y=\"136\" width=\"56\" height=\"56\" fill=\"#1E293B\" rx=\"8\"></rect> <rect x=\"8\" y=\"144\" width=\"40\" height=\"40\" fill=\"white\" rx=\"6\"></rect> <rect x=\"16\" y=\"152\" width=\"24\" height=\"24\" fill=\"#2DD4BF\" rx=\"4\"></rect><!-- Data pattern (stylized) --><g fill=\"#1E293B\"><!-- Row 1 --><rect x=\"72\" y=\"8\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"88\" y=\"8\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"104\" y=\"8\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"120\" y=\"8\" width=\"8\" height=\"8\" rx=\"2\"></rect><!-- Row 2 --><rect x=\"64\" y=\"24\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"80\" y=\"24\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"112\" y=\"24\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"128\" y=\"24\" width=\"8\" height=\"8\" rx=\"2\"></rect><!-- Center pattern --><rect x=\"72\" y=\"72\" width=\"48\" height=\"48\" fill=\"#1E293B\" rx=\"6\"></rect> <rect x=\"80\" y=\"80\" width=\"32\" height=\"32\" fill=\"white\" rx=\"4\"></rect> <rect x=\"88\" y=\"88\" width=\"16\" height=\"16\" fill=\"#2DD4BF\" rx=\"2\"></rect><!-- More data dots --><rect x=\"8\" y=\"72\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"24\" y=\"72\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"40\" y=\"72\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"8\" y=\"88\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"40\" y=\"88\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"8\" y=\"104\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"24\" y=\"104\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"40\" y=\"104\" width=\"8\" height=\"8\" rx=\"2\"></rect><!-- Right side data --><rect x=\"144\" y=\"72\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"160\" y=\"72\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"176\" y=\"72\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"144\" y=\"88\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"176\" y=\"88\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"144\" y=\"104\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"160\" y=\"104\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"176\" y=\"104\" width=\"8\" height=\"8\" rx=\"2\"></rect><!-- Bottom data --><rect x=\"72\" y=\"144\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"88\" y=\"144\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"104\" y=\"144\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"120\" y=\"144\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"64\" y=\"160\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"80\" y=\"160\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"112\" y=\"160\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"128\" y=\"160\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"72\" y=\"176\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"88\" y=\"176\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"104\" y=\"176\" width=\"8\" height=\"8\" rx=\"2\"></rect> <rect x=\"120\" y=\"176\" width=\"8\" height=\"8\" rx=\"2\"></rect></g></svg>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

// Helper functions

func getCopySuccessText(lang string) string {
	if lang == "ar" {
		return "تم النسخ!"
	}
	return "Copied!"
}

func getDownloadText(lang string) string {
	if lang == "ar" {
		return "تحميل رمز QR"
	}
	return "Download QR Code"
}

func getDownloadedText(lang string) string {
	if lang == "ar" {
		return "تم التحميل!"
	}
	return "Downloaded!"
}

var _ = templruntime.GeneratedTemplate
