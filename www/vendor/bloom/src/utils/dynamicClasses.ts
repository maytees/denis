import type { DynamicClasses } from "./types";

// dynamicClass("flex flex-row", {
// 	"asdf asdf": true,
// 	shouldHaveQuotes: false,
// });

// TODO: in the bloom compiler, if base isn't passed in,
// automatically transpile undefined, and assume dynamic is
// whats passed in. as long as its a Record<string, boolean>
// object ofc tho.
export const dynamicClass = (
	base: string | undefined,
	dynamic?: DynamicClasses,
): (() => string) => {
	return () => {
		let strBuilder = base ? base.trim() : "";

		if (!dynamic) return strBuilder;

		const entries = Object.entries(dynamic);

		for (const [className, val] of entries) {
			// const condition = typeof val === "function" ? val() : val;
			const condition = val();
			if (condition)
				strBuilder += strBuilder ? ` ${className.trim()}` : className.trim();
		}

		return strBuilder;
	};
};

export const $c = dynamicClass;
