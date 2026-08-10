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
insert(heading, "Testing");

const intro = el("p", { class: "my-4" });
insert(
	intro,
	'<code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">dig</code> is a command line tool for querying DNS servers. Since DENIS isn\'t your default DNS, pass <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">@addr</code> to point dig at it.',
);

const digCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});
const digCodeText = `dig @127.0.0.1 localhost

# on a port other than 53
dig @127.0.0.1 -p 5533 localhost`;
insert(digCode, highlight(digCodeText));

const behaviorText = el("p", { class: "my-4" });
insert(
	behaviorText,
	'If DENIS is working correctly, it grabs the record from your <a href="/docs/records" class="underline text-emerald-700">records.toml</a>, or falls back to your upstream DNS server. If it isn\'t, dig should hang and ultimately error.',
);

const goSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(goSection, "Go tests");

const goText = el("p", { class: "my-4" });
insert(goText, "Tests are a work in progress for DENIS. Run them with:");

const goCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});
const goCodeText = `# just & with gotestsum
just prettytest

# or without gotestsum
just test

# without just
go test ./...`;
insert(goCode, highlight(goCodeText));

const pager = el("div", {
	class:
		"flex flex-row justify-between border-t border-neutral-200 mt-10 pt-4 text-sm",
});
const prevLink = el("a", {
	href: "/docs/records",
	class: "underline text-emerald-700",
});
insert(prevLink, "← Records");
insert(pager, prevLink);

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
	digCode,
	behaviorText,
	goSection,
	goText,
	goCode,
	pager,
	footer,
];

for (const section of sections) {
	insert(main, section);
}

mount(main, "app");
