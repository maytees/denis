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
insert(heading, "Configuration");

const intro = el("p", { class: "my-4" });
insert(
	intro,
	'<code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">config.toml</code> controls how DENIS runs. Right now there is only a <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">[dns]</code> section — the reverse proxy, api, and web interface will get their own sections when they exist.',
);

const configCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});
const configCodeText = `[dns]
enabled = true # Toggle DNS
port = 53 # Default port for DNS
upstream = '8.8.8.8' # Fallback for unknown names`;
insert(configCode, highlight(configCodeText));

const fieldsSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(fieldsSection, "Fields");

const fieldsList = el("ul", {
	class: "list-disc pl-5 my-4 flex flex-col gap-2",
});

const fields: [string, string][] = [
	["enabled", "turns the DNS service on or off."],
	[
		"port",
		"the UDP port to listen on. 53 is the DNS default, but needs sudo — anything above 1024 doesn't.",
	],
	[
		"upstream",
		"where queries go when DENIS doesn't own the record. Google, Cloudflare, your router — any DNS server. Port defaults to 53 if left out.",
	],
];

for (const [label, description] of fields) {
	const field = el("li");
	const fieldLabel = el("b", { class: "font-semibold font-mono text-sm" });
	insert(fieldLabel, label);
	const fieldText = el("span");
	insert(fieldText, ` — ${description}`);
	insert(field, fieldLabel);
	insert(field, fieldText);
	insert(fieldsList, field);
}

const locationSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(locationSection, "Where it lives");

const locationText = el("p", { class: "my-4" });
insert(
	locationText,
	'If a <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">config</code> folder exists in your working directory, DENIS uses that (handy for local dev). Otherwise it falls back to your OS config directory, via <a href="https://github.com/kirsle/configdir?tab=readme-ov-file#configdir-for-go" target="_blank" class="underline text-emerald-700">configdir</a> — e.g. <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">~/Library/Application Support/denis</code> on macOS.',
);

const defaultsNote = el("p", {
	class: "border-l-2 border-emerald-700 pl-4 my-5 text-sm",
});
insert(
	defaultsNote,
	'<b class="block text-xs uppercase tracking-wide font-semibold text-emerald-700 mb-1">Note</b>If the config files don\'t exist, DENIS writes defaults on first run — DNS enabled, port 53, upstream <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">8.8.8.8:53</code>. A proper CLI with a config folder flag is in the works.',
);

const pager = el("div", {
	class:
		"flex flex-row justify-between border-t border-neutral-200 mt-10 pt-4 text-sm",
});
const prevLink = el("a", {
	href: "/docs/getting-started",
	class: "underline text-emerald-700",
});
insert(prevLink, "← Getting started");
const nextLink = el("a", {
	href: "/docs/records",
	class: "underline text-emerald-700",
});
insert(nextLink, "Records →");
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
	configCode,
	fieldsSection,
	fieldsList,
	locationSection,
	locationText,
	defaultsNote,
	pager,
	footer,
];

for (const section of sections) {
	insert(main, section);
}

mount(main, "app");
