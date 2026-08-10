// This file was written by AI (Claude). The site is built with Bloom,
// a hand-written signals library: https://github.com/maytees/bloom

import { el, insert, mount } from "bloom";
import { highlight } from "./highlight";

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
insert(heading, "DENIS");

const subheading = el("p", { class: "text-neutral-500 text-sm mt-1" });
insert(subheading, "A hand-written DNS server for your local network.");

const descParagraph = el("p", { class: "my-4" });
insert(
	descParagraph,
	"Give your local services names you'll actually remember — http://notes instead of localhost:5230. No DNS libraries doing the heavy lifting. A project for learning Go.",
);

const cta = el("div", { class: "flex flex-row gap-4" });

for (const [item, href] of Object.entries({
	"read the docs": "/docs",
	"view the source": "https://github.com/maytees/denis",
})) {
	const ctaLink = el("a", {
		href: href,
		class: "underline text-emerald-700",
		target: item === "view the source" ? "_blank" : "_self",
	});

	insert(ctaLink, item);
	insert(cta, ctaLink);
}

const separator = el("div", {
	class: "w-full bg-neutral-200 h-px rounded-full my-8",
});

const todaySection = el("h2", { class: "text-xl font-semibold" });
insert(todaySection, "Today - the DNS server");

const todayCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});

const todayCodeText = `# records.toml
[[records]]
name = 'notes'
type = 'A'
value = '127.0.0.1'
ttl = 300

# then ask DENIS about it
dig @127.0.0.1 notes`;

insert(todayCode, highlight(todayCodeText));

const planSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(planSection, "The plan - more than a DNS");

const planSubheading = el("p", { class: "text-neutral-500 text-sm mt-1" });
insert(
	planSubheading,
	"This doesn't exist yet. It's what the DNS server is building toward.",
);

const planCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});

const planCodeText = `http://notes
   -> DENIS resolves notes to 127.0.0.1
   -> a hand-written reverse proxy on port 80
   -> your app on localhost:5230

# managed from a web interface, backed by an
# api living inside the DENIS binary`;
insert(planCode, highlight(planCodeText));

const ideasSeparator = el("div", {
	class: "w-full bg-neutral-200 h-px rounded-full my-8",
});

const ideasSection = el("h2", { class: "text-xl font-semibold" });
insert(ideasSection, "Three ideas");

const ideasList = el("ul", {
	class: "list-disc pl-5 my-4 flex flex-col gap-2",
});

const ideas: [string, string][] = [
	[
		"Build your own.",
		"A hosts file or nginx could do this. The point is doing it anyway: the DNS server, the reverse proxy, and the HTTP server, all from scratch.",
	],
	[
		"Local names for local things.",
		"Every service in your homelab gets a real domain. No more memorizing ports.",
	],
	[
		"Straight from the spec.",
		"RFC 1035 implemented by hand — headers, names, and questions parsed byte by byte.",
	],
];

for (const [label, description] of ideas) {
	const idea = el("li");
	const ideaLabel = el("b", { class: "font-semibold" });
	insert(ideaLabel, label);
	const ideaText = el("span");
	insert(ideaText, ` ${description}`);
	insert(idea, ideaLabel);
	insert(idea, ideaText);
	insert(ideasList, idea);
}

const statusSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(statusSection, "Where it's at");

const statusList = el("ul", {
	class: "list-disc pl-5 my-4 flex flex-col gap-2",
});

const statuses: [string, string][] = [
	["done", "UDP server, listening for queries"],
	["done", "header, name & question parsing"],
	["wip", "the message pipeline — answers & upstream forwarding"],
	["next", "the reverse proxy, routing domains to ports"],
	["next", "api & web interface for managing records"],
];

for (const [status, description] of statuses) {
	const row = el("li");
	const label = el("b", { class: "font-semibold" });
	insert(label, status);
	const text = el("span");
	insert(text, ` — ${description}`);
	insert(row, label);
	insert(row, text);
	insert(statusList, row);
}

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
	descParagraph,
	cta,
	separator,
	todaySection,
	todayCode,
	planSection,
	planSubheading,
	planCode,
	ideasSeparator,
	ideasSection,
	ideasList,
	statusSection,
	statusList,
	footer,
];

for (const section of sections) {
	insert(main, section);
}

mount(main, "app");
