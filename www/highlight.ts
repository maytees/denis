// Tiny syntax highlighter for the code blocks on this site. Made with AI.
// Adapted from the Bloom website's highlighter — tweaked for the
// bash & TOML snippets in the DENIS docs (# comments, commands).

const KEYWORDS = new Set([
	// shell commands
	"sudo",
	"git",
	"go",
	"just",
	"air",
	"dig",
	"cp",
	"cd",
	"bun",
	"gotestsum",
	"docker",
	"networksetup",
	// toml booleans
	"true",
	"false",
]);

const TOKEN_RE =
	/(#[^\n]*|\/\/[^\n]*)|("(?:[^"\\]|\\.)*"|'(?:[^'\\]|\\.)*'|`(?:[^`\\]|\\.)*`)|(\b\d+(?:\.\d+)?\b)|([A-Za-z_$][A-Za-z0-9_$]*)/g;

function escapeHtml(text: string): string {
	return text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

export function highlight(code: string): string {
	let result = "";
	let lastIndex = 0;

	for (const match of code.matchAll(TOKEN_RE)) {
		const [full, comment, string, number, word] = match;
		const index = match.index ?? 0;

		result += escapeHtml(code.slice(lastIndex, index));

		if (comment) {
			result += `<span class="text-neutral-400 italic">${escapeHtml(comment)}</span>`;
		} else if (string) {
			result += `<span class="text-emerald-700">${escapeHtml(string)}</span>`;
		} else if (number) {
			result += `<span class="text-orange-700">${escapeHtml(number)}</span>`;
		} else if (word && KEYWORDS.has(word)) {
			result += `<span class="text-purple-700">${escapeHtml(word)}</span>`;
		} else {
			result += escapeHtml(full);
		}

		lastIndex = index + full.length;
	}

	result += escapeHtml(code.slice(lastIndex));
	return result;
}
