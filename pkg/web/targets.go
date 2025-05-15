package web

// Target is a struct representing a DOM target.
type Target struct {
	ID  string `json:"id"`
	Sel string `json:"sel"`
}

var (
	// TargetLogContainer is a target for the log container.
	// It is used to insert log entries into the DOM.
	TargetLogContainer = Target{
		ID:  "log-container",
		Sel: "#log-container",
	}

	// TargetStatusContent is a target for the status content.
	// It is used to update the status content in the DOM.
	TargetStatusContent = Target{
		ID:  "status-content",
		Sel: "#status-content",
	}
)

var (
	// LivePageTitle is the title of the live page.
	LivePageTitle = "Live Camera System"
)
