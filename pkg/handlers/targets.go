package handlers

// Target is a struct representing a dom target.
type Target struct {
	ID  string `json:"id"`
	Sel string `json:"sel"`
}

var (
	// LogContainer is a target for the log container.
	// It is used to insert log entries into the DOM.
	TargetLogContainer = Target{
		ID:  "log-container",
		Sel: "#log-container",
	}

	// StatusContent is a target for the status content.
	// It is used to update the status content in the DOM.
	TargetStatusContent = Target{
		ID:  "status-content",
		Sel: "#status-content",
	}
)

var (
	ManualPageTitle = "Manual Depth Map Generator"
	LivePageTitle   = "Live Camera System"
)
