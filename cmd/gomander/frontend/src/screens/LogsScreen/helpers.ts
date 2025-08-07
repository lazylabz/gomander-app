import AnsiToHtml from "ansi-to-html";

const ansiConverter = new AnsiToHtml({
  fg: "#000000",
  bg: "#ffffff",
  escapeXML: true,
  newline: true,
  stream: false,
});

const addSearchHighlight = (element: HTMLElement, searchQuery: string) => {
  if (
    !searchQuery ||
    !element.innerText.toLowerCase().includes(searchQuery.toLowerCase())
  ) {
    return element.innerHTML;
  }

  element.childNodes.forEach((node) => {
    if (node.nodeType === Node.TEXT_NODE) {
      const text = node.textContent || "";

      const replacedText = text.replace(
        new RegExp(`(${searchQuery})`, "gi"),
        (match) =>
          `<mark data-match="${crypto.randomUUID()}" class="bg-yellow-100 match">${match}</mark>`,
      );
      if (replacedText !== text) {
        const span = document.createElement("span");
        span.innerHTML = replacedText;
        node.replaceWith(span);
      }
    } else {
      // Recursively handle child nodes
      addSearchHighlight(node as HTMLElement, searchQuery);
    }
  });

  return element.innerHTML;
};

export const parseLog = (log: string, searchQuery: string) => {
  const html = ansiConverter.toHtml(log);

  if (!searchQuery) {
    return html;
  }

  const element = document.createElement("p");
  element.innerHTML = html;

  return addSearchHighlight(element, searchQuery);
};

export const extractMatchesIds = (log: string) => {
  const element = document.createElement("div");
  element.innerHTML = log;
  return Array.from(element.querySelectorAll("mark.match"))
    .map((e) => e.getAttribute("data-match") || "")
    .filter(Boolean);
};
