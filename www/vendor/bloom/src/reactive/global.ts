import type { Effect } from "../types";

type Global = {
	observer: Effect | null;
};

const global: Global = {
	observer: null,
};

export const Global = {
	get observer() {
		return global.observer;
	},
	set observer(effect: Effect | null) {
		global.observer = effect;
	},
};
