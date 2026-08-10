// args for signal can either be 1 or none
export type Signal<T = void> = (...args: T[] | [newValue: T]) => T;
// export type ExtractSignalType<S> = S extends Signal<infer T> ? T : never;

// is an object right now, even tho it holds 1 value.
export type ReactiveNode = { subscribers: Set<Effect> };

export type EffectCallback = () => void;
// dispose should be an object so it's clear,
// instead of a user set var name
export type EffectDispose = { dispose: () => void };
export type Effect = {
	readonly callbackFn: EffectCallback;
	readonly run: () => void;
	dependencies: Set<ReactiveNode>;
};
