export function insert(
	element: HTMLElement,
	value: string | number | Node | null,
) {
	switch (typeof value) {
		case "string":
			element.innerHTML = value;
			break;
		default:
			// TODO: idk what to do here instead of as node; fix later
			element.appendChild(value as Node);
	}
}
