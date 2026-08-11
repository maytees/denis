// This file was written by AI (Claude). The site is built with Bloom,
// a hand-written signals library: https://github.com/maytees/bloom

import { el, insert, mount } from "bloom";

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
insert(heading, "Documentation");

const subheading = el("p", { class: "text-neutral-500 text-sm mt-1" });
insert(
	subheading,
	"DENIS is a working DNS server, but it doesn't implement all of RFC 1035 yet. These docs cover getting it running on your machine.",
);

const docsList = el("ul", { class: "list-disc pl-5 my-4 flex flex-col gap-2" });

const docsLinks: [string, string, string][] = [
	["Getting started", "/docs/getting-started", "install, configure, run"],
	["Configuration", "/docs/configuration", "config.toml and where it lives"],
	["Records", "/docs/records", "records.toml and A records"],
	["Testing", "/docs/testing", "querying with dig, running the Go tests"],
	["Docker", "/docs/docker", "run DENIS in a container, use it as your system DNS"],
];

for (const [label, href, description] of docsLinks) {
	const item = el("li");
	insert(
		item,
		`<a href="${href}" class="underline text-emerald-700">${label}</a> — ${description}`,
	);
	insert(docsList, item);
}

const separator = el("div", {
	class: "w-full bg-neutral-200 h-px rounded-full my-8",
});

const visionSection = el("h2", { class: "text-xl font-semibold" });
insert(visionSection, "Working vs vision");

const visionText = el("p", { class: "my-4" });
insert(
	visionText,
	'Everything in these docs describes what works today: a DNS server that resolves your local records and forwards the rest upstream. The reverse proxy, api, and web interface described on the <a href="/" class="underline text-emerald-700">home page</a> are the vision — they do not exist yet.',
);

const pager = el("div", {
	class:
		"flex flex-row justify-end border-t border-neutral-200 mt-10 pt-4 text-sm",
});
const nextLink = el("a", {
	href: "/docs/getting-started",
	class: "underline text-emerald-700",
});
insert(nextLink, "Getting started →");
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
	subheading,
	docsList,
	separator,
	visionSection,
	visionText,
	pager,
	footer,
];

for (const section of sections) {
	insert(main, section);
}

mount(main, "app");
