#!/usr/bin/env node
/**
 * 构建末步：把产物路径打包成 zip（build 脚本最后以 && 串联调用）。
 *
 * 用法：node scripts/package-zip.mjs --out <输出zip名> [--rm <打包前删除的路径>...] <路径1> [路径2 ...]
 *   --rm：zip 之前先删除指定路径（如 next 的 dist/dev 开发缓存，避免被打进交付包）。
 *   - 不存在的路径跳过并告警；全部缺失则无事可做（跳过、exit 0）。
 *   - zip 工具按序探测：7z → zip → PowerShell Compress-Archive；
 *     全部不可用（如 node:20-slim Docker 构建容器，镜像里只 COPY dist 不需要 zip）
 *     时打 warning 跳过、exit 0，不阻断构建。
 *   - 工具存在但打包失败：带出原始错误、exit 1（不吞错）。
 */
import {execFileSync} from 'node:child_process';
import {existsSync, rmSync, statSync} from 'node:fs';

const argv = process.argv.slice(2);
const outFlag = argv.indexOf('--out');
const rmFlags = [];
for (let i = argv.indexOf('--rm'); i !== -1; i = argv.indexOf('--rm', i + 1)) {
  rmFlags.push(i);
}

// 所有要剔除出打包目标的位置：--out 及其值、--rm 及其值（--rm 的值是删除路径，不是打包路径）
const excluded = new Set(rmFlags.flatMap((i) => [i, i + 1]));
if (outFlag !== -1) excluded.add(outFlag), excluded.add(outFlag + 1);

const out = outFlag !== -1 ? argv[outFlag + 1] : 'dist.zip';
const rmPaths = rmFlags.map((i) => argv[i + 1]).filter(Boolean);
const targets = argv.filter((_, i) => !excluded.has(i));
if (!out || !targets.length) {
  console.error('[package-zip] 用法：node scripts/package-zip.mjs --out <zip名> [--rm <路径>...] <路径...>');
  process.exit(1);
}

for (const p of rmPaths) {
  if (existsSync(p)) {
    rmSync(p, {recursive: true, force: true});
    console.log(`[package-zip] 已删除 ${p}`);
  }
}

const present = targets.filter((p) => existsSync(p));
const missing = targets.filter((p) => !existsSync(p));
if (missing.length) console.warn(`[package-zip] 跳过不存在的路径：${missing.join(', ')}`);
if (!present.length) {
  console.warn('[package-zip] 无可打包内容，跳过');
  process.exit(0);
}

/** 执行一次打包；工具不存在返回 false，真实失败不吞错直接退出。 */
function tryZip(cmd, args) {
  try {
    execFileSync(cmd, args, {stdio: 'inherit', shell: false});
    return true;
  } catch (err) {
    if (err && err.code === 'ENOENT') return false; // 工具缺失 → 换下一个
    console.error(`[package-zip] ${cmd} 打包失败：`, err);
    process.exit(1);
  }
}

const ok =
  tryZip('7z', ['a', '-tzip', '-mx=6', '--', out, ...present]) ||
  tryZip('zip', ['-r', '-q', out, ...present]) ||
  tryZip('powershell', [
    '-NoProfile',
    '-Command',
    `Compress-Archive -Path ${present.map((p) => `'${p}'`).join(',')} -DestinationPath '${out}' -Force`,
  ]);

if (!ok) {
  console.warn('[package-zip] 未找到可用的 zip 工具（7z/zip/powershell），跳过打包（产物目录未受影响）');
  process.exit(0);
}

const sizeMb = (statSync(out).size / 1024 / 1024).toFixed(1);
console.log(`[package-zip] 已生成 ${out}（${sizeMb} MB）`);
