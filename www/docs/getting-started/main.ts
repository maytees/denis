// This file was written by AI (Claude). The site is built with Bloom,
// a hand-written signals library: https://github.com/maytees/bloom

import { el, insert, mount } from "bloom";
import { highlight } from "../../highlight";

const main = el("main", { class: "max-w-2xl mx-auto px-5 py-10" });

const nav = el("nav", { class: "flex flex-row items-center gap-4" });

const navItems = {
	denis: "/",
	docs: "/docs",
	github: "https://github.com/maytees/denis",
};

for (const [item, href] of Object.entries(navItems)) {
	const navLink = el("a", {
		href: href,
		class: "underline text-sm text-emerald-700",
		target: item === "github" ? "_blank" : "_self",
	});

	insert(navLink, item);
	insert(nav, navLink);
}

const heading = el("h1", { class: "text-2xl font-semibold mt-8" });
insert(heading, "Getting started");

const intro = el("p", { class: "my-4" });
insert(
	intro,
	'You need <a href="https://go.dev" target="_blank" class="underline text-emerald-700">go</a> installed. Optionally, grab <a href="https://github.com/casey/just" target="_blank" class="underline text-emerald-700">just</a> (task runner), <a href="https://github.com/air-verse/air" target="_blank" class="underline text-emerald-700">air</a> (hot reloading for dev), and <a href="https://github.com/gotestyourself/gotestsum" target="_blank" class="underline text-emerald-700">gotestsum</a> (pretty testing).',
);

const cloneCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});
const cloneCodeText = `git clone https://github.com/maytees/denis
cd denis`;
insert(cloneCode, highlight(cloneCodeText));

const configSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(configSection, "Create config files");

const configText = el("p", { class: "my-4" });
insert(
	configText,
	"DENIS reads two files: a config and a records file. From the repo folder, copy the examples to get going.",
);

const configCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});
const configCodeText = `# with just
just config-example

# without
cp ./config/config.example.toml ./config/config.toml
cp ./config/records.example.toml ./config/records.toml`;
insert(configCode, highlight(configCodeText));

const configNote = el("p", {
	class: "border-l-2 border-emerald-700 pl-4 my-5 text-sm",
});
insert(
	configNote,
	'<b class="block text-xs uppercase tracking-wide font-semibold text-emerald-700 mb-1">Note</b>If there is no <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">config</code> folder in your working directory, DENIS creates default config files in your OS config directory instead. See <a href="/docs/configuration" class="underline text-emerald-700">Configuration</a> for where that is.',
);

const runSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(runSection, "Run it");

const runCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});
const runCodeText = `# with just
just start

# without
go build -o dist/denis && sudo dist/denis

# or with hot reloading for dev
just dev  # runs: sudo air`;
insert(runCode, highlight(runCodeText));

const sudoNote = el("p", {
	class: "border-l-2 border-emerald-700 pl-4 my-5 text-sm",
});
insert(
	sudoNote,
	'<b class="block text-xs uppercase tracking-wide font-semibold text-emerald-700 mb-1">Note</b>DNS runs on port 53, a privileged port — that\'s why everything runs with <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">sudo</code>. Set a port above 1024 in your config to skip it.',
);

const pager = el("div", {
	class:
		"flex flex-row justify-between border-t border-neutral-200 mt-10 pt-4 text-sm",
});
const prevLink = el("a", {
	href: "/docs",
	class: "underline text-emerald-700",
});
insert(prevLink, "← Docs");
const nextLink = el("a", {
	href: "/docs/configuration",
	class: "underline text-emerald-700",
});
insert(nextLink, "Configuration →");
insert(pager, prevLink);
insert(pager, nextLink);

const footer = el("footer", {
	class: "border-t border-neutral-200 mt-14 pt-4 text-xs text-neutral-500",
});
const bloomAnchor = el("a", {
	href: "https://github.com/maytees/bloom",
	target: "_blank",
	class: "underline text-emerald-700",
});
insert(bloomAnchor, "Bloom");
const footerTail = el("span");
insert(
	footerTail,
	". Unlike the rest of DENIS, this site was written by AI (Claude).",
);
insert(footer, "This website is fully made with ");
insert(footer, bloomAnchor);
insert(footer, footerTail);

const sections = [
	nav,
	heading,
	intro,
	cloneCode,
	configSection,
	configText,
	configCode,
	configNote,
	runSection,
	runCode,
	sudoNote,
	pager,
	footer,
];

for (const section of sections) {
	insert(main, section);
}

mount(main, "app");
