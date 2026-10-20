"""登录策略执行链 + 系统参数域 E2E。

链路：
  1. 登录策略：黑名单 IP（TEST-NET 不命中）→ 放行；IP 白名单（不含本机）→ 拒绝；
     删除白名单 → 恢复；REGION 黑名单（北京市）→ 局域网归属不命中放行；
     REGION 白名单（局域网）→ geolite 解析命中放行；非法 IP 值 → 400
  2. 系统参数：内置三键已播种；创建自定义键 → 更新值 → 重复键 400 →
     删内置键 400 → 删自定义键成功
  3. 清理：删除全部测试策略/参数行

前置：core/admin 已启动（含本批二进制；本机登录 IP=127.0.0.1，geolite 解析省=局域网）。
用法：python test_login_policy_config.py
"""
import json

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


def try_login():
    try:
        admin_login("admin", "admin")
        return True, ""
    except AssertionError as e:
        return False, str(e)


def main():
    hdr = admin_login("admin", "admin")
    api = AdminApi(hdr)

    print("== 1. 登录策略执行链 ==")
    # 基线：当前无生效策略，登录应成功
    ok, _ = try_login()
    check("基线登录成功", ok)

    # 1a. 黑名单 TEST-NET IP（不命中本机 127.0.0.1）
    r = api("/admin/v1/login-policies", "POST", {
        "data": {"type": "BLACKLIST", "method": "IP", "value": "203.0.113.77", "reason": "e2e 黑名单"},
    })
    check("建黑名单 IP 策略", "error" not in r and "code" not in r, str(r)[:100])
    ok, msg = try_login()
    check("黑名单不命中 → 放行", ok, msg[:80])

    # 1b. IP 白名单（不含本机）→ 全局拒
    r = api("/admin/v1/login-policies", "POST", {
        "data": {"type": "WHITELIST", "method": "IP", "value": "10.99.99.99", "reason": "e2e 白名单"},
    })
    check("建 IP 白名单", "error" not in r and "code" not in r, str(r)[:100])
    ok, msg = try_login()
    check("白名单未命中 → 拒绝", not ok and "whitelist" in msg, msg[:100])

    # 删除白名单（用旧 token：token 未过期仍然有效）
    lst = api("/admin/v1/login-policies?page=1&pageSize=50")
    wl = [it for it in lst.get("items", []) if it.get("type") == "WHITELIST" and it.get("method") == "IP"]
    for it in wl:
        api(f"/admin/v1/login-policies/{it['id']}", "DELETE")
    ok, msg = try_login()
    check("删除白名单 → 恢复", ok, msg[:80])

    # 1c. REGION 黑名单（北京市；本机归属=局域网，不命中）
    r = api("/admin/v1/login-policies", "POST", {
        "data": {"type": "BLACKLIST", "method": "REGION", "value": "北京市", "reason": "e2e region 黑名单"},
    })
    check("建 REGION 黑名单", "error" not in r and "code" not in r, str(r)[:100])
    ok, msg = try_login()
    check("REGION 不命中 → 放行", ok, msg[:80])

    # 1d. REGION 白名单（局域网——geolite 对回环/内网 IP 的解析结果）
    r = api("/admin/v1/login-policies", "POST", {
        "data": {"type": "WHITELIST", "method": "REGION", "value": "局域网", "reason": "e2e region 白名单"},
    })
    check("建 REGION 白名单（局域网）", "error" not in r and "code" not in r, str(r)[:100])
    ok, msg = try_login()
    check("REGION 白名单命中（局域网）→ 放行", ok, msg[:80])

    # 1e. REGION 白名单改为他省 → 拒绝
    lst = api("/admin/v1/login-policies?page=1&pageSize=50")
    rwl = [it for it in lst.get("items", []) if it.get("method") == "REGION" and it.get("type") == "WHITELIST"]
    if rwl:
        api(f"/admin/v1/login-policies/{rwl[0]['id']}", "PUT", {
            "data": {"value": "北京市", "reason": "e2e region 白名单（他省）"},
            "updateMask": "value,reason",
        })
    ok, msg = try_login()
    check("REGION 白名单未命中（他省）→ 拒绝", not ok and "whitelist" in msg, msg[:100])

    # 1f. 非法 IP 值校验
    r = api("/admin/v1/login-policies", "POST", {
        "data": {"type": "BLACKLIST", "method": "IP", "value": "not-an-ip", "reason": "e2e"},
    })
    check("非法 IP 值被拒（400）", "400" in json.dumps(r), str(r)[:120])

    # 清理全部测试策略
    lst = api("/admin/v1/login-policies?page=1&pageSize=50")
    for it in lst.get("items", []):
        if (it.get("reason") or "").startswith("e2e"):
            api(f"/admin/v1/login-policies/{it['id']}", "DELETE")
    ok, msg = try_login()
    check("清理后登录恢复", ok, msg[:80])

    print("== 2. 系统参数域 ==")
    lst = api("/admin/v1/configs?page=1&pageSize=20")
    keys = [it.get("key") for it in (lst.get("items") or [])]
    check("内置三键已播种", all(k in keys for k in (
        "sys.password.minLen", "sys.password.maxAgeDays", "sys.password.historyCount")), str(keys))

    # 创建自定义键
    r = api("/admin/v1/configs", "POST", {
        "data": {"key": "e2e.config.dummy", "name": "E2E 测试参数", "value": "42", "valueType": "INT"},
    })
    check("创建自定义参数", "error" not in r and "code" not in r, str(r)[:100])

    # 重复键
    r = api("/admin/v1/configs", "POST", {
        "data": {"key": "e2e.config.dummy", "name": "dup", "value": "1", "valueType": "INT"},
    })
    check("重复键被拒", "already exists" in json.dumps(r), str(r)[:120])

    lst = api("/admin/v1/configs?page=1&pageSize=20")
    row = next((it for it in lst.get("items", []) if it.get("key") == "e2e.config.dummy"), None)
    check("自定义参数可检索", row is not None)

    # 更新值
    if row:
        r = api(f"/admin/v1/configs/{row['id']}", "PUT", {
            "data": {"value": "43"},
            "updateMask": "value",
        })
        got = api(f"/admin/v1/configs/{row['id']}")
        check("更新参数值生效", got.get("value") == "43", str(got.get("value")))

    # 删内置键 → 400
    builtin = next((it for it in lst.get("items", []) if it.get("key") == "sys.password.minLen"), None)
    if builtin:
        r = api(f"/admin/v1/configs/{builtin['id']}", "DELETE")
        check("内置键禁删（400）", "built-in" in json.dumps(r), str(r)[:120])

    # 删自定义键
    if row:
        r = api(f"/admin/v1/configs/{row['id']}", "DELETE")
        check("删除自定义参数", "error" not in r and "code" not in r, str(r)[:100])
        lst = api("/admin/v1/configs?page=1&pageSize=20")
        keys = [it.get("key") for it in (lst.get("items") or [])]
        check("删除后不再含", "e2e.config.dummy" not in keys)

    print(f"\n结果: {passed} passed, {failed} failed")
    raise SystemExit(1 if failed else 0)


if __name__ == "__main__":
    main()
