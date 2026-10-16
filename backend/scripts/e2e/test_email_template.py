"""邮件模板 CRUD E2E：创建→列表→更新→按 code 启用查询→删除。

前置：admin/core 服务已启动（email_templates 表由 core 自动迁移创建）。
用法：python test_email_template.py
"""
import sys

from common import AdminApi, admin_login

hdr = admin_login("admin", "admin")
api = AdminApi(hdr)
passed = failed = 0


def check(name, cond, extra=""):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✅ {name} {extra}")
    else:
        failed += 1
        print(f"  ❌ {name} {extra}")


# 1. 创建
r = api("/admin/v1/email-templates", "POST", {"data": {
    "code": "notify/system",
    "name": "系统通知邮件",
    "subject": "【{{.Title}}】风行内容中台通知",
    "bodyHtml": '<div><h2>{{.Title}}</h2><div>{{.Content}}</div>'
                '<a href="{{.SiteURL}}">{{.SiteURL}}</a></div>',
    "enabled": True,
    "description": "E2E 测试模板",
}})
tid = r.get("id")
check("创建模板", tid is not None, f"id={tid} code={r.get('code')}")

# 2. 重复 code 创建被拒
r = api("/admin/v1/email-templates", "POST", {"data": {
    "code": "notify/system", "name": "dup", "subject": "x", "enabled": True,
}})
check("重复 code 被拒(409)", r.get("code") == 409, f"code={r.get('code')}")

# 3. 列表含新模板
r = api("/admin/v1/email-templates?code=notify&pageSize=10")
items = r.get("items") or []
check("列表含模板", any(it.get("id") == tid for it in items), f"total={r.get('total')}")

# 4. 更新 subject
r = api(f"/admin/v1/email-templates/{tid}", "PUT", {
    "id": tid,
    "data": {"subject": "【{{.Title}}】已更新"},
    "updateMask": "subject",
})
check("更新主题", r.get("subject") == "【{{.Title}}】已更新")

# 5. FindEnabledByCode（内部 RPC 经 admin 未暴露——直接验证 List 里的 enabled 状态）
r = api("/admin/v1/email-templates?code=notify&enabled=true&pageSize=10")
check("启用状态可过滤", any(it.get("id") == tid for it in (r.get("items") or [])))

# 6. 删除
api(f"/admin/v1/email-templates/{tid}", "DELETE")
r = api("/admin/v1/email-templates?code=notify/system&pageSize=10")
check("删除后列表为空", not (r.get("items") or []))

print(f"邮件模板 E2E：{passed} 通过 / {failed} 失败")
sys.exit(1 if failed else 0)
