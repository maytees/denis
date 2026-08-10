export function mount(elementToMount: HTMLElement, id: string) {
	const mountEl = document.getElementById(id);
	// todo: throw or error log???; if throw then TODO: create custom bloom errors
	if (!mountEl) throw new Error(`Element mount point has no id.`);
	mountEl.appendChild(elementToMount);
}
