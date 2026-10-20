"""AI Agent 编排（切片①）E2E：应用 CRUD + 运行台账读取 + Run 未接入语义。

链路：
  1. 创建应用（租户盖章）→ 列表/详情回显 → 更新（掩码 name/enabled/graph）→ 删除
  2. 同租户重名 → 400；graph 非法 JSON 由 FE 层拦截（E2E 直发后端不校验结构——
     结构强校验在执行引擎切片的保存/运行路径）
  3. Run 未接入（切片②）：501 且不建运行行；ListRuns/GetRun 空集可读
  4. 清理

前置：core/admin 已启动（含本批二进制）。
用法：python test_ai_agent.py
"""
import json
import time
import uuid

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
    uname = f"e2e-agent-{uuid.uuid4().hex[:6]}"

    print("== 1. 创建与回显 ==")
    graph = json.dumps({
        "nodes": [
            {"id": "start", "type": "start", "data": {"inputVars": [{"name": "q"}]}},
            {"id": "end", "type": "end", "data": {"template": "答案：{{q}}"}},
        ],
        "edges": [{"source": "start", "target": "end"}],
    }, ensure_ascii=False)
    r = api("/admin/v1/ai/agents", "POST", {
        "data": {"name": uname, "description": "E2E 编排应用", "enabled": True, "graph": graph},
    })
    check("create 返回应用", r.get("name") == uname, str(r)[:120])
    app_id = r.get("id")
    check("租户盖章（admin→tenant 0）", r.get("tenantId") == 0, str(r.get("tenantId")))

    lst = api("/admin/v1/ai/agents?page=1&pageSize=50")
    names = [it.get("name") for it in (lst.get("items") or [])]
    check("列表含新应用", uname in names)
    check("total 为字符串", isinstance(lst.get("total"), str))

    got = api(f"/admin/v1/ai/agents/{app_id}")
    check("Get 回显 graph", json.loads(got.get("graph") or "{}").get("nodes") is not None)
    check("enabled 回填（值型 bool）", got.get("enabled") is True, str(got.get("enabled")))

    print("== 2. 更新与唯一约束 ==")
    r = api(f"/admin/v1/ai/agents/{app_id}", "PUT", {
        "data": {"name": uname + "-r", "enabled": False, "description": "更新后"},
        "updateMask": "name,enabled,description",
    })
    check("掩码更新受理", "error" not in r and "code" not in r, str(r)[:100])
    got = api(f"/admin/v1/ai/agents/{app_id}")
    check("更新生效", got.get("name") == uname + "-r" and got.get("enabled") is False, str(got.get("enabled")))

    r = api("/admin/v1/ai/agents", "POST", {
        "data": {"name": uname + "-r", "enabled": True},
    })
    check("租户内重名被拒（400）", "already exists" in json.dumps(r), str(r)[:120])

    print("== 3. Run 未接入（切片②）==")
    r = api(f"/admin/v1/ai/agents/{app_id}/run", "POST", {"input": {"q": "hi"}})
    check("Run 返回 501（未实现）", "501" in json.dumps(r) and "not implemented" in json.dumps(r), str(r)[:140])

    runs = api("/admin/v1/ai/agent-runs?page=1&pageSize=10")
    check("运行台账可读（空集）", runs.get("total") in ("0", 0, None), str(runs.get("total")))

    print("== 4. 删除 ==")
    r = api(f"/admin/v1/ai/agents/{app_id}", "DELETE")
    check("删除成功", "error" not in r and "code" not in r, str(r)[:100])
    lst = api("/admin/v1/ai/agents?page=1&pageSize=50")
    names = [it.get("name") for it in (lst.get("items") or [])]
    check("删除后不再含", uname + "-r" not in names)

    print(f"\n结果: {passed} passed, {failed} failed")
    raise SystemExit(1 if failed else 0)


if __name__ == "__main__":
    main()
