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

func EnrollmentTrackingPage(enrollments []models.Enrollment, stats models.EnrollmentStats, filter models.EnrollmentFilter, settings *models.ShopSettings, lang string, isRTL bool) templ.Component {
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
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 10, Col: 21}
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
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 10, Col: 49}
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
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "تتبع التسجيلات - Waqti.me")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "Enrollment Tracking - Waqti.me")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</title><script src=\"https://cdn.tailwindcss.com\"></script><script src=\"https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js\" defer></script><script src=\"https://unpkg.com/htmx.org@1.9.10\"></script><link rel=\"preconnect\" href=\"https://fonts.googleapis.com\"><link rel=\"preconnect\" href=\"https://fonts.gstatic.com\" crossorigin><link rel=\"stylesheet\" href=\"https://fonts.googleapis.com/css2?family=Cairo:wght@400;500;600;700&amp;display=swap\"><link rel=\"stylesheet\" href=\"https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&amp;display=swap\"><script>\n            tailwind.config = {\n                theme: {\n                    extend: {\n                        colors: {\n                            'gulf-teal': '#2DD4BF',\n                            'ivory-sand': '#FEFCE8',\n                            'slate-charcoal': '#1E293B'\n                        },\n                        fontFamily: {\n                            'cairo': ['Cairo', 'sans-serif'],\n                            'inter': ['Inter', 'sans-serif']\n                        }\n                    }\n                }\n            }\n        </script><style>\n            .font-primary {\n                font-family: { getFontFamily(isRTL) };\n            }\n\n            .gradient-bg {\n                background: linear-gradient(135deg, #F0FDFA 0%, #FEFCE8 100%);\n            }\n\n            .card-shadow {\n                box-shadow: 0 4px 20px rgba(45, 212, 191, 0.1);\n            }\n\n            .enrollment-row:hover {\n                background-color: rgba(45, 212, 191, 0.05);\n            }\n\n            .filter-select {\n                background-image: url(\"data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='m6 8 4 4 4-4'/%3e%3c/svg%3e\");\n                background-position: right 0.5rem center;\n                background-repeat: no-repeat;\n                background-size: 1.5em 1.5em;\n                padding-right: 2.5rem;\n            }\n        </style></head><body class=\"gradient-bg min-h-screen font-primary\"><!-- Header with Back Button --><header class=\"bg-white/80 backdrop-blur-sm border-b border-gulf-teal/20 sticky top-0 z-50\"><div class=\"max-w-md mx-auto px-4 py-4\"><div class=\"flex items-center space-x-3\"><a href=\"/dashboard\" class=\"p-2 hover:bg-gray-100 rounded-lg transition-colors\"><svg class=\"w-5 h-5 text-slate-charcoal\" fill=\"currentColor\" viewBox=\"0 0 24 24\">")
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
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "تتبع التسجيلات")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "Enrollment Tracking")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "</h1></div></div></header><!-- Main Content --><main class=\"max-w-md mx-auto px-4 py-6 space-y-6\"><div id=\"enrollment-content\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = EnrollmentContent(enrollments, stats, filter, settings, lang, isRTL).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "</div></main></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func EnrollmentContent(enrollments []models.Enrollment, stats models.EnrollmentStats, filter models.EnrollmentFilter, settings *models.ShopSettings, lang string, isRTL bool) templ.Component {
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
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, "<!-- Stats Overview --><div class=\"bg-white rounded-2xl p-6 card-shadow\"><h2 class=\"text-lg font-bold text-slate-charcoal mb-4\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, "إحصائيات التسجيلات")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 16, "Enrollment Statistics")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 17, "</h2><div class=\"grid grid-cols-3 gap-4\"><!-- Successful Sales --><div class=\"text-center\"><div class=\"text-2xl font-bold text-green-600\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", stats.SuccessfulSales))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 123, Col: 105}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 18, "</div><div class=\"text-xs text-gray-500\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 19, "مبيعات ناجحة")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 20, "Successful Sales")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 21, "</div></div><!-- Total Sales --><div class=\"text-center\"><div class=\"text-2xl font-bold text-gulf-teal\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%.0f", stats.TotalSales))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 135, Col: 102}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 22, "</div><div class=\"text-xs text-gray-500\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 23, "إجمالي المبيعات")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 24, "Total Sales")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 25, "</div></div><!-- Rejected Sales --><div class=\"text-center\"><div class=\"text-2xl font-bold text-red-600\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 string
		templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", stats.RejectedSales))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 147, Col: 101}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 26, "</div><div class=\"text-xs text-gray-500\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 27, "مبيعات مرفوضة")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 28, "Rejected Sales")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 29, "</div></div></div></div><!-- Filters --><div class=\"bg-white rounded-2xl p-4 card-shadow\"><form hx-post=\"/enrollments/filter\" hx-target=\"#enrollment-content\" hx-swap=\"innerHTML\" hx-trigger=\"change\" class=\"space-y-4\"><!-- Time Range Filter --><div><label class=\"block text-sm font-medium text-gray-700 mb-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 30, "الفترة الزمنية")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 31, "Time Range")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 32, "</label> <select name=\"time_range\" class=\"w-full p-2 border border-gray-300 rounded-lg filter-select appearance-none bg-white\"><option value=\"days\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if filter.TimeRange == "days" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 33, " selected")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 34, ">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 35, "آخر 30 يوم")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 36, "Last 30 Days")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 37, "</option> <option value=\"months\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if filter.TimeRange == "months" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 38, " selected")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 39, ">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 40, "آخر 12 شهر")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 41, "Last 12 Months")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 42, "</option> <option value=\"year\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if filter.TimeRange == "year" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 43, " selected")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 44, ">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 45, "آخر 5 سنوات")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 46, "Last 5 Years")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 47, "</option></select></div><!-- Sort Options --><div class=\"flex space-x-4\"><div class=\"flex-1\"><label class=\"block text-sm font-medium text-gray-700 mb-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 48, "ترتيب حسب")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 49, "Sort By")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 50, "</label> <select name=\"order_by\" class=\"w-full p-2 border border-gray-300 rounded-lg filter-select appearance-none bg-white\"><option value=\"date\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if filter.OrderBy == "date" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 51, " selected")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 52, ">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 53, "التاريخ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 54, "Date")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 55, "</option> <option value=\"price\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if filter.OrderBy == "price" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 56, " selected")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 57, ">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 58, "السعر")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 59, "Price")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 60, "</option> <option value=\"name\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if filter.OrderBy == "name" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 61, " selected")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 62, ">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 63, "الاسم")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 64, "Name")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 65, "</option></select></div><div class=\"flex-1\"><label class=\"block text-sm font-medium text-gray-700 mb-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 66, "الاتجاه")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 67, "Direction")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 68, "</label> <select name=\"order_dir\" class=\"w-full p-2 border border-gray-300 rounded-lg filter-select appearance-none bg-white\"><option value=\"asc\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if filter.OrderDir == "asc" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 69, " selected")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 70, ">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 71, "تصاعدي")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 72, "Ascending")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 73, "</option> <option value=\"desc\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if filter.OrderDir == "desc" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 74, " selected")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 75, ">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 76, "تنازلي")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 77, "Descending")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 78, "</option></select></div></div></form></div><!-- Enrollments Table --><div class=\"bg-white rounded-2xl card-shadow overflow-hidden\"><div class=\"p-4 border-b border-gray-100\"><h3 class=\"font-semibold text-slate-charcoal\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 79, "التسجيلات")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 80, "Enrollments")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 81, "</h3></div><div class=\"divide-y divide-gray-50\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(enrollments) == 0 {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 82, "<div class=\"p-8 text-center text-gray-500\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if lang == "ar" {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 83, "لا توجد تسجيلات")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 84, "No enrollments found")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 85, "</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			for _, enrollment := range enrollments {
				templ_7745c5c3_Err = EnrollmentRow(enrollment, filter, settings, lang, isRTL).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 86, "</div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func EnrollmentRow(enrollment models.Enrollment, filter models.EnrollmentFilter, settings *models.ShopSettings, lang string, isRTL bool) templ.Component {
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
		templ_7745c5c3_Var8 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var8 == nil {
			templ_7745c5c3_Var8 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 87, "<div class=\"enrollment-row p-4\"><div class=\"flex items-center justify-between\"><!-- Enrollment Info --><div class=\"flex-1\"><div class=\"flex items-center space-x-2 mb-1\"><h4 class=\"font-medium text-slate-charcoal text-sm\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			var templ_7745c5c3_Var9 string
			templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(enrollment.WorkshopNameAr)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 304, Col: 55}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			var templ_7745c5c3_Var10 string
			templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(enrollment.WorkshopName)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 306, Col: 53}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 88, "</h4>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var11 = []any{fmt.Sprintf("px-2 py-1 text-xs rounded %s", getStatusBadgeClass(enrollment.Status))}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var11...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 89, "<span class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var12 string
		templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var11).String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 1, Col: 0}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 90, "\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if lang == "ar" {
			var templ_7745c5c3_Var13 string
			templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(enrollment.StatusAr)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 311, Col: 49}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			var templ_7745c5c3_Var14 string
			templ_7745c5c3_Var14, templ_7745c5c3_Err = templ.JoinStringErrs(enrollment.Status)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 313, Col: 47}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var14))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 91, "</span></div><div class=\"flex items-center space-x-4 text-xs text-gray-500 mb-1\"><span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var15 string
		templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(enrollment.StudentName)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 318, Col: 50}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 92, "</span> <span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var16 string
		templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs(enrollment.EnrollmentDate.Format("2006/01/02"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 319, Col: 74}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 93, "</span></div><div class=\"text-sm font-semibold text-gulf-teal\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var17 string
		templ_7745c5c3_Var17, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%.2f %s", enrollment.TotalPrice, getCurrencySymbol(settings, lang)))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 322, Col: 102}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var17))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 94, "</div></div><!-- Actions --><div class=\"flex items-center space-x-2\"><!-- Edit Button --><button class=\"p-2 text-gulf-teal hover:bg-gulf-teal/10 rounded-lg transition-colors\"><svg class=\"w-4 h-4\" fill=\"currentColor\" viewBox=\"0 0 24 24\"><path d=\"M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z\"></path></svg></button><!-- Delete Button --><button class=\"p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors\" hx-post=\"/enrollments/delete\" hx-values=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var18 string
		templ_7745c5c3_Var18, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf(`{"enrollment_id": "%s", "current_time_range": "%s", "current_order_by": "%s", "current_order_dir": "%s"}`, enrollment.ID.String(), filter.TimeRange, filter.OrderBy, filter.OrderDir))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 339, Col: 226}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var18))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 95, "\" hx-target=\"#enrollment-content\" hx-swap=\"innerHTML\" hx-confirm=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var19 string
		templ_7745c5c3_Var19, templ_7745c5c3_Err = templ.JoinStringErrs(getDeleteConfirmText(lang))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/templates/enrollment_tracking.templ`, Line: 342, Col: 59}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var19))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 96, "\"><svg class=\"w-4 h-4\" fill=\"currentColor\" viewBox=\"0 0 24 24\"><path d=\"M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z\"></path></svg></button></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

// Helper functions
func getStatusBadgeClass(status string) string {
	switch status {
	case "successful":
		return "bg-green-100 text-green-800"
	case "rejected":
		return "bg-red-100 text-red-800"
	case "pending":
		return "bg-yellow-100 text-yellow-800"
	default:
		return "bg-gray-100 text-gray-800"
	}
}

func getDeleteConfirmText(lang string) string {
	if lang == "ar" {
		return "هل أنت متأكد من حذف هذا التسجيل؟"
	}
	return "Are you sure you want to delete this enrollment?"
}

var _ = templruntime.GeneratedTemplate
