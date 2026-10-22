/**
 * 生成 tabBar PNG 图标（纯 Node 实现，无第三方依赖）。
 *
 * 用法：node scripts/gen-tab-icons.mjs
 * 输出：src/assets/tabbar/{home,post,favorites,user}{,-active}.png
 *
 * 小程序 tabBar 图标要求：PNG，建议 81x81，≤40KB；H5 端同样适用。
 * 图形用 SDF（有向距离场）定义，4x4 超采样抗锯齿。
 */
import {deflateSync} from 'node:zlib';
import {mkdirSync, writeFileSync} from 'node:fs';
import {dirname, join} from 'node:path';
import {fileURLToPath} from 'node:url';

const SIZE = 81; // 输出尺寸
const SS = 4; // 每轴超采样倍数

// ---------- SDF 基元（设计坐标 81x81） ----------
const circle = (cx, cy, r) => (x, y) => Math.hypot(x - cx, y - cy) - r;
const segment = (ax, ay, bx, by, r) => (x, y) => {
  const pax = x - ax, pay = y - ay, bax = bx - ax, bay = by - ay;
  const h = Math.max(0, Math.min(1, (pax * bax + pay * bay) / (bax * bax + bay * bay)));
  return Math.hypot(pax - bax * h, pay - bay * h) - r;
};
const roundRect = (cx, cy, hw, hh, r) => (x, y) => {
  const qx = Math.abs(x - cx) - hw + r, qy = Math.abs(y - cy) - hh + r;
  return Math.min(Math.max(qx, qy), 0) + Math.hypot(Math.max(qx, 0), Math.max(qy, 0)) - r;
};
const triangle = (x1, y1, x2, y2, x3, y3) => (x, y) => {
  const d = (ax, ay, bx, by) => {
    const pax = x - ax, pay = y - ay, bax = bx - ax, bay = by - ay;
    const h = Math.max(0, Math.min(1, (pax * bax + pay * bay) / (bax * bax + bay * bay)));
    return Math.hypot(pax - bax * h, pay - bay * h);
  };
  const edge = Math.min(d(x1, y1, x2, y2), d(x2, y2, x3, y3), d(x3, y3, x1, y1));
  const s1 = (x2 - x1) * (y - y1) - (y2 - y1) * (x - x1);
  const s2 = (x3 - x2) * (y - y2) - (y3 - y2) * (x - x2);
  const s3 = (x1 - x3) * (y - y3) - (y1 - y3) * (x - x3);
  const inside = (s1 >= 0 && s2 >= 0 && s3 >= 0) || (s1 <= 0 && s2 <= 0 && s3 <= 0);
  return inside ? -edge : edge;
};
const union = (...fs) => (x, y) => Math.min(...fs.map(f => f(x, y)));
const subtract = (a, b) => (x, y) => Math.max(a(x, y), -b(x, y));

// ---------- 四个图标（CMS：首页 / 文章 / 收藏 / 我的） ----------
const C = SIZE / 2;

// 旋转辅助：把判定点绕原点反转 deg，即等价于图形被旋转 deg
const rotate = (f, deg) => (x, y) => {
  const a = (deg * Math.PI) / 180;
  const c = Math.cos(a), s = Math.sin(a);
  return f(x * c + y * s, -x * s + y * c);
};

const icons = {
  // 首页：屋顶三角 + 房身 + 镂空门
  home: subtract(
    union(
      triangle(11, 41, C, 13, 70, 41),
      roundRect(C, 52, 24, 19, 3),
    ),
    roundRect(C, 62, 7.5, 10, 3),
  ),
  // 文章：圆角文档 + 三行镂空正文
  post: subtract(
    roundRect(C, C, 21, 26, 4),
    union(
      segment(32, 33, 52, 33, 2.5),
      segment(32, 43, 52, 43, 2.5),
      segment(32, 53, 45, 53, 2.5),
    ),
  ),
  // 收藏：心形（双圆 + 45° 菱形）
  favorites: union(
    circle(C - 10, 33, 11),
    circle(C + 10, 33, 11),
    rotate(roundRect(C, 44, 13.5, 13.5, 3), 45),
  ),
  // 我的：头 + 肩
  user: union(
    circle(C, 27, 11.5),
    roundRect(C, 61, 17, 12.5, 10),
  ),
};

// ---------- PNG 编码 ----------
const crcTable = (() => {
  const t = new Int32Array(256);
  for (let n = 0; n < 256; n++) {
    let c = n;
    for (let k = 0; k < 8; k++) c = c & 1 ? 0xedb88320 ^ (c >>> 1) : c >>> 1;
    t[n] = c;
  }
  return t;
})();
function crc32(buf) {
  let c = 0xffffffff;
  for (const b of buf) c = crcTable[(c ^ b) & 0xff] ^ (c >>> 8);
  return (c ^ 0xffffffff) >>> 0;
}
function chunk(type, data) {
  const len = Buffer.alloc(4);
  len.writeUInt32BE(data.length);
  const body = Buffer.concat([Buffer.from(type, 'ascii'), data]);
  const crc = Buffer.alloc(4);
  crc.writeUInt32BE(crc32(body));
  return Buffer.concat([len, body, crc]);
}
function encodePNG(rgba, w, h) {
  const ihdr = Buffer.alloc(13);
  ihdr.writeUInt32BE(w, 0);
  ihdr.writeUInt32BE(h, 4);
  ihdr[8] = 8; // bit depth
  ihdr[9] = 6; // RGBA
  const raw = Buffer.alloc((w * 4 + 1) * h);
  for (let y = 0; y < h; y++) {
    raw[y * (w * 4 + 1)] = 0; // filter: none
    rgba.copy(raw, y * (w * 4 + 1) + 1, y * w * 4, (y + 1) * w * 4);
  }
  return Buffer.concat([
    Buffer.from([0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a]),
    chunk('IHDR', ihdr),
    chunk('IDAT', deflateSync(raw)),
    chunk('IEND', Buffer.alloc(0)),
  ]);
}

// ---------- 渲染 ----------
function hexRGB(hex) {
  return [parseInt(hex.slice(1, 3), 16), parseInt(hex.slice(3, 5), 16), parseInt(hex.slice(5, 7), 16)];
}
function render(sdf, hex) {
  const [R, G, B] = hexRGB(hex);
  const rgba = Buffer.alloc(SIZE * SIZE * 4);
  for (let py = 0; py < SIZE; py++) {
    for (let px = 0; px < SIZE; px++) {
      let hit = 0;
      for (let sy = 0; sy < SS; sy++) {
        for (let sx = 0; sx < SS; sx++) {
          const x = px + (sx + 0.5) / SS;
          const y = py + (sy + 0.5) / SS;
          if (sdf(x, y) < 0) hit++;
        }
      }
      const i = (py * SIZE + px) * 4;
      rgba[i] = R;
      rgba[i + 1] = G;
      rgba[i + 2] = B;
      rgba[i + 3] = Math.round((hit / (SS * SS)) * 255);
    }
  }
  return rgba;
}

const outDir = join(dirname(fileURLToPath(import.meta.url)), '..', 'src', 'assets', 'tabbar');
mkdirSync(outDir, {recursive: true});

const NORMAL = '#7a7e83';
const ACTIVE = '#006be6';
for (const [name, sdf] of Object.entries(icons)) {
  for (const [suffix, color] of [['', NORMAL], ['-active', ACTIVE]]) {
    const png = encodePNG(render(sdf, color), SIZE, SIZE);
    const file = join(outDir, `${name}${suffix}.png`);
    writeFileSync(file, png);
    console.log('written', file, `${png.length}B`);
  }
}
