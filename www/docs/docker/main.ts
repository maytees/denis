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
insert(heading, "Docker");

const intro = el("p", { class: "my-4" });
insert(
	intro,
	'Instead of installing Go, you can run DENIS in a container. You need <a href="https://docs.docker.com/get-docker/" target="_blank" class="underline text-emerald-700">Docker</a> (Docker Desktop on Mac/Windows) — nothing else.',
);

const configSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(configSection, "Create config files");

const configText = el("p", { class: "my-4" });
insert(
	configText,
	'Same two files as a native install — see <a href="/docs/configuration" class="underline text-emerald-700">Configuration</a>. They are mounted into the container, not baked into the image, so you can edit them without rebuilding.',
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

const hostNote = el("p", {
	class: "border-l-2 border-emerald-700 pl-4 my-5 text-sm",
});
insert(
	hostNote,
	'<b class="block text-xs uppercase tracking-wide font-semibold text-emerald-700 mb-1">Important</b>In <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">config/config.toml</code>, set <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">host = \'0.0.0.0\'</code>. Inside a container, binding 127.0.0.1 means Docker can\'t route traffic to DENIS. This does not expose DENIS to your network — the compose file only publishes the port on your machine\'s loopback.',
);

const runSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(runSection, "Build and run");

const runCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});
const runCodeText = `# with just
just docker-compose

# without
docker compose up -d --build

# watch the logs (Ctrl+C stops watching, not the container)
just docker-logs  # runs: docker compose logs -f`;
insert(runCode, highlight(runCodeText));

const sudoNote = el("p", {
	class: "border-l-2 border-emerald-700 pl-4 my-5 text-sm",
});
insert(
	sudoNote,
	'<b class="block text-xs uppercase tracking-wide font-semibold text-emerald-700 mb-1">Note</b>No <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">sudo</code> needed — port 53 is only privileged on the host, not inside the container.',
);

const testSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(testSection, "Test it");

const testText = el("p", { class: "my-4" });
insert(
	testText,
	'Same as <a href="/docs/testing" class="underline text-emerald-700">Testing</a>: one query for a record you own, one for a domain that forwards upstream. Both must answer before you point your OS at DENIS.',
);

const testCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});
const testCodeText = `# with just
just docker-check

# without
dig @127.0.0.1 localhost    # a record from records.toml
dig @127.0.0.1 google.com   # forwarded upstream`;
insert(testCode, highlight(testCodeText));

const updateSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(updateSection, "Updating");

const updateText = el("p", { class: "my-4" });
insert(
	updateText,
	"Code lives in the image; config lives on your disk. Changed Go code (or pulled new commits)? Rebuild. Changed a TOML file? Just restart — records are read once at startup.",
);

const updateCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});
const updateCodeText = `# after code changes — builds first, then swaps containers,
# so a failed build never takes DENIS down
just docker-compose

# after config/records changes
just docker-restart  # runs: docker compose restart

# stop it (stays stopped until you start it again)
just docker-down  # runs: docker compose down`;
insert(updateCode, highlight(updateCodeText));

const systemSection = el("h2", { class: "text-xl font-semibold mt-8" });
insert(systemSection, "Using it as your system DNS");

const systemText = el("p", { class: "my-4" });
insert(
	systemText,
	"To make your OS resolve every name through DENIS, point your system's DNS server at 127.0.0.1. On macOS:",
);

const systemCode = el("pre", {
	class:
		"bg-neutral-100 border-l-2 border-neutral-300 p-4 my-4 text-sm font-mono leading-relaxed overflow-x-auto",
});
const systemCodeText = `# point macOS at DENIS (swap "Wi-Fi" for your active service)
sudo networksetup -setdnsservers "Wi-Fi" 127.0.0.1

# revert to DHCP-provided DNS
sudo networksetup -setdnsservers "Wi-Fi" "Empty"`;
insert(systemCode, highlight(systemCodeText));

const escapeNote = el("p", {
	class: "border-l-2 border-emerald-700 pl-4 my-5 text-sm",
});
insert(
	escapeNote,
	'<b class="block text-xs uppercase tracking-wide font-semibold text-emerald-700 mb-1">Warning</b>While DENIS is your system resolver, stopping the container means your machine can\'t resolve anything. Keep the revert command handy. The compose file sets <code class="rounded bg-neutral-100 px-1 py-0.5 font-mono text-[0.85em]">restart: unless-stopped</code>, so DENIS comes back after crashes and reboots on its own — but an explicit stop stays stopped.',
);

const pager = el("div", {
	class:
		"flex flex-row justify-between border-t border-neutral-200 mt-10 pt-4 text-sm",
});
const prevLink = el("a", {
	href: "/docs/testing",
	class: "underline text-emerald-700",
});
insert(prevLink, "← Testing");
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
	configSection,
	configText,
	configCode,
	hostNote,
	runSection,
	runCode,
	sudoNote,
	testSection,
	testText,
	testCode,
	updateSection,
	updateText,
	updateCode,
	systemSection,
	systemText,
	systemCode,
	escapeNote,
	pager,
	footer,
];

for (const section of sections) {
	insert(main, section);
}

mount(main, "app");
