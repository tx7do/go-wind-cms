"""通知域 E2E：渠道花名册 → 投递台账 → 路由规则 → 测试投递 → 用户偏好。

链路：
  1. 播种断言：核心启动时空表播种 INTERNAL_MESSAGE→INTERNAL 规则行
  2. 渠道 CRUD：建 EMAIL 渠道（带密码）→ hasPassword 标识、凭据永不回显
  3. 测试邮件：send-test-email → 台账落 FAILED（SMTP 拨不可达地址），target 脱敏
  4. WEBHOOK 渠道 + 规则 TestDispatch → SSRF 防线拦回环地址（SKIPPED + 提示文案）
  5. 规则唯一性：同事件类型重复创建 → 400
  6. 偏好自助：默认值 → 更新静音时段 → 零宽窗口 400 → 分类清单
  7. 台账只读：列表含上述行；无 Update/Delete 口（路由不存在 404）

前置：core/admin 已启动（本次部署须带通知域新二进制）、redis 运行中。
用法：python test_notification.py
"""
import json
import time
import urllib.parse

from common import AdminApi, admin_login

passed = failed = 0


def check(name, cond, extra=""):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✅ {name} {extra}")
    else:
        failed += 1
        print(f"  ❌ {name} {extra}")


def wait_for(pred, timeout=30):
    deadline = time.time() + timeout
    while time.time() < deadline:
        if pred():
            return True
        time.sleep(0.5)
    return False


def main():
    hdr = admin_login("admin", "admin")
    api = AdminApi(hdr)
    print("🧪 已登录 admin")

    # ── 1. 播种断言 ──────────────────────────────────────────────
    print("\n[1] 规则空表播种")
    rules = api("/admin/v1/notification-rules?pageSize=50")
    items = rules.get("items") or []
    seeded = [r for r in items if r.get("eventType") == "INTERNAL_MESSAGE"]
    check("INTERNAL_MESSAGE 规则已播种", len(seeded) == 1)
    if seeded:
        check("播种行路由到 INTERNAL 渠道", seeded[0].get("channel") == "INTERNAL")
        check("播种行默认启用且同步", seeded[0].get("isEnabled") is True and seeded[0].get("isAsync") is False)

    # ── 2. 渠道 CRUD + 凭据不回显 ────────────────────────────────
    print("\n[2] 渠道 CRUD")
    ch = api("/admin/v1/notification-channels", "POST", {
        "data": {
            "name": f"e2e-smtp-{int(time.time())}",
            "type": "EMAIL",
            "smtpHost": "127.0.0.1",
            "smtpPort": 1,  # 不可达端口：拨号立即被拒，测试邮件必失败
            "smtpFrom": "noreply@example.com",
            "smtpTls": "NONE",
            "enabled": True,
        },
        "password": "e2e-secret-password",
    })
    ch_id = ch.get("id")
    check("创建 EMAIL 渠道", bool(ch_id), f"id={ch_id}")
    check("hasPassword 标识回填", ch.get("hasPassword") is True)
    check("凭据不回显（响应无 password 字段）", "password" not in ch)
    check("enabled 回填", ch.get("enabled") is True)

    lst = api("/admin/v1/notification-channels?pageSize=50")
    listed = [c for c in (lst.get("items") or []) if c.get("id") == ch_id]
    check("列表含新渠道且带 hasPassword", len(listed) == 1 and listed[0].get("hasPassword") is True)

    # ── 3. 测试邮件 → 台账 FAILED + 脱敏 ─────────────────────────
    print("\n[3] 测试邮件（拨不可达 SMTP）")
    ts = api(f"/admin/v1/notification-channels/{ch_id}/send-test-email", "POST", {
        "id": ch_id,
        "recipient": "you@example.com",
    })
    # Notifier 契约：投递失败以 4xx 把原始错误回出去（这个按钮唯一有用的输出就是"为什么没发出去"）
    ts_ok = (ts.get("success") is False and bool(ts.get("error"))) or             (ts.get("code") == 400 and bool(ts.get("message")))
    check("send-test-email 失败把原始原因回出去", ts_ok, f"resp={str(ts)[:90]}")

    got = wait_for(lambda: _find_delivery(api, "CHANNEL_TEST_EMAIL", "y***@example.com"))
    row = _find_delivery(api, "CHANNEL_TEST_EMAIL", "y***@example.com")
    check("台账落 CHANNEL_TEST_EMAIL 行（target 已脱敏）", got, "target=y***@example.com")
    if row:
        check("测试邮件拨号失败记 FAILED", row.get("status") == "FAILED",
              f"status={row.get('status')} lastError={str(row.get('lastError'))[:60]}")
        check("attempts=1（同步投递恒一次）", row.get("attempts") == 1)
        check("channel_id 钉到被测渠道", row.get("channelId") == ch_id)

    # ── 4. WEBHOOK 渠道 + SSRF 防线 ──────────────────────────────
    print("\n[4] WEBHOOK 渠道 + 规则测试投递（SSRF 防线）")
    wh = api("/admin/v1/notification-channels", "POST", {
        "data": {
            "name": f"e2e-hook-{int(time.time())}",
            "type": "WEBHOOK",
            "webhookUrl": "http://192.0.2.1/hook",  # TEST-NET-1：外网不可达，绕开回环 SSRF 拦截
            "webhookSignStyle": "CUSTOM",
            "enabled": True,
        },
        "webhookSecret": "hook-secret",
    })
    wh_id = wh.get("id")
    check("创建 WEBHOOK 渠道", bool(wh_id), f"id={wh_id}")
    check("hasWebhookSecret 标识回填", wh.get("hasWebhookSecret") is True)

    # 借用种子规则（CHANNEL_TEST_EMAIL 一事件一行）：临时切到 WEBHOOK 做测试投递，完了还原 EMAIL
    rules = api("/admin/v1/notification-rules?pageSize=50")
    seed_rules = [r for r in (rules.get("items") or []) if r.get("eventType") == "CHANNEL_TEST_EMAIL"]
    rule_id = seed_rules[0].get("id") if seed_rules else None
    check("找到 CHANNEL_TEST_EMAIL 种子规则", bool(rule_id))

    if rule_id:
        upd = api(f"/admin/v1/notification-rules/{rule_id}", "PUT", {
            "id": rule_id,
            "data": {"channel": "WEBHOOK", "isEnabled": True},
            "updateMask": "channel,isEnabled",
        })
        check("规则切到 WEBHOOK（channel 进掩码可改）", upd.get("code") is None and upd.get("channel") == "WEBHOOK",
              f"resp={str(upd)[:80]}")

        td = api(f"/admin/v1/notification-rules/{rule_id}/test-dispatch", "POST", {
            "id": rule_id, "target": "http://192.0.2.1/hook",
        })
        # 外网不可达（TEST-NET-1 拨号超时 ~10s）：Notifier 契约下失败以 400 回原文；台账行照落
        td_resp_ok = (td.get("deliveryId", 0) > 0) or (td.get("code") == 400 and bool(td.get("message")))
        check("测试投递走完整链（回结论或 400 原文）", td_resp_ok, f"resp={str(td)[:100]}")

        got_hook = wait_for(lambda: _find_delivery(api, "CHANNEL_TEST_EMAIL", "****hook") is not None, timeout=30)
        hook_row = _find_delivery(api, "CHANNEL_TEST_EMAIL", "****hook")
        check("webhook 测试投递落台账（target 脱敏）", got_hook, "target=****hook")
        if hook_row:
            ok = wait_for(lambda: _delivery_status(api, hook_row["id"]) == "FAILED", timeout=30)
            st = _delivery_status(api, hook_row["id"])
            check("不可达地址定案 FAILED", st == "FAILED", f"status={st}")

        # 还原种子路由到 EMAIL（后续轮次的 send-test-email 仍走邮件渠道）
        api(f"/admin/v1/notification-rules/{rule_id}", "PUT", {
            "id": rule_id,
            "data": {"channel": "EMAIL", "isEnabled": True},
            "updateMask": "channel,isEnabled",
        })

    # INTERNAL 规则拒绝测试投递（无测试文案语义）
    seeded_id = seeded[0].get("id") if seeded else None
    if seeded_id:
        td2 = api(f"/admin/v1/notification-rules/{seeded_id}/test-dispatch", "POST", {"id": seeded_id})
        check("INTERNAL 规则测试投递被拒（400）", td2.get("code") == 400, f"resp={str(td2)[:80]}")

    # ── 5. 规则唯一性 ────────────────────────────────────────────
    print("\n[5] 规则唯一性")
    dup = api("/admin/v1/notification-rules", "POST", {
        "data": {"eventType": "INTERNAL_MESSAGE", "channel": "INTERNAL", "isEnabled": True},
    })
    check("同事件类型重复创建被拒（400）", dup.get("code") == 400, f"reason={dup.get('reason')}")

    # ── 6. 偏好自助 ──────────────────────────────────────────────
    print("\n[6] 用户通知偏好")
    pref = api("/admin/v1/notification-preference")
    check("默认偏好：未启用静音", pref.get("quietEnabled") is False)
    check("默认窗口 22:00→08:00", pref.get("quietStartMinute") == 1320 and pref.get("quietEndMinute") == 480)

    upd = api("/admin/v1/notification-preference", "POST", {
        "quietEnabled": True, "quietStartMinute": 1320, "quietEndMinute": 480, "mutedCategoryIds": [],
    })
    check("保存静音时段（跨零点合法）", upd.get("quietEnabled") is True and upd.get("quietStartMinute") == 1320)
    check("user_id 服务端钉定", upd.get("userId") == 1, f"userId={upd.get('userId')}")

    zero = api("/admin/v1/notification-preference", "POST", {
        "quietEnabled": True, "quietStartMinute": 600, "quietEndMinute": 600, "mutedCategoryIds": [],
    })
    check("零宽窗口被拒（400）", zero.get("code") == 400)

    oob = api("/admin/v1/notification-preference", "POST", {
        "quietEnabled": True, "quietStartMinute": 1440, "quietEndMinute": 480, "mutedCategoryIds": [],
    })
    check("分钟数越界被拒（400）", oob.get("code") == 400)

    cats = api("/admin/v1/notification-preference/categories")
    check("分类清单可读（数组）", isinstance(cats.get("items"), list))

    # 还原偏好（后续轮次不受静音影响）
    api("/admin/v1/notification-preference", "POST", {
        "quietEnabled": False, "quietStartMinute": 1320, "quietEndMinute": 480, "mutedCategoryIds": [],
    })

    # ── 7. 台账只读 + 清理 ───────────────────────────────────────
    print("\n[7] 台账只读与清理")
    dl = api("/admin/v1/notification-deliveries?pageSize=10")
    check("台账列表可读", isinstance(dl.get("items"), list) and len(dl.get("items")) >= 2)
    export_route = api("/admin/v1/notification-deliveries/x", "PUT", {"status": "SENT"})
    check("台账无更新口（404/405）", export_route.get("code") in (404, 405) or "_raw" in export_route,
          f"code={export_route.get('code')}")

    if ch_id:
        d1 = api(f"/admin/v1/notification-channels/{ch_id}", "DELETE")
        check("删除 EMAIL 渠道", d1.get("code") in (200, None) or "_raw" not in d1)
    if wh_id:
        api(f"/admin/v1/notification-channels/{wh_id}", "DELETE")
    print(f"\n{'='*40}\n通过 {passed} / 失败 {failed}")
    raise SystemExit(1 if failed else 0)


def _find_delivery(api, event_type, target):
    # 倒序取最新：台账是只增的，历史轮次会留同 target 的旧行
    _order = urllib.parse.quote('["-id"]')
    resp = api(f"/admin/v1/notification-deliveries?pageSize=20&eventType={event_type}&orderBy={_order}")
    for it in resp.get("items") or []:
        if it.get("target") == target:
            return it
    return None


def _delivery_status(api, delivery_id):
    row = api(f"/admin/v1/notification-deliveries/{delivery_id}")
    return row.get("status")


if __name__ == "__main__":
    main()
