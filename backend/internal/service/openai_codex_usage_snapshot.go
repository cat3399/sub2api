package service

// buildCodexUsageExtraUpdates converts an upstream Codex usage snapshot (x-codex-* headers)
// into account.extra fields.
//
// Note: This intentionally mirrors the logic in OpenAIGatewayService.updateCodexUsageSnapshot
// so that other code paths (e.g. admin-triggered refresh) can reuse the same mapping.
func buildCodexUsageExtraUpdates(snapshot *OpenAICodexUsageSnapshot) map[string]any {
	if snapshot == nil {
		return nil
	}

	updates := make(map[string]any)
	if snapshot.PrimaryUsedPercent != nil {
		updates["codex_primary_used_percent"] = *snapshot.PrimaryUsedPercent
	}
	if snapshot.PrimaryResetAfterSeconds != nil {
		updates["codex_primary_reset_after_seconds"] = *snapshot.PrimaryResetAfterSeconds
	}
	if snapshot.PrimaryWindowMinutes != nil {
		updates["codex_primary_window_minutes"] = *snapshot.PrimaryWindowMinutes
	}
	if snapshot.SecondaryUsedPercent != nil {
		updates["codex_secondary_used_percent"] = *snapshot.SecondaryUsedPercent
	}
	if snapshot.SecondaryResetAfterSeconds != nil {
		updates["codex_secondary_reset_after_seconds"] = *snapshot.SecondaryResetAfterSeconds
	}
	if snapshot.SecondaryWindowMinutes != nil {
		updates["codex_secondary_window_minutes"] = *snapshot.SecondaryWindowMinutes
	}
	if snapshot.PrimaryOverSecondaryPercent != nil {
		updates["codex_primary_over_secondary_percent"] = *snapshot.PrimaryOverSecondaryPercent
	}
	if snapshot.UpdatedAt != "" {
		updates["codex_usage_updated_at"] = snapshot.UpdatedAt
	}

	// Normalize to canonical 5h/7d fields based on window_minutes.
	//
	// IMPORTANT: We can only reliably determine window type from window_minutes field.
	// The reset_after_seconds is remaining time, not window size.
	var primaryWindowMins, secondaryWindowMins int
	var hasPrimaryWindow, hasSecondaryWindow bool

	if snapshot.PrimaryWindowMinutes != nil {
		primaryWindowMins = *snapshot.PrimaryWindowMinutes
		hasPrimaryWindow = true
	}
	if snapshot.SecondaryWindowMinutes != nil {
		secondaryWindowMins = *snapshot.SecondaryWindowMinutes
		hasSecondaryWindow = true
	}

	var use5hFromPrimary, use7dFromPrimary bool
	var use5hFromSecondary, use7dFromSecondary bool

	if hasPrimaryWindow && hasSecondaryWindow {
		// Both window sizes known: compare and assign smaller to 5h, larger to 7d.
		if primaryWindowMins < secondaryWindowMins {
			use5hFromPrimary = true
			use7dFromSecondary = true
		} else {
			use5hFromSecondary = true
			use7dFromPrimary = true
		}
	} else if hasPrimaryWindow {
		// Only primary window size known: classify by absolute threshold.
		if primaryWindowMins <= 360 {
			use5hFromPrimary = true
		} else {
			use7dFromPrimary = true
		}
	} else if hasSecondaryWindow {
		// Only secondary window size known: classify by absolute threshold.
		if secondaryWindowMins <= 360 {
			use5hFromSecondary = true
		} else {
			use7dFromSecondary = true
		}
	} else {
		// No window_minutes available: cannot reliably determine window types.
		// Fall back to legacy assumption (may be incorrect): primary=7d, secondary=5h.
		if snapshot.SecondaryUsedPercent != nil || snapshot.SecondaryResetAfterSeconds != nil || snapshot.SecondaryWindowMinutes != nil {
			use5hFromSecondary = true
		}
		if snapshot.PrimaryUsedPercent != nil || snapshot.PrimaryResetAfterSeconds != nil || snapshot.PrimaryWindowMinutes != nil {
			use7dFromPrimary = true
		}
	}

	// Write canonical 5h fields.
	if use5hFromPrimary {
		if snapshot.PrimaryUsedPercent != nil {
			updates["codex_5h_used_percent"] = *snapshot.PrimaryUsedPercent
		}
		if snapshot.PrimaryResetAfterSeconds != nil {
			updates["codex_5h_reset_after_seconds"] = *snapshot.PrimaryResetAfterSeconds
		}
		if snapshot.PrimaryWindowMinutes != nil {
			updates["codex_5h_window_minutes"] = *snapshot.PrimaryWindowMinutes
		}
	} else if use5hFromSecondary {
		if snapshot.SecondaryUsedPercent != nil {
			updates["codex_5h_used_percent"] = *snapshot.SecondaryUsedPercent
		}
		if snapshot.SecondaryResetAfterSeconds != nil {
			updates["codex_5h_reset_after_seconds"] = *snapshot.SecondaryResetAfterSeconds
		}
		if snapshot.SecondaryWindowMinutes != nil {
			updates["codex_5h_window_minutes"] = *snapshot.SecondaryWindowMinutes
		}
	}

	// Write canonical 7d fields.
	if use7dFromPrimary {
		if snapshot.PrimaryUsedPercent != nil {
			updates["codex_7d_used_percent"] = *snapshot.PrimaryUsedPercent
		}
		if snapshot.PrimaryResetAfterSeconds != nil {
			updates["codex_7d_reset_after_seconds"] = *snapshot.PrimaryResetAfterSeconds
		}
		if snapshot.PrimaryWindowMinutes != nil {
			updates["codex_7d_window_minutes"] = *snapshot.PrimaryWindowMinutes
		}
	} else if use7dFromSecondary {
		if snapshot.SecondaryUsedPercent != nil {
			updates["codex_7d_used_percent"] = *snapshot.SecondaryUsedPercent
		}
		if snapshot.SecondaryResetAfterSeconds != nil {
			updates["codex_7d_reset_after_seconds"] = *snapshot.SecondaryResetAfterSeconds
		}
		if snapshot.SecondaryWindowMinutes != nil {
			updates["codex_7d_window_minutes"] = *snapshot.SecondaryWindowMinutes
		}
	}

	return updates
}

