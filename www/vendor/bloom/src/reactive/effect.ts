import type { Effect, EffectCallback, EffectDispose, Node } from "../types";
import { Global } from "./global";

export function effect(callbackFn: EffectCallback): EffectDispose {
	const dependencies = new Set<Node>();
	const effectObj: Effect = {
		callbackFn,
		dependencies,
		run: () => {
			cleanEffect(dependencies, effectObj);
			Global.observer = effectObj;
			callbackFn();
			Global.observer = null;
		},
	};

	effectObj.run();

	return {
		dispose: () => {
			cleanEffect(dependencies, effectObj);
		},
	};
}

function cleanEffect(dependencies: Set<Node>, effectObj: Effect) {
	for (const node of dependencies) node.subscribers.delete(effectObj);
	dependencies.clear();
}
