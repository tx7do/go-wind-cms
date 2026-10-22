/**
 * 图标子集生成器
 *
 * 扫 src 下用到的 `carbon:<name>`，产出两份：
 *   1. src/plugins/generated/iconfont.css    图标字体 @font-face（全端唯一的渲染载体）
 *   2. src/plugins/generated/icon-glyphs.ts  图标名 → 码点
 *
 * 为什么必须裁：@iconify-json/carbon 全量 2776 个图标 / 1167 KB JSON，全项目只用到几十个；
 * 之前它整份进 vendors.js，dist 2.3 MB 直接顶穿微信主包 2 MB 上限。
 * 为什么全端走字体而不是 svg：小程序端 svg 进不了渲染树（Taro 的模板选择只认固定节点名
 * 白名单，编译产物 *.wxml 里 svg 出现 0 次），H5 若留 svg 就成了预览面 ≠ 出货面。
 *
 * 名字写错（图标集里没有）一律 fail，不留到运行时变空白。
 */
import {createRequire} from 'node:module';
import {fileURLToPath} from 'node:url';
import {dirname, join, relative} from 'node:path';
import {readFileSync, writeFileSync, mkdirSync, readdirSync, statSync} from 'node:fs';
import {Readable} from 'node:stream';
import {SVGIcons2SVGFontStream} from 'svgicons2svgfont';
import svg2ttf from 'svg2ttf';

const require = createRequire(import.meta.url);
const here = dirname(fileURLToPath(import.meta.url));
const projectRoot = join(here, '..');
const srcRoot = join(projectRoot, 'src');
const outDir = join(srcRoot, 'plugins', 'generated');

const FONT_NAME = 'gw-icon';
const FIRST_CODEPOINT = 0xe001;
const ICON_PREFIX = 'carbon:';

const carbon = require('@iconify-json/carbon/icons.json');
/** carbon 是 32×32 栅格（不是 16）：viewBox 写错会把图形放大一倍切掉，字体和 svg 两条路都得用它 */
const GRID_W = carbon.width || 16;
const GRID_H = carbon.height || 16;

function walk(dir, files = []) {
    for (const entry of readdirSync(dir)) {
        const full = join(dir, entry);
        if (statSync(full).isDirectory()) {
            if (full === outDir) continue;
            walk(full, files);
        } else if (/\.(t|j)sx?$/.test(entry)) {
            files.push(full);
        }
    }
    return files;
}

function collectUsedNames() {
    const found = new Set();
    for (const file of walk(srcRoot)) {
        const text = readFileSync(file, 'utf8');
        for (const match of text.matchAll(/carbon:([a-z][a-z0-9-]*)/g)) {
            found.add(match[1]);
        }
    }
    return [...found].sort();
}

/** iconify body 里的占位符按 16x16 网格展开，字体和 svg 两条路都用同一份 body */
function resolveBody(name) {
    const icon = carbon.icons[name];
    return icon.body
        .replace(/\{width\}/g, String(icon.width || GRID_W))
        .replace(/\{height\}/g, String(icon.height || GRID_H));
}

/**
 * 字体只认「填出来的轮廓」，svg 的绘制属性一律无效：
 * fill="none" 的子路径在 svg 里本来就不落地，转字体前必须删掉，否则它会和主路径
 * 按非零环绕并成一坨（carbon:checkmark-filled 就是靠它盖住勾形的）；
 * fill="currentColor" 同理要摘掉，字体颜色由 color 决定。
 */
function toFontBody(body) {
    return body
        .replace(/<path[^>]*fill="none"[^>]*\/>/g, '')
        .replace(/ fill="currentColor"/g, '');
}

function buildSvgFont(glyphs) {
    return new Promise((resolve, reject) => {
        const chunks = [];
        const fontStream = new SVGIcons2SVGFontStream({
            fontName: FONT_NAME,
            fontHeight: 1024,
            descent: 0,
            ascent: 1024,
            normalize: true,
            centerHorizontally: true,
            fontId: FONT_NAME,
        });
        fontStream.on('data', (c) => chunks.push(String(c)));
        fontStream.on('end', () => resolve(chunks.join('')));
        fontStream.on('error', reject);
        for (const glyph of glyphs) {
            const stream = Readable.from([
                `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${GRID_W} ${GRID_H}" width="${GRID_W}" height="${GRID_H}">${glyph.body}</svg>`,
            ]);
            stream.metadata = {name: glyph.name, unicode: [glyph.char]};
            fontStream.write(stream);
        }
        fontStream.end();
    });
}

async function main() {
    const used = collectUsedNames();
    const missing = used.filter((n) => !carbon.icons[n]);
    if (missing.length) {
        console.error(`[gen-icons] carbon 图标集里不存在：${missing.map((m) => ICON_PREFIX + m).join(', ')}`);
        console.error(`[gen-icons] 可选名参考：${Object.keys(carbon.icons).filter((n) => n.startsWith(missing[0].slice(0, 4))).slice(0, 8).join(', ')}`);
        process.exit(1);
    }

    const stroked = used.filter((n) => /stroke=/.test(carbon.icons[n].body));
    if (stroked.length) {
        console.warn(`[gen-icons] 以下图标带描边，转字体后会丢线宽：${stroked.join(', ')}`);
    }

    const glyphs = used.map((name, i) => ({name, char: String.fromCodePoint(FIRST_CODEPOINT + i), body: toFontBody(resolveBody(name))}));
    const svgFont = await buildSvgFont(glyphs);
    const ttf = svg2ttf(svgFont, {});
    const base64 = Buffer.from(ttf.buffer).toString('base64');

    mkdirSync(outDir, {recursive: true});
    writeFileSync(join(outDir, 'icon-glyphs.ts'), [
        '/* 生成文件，勿手改：node scripts/gen-icons.mjs */',
        'export const iconGlyphs: Record<string, string> = {',
        ...glyphs.map((g) => `    '${ICON_PREFIX}${g.name}': '${g.char}',`),
        '};',
        '',
    ].join('\n'));
    writeFileSync(join(outDir, 'iconfont.css'), [
        '/* 生成文件，勿手改：node scripts/gen-icons.mjs */',
        '@font-face {',
        `  font-family: '${FONT_NAME}';`,
        `  src: url('data:font/ttf;base64,${base64}') format('truetype');`,
        '}',
        '',
        `.${FONT_NAME} {`,
        `  font-family: '${FONT_NAME}';`,
        '  font-style: normal;',
        '  font-weight: normal;',
        '  line-height: 1;',
        '  text-rendering: optimizeLegibility;',
        '}',
        '',
    ].join('\n'));

    const kb = (n) => `${(n / 1024).toFixed(1)} KB`;
    console.log(`[gen-icons] ${used.length} 个图标 | 字体 ${kb(ttf.buffer.byteLength)}（base64 ${kb(base64.length)}）｜对照：全量 carbon JSON ${kb(readFileSync(require.resolve('@iconify-json/carbon/icons.json'), 'utf8').length)}`);
    console.log(`[gen-icons] 产物：${relative(projectRoot, join(outDir, 'iconfont.css'))} / icon-glyphs.ts`);
}

await main();
