// Renders README.md into the Pages template.
// README stays the single source of truth; the hero and page shell live in
// template.html, and the README's own title/tagline/footer feed into it.

const fs = require("fs");
const path = require("path");
const { marked } = require("marked");

const root = path.resolve(__dirname, "..", "..");
const pagesDir = __dirname;
const outDir = path.join(root, "_site");

const readme = fs.readFileSync(path.join(root, "README.md"), "utf8");
let html = marked.parse(readme).trim();

// The template supplies its own <h1>; drop the README's.
html = html.replace(/^<h1[^>]*>[\s\S]*?<\/h1>\s*/, "");

// First paragraph becomes the hero tagline.
let tagline = "";
html = html.replace(/^<p>([\s\S]*?)<\/p>\s*/, (_, inner) => {
  tagline = inner.trim();
  return "";
});

// Everything after the final <hr> is the footer (source link + credit).
let footer = "";
const lastHr = html.lastIndexOf("<hr>");
if (lastHr !== -1) {
  footer = html.slice(lastHr + "<hr>".length).trim();
  html = html.slice(0, lastHr).trim();
}

const template = fs.readFileSync(path.join(pagesDir, "template.html"), "utf8");
// Nothing is re-indented here: <pre> content is whitespace-significant, so
// cosmetic indentation would corrupt multi-line code blocks.
const page = template
  .replace("{{TAGLINE_TEXT}}", stripTags(tagline))
  .replace("{{TAGLINE}}", tagline)
  .replace("{{CONTENT}}", html)
  .replace("{{FOOTER}}", footer);

fs.mkdirSync(outDir, { recursive: true });
fs.writeFileSync(path.join(outDir, "index.html"), page);
for (const asset of ["style.css", "copy.js"]) {
  fs.copyFileSync(path.join(pagesDir, asset), path.join(outDir, asset));
}

console.log(`built _site/index.html (${page.length} bytes)`);

function stripTags(s) {
  return s.replace(/<[^>]+>/g, "").replace(/"/g, "&quot;");
}
