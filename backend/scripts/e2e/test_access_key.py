"""AccessKey（OpenAPI AK/SK 凭证）E2E：管理 CRUD + 令牌交换 + 轮换 + 防爆破限流。

链路：
  1. 创建凭证（secret 仅本次明文返回）→ 列表/详情回显
  2. IssueToken 正确 AK/SK → 机器 JWT → 以 Bearer 调 admin API（posts 列表）
  3. Update 改名生效；掩码带 accessKey（不可变字段）→ AK 不被改写
  4. ResetSecret → 新 SK 可交换、旧 SK 立即失效（成功交换清零限流计数）
  5. 停用 → 交换拒绝（disabled）；重新启用恢复
  6. 连续错误 SK → 第 6 次起锁定（IP+AK 双维度，5 次/15 分钟）
  7. 清理：删除凭证

前置：core/admin 已启动（含本批二进制），redis 运行中。
注意：限流计数带 IP 维度（本机 127.0.0.1），失败用例须全部排在成功用例之后，
且中途任意一次成功交换会把计数清零——顺序不能倒。
用法：python test_access_key.py
"""
import json
import time
import urllib.error
import urllib.request

from common import ADMIN_BASE, AdminApi, _admin_api_raw, admin_login

# 限流按调用方 IP 分维度计数，而 urllib 对 localhost 会在 IPv4/IPv6 间摇摆
# （::1 与 127.0.0.1 各一把锁）。交换请求全部钉死 IPv4，保证计数落在同一把锁上。
TOKEN_BASE = "http://127.0.0.1:6600"


def issue_token(ak, sk):
    r = urllib.request.Request(
        TOKEN_BASE + "/admin/v1/access-keys/token",
        data=json.dumps({"accessKey": ak, "secret": sk}).encode(),
        method="POST",
        headers={"Content-Type": "application/json"},
    )
    try:
        return json.loads(urllib.request.urlopen(r).read().decode())
    except urllib.error.HTTPError as e:
        return json.loads(e.read().decode())

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

    print("== 1. 创建与回显 ==")
    r = api("/admin/v1/access-keys", "POST", {"data": {"name": "e2e-ak-A"}})
    check("create 返回结构", "secret" in r and r.get("data", {}).get("accessKey", "").startswith("ak-"), str(r)[:120])
    ak_a = r["data"]["accessKey"]
    sk_a = r["secret"]
    id_a = r["data"]["id"]
    check("SK 前缀 sk-", sk_a.startswith("sk-"))
    check("状态默认启用", r["data"].get("status") == "ON")
    check("secret 不回显于 data", "secret" not in (r["data"] or {}))

    r = api("/admin/v1/access-keys", "POST", {"data": {"name": "e2e-ak-B"}})
    ak_b, sk_b, id_b = r["data"]["accessKey"], r["secret"], r["data"]["id"]

    lst = api("/admin/v1/access-keys?page=1&pageSize=50")
    aks = [it.get("accessKey") for it in (lst.get("items") or [])]
    check("列表含新凭证", ak_a in aks and ak_b in aks)
    check("列表 total 为字符串(uint64 JSON)", isinstance(lst.get("total"), str))

    got = api(f"/admin/v1/access-keys/{id_a}")
    check("Get 按 id 回显 AK", got.get("accessKey") == ak_a)
    check("Get 不泄露 secret 摘要", "secretHash" not in got and "secret_hash" not in got)

    print("== 2. 令牌交换与机器令牌调用 ==")
    r = issue_token(ak_a, sk_a)
    check("交换返回 bearer", r.get("tokenType") == "bearer" and bool(r.get("accessToken")), str(r)[:120])
    check("expiresIn > 0", int(r.get("expiresIn") or 0) > 0)
    machine_tok = r.get("accessToken", "")

    raw, _ = _admin_api_raw("/admin/v1/posts?page=1&pageSize=1", headers={"Authorization": "Bearer " + machine_tok})
    check("机器令牌可调 admin API", raw and '"items"' in raw, raw[:120])

    raw, _ = _admin_api_raw("/admin/v1/posts?page=1&pageSize=1")
    check("无 token 调用被拒", "accessToken" not in raw and "未认证" not in raw and ("401" in raw[:40] or "Unauthorized" in raw or "unauthorized" in raw.lower()), raw[:120])

    got = api(f"/admin/v1/access-keys/{id_a}")
    check("lastUsedAt 已刷新", bool(got.get("lastUsedAt")), str(got.get("lastUsedAt")))

    print("== 3. Update 与不可变字段 ==")
    r = api(f"/admin/v1/access-keys/{id_a}", "PUT", {
        "data": {"name": "e2e-ak-A-renamed"},
        "updateMask": "name",
    })
    got = api(f"/admin/v1/access-keys/{id_a}")
    check("改名生效", got.get("name") == "e2e-ak-A-renamed", str(got.get("name")))

    r = api(f"/admin/v1/access-keys/{id_a}", "PUT", {
        "data": {"accessKey": "ak-hacked", "name": "e2e-ak-A-renamed"},
        "updateMask": "accessKey,name",
    })
    got = api(f"/admin/v1/access-keys/{id_a}")
    check("掩码带 accessKey 不改写 AK（不可变白名单）", got.get("accessKey") == ak_a, str(got.get("accessKey")))

    print("== 4. 密钥轮换 ==")
    r = api(f"/admin/v1/access-keys/{id_a}/secret", "PUT", {})
    check("ResetSecret 返回新明文", (r.get("secret") or "").startswith("sk-"), str(r)[:120])
    sk_a2 = r["secret"]

    r = issue_token(ak_a, sk_a)
    check("旧 SK 立即失效", r.get("code") == "400" or "invalid" in str(r.get("error", "")) or "invalid" in str(r.get("message", "")), str(r)[:120])

    r = issue_token(ak_a, sk_a2)
    check("新 SK 可交换", bool(r.get("accessToken")), str(r)[:120])
    # 成功交换清零限流计数，后续失败用例从干净状态起算

    print("== 5. 停用语义 ==")
    api(f"/admin/v1/access-keys/{id_b}", "PUT", {"data": {"status": "OFF"}, "updateMask": "status"})
    r = issue_token(ak_b, sk_b)
    check("停用后交换被拒（disabled）", "disabled" in json.dumps(r), str(r)[:120])
    api(f"/admin/v1/access-keys/{id_b}", "PUT", {"data": {"status": "ON"}, "updateMask": "status"})

    print("== 6. 防爆破限流（IP+AK，5 次锁 15 分钟）==")
    # 第 4 节的成功交换已清零限流计数；连错 5 次 → 第 6 次起锁定
    for i in range(4):
        r = issue_token(ak_a, "sk-" + "0" * 64)
        check(f"错误 SK 第 {i + 1} 次仅报无效（未锁定）", "invalid access key" in json.dumps(r), str(r)[:100])
    issue_token(ak_a, "sk-" + "0" * 64)  # 第 5 次失败：计数达阈值
    r = issue_token(ak_a, sk_a2)
    check("第 6 次尝试被锁定", "too many attempts" in json.dumps(r), str(r)[:120])
    r = issue_token(ak_b, sk_b)
    check("同 IP 其他凭证也被锁（IP 维度）", "too many attempts" in json.dumps(r), str(r)[:120])

    print("== 7. 清理 ==")
    r = api(f"/admin/v1/access-keys/{id_a}", "DELETE")
    r2 = api(f"/admin/v1/access-keys/{id_b}", "DELETE")
    lst = api("/admin/v1/access-keys?page=1&pageSize=50")
    aks = [it.get("accessKey") for it in (lst.get("items") or [])]
    check("删除后列表不再含", ak_a not in aks and ak_b not in aks, str(r)[:80] + str(r2)[:80])

    print(f"\n结果: {passed} passed, {failed} failed")
    raise SystemExit(1 if failed else 0)


if __name__ == "__main__":
    main()
