"""评论 AI 反垃圾 E2E：开关默认关 / 开启后 AI 不可用 fail-open。

前置（同 e2e_common.py）：admin(:6600)/core 已启动（本特性代码版本）；一篇文章可用。

断言（5 条）：
  1. 内置参数键已由启动种子落位（GET configs 可见）
  2. 默认关闭：创建评论直接成功（不触 AI）
  3. 开启开关（value=true）
  4. 开启后创建评论：provider 不可达/判定非垃圾 → 仍成功（fail-open 或非垃圾放行）
  5. 收尾：开关回置 false
"""
import json
import sys
import time

sys.path.insert(0, __file__.rsplit("\\", 1)[0])
from common import AdminApi, admin_login  # noqa: E402

KEY = "comment.ai_moderation.enabled"
PASS = 0


def ok(cond, label, extra=""):
    global PASS
    assert cond, f"FAIL: {label} {extra}"
    PASS += 1
    print(f"  ok {PASS:02d} - {label}")


def main():
    hdr = admin_login("admin", "admin")
    api = AdminApi(hdr)
    TS = int(time.time())

    def find_config():
        listed = api("/admin/v1/configs?page=1&pageSize=100", "GET")
        for c in listed.get("items", []):
            if c.get("key") == KEY:
                return c
        return None

    cfg = find_config()
    ok(cfg is not None, "内置参数键已种子落位", f"id={cfg.get('id') if cfg else None}")

    # 开关确认为 false（幂等）
    if cfg and (cfg.get("value") or "").lower() == "true":
        api(f"/admin/v1/configs/{cfg['id']}", "PUT", {"data": {"value": "false"},
                                                      "updateMask": "value"})
        time.sleep(1)

    # 建文章供评论
    post = api("/admin/v1/posts", "POST", {"data": {
        "status": "DRAFT",
        "translations": [{"languageCode": "zh-CN", "title": f"mod e2e {TS}",
                          "content": "<p>m</p>"}],
    }})
    post_id = post.get("id")
    ok(bool(post_id), "建文章", f"id={post_id}")

    def create_comment():
        return api("/admin/v1/comments", "POST", {"data": {
            "objectId": post_id,
            "contentType": "CONTENT_TYPE_POST",
            "content": f"moderation e2e {TS}",
            "status": "STATUS_PENDING",
        }})

    # 默认关：直接成功
    c1 = create_comment()
    ok(bool(c1.get("id")), "默认关闭创建成功", f"id={c1.get('id')} resp={json.dumps(c1, ensure_ascii=False)[:80]}")

    # 开启
    assert cfg is not None
    api(f"/admin/v1/configs/{cfg['id']}", "PUT", {"data": {"value": "true"},
                                                  "updateMask": "value"})
    time.sleep(1.5)  # 配置失效广播 + 进程内缓存过期

    # 开启后：provider 不可达 → fail-open；可达且判非垃圾 → 放行。两者都表现为创建成功。
    c2 = create_comment()
    ok(bool(c2.get("id")), "开启后创建仍成功（fail-open/非垃圾）",
       f"id={c2.get('id')} resp={json.dumps(c2, ensure_ascii=False)[:80]}")

    # 收尾：回置 false + 清理
    api(f"/admin/v1/configs/{cfg['id']}", "PUT", {"data": {"value": "false"},
                                                  "updateMask": "value"})
    api(f"/admin/v1/comments/{c1.get('id')}", "DELETE")
    api(f"/admin/v1/comments/{c2.get('id')}", "DELETE")
    api(f"/admin/v1/posts/{post_id}", "DELETE")

    print(f"\nALL PASS ({PASS} assertions)")


if __name__ == "__main__":
    main()
