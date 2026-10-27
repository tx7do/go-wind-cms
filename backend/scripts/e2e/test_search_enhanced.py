"""搜索深化 E2E：高亮 / 建议词 / 同义词（查询侧扩展）。

前置：core/admin/app 三服务已启动（本特性代码版本）；OpenSearch :9200 运行中且
demo 文章已重建索引（go run ./cmd/reindex，demo 文章属租户 1）；sys_tenants id=1
已绑域名 127.0.0.1:6700（app 按 Host 解析租户；E2E 请求显式带该 Host）。

demo 锚点：官方演示文章标题含「内容中台」（tenant 1, PUBLISHED）。
同义词断言用临时触发词映射到「内容中台」——扩展词命中即扩展生效。

断言（9 条）：
  1. 索引有已发布文档：直接搜「内容中台」命中
  2. 建同义词（平台级，全局生效）：临时触发词 → 内容中台
  3. 同义词回读一致
  4. 同义词扩展：搜触发词命中 demo 文章
  5. highlightFragments 含 <mark>
  6. 响应携带 summary 字段
  7. 建议词端点：前缀「内容中台」返回标题列表
  8. 删同义词后：搜触发词不再命中（扩展撤除）
  9. 全部断言完成
"""
import json
import sys
import time
import urllib.parse
import urllib.request

sys.path.insert(0, __file__.rsplit("\\", 1)[0])
from common import APP_BASE, AdminApi, admin_login  # noqa: E402

PASS = 0


def ok(cond, label, extra=""):
    global PASS
    assert cond, f"FAIL: {label} {extra}"
    PASS += 1
    print(f"  ok {PASS:02d} - {label}")


def app_get(path):
    # 搜索租户按 Host 解析：显式指定已绑域名的 Host
    req = urllib.request.Request(APP_BASE + path, headers={"Host": "127.0.0.1:6700"})
    return json.loads(urllib.request.urlopen(req).read().decode())


def search_hits(q):
    enc = urllib.parse.quote(q)
    return app_get(f"/app/v1/posts/search?query={enc}&language=zh-CN").get("items") or []


def wait_hits(q, min_n=1, timeout=15):
    deadline = time.time() + timeout
    while time.time() < deadline:
        if len(search_hits(q)) >= min_n:
            return True
        time.sleep(1)
    return False


def main():
    hdr = admin_login("admin", "admin")
    api = AdminApi(hdr)
    TS = int(time.time())

    # ---- 清理上次运行残留 ----
    for s in (api("/admin/v1/search-synonyms", "GET").get("items") or []):
        if str(s.get("term", "")).startswith("中枢宝典") or str(s.get("term", "")).startswith("e2esyn"):
            api(f"/admin/v1/search-synonyms/{s.get('id')}", "DELETE")

    # ---- 1. 索引基线：触发词本身应 0 命中（文档不含该词）----
    # 触发词用真实中文词（smartcn 正常分词）；随机长串会被二元语法切成
    # bi-gram token 在 OR 语义下噪声全中，不能做扩展断言的载体。
    trigger = "中枢宝典"
    ok(len(search_hits(trigger)) == 0, "触发词基线 0 命中（文档不含）", "")

    # ---- 2. 索引有已发布文档：直接搜 demo 关键词 ----
    hits = search_hits("内容中台")
    ok(len(hits) >= 1, "索引有已发布文档：搜「内容中台」命中", f"n={len(hits)}")
    syn = api("/admin/v1/search-synonyms", "POST", {"data": {
        "term": trigger,
        "synonyms": ["内容中台"],
        "enabled": True,
    }})
    syn_id = syn.get("id")
    ok(bool(syn_id), "建同义词", f"id={syn_id}")
    listed = api("/admin/v1/search-synonyms", "GET").get("items") or []
    back = next((x for x in listed if x.get("id") == syn_id), None)
    ok(back is not None and back.get("synonyms") == ["内容中台"], "同义词回读一致", str(back)[:80])

    # ---- 3. 同义词扩展：搜触发词命中 demo 文章 ----
    ok(wait_hits(trigger, min_n=1), "同义词扩展：搜触发词命中", "")

    # ---- 4/5. 高亮与 summary ----
    hits = search_hits("内容中台")
    hit = hits[0] if hits else {}
    frags = hit.get("highlightFragments") or []
    ok(any("<mark>" in f for f in frags), "highlightFragments 含 <mark>", str(frags[:1])[:60])
    ok("summary" in hit, "响应携带 summary 字段", f"summary={str(hit.get('summary'))[:30]}")

    # ---- 6. 建议词 ----
    enc = urllib.parse.quote("内容中台")
    sug = app_get(f"/app/v1/posts/search/suggest?query={enc}&language=zh-CN")
    titles = sug.get("titles") or []
    ok(len(titles) >= 1, "建议词返回标题列表", f"titles={str(titles)[:80]}")

    # ---- 7. 同义词扩展再验（时间点独立于 3）----
    ok(wait_hits(trigger, min_n=1), "同义词扩展再验", "")

    # ---- 8. 删同义词 → 扩展撤除 ----
    api(f"/admin/v1/search-synonyms/{syn_id}", "DELETE")
    ok(not wait_hits(trigger, min_n=1, timeout=6), "删词后搜触发词不再命中（扩展撤除）", "")

    print(f"\nALL PASS ({PASS} assertions)")


if __name__ == "__main__":
    main()
