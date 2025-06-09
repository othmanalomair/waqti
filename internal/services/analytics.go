package services

import (
	"fmt"
	"waqti/internal/database"
	"waqti/internal/models"

	"github.com/google/uuid"
)

type AnalyticsService struct{}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

func (s *AnalyticsService) GetClicksByCreatorID(creatorID uuid.UUID, filter models.AnalyticsFilter) []models.AnalyticsClick {
	// Get clicks from database
	dbClicks, err := database.Instance.GetAnalyticsClicksByCreatorID(creatorID, 100)
	if err != nil {
		// Return empty slice on error
		return []models.AnalyticsClick{}
	}

	// Convert database clicks to models
	var clicks []models.AnalyticsClick
	for _, dbClick := range dbClicks {
		click := models.AnalyticsClick{
			ID:         dbClick.ID,
			CreatorID:  dbClick.CreatorID,
			Country:    dbClick.Country,
			CountryAr:  dbClick.CountryAr,
			Device:     dbClick.Device,
			DeviceAr:   dbClick.DeviceAr,
			OS:         dbClick.OS,
			OSAr:       dbClick.OSAr,
			Platform:   dbClick.Platform,
			PlatformAr: dbClick.PlatformAr,
			ClickedAt:  dbClick.ClickedAt,
			CreatedAt:  dbClick.ClickedAt, // Use clicked_at as created_at for compatibility
		}

		// Handle optional fields
		if dbClick.Browser != nil {
			click.Browser = *dbClick.Browser
		} else {
			click.Browser = "Unknown"
		}
		if dbClick.BrowserAr != nil {
			click.BrowserAr = *dbClick.BrowserAr
		} else {
			click.BrowserAr = "غير محدد"
		}

		clicks = append(clicks, click)
	}

	return clicks
}

func (s *AnalyticsService) GetAnalyticsStats(creatorID uuid.UUID, filter models.AnalyticsFilter) models.AnalyticsStats {
	clicks := s.GetClicksByCreatorID(creatorID, filter)

	// Get total clicks from database
	totalClicks, err := database.Instance.GetAnalyticsStats(creatorID)
	if err != nil {
		totalClicks = len(clicks)
	}

	stats := models.AnalyticsStats{
		TotalClicks:       totalClicks,
		DateRange:         filter.DateRange,
		CountryBreakdown:  make(map[string]int),
		DeviceBreakdown:   make(map[string]int),
		OSBreakdown:       make(map[string]int),
		PlatformBreakdown: make(map[string]int),
	}

	for _, click := range clicks {
		stats.CountryBreakdown[click.Country]++
		stats.DeviceBreakdown[click.Device]++
		stats.OSBreakdown[click.OS]++
		stats.PlatformBreakdown[click.Platform]++
	}

	return stats
}

func (s *AnalyticsService) GenerateCompleteScript(clicks []models.AnalyticsClick, stats models.AnalyticsStats) string {
	jsCode := fmt.Sprintf(`<script type="text/javascript">
let currentFilter = 'hours';

// Analytics data from server
const totalClicksData = %d;
const realClicksData = [`, stats.TotalClicks)
	
	for i, click := range clicks {
		if i > 0 {
			jsCode += ",\n"
		}
		jsCode += fmt.Sprintf("    { time: new Date('%s'), platform: '%s' }", 
			click.ClickedAt.Format("2006-01-02T15:04:05Z07:00"), 
			click.Platform)
	}
	
	jsCode += `];

let chartData = {
    hours: generateRealHoursData(),
    days: generateRealDaysData(),
    months: generateRealMonthsData()
};

function switchTimeFilter(filter) {
    currentFilter = filter;
    
    // Update tab states
    document.querySelectorAll('[id^="tab-"]').forEach(tab => {
        if (tab.id === 'tab-' + filter) {
            tab.className = tab.className.replace('tab-inactive', 'tab-active');
        } else {
            tab.className = tab.className.replace('tab-active', 'tab-inactive');
        }
    });
    
    // Update chart
    updateChart(filter);
}

function generateRealHoursData() {
    const hours = Array.from({length: 24}, (_, i) => ({
        label: i + ':00',
        value: 0,
        height: 5
    }));
    
    // Count real clicks by hour for today
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    
    realClicksData.forEach(click => {
        if (click.time >= today) {
            const hour = click.time.getHours();
            hours[hour].value++;
        }
    });
    
    // Calculate heights based on max value (ensure minimum height for visibility)
    const maxValue = Math.max(...hours.map(h => h.value), 1);
    hours.forEach(hour => {
        hour.height = Math.max((hour.value / maxValue) * 100, 5);
    });
    
    return hours;
}

function generateRealDaysData() {
    const days = [];
    const now = new Date();
    
    for (let i = 29; i >= 0; i--) {
        const date = new Date(now);
        date.setDate(date.getDate() - i);
        date.setHours(0, 0, 0, 0);
        
        const nextDay = new Date(date);
        nextDay.setDate(nextDay.getDate() + 1);
        
        const count = realClicksData.filter(click => 
            click.time >= date && click.time < nextDay
        ).length;
        
        days.push({
            label: date.toLocaleDateString('en', { month: 'short', day: 'numeric' }),
            value: count,
            height: 5
        });
    }
    
    // Calculate heights
    const maxValue = Math.max(...days.map(d => d.value), 1);
    days.forEach(day => {
        day.height = Math.max((day.value / maxValue) * 100, 5);
    });
    
    return days;
}

function generateRealMonthsData() {
    const months = [];
    const now = new Date();
    
    for (let i = 11; i >= 0; i--) {
        const startDate = new Date(now.getFullYear(), now.getMonth() - i, 1);
        const endDate = new Date(now.getFullYear(), now.getMonth() - i + 1, 1);
        
        const count = realClicksData.filter(click => 
            click.time >= startDate && click.time < endDate
        ).length;
        
        months.push({
            label: startDate.toLocaleDateString('en', { month: 'short' }),
            value: count,
            height: 5
        });
    }
    
    // Calculate heights
    const maxValue = Math.max(...months.map(m => m.value), 1);
    months.forEach(month => {
        month.height = Math.max((month.value / maxValue) * 100, 5);
    });
    
    return months;
}

function updateChart(filter) {
    const chart = document.getElementById('traffic-chart');
    const labels = document.getElementById('chart-labels');
    const data = chartData[filter];
    
    // Clear existing content
    chart.innerHTML = '';
    labels.innerHTML = '';
    
    if (data && data.length > 0) {
        // Create bars
        data.forEach(item => {
            const bar = document.createElement('div');
            bar.className = 'chart-bar flex-1';
            bar.style.height = item.height + '%';
            bar.setAttribute('data-value', item.value);
            chart.appendChild(bar);
        });
        
        // Create labels (show every nth label to avoid crowding)
        const labelStep = filter === 'hours' ? 3 : (filter === 'days' ? 5 : 1);
        data.forEach((item, index) => {
            const label = document.createElement('div');
            label.className = 'flex-1 text-center';
            label.textContent = (index % labelStep === 0) ? item.label : '';
            labels.appendChild(label);
        });
    } else {
        // Show placeholder when no data
        chart.innerHTML = '<div class="w-full text-center text-gray-400 py-8">No data available for this time period</div>';
    }
}

// Initialize chart on load
document.addEventListener('DOMContentLoaded', function() {
    updateChart('hours');
});
</script>`
	
	return jsCode
}
