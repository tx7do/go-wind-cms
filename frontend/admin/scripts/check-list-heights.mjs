/**
 * 列表页高度契约巡检（Layer 3 视觉回归网，2026-10-22）。
 *
 * 契约：每个列表页的表格卡（最后一个 .ant-pro-card）底边应贴合内容容器
 * （.content-container）底边，间距 = 容器 padding(16px) + 边框 ≈ 19px（容差 ±4）；
 * fixed 容器内不得有可见内容被 overflow 裁掉。
 *
 * 用法（先起 dev server 与后端）：
 *   pnpm exec node scripts/check-list-heights.mjs
 * 环境变量：
 *   BASE_URL        默认 http://localhost:15999
 *   ADMIN_USER/PASS 默认 admin/admin
 *   REDIS_CONTAINER 默认 redis-server（登录验证码从该容器 Redis 直读，需 docker 可用）
 *   REDIS_PASSWORD  默认 *Abcd123456（backend/app/admin/service/configs/data.yaml）
 *   HEADLESS        默认 true，设 false 可看浏览器过程
 * 退出码：全部通过 0；任一失败 1（失败清单打印在末尾）。
 */
import { execFileSync } from 'node:child_process';
import { chromium } from 'playwright-core';

const BASE_URL = process.env.BASE_URL ?? 'http://localhost:15999';
const ADMIN_USER = process.env.ADMIN_USER ?? 'admin';
const ADMIN_PASS = process.env.ADMIN_PASS ?? 'admin';
const REDIS_CONTAINER = process.env.REDIS_CONTAINER ?? 'redis-server';
const REDIS_PASSWORD = process.env.REDIS_PASSWORD ?? '*Abcd123456';

/** 有 ListTable 的列表路由（对照 src/router/modules/*.tsx；新增列表页请同步此表） */
const LIST_ROUTES = [
  '/ai/agents', '/ai/knowledge', '/ai/providers', '/ai/usage',
  '/content/categories', '/content/components', '/content/models', '/content/pages',
  '/content/posts', '/content/tags', '/content/form-submissions', '/content/experiments',
  '/log/api-audit-logs', '/log/data-access-audit-logs', '/log/login-audit-logs',
  '/log/operation-audit-logs', '/log/permission-audit-logs', '/log/policy-evaluation-logs',
  '/engagement/comments', '/tenant/members',
  '/system/access-keys', '/system/configs', '/system/dict', '/system/files', '/system/flow-rule',
  '/system/languages', '/system/login-policies', '/system/monitor-alerts', '/system/online-sessions',
  '/system/tasks', '/system/tenant-quotas', '/system/webhooks',
  '/internal-message/categories', '/internal-message/inbox', '/internal-message/messages',
  '/media/media-assets',
  '/notification/channels', '/notification/deliveries', '/notification/rules', '/notification/templates',
  '/opm/org-units', '/opm/positions', '/opm/users',
  '/permission/apis', '/permission/codes', '/permission/menus', '/permission/roles',
  '/site-setting/navigations', '/site-setting/redirects', '/site-setting/sites', '/site-setting/site-settings',
  '/social/accounts', '/social/distributions',
];

const redis = (args) =>
  execFileSync('docker', ['exec', REDIS_CONTAINER, 'redis-cli', '-a', REDIS_PASSWORD, '--no-auth-warning', ...args], {
    encoding: 'utf8',
  });

/** 从 Redis 直读验证码答案（键 gowind-cms:captcha:<id>，id 取登录页验证码请求响应头 X-Captcha-Id） */
function readCaptchaFromRedis(captchaId) {
  return redis(['GET', `gowind-cms:captcha:${captchaId}`]).trim();
}

async function login(page) {
  // 监听要先于 goto 挂上：登录页挂载即拉验证码。captchaId 在 GenerateCaptcha 的
  // JSON 响应体里；答案只在 Redis（键 gowind-cms:captcha:<id>），须 docker 可用。
  let captchaIdResolve;
  const captchaId = new Promise((resolve) => (captchaIdResolve = resolve));
  page.on('response', async (res) => {
    try {
      if (!(res.headers()['content-type'] ?? '').includes('json')) return;
      const body = await res.json();
      if (body && typeof body.captchaId === 'string' && body.captchaId) captchaIdResolve(body.captchaId);
    } catch {
      // 非 JSON 响应或流已消费，忽略
    }
  });
  await page.goto(`${BASE_URL}/auth/login`, { waitUntil: 'domcontentloaded' });
  await page.waitForTimeout(1200);
  if (!page.url().includes('/auth/login')) return; // 已有登录态
  const id = await Promise.race([captchaId, new Promise((r) => setTimeout(() => r(null), 8000))]);
  if (!id) throw new Error('未捕获 GenerateCaptcha 响应（captchaId）');
  await page.getByPlaceholder('用户名 / 邮箱 / 手机号').fill(ADMIN_USER);
  await page.getByPlaceholder('请输入密码').fill(ADMIN_PASS);
  await page.getByPlaceholder('请输入验证码').fill(readCaptchaFromRedis(id));
  await page.locator('button.login-submit-btn').click();
  await page.waitForURL((u) => !String(u).includes('/auth/login'), { timeout: 15000 });
}

/** 页内度量：卡底贴合度 + fixed 容器内容裁剪。
 *  容器两种合法形态：.content-container（ContentContainer 包装，卡底 gap≈19=16padding+边框）
 *  或裸 .page-container-content（历史页面直包，hook 把分页器钉到容器底，gap≈2）。 */
const MEASURE = () => {
  const cc = document.querySelector('.content-container');
  const container = cc ?? document.querySelector('.page-container-content');
  if (!container) return { hasContainer: false, fixed: false };
  const r = (el) => el.getBoundingClientRect();
  const out = { hasContainer: true, fixed: !!cc && cc.className.includes('--fixed'), kind: cc ? 'cc' : 'pcc' };
  const cards = document.querySelectorAll('.ant-pro-card');
  const lastCard = cards[cards.length - 1];
  if (lastCard) out.gap = Math.round(r(container).bottom - r(lastCard).bottom);
  if (out.fixed) {
    const cRect = r(container);
    const body = document.querySelector('.ant-table-body');
    let clipped = 0;
    container.querySelectorAll('*').forEach((el) => {
      if (el.children.length > 0) return;
      const tag = el.tagName;
      if (tag === 'COL' || tag === 'COLGROUP') return;
      if (body && body.contains(el)) return; // 表体内部滚动属正常
      const cls = String(el.className ?? '');
      if (cls.includes('sr-only') || el.getAttribute('aria-hidden') === 'true') return;
      const st = getComputedStyle(el);
      if (st.display === 'none' || st.visibility === 'hidden' || parseFloat(st.opacity) === 0) return;
      if (Math.round(r(el).bottom) > cRect.bottom + 6) clipped += 1;
    });
    out.clipped = clipped;
  }
  return out;
};

/** 期望 gap 基线（按容器形态） */
const GAP_BASELINE = { cc: 19, pcc: 2 };

const failures = [];
const browser = await chromium.launch({ channel: 'msedge', headless: process.env.HEADLESS !== 'false' });
try {
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
  await login(page);
  for (const route of LIST_ROUTES) {
    await page.goto(BASE_URL + route, { waitUntil: 'domcontentloaded' });
    // 等表格出现后连续采样，gap 连续两次一致才算稳定（hook 的 1.5s 兜底计时器/
    // 主从双表异步加载都会让早采样读到过渡值）
    let m = null;
    let prevGap;
    for (let i = 0; i < 14; i++) {
      await page.waitForTimeout(600);
      m = await page.evaluate(MEASURE);
      if (m.gap !== undefined && m.gap === prevGap) break;
      prevGap = m.gap;
    }
    const problems = [];
    if (!m.hasContainer) problems.push('无 .content-container/.page-container-content');
    const baseline = GAP_BASELINE[m.kind] ?? 19;
    if (m.gap !== undefined && Math.abs(m.gap - baseline) > 4) {
      problems.push(`gap=${m.gap}（期望 ${baseline}±4）`);
    }
    if (m.fixed && m.clipped > 0) problems.push(`${m.clipped} 个可见元素被容器裁剪`);
    if (problems.length > 0) {
      failures.push({ route, problems });
      console.error(`✗ ${route}  ${problems.join('；')}`);
    } else {
      console.log(`✓ ${route}  [${m.kind}] gap=${m.gap ?? '-'} clipped=${m.clipped ?? 0}`);
    }
  }
} finally {
  await browser.close();
}
if (failures.length > 0) {
  console.error(`\n${failures.length}/${LIST_ROUTES.length} 个列表页违反高度契约`);
  process.exit(1);
}
console.log(`\n全部 ${LIST_ROUTES.length} 个列表页高度契约通过`);
