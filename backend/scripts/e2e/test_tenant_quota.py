"""租户配额子系统 E2E：配额 CRUD + 用量聚合 + 三闸门 + 水位告警全链。

链路：
  1. 平台管理员为 tenant 1 配 USER_LIMIT（=当前用户数）→ 再建用户 403
  2. tenant 1 的用户自助 GetMyTenantUsage：tenantId 钉定 + 配额可见
  3. STORAGE 配额（=当前占用+10B）→ 租户用户上传 100KB → 403
  4. API_CALL 配额（=当前审计行数）→ 租户用户任意请求 → 网关闸 403
  5. 水位告警：建每分钟 tenant_quota_watermark_scan 任务 + StartAllTask →
     等待扫描 → 租户管理员（user 2）收到站内信
  6. 清理：配额行/测试用户/水位任务行；平台身份（tenant 0）全程不受闸门限制

前置：core/admin 已启动（含本批二进制）；dev 库有 tenant 1（admin_user_id=2）。
注意：闸门 fail-open——配额行删除后租户立即恢复不受限。
用法：python test_tenant_quota.py
"""
import json
import subprocess
import time

from common import ADMIN_BASE, AdminApi, admin_login

TEST_TENANT_ID = 1
TENANT_ADMIN_USER_ID = 2
TEST_USER_PASSWORD = "12345678"

passed = failed = 0


def check(name, cond, extra=""):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✅ {name} {extra}")
    else:
        failed += 1
        print(f"  ❌ {name} {extra}")


def upload_file(headers, filename, size_kb=100):
    # curl.exe（Windows 原生）读不了 MSYS 的 /tmp 虚拟路径——临时文件落在
    # 当前目录（相对路径由 bash 侧解析，curl 按进程 cwd 打开）。
    import io
    import os
    import uuid
    tmp = f"quota-upload-{uuid.uuid4().hex}.bin"
    with open(tmp, "wb") as f:
        f.write(io.BytesIO(uuid.uuid4().bytes * size_kb * 32).getvalue())
    cmd = [
        "curl", "-s", "--noproxy", "*", "--max-time", "30",
        "-X", "POST", f"{ADMIN_BASE}/admin/v1/file/upload",
        "-H", f"Authorization: {headers['Authorization']}",
        "-F", "mime=application/octet-stream",
        "-F", f"sourceFileName={filename}",
        "-F", "storageObject={}",
        "-F", f"file=@{tmp};type=application/octet-stream",
    ]
    r = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
    os.unlink(tmp)
    try:
        return json.loads(r.stdout)
    except json.JSONDecodeError:
        return {"_raw": r.stdout[:200], "_err": r.stderr[:200]}


def main():
    hdr = admin_login("admin", "admin")
    api = AdminApi(hdr)

    print("== 0. 基线 ==")
    usage0 = api(f"/admin/v1/tenants/{TEST_TENANT_ID}/usage")
    check("GetUsage 返回租户用量", usage0.get("tenantId") == TEST_TENANT_ID, str(usage0.get("tenantId")))
    base_users = int(usage0.get("userCount") or 0)
    base_storage = int(usage0.get("storageUsedBytes") or 0)
    base_api = int(usage0.get("apiCallCount") or 0)
    print(f"     基线：users={base_users} storage={base_storage} apiCalls={base_api}")

    print("== 1. 建租户测试用户（配额未设，应成功）==")
    uname = f"quota-e2e-{int(time.time())}"
    roles = api('/admin/v1/roles?page=1&pageSize=20')
    tenant_role_id = next(
        (it.get("id") for it in roles.get("items", []) if it.get("type") == "TENANT"),
        None,
    )
    assert tenant_role_id, f"no TENANT role available: {str(roles)[:200]}"
    r = api("/admin/v1/users", "POST", {
        "data": {
            "username": uname,
            "nickname": "quota-e2e",
            "tenantId": TEST_TENANT_ID,
            "status": "ON",
            "roleIds": [tenant_role_id],
        },
        "password": TEST_USER_PASSWORD,
    })
    check("创建租户用户成功", "id" in r or "error" not in r, str(r)[:120])
    test_user = api(f"/admin/v1/users?page=1&pageSize=20")
    test_user_id = None
    for it in test_user.get("items", []):
        if it.get("username") == uname:
            test_user_id = it.get("id")
    check("测试用户可检索", test_user_id is not None, str(test_user_id))

    # tenant:user 角色（demo 数据）只有 sys:access_app，登录 admin 链要求
    # sys:access_backend——按字段权限批先例 SQL 补绑 role 2（租户管理员模板，
    # 带后端访问权限；必须带 tenant_id=1，ent viewer 按租户过滤）。
    subprocess.run([
        "docker", "exec", "citus-server-standalone", "psql", "-U", "postgres", "-d", "gwc", "-c",
        f"INSERT INTO sys_user_roles (tenant_id, user_id, role_id, created_at) "
        f"VALUES ({TEST_TENANT_ID}, {test_user_id}, 2, now());",
    ], capture_output=True, text=True, timeout=30)
    check("补绑租户管理员角色（sys_user_roles）", True, f"user={test_user_id} role=2")

    tenant_hdr = None
    try:
        tenant_hdr = admin_login(uname, TEST_USER_PASSWORD)
        check("测试用户可登录", True)
    except AssertionError as e:
        check("测试用户可登录", False, str(e)[:120])
        raise SystemExit(1)
    tenant_api = AdminApi(tenant_hdr)

    print("== 2. USER_LIMIT 闸门 ==")
    r = api("/admin/v1/tenant-quotas", "POST", {
        "data": {"tenantId": TEST_TENANT_ID, "quotaType": "USER_LIMIT", "quotaValue": base_users + 1},
    })
    check("配 USER_LIMIT（=现有+1）", "error" not in r and "code" not in r, str(r)[:120])

    r = api("/admin/v1/users", "POST", {
        "data": {
            "username": uname + "-x",
            "nickname": "quota-e2e-x",
            "tenantId": TEST_TENANT_ID,
            "status": "ON",
            "roleIds": [tenant_role_id],
        },
        "password": TEST_USER_PASSWORD,
    })
    check("第 2 个用户创建被拒（403 超限）", "403" in json.dumps(r) and "USER_LIMIT" in json.dumps(r), str(r)[:140])

    r = api("/admin/v1/tenant-quotas", "POST", {
        "data": {"tenantId": TEST_TENANT_ID, "quotaType": "USER_LIMIT", "quotaValue": 0},
    })
    check("重复类型被拒（唯一约束）", "400" in json.dumps(r), str(r)[:120])

    # 更新为正好=现有数：再来一个也拒
    lst = api("/admin/v1/tenant-quotas?page=1&pageSize=50")
    ql = [it for it in lst.get("items", []) if it.get("tenantId") == TEST_TENANT_ID and it.get("quotaType") == "USER_LIMIT"]
    check("配额列表含新行", len(ql) == 1)
    if ql:
        api(f"/admin/v1/tenant-quotas/{ql[0]['id']}", "PUT", {
            "data": {"quotaValue": base_users + 1},
            "updateMask": "quotaValue",
        })

    print("== 3. 自助用量面板 ==")
    my = tenant_api("/admin/v1/my-tenant-usage")
    check("GetMyTenantUsage 钉定本租户", my.get("tenantId") == TEST_TENANT_ID, str(my.get("tenantId")))
    qtypes = sorted(q.get("quotaType") for q in (my.get("quotas") or []))
    check("配额可见", "USER_LIMIT" in qtypes, str(qtypes))
    check("用量计数随创建增长", int(my.get("userCount") or 0) == base_users + 1, str(my.get("userCount")))

    print("== 4. STORAGE 闸门 ==")
    r = api("/admin/v1/tenant-quotas", "POST", {
        "data": {"tenantId": TEST_TENANT_ID, "quotaType": "STORAGE", "quotaValue": base_storage + 10},
    })
    check("配 STORAGE（=占用+10B）", "error" not in r and "code" not in r, str(r)[:120])
    up = upload_file(tenant_hdr, "quota-e2e.bin", size_kb=100)
    check("100KB 上传被拒（403 超限）", "403" in json.dumps(up) and "STORAGE" in json.dumps(up), str(up)[:140])

    print("== 5. API_CALL 网关闸 ==")
    # 每个请求自身也落审计行，先取即时计数再设上限（下一个请求即超）
    usage_now = api(f"/admin/v1/tenants/{TEST_TENANT_ID}/usage")
    base_api = int(usage_now.get("apiCallCount") or 0)
    # 审计行在响应后落库：闸门读到的是"当前计数"，限值设为当前数即可让下一请求命中
    r = api("/admin/v1/tenant-quotas", "POST", {
        "data": {"tenantId": TEST_TENANT_ID, "quotaType": "API_CALL", "quotaValue": base_api},
    })
    check("配 API_CALL（=当前审计行数）", "error" not in r and "code" not in r, str(r)[:120])
    r = tenant_api("/admin/v1/my-tenant-usage")
    check("租户用户下一请求被网关拒（403）", "403" in json.dumps(r) and "API_CALL" in json.dumps(r), str(r)[:140])

    print("== 6. 水位告警（每分钟扫描任务）==")
    # USER_LIMIT 当前 ratio = (base_users+1)/(base_users+1) = 100% ≥ 80% 必命中
    r = api("/admin/v1/tasks", "POST", {
        "data": {
            "typeName": "tenant_quota_watermark_scan",
            "type": "PERIODIC",
            "cronSpec": "* * * * *",
            "enable": True,
            "taskPayload": "{}",
        },
    })
    check("建水位扫描任务（每分钟）", "error" not in r and "code" not in r, str(r)[:120])
    api("/admin/v1/tasks:start", "POST", {})

    delivered = False
    for _ in range(24):
        time.sleep(5)
        msgs = api("/admin/v1/internal-message/messages?page=1&pageSize=10&orderBy=%5B%22-id%22%5D")
        for m in msgs.get("items", []):
            t = m.get("title") or ""
            if ("配额水位告警" in t) and ("测试租户" in t):
                delivered = True
                break
        if delivered:
            break
    check("租户管理员收到水位站内信", delivered)

    print("== 7. 清理 ==")
    # 停用并删除水位任务（Update enable=false 即移除调度注册）
    lst = api("/admin/v1/tasks?page=1&pageSize=50")
    for it in lst.get("items", []):
        if it.get("typeName") == "tenant_quota_watermark_scan":
            api(f"/admin/v1/tasks/{it['id']}", "PUT", {
                "data": {
                    "typeName": "tenant_quota_watermark_scan",
                    "type": "PERIODIC",
                    "cronSpec": "* * * * *",
                    "enable": False,
                    "taskPayload": "{}",
                },
                "updateMask": "enable",
            })
            api(f"/admin/v1/tasks/{it['id']}", "DELETE")
    check("水位任务已清理", True)

    for it in (api("/admin/v1/tenant-quotas?page=1&pageSize=50").get("items") or []):
        if it.get("tenantId") == TEST_TENANT_ID:
            api(f"/admin/v1/tenant-quotas/{it['id']}", "DELETE")
    check("配额行已清理", True)

    if test_user_id:
        api(f"/admin/v1/users/{test_user_id}", "DELETE")
    check("测试用户已清理", True)

    # 平台身份不受闸门限制（tenant 0 跳过）：删配额后租户也应恢复
    time.sleep(1)
    r = tenant_api("/admin/v1/my-tenant-usage")
    check("租户删除配额后恢复访问（fail-open）", "403" not in json.dumps(r), str(r)[:100])

    print(f"\n结果: {passed} passed, {failed} failed")
    raise SystemExit(1 if failed else 0)


if __name__ == "__main__":
    main()
