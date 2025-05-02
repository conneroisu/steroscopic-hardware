
let polling = true;
let searchTerm = '';

// Toggle polling state
function togglePolling() {
	polling = !polling;
	const btn = document.getElementById('pause-btn');
	const status = document.getElementById('connection-status');

	if (polling) {
		btn.textContent = 'Pause';
		btn.classList.add('pause');
		status.textContent = 'Connected';
		status.classList.add('connected');
		status.classList.remove('disconnected');
		// Trigger an immediate poll
		htmx.trigger('#log-container', 'newLogsLoaded');
	} else {
		btn.textContent = 'Resume';
		btn.classList.remove('pause');
		status.textContent = 'Paused';
		status.classList.remove('connected');
		status.classList.add('disconnected');
	}
}

// Clear all logs
function clearLogs() {
	const container = document.getElementById('log-container');
	// Keep the HTMX attributes
	const attrs = {
		'hx-get': container.getAttribute('hx-get'),
		'hx-trigger': container.getAttribute('hx-trigger'),
		'hx-swap': container.getAttribute('hx-swap')
	};

	// Clear content
	container.innerHTML = '<div class="empty-state"><p>No logs available</p></div>';

	// Restore HTMX attributes
	for (const [key, value] of Object.entries(attrs)) {
		container.setAttribute(key, value);
	}

	// Update the lastIndex parameter to the latest
	const urlParams = new URLSearchParams(attrs['hx-get']);
	const lastIndex = urlParams.get('lastIndex');
	if (lastIndex) {
		container.setAttribute('hx-get', '?lastIndex=' + lastIndex);
	}
}

// Filter logs based on search term
function filterLogs() {
	searchTerm = document.getElementById('search-input').value.toLowerCase();
	document.querySelectorAll('.log-entry').forEach(entry => {
		const content = entry.getAttribute('data-log-content').toLowerCase();
		const attrs = entry.querySelectorAll('.log-attr-value');

		let found = content.includes(searchTerm);

		// Also search through attributes
		if (!found && attrs.length > 0) {
			attrs.forEach(attr => {
				if (attr.textContent.toLowerCase().includes(searchTerm)) {
					found = true;
				}
			});
		}

		entry.style.display = found ? 'block' : 'none';
	});
}

// Auto-scroll to bottom when new logs arrive
function scrollToBottom() {
	if (document.getElementById('auto-scroll').checked) {
		window.scrollTo(0, document.body.scrollHeight);
	}
}

// Intercept HTMX before request to check if polling is enabled
document.addEventListener('htmx:beforeRequest', function(evt) {
	if (!polling && evt.detail.triggerSpec.includes('every')) {
		evt.detail.xhr.abort();
	}
});

// Apply filtering after new content is added
document.addEventListener('htmx:afterSwap', function() {
	if (searchTerm) {
		filterLogs();
	}
	scrollToBottom();
});

// Setup search filtering
document.getElementById('search-input').addEventListener('input', filterLogs);

// Update URL in hx-get when lastIndex changes
document.addEventListener('htmx:configRequest', function(evt) {
	// Extract the current lastIndex from the URL
	const urlStr = evt.detail.path;
	const url = new URL(urlStr, window.location.href);
	const params = url.searchParams;

	// Update the lastIndex in the hx-get attribute
	if (params.has('lastIndex')) {
		const container = document.getElementById('log-container');
		container.setAttribute('hx-get', '?lastIndex=' + params.get('lastIndex'));
	}
});
