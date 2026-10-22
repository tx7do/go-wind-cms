// i18n 完整性检查：代码用到的 key vs zh-CN/en-US 语言包定义（双向：缺失 + 未使用）
// 用法：node scripts/i18n-check.mjs
//
// 扫描策略（保守优先：宁可漏报未使用，不可误删在用 key）：
// - 剥离注释后，源码里所有「以已知 ns 开头的字符串字面量」都计为已使用
//   （覆盖 t() / usePageTitle() / key 构造函数返回值 / 三元表达式等）
// - next-intl-compat 的 useTranslations('ns') 上下文里的短 key 补全为 `ns.key`
// - 模板字符串 `` t(`edu.wrong.reason_${x}`) `` 记录静态前缀，前缀命中的 key 视为已使用
import fs from 'node:fs';
import path from 'node:path';

const ROOT = path.resolve(process.cwd(), 'src');
const LOCALES = path.resolve(process.cwd(), 'messages');

// ---------- 1. 收集语言包 key（全站扁平成一个命名空间，full key = `文件名.路径`） ----------
const flatten = (obj, prefix = '', out = new Set()) => {
    for (const [k, v] of Object.entries(obj)) {
        const key = prefix ? `${prefix}.${k}` : k;
        if (v && typeof v === 'object') flatten(v, key, out);
        else out.add(key);
    }
    return out;
};

const loadLocale = (lang) => {
    const dir = path.join(LOCALES, lang);
    const keys = new Set();
    for (const f of fs.readdirSync(dir)) {
        if (!f.endsWith('.json')) continue;
        const ns = f.replace(/\.json$/, '');
        for (const k of flatten(JSON.parse(fs.readFileSync(path.join(dir, f), 'utf8')))) {
            keys.add(`${ns}.${k}`);
        }
    }
    return keys;
};

const zh = loadLocale('zh-CN');
const en = loadLocale('en-US');
const namespaces = new Set([...zh, ...en].map((k) => k.split('.')[0]));
const nsPrefixRe = new RegExp(`^(?:${[...namespaces].map(escapeRe).join('|')})\\.`);

function escapeRe(s) {
    return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

// ---------- 2. 扫描源码 key 引用 ----------
const walk = (dir, out = []) => {
    for (const e of fs.readdirSync(dir, { withFileTypes: true })) {
        const p = path.join(dir, e.name);
        if (e.isDirectory()) walk(p, out);
        else if (/\.(ts|tsx|js|jsx)$/.test(e.name)) out.push(p);
    }
    return out;
};

// 剥离注释（保留字符串字面量内容），string-aware，够用即可
function stripComments(code) {
    let out = '';
    let i = 0;
    let strCh = null; // 当前所处字符串的引号
    while (i < code.length) {
        const c = code[i];
        if (strCh) {
            out += c;
            if (c === '\\') {
                out += code[i + 1] ?? '';
                i += 2;
                continue;
            }
            if (c === strCh) strCh = null;
            i += 1;
            continue;
        }
        if (c === '"' || c === "'" || c === '`') {
            strCh = c;
            out += c;
            i += 1;
            continue;
        }
        if (c === '/' && code[i + 1] === '/') {
            while (i < code.length && code[i] !== '\n') i += 1;
            continue;
        }
        if (c === '/' && code[i + 1] === '*') {
            i += 2;
            while (i < code.length && !(code[i] === '*' && code[i + 1] === '/')) i += 1;
            i += 2;
            continue;
        }
        out += c;
        i += 1;
    }
    return out;
}

// 提取字符串字面量：返回 {value, dynamic}，dynamic=含 ${} 的模板串取 `${` 前静态前缀
function extractStrings(code) {
    const out = [];
    let i = 0;
    while (i < code.length) {
        const c = code[i];
        if (c === '"' || c === "'" || c === '`') {
            const ch = c;
            let v = '';
            i += 1;
            while (i < code.length && code[i] !== ch) {
                if (code[i] === '\\') {
                    v += code[i] + (code[i + 1] ?? '');
                    i += 2;
                    continue;
                }
                v += code[i];
                i += 1;
            }
            i += 1;
            const dyn = v.indexOf('${');
            out.push(dyn === -1 ? {value: v, dynamic: false} : {value: v.slice(0, dyn), dynamic: true});
            continue;
        }
        i += 1;
    }
    return out;
}

const files = walk(ROOT);
const used = new Map(); // key -> Set(file)
const dynamicPrefixes = []; // {prefix, rel}

for (const file of files) {
    const code = stripComments(fs.readFileSync(file, 'utf8'));
    const rel = path.relative(ROOT, file).replace(/\\/g, '/');

    // useTranslations('ns') 上下文：返回值变量名 -> ns（短 key 补全为 ns.key）
    const nsByVar = new Map();
    const declRe = /(?:const|let|var)\s+(\w+)\s*=\s*useTranslations\(\s*['"]([\w-]+)['"]\s*\)/g;
    let m;
    while ((m = declRe.exec(code))) nsByVar.set(m[1], m[2]);

    for (const {value, dynamic} of extractStrings(code)) {
        if (dynamic) {
            if (/^[\w.-]+$/.test(value)) dynamicPrefixes.push({prefix: value, rel});
            continue;
        }
        if (!/^[\w][\w.-]*$/.test(value)) continue;
        // 带 ns 前缀的直接记；不带前缀的短 key 只可能出现在 useTranslations 上下文，下方单独处理
        if (nsPrefixRe.test(value)) mark(used, value, rel);
    }

    // useTranslations 上下文里的短 key（next-intl 语义：key 一律相对 ns）
    for (const [varName, ns] of nsByVar) {
        const callRe = new RegExp(`\\b${escapeRe(varName)}\\(\\s*['"]([\\w.-]+)['"]`, 'g');
        while ((m = callRe.exec(code))) {
            mark(used, `${ns}.${m[1]}`, rel);
        }
    }
}

function mark(map, key, rel) {
    if (!map.has(key)) map.set(key, new Set());
    map.get(key).add(rel);
}

const usedKeys = new Set(used.keys());
for (const {prefix} of dynamicPrefixes) {
    for (const k of zh) if (k.startsWith(prefix)) usedKeys.add(k);
    for (const k of en) if (k.startsWith(prefix)) usedKeys.add(k);
}

// ---------- 3. 比对 ----------
const missingZh = [...usedKeys].filter((k) => !zh.has(k)).sort();
const missingEn = [...usedKeys].filter((k) => !en.has(k)).sort();

const zhOnly = [...zh].filter((k) => !en.has(k)).sort();
const enOnly = [...en].filter((k) => !zh.has(k)).sort();

const unusedZh = [...zh].filter((k) => !usedKeys.has(k)).sort();
const unusedEn = [...en].filter((k) => !usedKeys.has(k)).sort();

const srcOf = (k) => [...(used.get(k) ?? [])].join(', ');

console.log('=== 代码用到但语言包缺失（zh 或 en 任一）:', new Set([...missingZh, ...missingEn]).size);
for (const k of new Set([...missingZh, ...missingEn])) {
    console.log(`  ${k}  (zh:${missingZh.includes(k) ? '缺' : '有'} en:${missingEn.includes(k) ? '缺' : '有'})  <- ${srcOf(k)}`);
}
console.log('\n=== zh-CN 有但 en-US 没有:', zhOnly.length);
for (const k of zhOnly) console.log('  ' + k);
console.log('\n=== en-US 有但 zh-CN 没有:', enOnly.length);
for (const k of enOnly) console.log('  ' + k);
console.log('\n=== 语言包定义但代码未使用:', unusedZh.length, '(以 zh-CN 计)');
for (const k of unusedZh) console.log('  ' + k);
console.log('\n=== 动态 key 前缀:', dynamicPrefixes.length);
for (const d of dynamicPrefixes) console.log(`  ${d.prefix}  <- ${d.rel}`);
