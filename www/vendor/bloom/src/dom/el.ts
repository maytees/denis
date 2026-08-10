import { effect } from "../reactive/effect";
import type { IntrinsicElements } from "../types";

export function el<K extends keyof IntrinsicElements>(
	tag: K,
	props?: IntrinsicElements[K],
) {
	const element = document.createElement(tag);

	const entries = Object.entries(props ?? {});

	for (const [key, value] of entries) {
		if (!value) continue;

		if (
			key.startsWith("on") &&
			key[2] === key[2]?.toUpperCase() &&
			typeof value === "function"
		) {
			const eventName = key.slice(2).toLowerCase();
			element.addEventListener(eventName, value);
			continue;
		}

		if (typeof value === "function") {
			effect(() => {
				const evalVal = value();
				element.setAttribute(key, evalVal);
			});
			continue;
		}
		element.setAttribute(key, String(value));
	}

	return element;
}
