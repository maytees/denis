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
insert(heading, "Records");

const intro = el("p", { class: "my-4" });
insert(
	intro,
	'<code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">records.toml</code> holds the records DENIS answers for. It lives next to <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">config.toml</code>. Each record maps a name to an address.',
);

const recordsCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});
const recordsCodeText = `[[records]]
name = 'localhost' # Domain
type = 'A' # Record type (only A available)
value = '127.0.0.1' # Address
ttl = 300 # Cache time

[[records]]
name = 'my.google'
type = 'A'
value = '142.251.16.138'
ttl = 0`;
insert(recordsCode, highlight(recordsCodeText));

const fieldsSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(fieldsSection, "Fields");

const fieldsList = el("ul", {
	class: "list-disc pl-5 my-4 flex flex-col gap-2",
});

const fields: [string, string][] = [
	[
		"name",
		"the domain to answer for. Lookups are case-insensitive; store names in lowercase.",
	],
	["type", "the record type. Only A records work for now."],
	[
		"value",
		"an IPv4 address, without a port — a DNS answer is just an IP. Routing ports is the reverse proxy's job, later.",
	],
	["ttl", "how long resolvers may cache the answer, in seconds."],
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

const fallbackText = el("p", { class: "my-4" });
insert(
	fallbackText,
	'Queries for names not in this file are forwarded to your configured <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">upstream</code> DNS server, so normal browsing keeps working with DENIS in front.',
);

const macosNote = el("p", {
	class: "border-l-2 border-emerald-700 pl-4 my-5 text-sm",
});
insert(
	macosNote,
	'<b class="block text-xs uppercase tracking-wide font-semibold text-emerald-700 mb-1">macOS</b>Apple uses <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">.local</code> domains for multicast DNS (on port 5353), so records ending in <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">.local</code> won\'t work on macOS. Pick another ending.',
);

const pager = el("div", {
	class:
		"flex flex-row justify-between border-t border-neutral-200 mt-10 pt-4 text-sm",
});
const prevLink = el("a", {
	href: "/docs/configuration",
	class: "underline text-emerald-700",
});
insert(prevLink, "← Configuration");
const nextLink = el("a", {
	href: "/docs/testing",
	class: "underline text-emerald-700",
});
insert(nextLink, "Testing →");
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
	recordsCode,
	fieldsSection,
	fieldsList,
	fallbackText,
	macosNote,
	pager,
	footer,
];

for (const section of sections) {
	insert(main, section);
}

mount(main, "app");
