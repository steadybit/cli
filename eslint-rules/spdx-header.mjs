// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

const LICENSE_LINE = ' SPDX-License-Identifier: MIT';
const COPYRIGHT_PATTERN = /^ SPDX-FileCopyrightText: \d{4} Steadybit GmbH$/;

// A local rule rather than a plugin: the repo dropped eslint-plugin-header when it moved
// to flat config, and enforcing two fixed comment lines does not warrant a dependency.
export const spdxHeader = {
  meta: {
    type: 'layout',
    fixable: 'code',
    schema: [],
    messages: {
      missing: 'Missing the SPDX copyright header.',
      malformed:
        'The SPDX copyright header must be "// SPDX-License-Identifier: MIT" and ' +
        '"// SPDX-FileCopyrightText: <year> Steadybit GmbH" on the two topmost lines.',
    },
  },
  create(context) {
    const sourceCode = context.sourceCode;

    return {
      Program(node) {
        const comments = sourceCode.getAllComments();
        // Both spellings are needed: espree reports the `#!` line of a .mjs file as
        // "Hashbang", while typescript-eslint reports it as "Shebang" in the .ts CLI
        // entry points.
        const hashbang = comments.find(comment => comment.type === 'Hashbang' || comment.type === 'Shebang');
        const header = comments.filter(comment => comment !== hashbang).slice(0, 2);
        const startsAfter = hashbang ? hashbang.range[1] : 0;

        const isPresent =
          header.length === 2 &&
          header.every(comment => comment.type === 'Line') &&
          header[0].value === LICENSE_LINE &&
          COPYRIGHT_PATTERN.test(header[1].value) &&
          header[0].loc.start.line + 1 === header[1].loc.start.line &&
          sourceCode.text.slice(startsAfter, header[0].range[0]).trim() === '';

        if (isPresent) {
          return;
        }

        // Only an outright absent header is auto-fixed. Rewriting a header that is
        // present but wrong risks stacking a second one on top of the first.
        const hasSomeHeader = header.some(comment => comment.value.includes('SPDX-'));

        context.report({
          node,
          messageId: hasSomeHeader ? 'malformed' : 'missing',
          fix: hasSomeHeader
            ? undefined
            : fixer => {
                const text = `//${LICENSE_LINE}\n// SPDX-FileCopyrightText: ${new Date().getFullYear()} Steadybit GmbH\n`;
                return hashbang
                  ? fixer.insertTextAfterRange(hashbang.range, `\n${text}`)
                  : fixer.insertTextBeforeRange([0, 0], `${text}\n`);
              },
        });
      },
    };
  },
};
