import type { ReactiveNode, Signal } from "../types";
import { Global } from "./global";

export function signal<T = undefined>(defaultValue?: T): Signal<T> {
	let value = defaultValue as T;
	// a node is just an obj holding a signal's state.
	// effect's use a node to identify a signal
	const node: ReactiveNode = { subscribers: new Set() };

	return (...args) => {
		if (args.length === 0) {
			if (Global.observer && !node.subscribers.has(Global.observer)) {
				// register subscribers for this signal
				node.subscribers.add(Global.observer);

				// register dependencies for the effect
				Global.observer.dependencies.add(node);
			}
			return value;
		}

		if (Object.is(value, args[0])) return value;

		value = args[0];

		// you have to spread this into an array to prevent an
		// infinite loop
		[...node.subscribers].forEach((effect) => {
			effect.run();
		});

		return value;
	};
}

// Shorthand alias for signals
export const $ = signal;
