// 读取 .env.development 后再启动 nuxt dev。
// Nuxt CLI 解析 dev server 端口（读 process.env.PORT）早于 --dotenv 加载，
// 因此 PORT 只能从进程环境变量读取——在这里提前注入，使 dev server 端口可在 .env.development 中配置。
// 若首选端口被占用，则向上顺延直至找到可用端口。
import { readFileSync } from 'node:fs';
import net from 'node:net';

const envPath = new URL('../.env.development', import.meta.url);
for (const line of readFileSync(envPath, 'utf8').split(/\r?\n/)) {
  const entry = line.trim();
  if (!entry || entry.startsWith('#')) continue;
  const eq = entry.indexOf('=');
  if (eq <= 0) continue;
  const key = entry.slice(0, eq).trim();
  let value = entry.slice(eq + 1).trim();
  if (value.length > 1 && ((value.startsWith('"') && value.endsWith('"')) || (value.startsWith("'") && value.endsWith("'")))) {
    value = value.slice(1, -1);
  }
  process.env[key] ??= value;
}

// 探测端口在指定地址上可否绑定。nuxt dev 默认 hostname 0.0.0.0，但实际
// 监听可能落到 IPv6 (::)，故 IPv4 与 IPv6 任意一端占用即视为不可用。
function tryBind(port, host) {
  return new Promise((resolve) => {
    const server = net.createServer();
    server.unref();
    server.on('error', () => resolve(false));
    server.listen(port, host, () => {
      server.close(() => resolve(true));
    });
  });
}
async function isPortAvailable(port) {
  return (await tryBind(port, '0.0.0.0')) && (await tryBind(port, '::'));
}

const MAX_PORT = 65535;
const requestedPort = Number.parseInt(process.env.PORT, 10);
let port = Number.isInteger(requestedPort) && requestedPort >= 1 && requestedPort <= MAX_PORT ? requestedPort : 3000;
while (!(await isPortAvailable(port))) {
  if (port >= MAX_PORT) {
    console.error(`[dev.mjs] 从 ${requestedPort} 顺延至 ${MAX_PORT} 仍无可用端口，退出`);
    process.exit(1);
  }
  port++;
}
if (Number.isInteger(requestedPort) && port !== requestedPort) {
  console.warn(`[dev.mjs] 端口 ${requestedPort} 被占用，改用 ${port}`);
}
process.argv = [process.argv[0], 'nuxt', 'dev', '--dotenv', '.env.development', '--port', String(port)];
await import('../node_modules/nuxt/bin/nuxt.mjs');
