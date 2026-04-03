import { useCallback } from "react";

// Analytics event types
export type AnalyticsEvent =
  | "result_viewed"
  | "result_exported"
  | "result_printed"
  | "violation_expanded"
  | "violation_collapsed"
  | "help_tooltip_viewed"
  | "tab_switched";

interface AnalyticsEventData {
  event: AnalyticsEvent;
  properties?: Record<string, string | number | boolean>;
  timestamp?: string;
}

/**
 * Custom hook for tracking user interactions and analytics events
 * This provides a centralized way to track user behavior across the application
 */
export function useAnalytics() {
  const trackEvent = useCallback((event: AnalyticsEvent, properties?: Record<string, string | number | boolean>) => {
    const eventData: AnalyticsEventData = {
      event,
      properties,
      timestamp: new Date().toISOString(),
    };

    // Log to console in development
    if (process.env.NODE_ENV === "development") {
      console.log("[Analytics]", eventData);
    }

    // In production, this would send to your analytics service
    // Example integrations:
    // - Google Analytics: gtag('event', event, properties)
    // - Mixpanel: mixpanel.track(event, properties)
    // - Segment: analytics.track(event, properties)
    // - Custom API: fetch('/api/analytics', { method: 'POST', body: JSON.stringify(eventData) })

    try {
      // Store in localStorage for debugging/testing
      const events = JSON.parse(localStorage.getItem("analytics_events") || "[]");
      events.push(eventData);
      // Keep only last 100 events
      if (events.length > 100) {
        events.shift();
      }
      localStorage.setItem("analytics_events", JSON.stringify(events));
    } catch (error) {
      console.error("Failed to store analytics event:", error);
    }
  }, []);

  return { trackEvent };
}
