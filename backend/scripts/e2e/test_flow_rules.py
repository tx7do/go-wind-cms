"""自动化 Flow 规则 E2E：事件触发 → 条件判定 → webhook 动作执行。

链路：
  1. 建 mock 接收端（:6901，捕获 POST body 与签名头）
  2. 建规则：post.published + 条件 POST_STATUS_PUBLISHED + webhook 动作
  3. 文章草稿态更新（post.updated）→ 不触发
  4. 文章转 PUBLISHED（post.published）→ 触发：body/envelope、HMAC 签名、match_count+1
  5. 条件不满足的事件（文章回 DRAFT 的 post.updated）→ 不触发
  6. 通知动作：规则触发 → 站内信落库 + 收件人收件箱
  7. 延时动作：webhook → delay 5s → webhook，第二跳晚于第一跳 ≥5s
  8. 定时触发：triggerCron 每分钟 → envelope event=schedule/entity=0

前置：core/admin 已启动（taskScheduler 可用，redis 运行中）。
用法：python test_flow_rules.py
"""
import hashlib
import hmac as hmac_mod
import json
import threading
import time
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

from common import AdminApi, admin_login

RECEIVER_PORT = 6901
RECEIVER_SECRET = "flow-e2e-secret"

passed = failed = 0
received = []  # [{path, body, signature}]


def check(name, cond, extra=""):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✅ {name} {extra}")
    else:
        failed += 1
        print(f"  ❌ {name} {extra}")


class ReceiverHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(length).decode()
        received.append({
            "path": self.path,
            "body": body,
            "signature": self.headers.get("X-Webhook-Signature", ""),
            "event": self.headers.get("X-Flow-Event", ""),
            "ts": time.time(),
        })
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"ok")

    def log_message(self, *args):
        pass


def start_receiver():
    server = ThreadingHTTPServer(("127.0.0.1", RECEIVER_PORT), ReceiverHandler)
    threading.Thread(target=server.serve_forever, daemon=True).start()
    return server


def wait_for(pred, timeout=60):
    deadline = time.time() + timeout
    while time.time() < deadline:
        if pred():
            return True
        time.sleep(0.5)
    return False


def main():
    receiver = start_receiver()
    print("🧪 mock 接收端已启动 :6901")
    hdr = admin_login("admin", "admin")
    api = AdminApi(hdr)
    TS = int(time.time())

    # ===== 1. 建规则 =====
    rule = api("/admin/v1/flow-rules", "POST", {"data": {
        "name": f"e2e-flow-{TS}",
        "triggerEvent": "post.published",
        "conditionStatus": "POST_STATUS_PUBLISHED",
        "actions": json.dumps([{
            "type": "webhook",
            "url": f"http://127.0.0.1:{RECEIVER_PORT}/flow",
            "secret": RECEIVER_SECRET,
        }]),
        "enabled": True,
    }})
    rule_id = rule.get("id")
    check("建规则", bool(rule_id), f"id={rule_id} resp={json.dumps(rule, ensure_ascii=False)[:100]}")

    # 动作 JSON 校验拒绝（坏动作）
    bad = api("/admin/v1/flow-rules", "POST", {"data": {
        "name": f"e2e-flow-bad-{TS}",
        "triggerEvent": "post.published",
        "actions": json.dumps([{"type": "unknown"}]),
    }})
    check("坏动作类型被拒", "_raw" in bad or bad.get("code") == 400,
          f"resp={json.dumps(bad, ensure_ascii=False)[:80]}")

    # ===== 2. 建文章（DRAFT）并更新 → post.updated，条件是 PUBLISHED，不应触发 =====
    before = len(received)
    post = api("/admin/v1/posts", "POST", {"data": {
        "status": "DRAFT",
        "translations": [{"languageCode": "zh-CN", "title": f"flow e2e {TS}",
                          "content": "<p>x</p>"}],
    }})
    post_id = post.get("id")
    check("建草稿文章", bool(post_id), f"id={post_id}")
    api(f"/admin/v1/posts/{post_id}", "PUT", {"data": {"title": f"flow e2e {TS} v2"}})
    time.sleep(3)
    check("草稿更新不触发 published 规则", len(received) == before,
          f"received={len(received)}")

    # ===== 3. 转 PUBLISHED → 触发 =====
    # 直接 PUT status 会被发布审批工作流拦截（status 不落库），发布走审核链：
    # submit（DRAFT→审核中）→ approve（→PUBLISHED，ApproveReview 补发 published 事件）
    api(f"/admin/v1/posts/{post_id}/review/submit", "POST", {"postId": post_id})
    api(f"/admin/v1/posts/{post_id}/review/approve", "POST", {"postId": post_id})
    ok = wait_for(lambda: len(received) > before)
    check("发布触发 webhook", ok, f"received={len(received)}")
    if ok and len(received) > before:
        hit = received[-1]
        body = json.loads(hit["body"])
        check("envelope 事件正确", body.get("event") == "post.published", f"{body.get('event')}")
        check("envelope 实体正确", body.get("entity_id") == post_id,
              f"{body.get('entity_id')} vs {post_id}")
        check("envelope 规则正确", body.get("rule_id") == rule_id,
              f"{body.get('rule_id')} vs {rule_id}")
        expect_sig = "sha256=" + hmac_mod.new(
            RECEIVER_SECRET.encode(), hit["body"].encode(), hashlib.sha256).hexdigest()
        check("HMAC 签名验证", hit["signature"] == expect_sig,
              f"sig={hit['signature'][:24]}…")

    # ===== 4. match_count 递增 =====
    ok = wait_for(lambda: (lambda r: int((r.get("matchCount") or 0)) >= 1)(
        api(f"/admin/v1/flow-rules/{rule_id}", "GET")), timeout=15)
    r = api(f"/admin/v1/flow-rules/{rule_id}", "GET")
    check("命中次数递增", int(r.get("matchCount") or 0) >= 1,
          f"matchCount={r.get('matchCount')}")

    # ===== 5. 回 DRAFT → post.updated，条件不满足不触发 =====
    before = len(received)
    api(f"/admin/v1/posts/{post_id}/review/reject", "POST", {"postId": post_id, "reason": "e2e"})
    time.sleep(3)
    check("条件不满足不触发", len(received) == before, f"received={len(received)}")

    # ===== 6. 通知动作：规则触发 → 站内信落库 + 收件人收件箱 =====
    title_mark = f"flow-notify-{TS}"
    nrule = api("/admin/v1/flow-rules", "POST", {"data": {
        "name": f"e2e-flow-notify-{TS}",
        "triggerEvent": "post.published",
        "conditionStatus": "",
        "actions": json.dumps([{
            "type": "notification",
            "title": title_mark + " {event}",
            "content": "post " + str(0) + " {entity_id} fired",
            "target_user_ids": [1],
        }]),
        "enabled": True,
    }})
    nrule_id = nrule.get("id")
    check("建通知规则", bool(nrule_id), f"id={nrule_id}")

    # 校验：通知动作缺收件人被拒
    bad_n = api("/admin/v1/flow-rules", "POST", {"data": {
        "name": f"e2e-flow-notify-bad-{TS}",
        "triggerEvent": "post.published",
        "actions": json.dumps([{"type": "notification", "title": "t", "content": "c"}]),
    }})
    check("通知动作缺收件人被拒", bad_n.get("code") == 400,
          f"code={bad_n.get('code')}")

    npost = api("/admin/v1/posts", "POST", {"data": {
        "status": "DRAFT",
        "translations": [{"languageCode": "zh-CN", "title": f"flow notify e2e {TS}",
                          "content": "<p>x</p>"}],
    }})
    npost_id = npost.get("id")
    api(f"/admin/v1/posts/{npost_id}/review/submit", "POST", {"postId": npost_id})
    api(f"/admin/v1/posts/{npost_id}/review/approve", "POST", {"postId": npost_id})
    ok = wait_for(lambda: any(
        title_mark in (m.get("title") or "")
        for m in (api("/admin/v1/internal-message/messages", "GET").get("items") or [])
    ), timeout=90)
    check("通知动作落库（消息标题命中模板）", ok, f"mark={title_mark}")

    api(f"/admin/v1/flow-rules/{nrule_id}", "DELETE")
    api(f"/admin/v1/posts/{npost_id}", "DELETE")

    # ===== 7. 延时动作：webhook → delay 5s → webhook =====
    drule = api("/admin/v1/flow-rules", "POST", {"data": {
        "name": f"e2e-flow-delay-{TS}",
        "triggerEvent": "post.published",
        "conditionStatus": "",
        "actions": json.dumps([
            {"type": "webhook", "url": f"http://127.0.0.1:{RECEIVER_PORT}/flow-fast"},
            {"type": "delay", "seconds": 5},
            {"type": "webhook", "url": f"http://127.0.0.1:{RECEIVER_PORT}/flow-slow"},
        ]),
        "enabled": True,
    }})
    drule_id = drule.get("id")
    check("建延时规则", bool(drule_id), f"id={drule_id}")

    bad_d = api("/admin/v1/flow-rules", "POST", {"data": {
        "name": f"e2e-flow-delay-bad-{TS}",
        "triggerEvent": "post.published",
        "actions": json.dumps([{"type": "delay", "seconds": 0}]),
    }})
    check("延时秒数非法被拒", bad_d.get("code") == 400, f"code={bad_d.get('code')}")

    dpost = api("/admin/v1/posts", "POST", {"data": {
        "status": "DRAFT",
        "translations": [{"languageCode": "zh-CN", "title": f"flow delay e2e {TS}",
                          "content": "<p>x</p>"}],
    }})
    dpost_id = dpost.get("id")
    api(f"/admin/v1/posts/{dpost_id}/review/submit", "POST", {"postId": dpost_id})
    api(f"/admin/v1/posts/{dpost_id}/review/approve", "POST", {"postId": dpost_id})

    fast_idx = len(received)
    ok = wait_for(lambda: any(r["path"] == "/flow-fast" for r in received[fast_idx:]),
                  timeout=90)
    check("延时第一跳（立即 webhook）到达", ok)
    fast_ts = next((r["ts"] for r in received[fast_idx:] if r["path"] == "/flow-fast"), None)

    slow_idx = fast_idx
    ok = wait_for(lambda: any(r["path"] == "/flow-slow" for r in received[slow_idx:]),
                  timeout=120)
    slow_ts = next((r["ts"] for r in received[slow_idx:] if r["path"] == "/flow-slow"), None)
    check("延时第二跳（续跑 webhook）到达", ok and slow_ts is not None)
    if fast_ts and slow_ts:
        gap = slow_ts - fast_ts
        check("第二跳晚于第一跳 ≥5s（delay 生效）", gap >= 5, f"gap={gap:.1f}s")

    # ===== 8. 定时触发：triggerCron 每分钟 =====
    bad_c = api("/admin/v1/flow-rules", "POST", {"data": {
        "name": f"e2e-flow-cron-bad-{TS}",
        "triggerEvent": "post.published",
        "triggerCron": "not-a-cron",
        "actions": json.dumps([{"type": "webhook",
                                "url": f"http://127.0.0.1:{RECEIVER_PORT}/flow-cron"}]),
        "enabled": True,
    }})
    check("非法 cron 被拒", bad_c.get("code") == 400, f"code={bad_c.get('code')}")

    crule = api("/admin/v1/flow-rules", "POST", {"data": {
        "name": f"e2e-flow-cron-{TS}",
        "triggerEvent": "post.published",
        "triggerCron": "* * * * *",
        "actions": json.dumps([{"type": "webhook",
                                "url": f"http://127.0.0.1:{RECEIVER_PORT}/flow-cron"}]),
        "enabled": True,
    }})
    crule_id = crule.get("id")
    check("建定时规则", bool(crule_id), f"id={crule_id}")

    cron_idx = len(received)
    ok = wait_for(lambda: any(r["path"] == "/flow-cron" for r in received[cron_idx:]),
                  timeout=150)
    check("定时规则触发（≤下一分钟+队列滞后）", ok)
    if ok:
        hit = next(r for r in received[cron_idx:] if r["path"] == "/flow-cron")
        body = json.loads(hit["body"])
        check("envelope event=schedule", body.get("event") == "schedule",
              f"event={body.get('event')}")
        check("envelope entity_id=0", body.get("entity_id") == 0,
              f"entity={body.get('entity_id')}")
        check("envelope rule_id 命中", body.get("rule_id") == crule_id,
              f"rule={body.get('rule_id')}")

    api(f"/admin/v1/flow-rules/{drule_id}", "DELETE")
    api(f"/admin/v1/posts/{dpost_id}", "DELETE")
    api(f"/admin/v1/flow-rules/{crule_id}", "DELETE")

    # ===== 清理 =====
    api(f"/admin/v1/flow-rules/{rule_id}", "DELETE")
    api(f"/admin/v1/posts/{post_id}", "DELETE")
    receiver.shutdown()
    print(f"\n通过 {passed} / 失败 {failed}")
    raise SystemExit(1 if failed else 0)


if __name__ == "__main__":
    main()
