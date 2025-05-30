// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.865
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"waqti/internal/models"
)

func ReorderWorkshopsPage(creator *models.Creator, workshops []models.Workshop, stats models.DashboardStats, lang string, isRTL bool) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(lang)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 10, Col: 21}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\" dir=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(getDirection(isRTL))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 10, Col: 49}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "إعادة ترتيب الورش - Waqti.me")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "Reorder Workshops - Waqti.me")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</title><script src=\"https://cdn.tailwindcss.com\"></script><script src=\"https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js\" defer></script><script src=\"https://unpkg.com/htmx.org@1.9.10\"></script><link rel=\"preconnect\" href=\"https://fonts.googleapis.com\"><link rel=\"preconnect\" href=\"https://fonts.gstatic.com\" crossorigin><link rel=\"stylesheet\" href=\"https://fonts.googleapis.com/css2?family=Cairo:wght@400;500;600;700&amp;display=swap\"><link rel=\"stylesheet\" href=\"https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&amp;display=swap\"><script>\n            tailwind.config = {\n                theme: {\n                    extend: {\n                        colors: {\n                            'gulf-teal': '#2DD4BF',\n                            'ivory-sand': '#FEFCE8',\n                            'slate-charcoal': '#1E293B'\n                        },\n                        fontFamily: {\n                            'cairo': ['Cairo', 'sans-serif'],\n                            'inter': ['Inter', 'sans-serif']\n                        }\n                    }\n                }\n            }\n        </script><style>\n            .font-primary {\n                font-family: { getFontFamily(isRTL) };\n            }\n\n            .gradient-bg {\n                background: linear-gradient(135deg, #F0FDFA 0%, #FEFCE8 100%);\n            }\n\n            .card-shadow {\n                box-shadow: 0 4px 20px rgba(45, 212, 191, 0.1);\n            }\n\n            .badge {\n                background: linear-gradient(45deg, #2DD4BF, #06B6D4);\n            }\n\n            .workshop-row:hover {\n                background-color: rgba(45, 212, 191, 0.05);\n            }\n\n            .reorder-btn:hover {\n                background-color: rgba(45, 212, 191, 0.1);\n            }\n        </style></head><body class=\"gradient-bg min-h-screen font-primary\"><!-- Header with Back Button --><header class=\"bg-white/80 backdrop-blur-sm border-b border-gulf-teal/20 sticky top-0 z-50\"><div class=\"max-w-md mx-auto px-4 py-4\"><div class=\"flex items-center space-x-3\"><a href=\"/dashboard\" class=\"p-2 hover:bg-gray-100 rounded-lg transition-colors\"><svg class=\"w-5 h-5 text-slate-charcoal\" fill=\"currentColor\" viewBox=\"0 0 24 24\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if isRTL {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "<path d=\"M8.59 16.34l4.58-4.59-4.58-4.59L10 5.75l6 6-6 6z\"></path>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "<path d=\"M15.41 16.34L10.83 11.75l4.58-4.59L14 5.75l-6 6 6 6z\"></path>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "</svg></a><h1 class=\"text-lg font-bold text-slate-charcoal\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "إعادة ترتيب الورش")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "Reorder Workshops")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "</h1></div></div></header><!-- Main Content --><main class=\"max-w-md mx-auto px-4 py-6 space-y-6\"><!-- Workshop Overview --><div class=\"bg-white rounded-2xl p-6 card-shadow\"><h2 class=\"text-lg font-bold text-slate-charcoal mb-4\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "نظرة عامة على الورش")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, "Workshop Overview")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, "</h2><div class=\"grid grid-cols-3 gap-4\"><!-- Total Workshops --><div class=\"text-center\"><div class=\"text-2xl font-bold text-gulf-teal\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", stats.TotalWorkshops))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 115, Col: 112}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 16, "</div><div class=\"text-xs text-gray-500\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 17, "إجمالي الورش")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 18, "Total Workshops")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 19, "</div></div><!-- Projected Sales --><div class=\"text-center\"><div class=\"text-2xl font-bold text-green-600\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%.0f", stats.ProjectedSales))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 127, Col: 114}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 20, "</div><div class=\"text-xs text-gray-500\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 21, "المبيعات المتوقعة")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 22, "Projected Sales")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 23, "</div></div><!-- Remaining Seats --><div class=\"text-center\"><div class=\"text-2xl font-bold text-orange-600\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", stats.RemainingSeats))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 139, Col: 113}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 24, "</div><div class=\"text-xs text-gray-500\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 25, "المقاعد المتبقية")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 26, "Remaining Seats")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 27, "</div></div></div></div><!-- Add New Workshop Button - UPDATED WITH PROPER LINK --><a href=\"/workshops/add\" class=\"w-full bg-gulf-teal text-white py-3 rounded-xl font-medium hover:bg-teal-600 transition-colors flex items-center justify-center space-x-2 block\"><svg class=\"w-5 h-5\" fill=\"currentColor\" viewBox=\"0 0 24 24\"><path d=\"M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z\"></path></svg> <span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 28, "إضافة ورشة جديدة")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 29, "Add New Workshop")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 30, "</span></a><!-- Workshops List --><div class=\"bg-white rounded-2xl card-shadow overflow-hidden\"><div class=\"p-4 border-b border-gray-100\"><h3 class=\"font-semibold text-slate-charcoal\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 31, "ورشاتي")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 32, "My Workshops")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 33, "</h3></div><div id=\"workshops-list\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = WorkshopsList(workshops, lang, isRTL).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 34, "</div></div><!-- Bottom Spacing --><div class=\"h-6\"></div></main></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func WorkshopsList(workshops []models.Workshop, lang string, isRTL bool) templ.Component {
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
		templ_7745c5c3_Var7 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var7 == nil {
			templ_7745c5c3_Var7 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		for i, workshop := range workshops {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 35, "<div class=\"workshop-row border-b border-gray-50 last:border-b-0\"><div class=\"p-4 flex items-center justify-between\"><!-- Workshop Info --><div class=\"flex-1\"><div class=\"flex items-center space-x-2 mb-1\"><h4 class=\"font-medium text-slate-charcoal text-sm\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if lang == "ar" {
				var templ_7745c5c3_Var8 string
				templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(workshop.TitleAr)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 198, Col: 50}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				var templ_7745c5c3_Var9 string
				templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(workshop.Title)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 200, Col: 48}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 36, "</h4>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !workshop.IsActive {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 37, "<span class=\"px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if lang == "ar" {
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 38, "غير نشط")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 39, "Inactive")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 40, "</span>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 41, "</div><div class=\"flex items-center space-x-4 text-xs text-gray-500\"><span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var10 string
			templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%.0f KD", workshop.Price))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 214, Col: 70}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 42, "</span> <span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var11 string
			templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", workshop.MaxStudents))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 215, Col: 71}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 43, " ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if lang == "ar" {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 44, "مقعد")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 45, "seats")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 46, "</span></div></div><!-- Actions --><div class=\"flex items-center space-x-2\"><!-- Reorder Buttons --><div class=\"flex flex-col space-y-1\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if i > 0 {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 47, "<button class=\"reorder-btn p-1 rounded hover:bg-gray-100 transition-colors\" hx-post=\"/workshops/reorder\" hx-values=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var12 string
				templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf(`{"workshop_id": "%d", "direction": "up"}`, workshop.ID))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 233, Col: 112}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 48, "\" hx-target=\"#workshops-list\" hx-swap=\"innerHTML\"><svg class=\"w-4 h-4 text-gray-600\" fill=\"currentColor\" viewBox=\"0 0 24 24\"><path d=\"M7 14l5-5 5 5z\"></path></svg></button> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if i < len(workshops)-1 {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 49, "<button class=\"reorder-btn p-1 rounded hover:bg-gray-100 transition-colors\" hx-post=\"/workshops/reorder\" hx-values=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var13 string
				templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf(`{"workshop_id": "%d", "direction": "down"}`, workshop.ID))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 246, Col: 114}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 50, "\" hx-target=\"#workshops-list\" hx-swap=\"innerHTML\"><svg class=\"w-4 h-4 text-gray-600\" fill=\"currentColor\" viewBox=\"0 0 24 24\"><path d=\"M7 10l5 5 5-5z\"></path></svg></button>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 51, "</div><!-- Status Toggle -->")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var14 = []any{fmt.Sprintf("p-2 rounded-lg transition-colors %s", getStatusButtonClass(workshop.IsActive))}
			templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var14...)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 52, "<button class=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var15 string
			templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var14).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 1, Col: 0}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 53, "\" hx-post=\"/workshops/toggle-status\" hx-values=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var16 string
			templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf(`{"workshop_id": "%d"}`, workshop.ID))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/reorder_workshops.templ`, Line: 261, Col: 85}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 54, "\" hx-target=\"#workshops-list\" hx-swap=\"innerHTML\"><svg class=\"w-4 h-4\" fill=\"currentColor\" viewBox=\"0 0 24 24\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if workshop.IsActive {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 55, "<path d=\"M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z\"></path>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 56, "<path d=\"M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z\"></path>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 57, "</svg></button><!-- Edit Button --><button class=\"p-2 text-gulf-teal hover:bg-gulf-teal/10 rounded-lg transition-colors\"><svg class=\"w-4 h-4\" fill=\"currentColor\" viewBox=\"0 0 24 24\"><path d=\"M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z\"></path></svg></button></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		return nil
	})
}

// Helper function for status button styling
func getStatusButtonClass(isActive bool) string {
	if isActive {
		return "text-green-600 hover:bg-green-50"
	}
	return "text-gray-400 hover:bg-gray-50"
}

var _ = templruntime.GeneratedTemplate
