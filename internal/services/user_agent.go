package services

import (
	"strings"
)

// UserAgentInfo contains parsed user agent information
type UserAgentInfo struct {
	Device     string
	DeviceAr   string
	OS         string
	OSAr       string
	Browser    string
	BrowserAr  string
	Platform   string
	PlatformAr string
}

// ParseUserAgent parses user agent string and extracts device information
func ParseUserAgent(userAgent, referrer string) UserAgentInfo {
	ua := strings.ToLower(userAgent)
	ref := strings.ToLower(referrer)

	info := UserAgentInfo{
		Device:     "Desktop",
		DeviceAr:   "سطح المكتب",
		OS:         "Unknown",
		OSAr:       "غير محدد",
		Browser:    "Unknown",
		BrowserAr:  "غير محدد",
		Platform:   "Direct",
		PlatformAr: "مباشر",
	}

	// Detect device
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		info.Device = "Mobile"
		info.DeviceAr = "جوال"
	}

	// Detect OS
	if strings.Contains(ua, "windows") {
		info.OS = "Windows"
		info.OSAr = "ويندوز"
	} else if strings.Contains(ua, "mac") || strings.Contains(ua, "darwin") {
		info.OS = "macOS"
		info.OSAr = "ماك أو إس"
	} else if strings.Contains(ua, "linux") {
		info.OS = "Linux"
		info.OSAr = "لينوكس"
	} else if strings.Contains(ua, "android") {
		info.OS = "Android"
		info.OSAr = "أندرويد"
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ios") {
		info.OS = "iOS"
		info.OSAr = "آي أو إس"
	}

	// Detect browser
	if strings.Contains(ua, "chrome") && !strings.Contains(ua, "edg") {
		info.Browser = "Chrome"
		info.BrowserAr = "كروم"
	} else if strings.Contains(ua, "firefox") {
		info.Browser = "Firefox"
		info.BrowserAr = "فايرفوكس"
	} else if strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") {
		info.Browser = "Safari"
		info.BrowserAr = "سفاري"
	} else if strings.Contains(ua, "edg") {
		info.Browser = "Edge"
		info.BrowserAr = "إيدج"
	}

	// Detect platform/referrer
	if strings.Contains(ref, "instagram") || strings.Contains(ref, "ig.me") {
		info.Platform = "Instagram"
		info.PlatformAr = "إنستغرام"
	} else if strings.Contains(ref, "snapchat") || strings.Contains(ref, "snap.com") {
		info.Platform = "Snapchat"
		info.PlatformAr = "سناب شات"
	} else if strings.Contains(ref, "whatsapp") || strings.Contains(ref, "wa.me") {
		info.Platform = "WhatsApp"
		info.PlatformAr = "واتساب"
	} else if strings.Contains(ref, "twitter") || strings.Contains(ref, "t.co") {
		info.Platform = "Twitter"
		info.PlatformAr = "تويتر"
	} else if strings.Contains(ref, "facebook") || strings.Contains(ref, "fb.com") {
		info.Platform = "Facebook"
		info.PlatformAr = "فيسبوك"
	} else if strings.Contains(ref, "tiktok") {
		info.Platform = "TikTok"
		info.PlatformAr = "تيك توك"
	} else if strings.Contains(ref, "youtube") {
		info.Platform = "YouTube"
		info.PlatformAr = "يوتيوب"
	} else if strings.Contains(ref, "linkedin") {
		info.Platform = "LinkedIn"
		info.PlatformAr = "لينكد إن"
	} else if ref != "" {
		info.Platform = "Other"
		info.PlatformAr = "أخرى"
	}

	return info
}

// GetCountryFromIP returns country information (simplified for now)
func GetCountryFromIP(ipAddress string) (string, string) {
	// For now, return Kuwait as default
	// In production, you'd use a GeoIP service
	return "Kuwait", "الكويت"
}