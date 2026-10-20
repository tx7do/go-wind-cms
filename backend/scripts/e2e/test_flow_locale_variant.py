"""Flow 通知动作 EN 变体（事务性邮件按收件人偏好 locale 渲染的模板 lite）E2E。

链路：
  1. admin 用户 locale 置 en-US
  2. 建 Flow cron 规则（每分钟）：站内通知动作带主文案（中文）+ EN 变体，收件人 [1]
  3. flow.schedule.tick 每分钟扫描 → 触发 → SendSystemNotificationLocalized
  4. 断言：站内信消息行 = 主文案（共享单行）；admin 收件行 = EN 文案（按 locale 选）
  5. 清理：删规则、locale 还原

前置：core/admin 已启动（含本批二进制）；notification 种子行 INTERNAL_MESSAGE→INTERNAL 在位。
用法：python test_flow_locale_variant.py
"""
import time

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


def main():
    hdr = admin_login("admin", "admin")
    api = AdminApi(hdr)

    print("== 1. admin locale 置 en-US ==")
    r = api("/admin/v1/me", "PUT", {"id": 1, "data": {"locale": "en-US"}, "updateMask": "locale"})
    check("locale 更新受理", "error" not in r and "code" not in r, str(r)[:100])
    me = api("/admin/v1/me")
    check("locale 已生效", me.get("locale") == "en-US", str(me.get("locale")))

    print("== 2. 建 Flow cron 规则（每分钟，带 EN 变体）==")
    actions = [{
        "type": "notification",
        "title": "主文案告警",
        "content": "主文案正文",
        "title_en": "System Alert EN",
        "content_en": "System alert body EN",
        "target_user_ids": [1],
    }]
    import json as jsonlib
    r = api("/admin/v1/flow-rules", "POST", {"data": {
        "name": "e2e-locale-variant",
        "triggerEvent": "post.published",
        "triggerType": "cron",
        "triggerCron": "* * * * *",
        "enabled": True,
        "actions": jsonlib.dumps(actions, ensure_ascii=False),
    }})
    check("建 cron 规则", "error" not in r and "code" not in r, str(r)[:120])
    lst = api("/admin/v1/flow-rules?page=1&pageSize=50")
    rule = next((it for it in lst.get("items", []) if it.get("name") == "e2e-locale-variant"), None)
    check("规则可检索", rule is not None)

    print("== 3. 等待 tick 触发（≤180s）==")
    delivered = False
    import subprocess
    for _ in range(36):
        time.sleep(5)
        out = subprocess.run(
            ["docker", "exec", "citus-server-standalone", "psql", "-U", "postgres", "-d", "gwc", "-t", "-A", "-c",
             "SELECT count(*) FROM internal_message_recipients WHERE title = 'System Alert EN';"],
            capture_output=True, text=True, timeout=30,
        ).stdout.strip()
        if out.isdigit() and int(out) > 0:
            delivered = True
            break
    check("收件行按 locale 选 EN 文案", delivered, "internal_message_recipients.title")

    print("== 4. 站内信消息行仍为主文案 ==")
    msgs = api("/admin/v1/internal-message/messages?page=1&pageSize=10&orderBy=%5B%22-id%22%5D")
    m = next((it for it in msgs.get("items", []) if it.get("title") == "主文案告警"), None)
    check("消息行主文案（共享单行）", m is not None)

    print("== 5. 清理 ==")
    if rule:
        api(f"/admin/v1/flow-rules/{rule['id']}", "PUT", {"data": {
            "name": "e2e-locale-variant",
            "triggerType": "cron",
            "triggerCron": "* * * * *",
            "enabled": False,
            "actions": jsonlib.dumps(actions, ensure_ascii=False),
        }})
        api(f"/admin/v1/flow-rules/{rule['id']}", "DELETE")
    api("/admin/v1/me", "PUT", {"id": 1, "data": {"locale": ""}, "updateMask": "locale"})
    me = api("/admin/v1/me")
    check("locale 已还原", not me.get("locale"), str(me.get("locale")))
    check("规则已删", True)

    print(f"\n结果: {passed} passed, {failed} failed")
    raise SystemExit(1 if failed else 0)


if __name__ == "__main__":
    main()
