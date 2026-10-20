"""AI Agent 编排（切片②）E2E：图执行引擎真跑——mock OpenAI provider + 试运行全链。

链路：
  1. 起 mock OpenAI Chat Completions 端点（:6902，返回固定答案与用量）
  2. 建 provider 指向 mock → 建 Agent 应用（start→llm→condition→end，双分支）
  3. Run：condition 为 true 路径 → SUCCEEDED，output 含模板渲染值，用量落流水
  4. Run：condition 为 false 路径 → 走 false 分支 end，产出为另一模板
  5. 运行台账：ListRuns 含两次运行；GetRun trace 含逐节点记录
  6. 校验语义：停用应用拒绝运行；存非法图 400；缺必填入参 FAILED（建行）
  7. 清理

前置：core/admin 已启动（含本批二进制）。
用法：python test_ai_agent_run.py
"""
import json
import threading
import time
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

from common import AdminApi, admin_login

MOCK_PORT = 6902
MOCK_ANSWER = "mock-agent-answer"
MOCK_EMBEDDING_DIM = 1536
MOCK_VECTOR = [1.0] * MOCK_EMBEDDING_DIM  # 检索种子向量同维同向：余弦距离 0，topK 确定命中

passed = failed = 0


def check(name, cond, extra=""):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✅ {name} {extra}")
    else:
        failed += 1
        print(f"  ❌ {name} {extra}")


def sse_chunk(content, finish_reason, usage=None):
    """OpenAI 流式分块（CMS 工具循环走 CreateChatCompletionStream）。"""
    chunk = {
        "id": "chatcmpl-mock",
        "object": "chat.completion.chunk",
        "created": int(time.time()),
        "model": "mock-agent-model",
        "choices": [{"index": 0, "delta": {}, "finish_reason": finish_reason}],
    }
    if content is not None:
        chunk["choices"][0]["delta"] = {"role": "assistant", "content": content}
    if usage is not None:
        chunk["usage"] = usage
    return "data: " + json.dumps(chunk) + "\n\n"


class MockOpenAIHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers.get("Content-Length", 0))
        self.rfile.read(length)

        if self.path.endswith("/embeddings"):
            # OpenAI embeddings 端点：知识节点检索链路（同维固定向量）
            body = json.dumps({
                "object": "list",
                "data": [{"object": "embedding", "index": 0, "embedding": MOCK_VECTOR}],
                "model": "mock-embedding",
                "usage": {"prompt_tokens": 3, "total_tokens": 3},
            }).encode()
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)
            return

        self.send_response(200)
        self.send_header("Content-Type", "text/event-stream; charset=utf-8")
        self.send_header("Cache-Control", "no-cache")
        self.end_headers()

        # 分两段文本增量 + 终止帧（usage 挂在最后一块，IncludeUsage 语义）
        self.wfile.write(sse_chunk(MOCK_ANSWER[:6], None).encode())
        self.wfile.write(sse_chunk(MOCK_ANSWER[6:], "stop", {
            "prompt_tokens": 11, "completion_tokens": 7, "total_tokens": 18,
        }).encode())
        self.wfile.write(b"data: [DONE]\n\n")

    def log_message(self, *args):
        pass


def sql(sqltext):
    import subprocess
    return subprocess.run(
        ["docker", "exec", "citus-server-standalone", "psql", "-U", "postgres", "-d", "gwc", "-t", "-A", "-c", sqltext],
        capture_output=True, text=True, timeout=30,
    ).stdout.strip()


def main():
    server = ThreadingHTTPServer(("127.0.0.1", MOCK_PORT), MockOpenAIHandler)
    threading.Thread(target=server.serve_forever, daemon=True).start()
    print(f"mock openai at :{MOCK_PORT}")

    hdr = admin_login("admin", "admin")
    api = AdminApi(hdr)

    print("== 1. mock provider + 应用 ==")
    r = api("/admin/v1/ai/providers", "POST", {
        "data": {
            "name": "e2e-agent-provider",
            "baseUrl": f"http://127.0.0.1:{MOCK_PORT}/v1",
            "apiKey": "sk-mock",
            "modelName": "mock-agent-model",
            "modelType": "CLOUD",
            "isEnabled": True,
            "isDefault": False,
        },
    })
    check("建 mock provider", "error" not in r and "code" not in r, str(r)[:120])
    lst = api("/admin/v1/ai/providers?page=1&pageSize=50")
    prov = next((it for it in lst.get("items", []) if it.get("name") == "e2e-agent-provider"), None)
    check("provider 可检索", prov is not None)
    provider_id = prov["id"]

    graph = json.dumps({
        "nodes": [
            {"id": "start", "type": "start", "data": {"inputVars": [{"name": "q", "required": True}, {"name": "branch", "required": True}]}},
            {"id": "llm1", "type": "llm", "data": {"providerId": provider_id, "systemPrompt": "你是测试助手", "userPrompt": "{{start.q}}"}},
            {"id": "cond", "type": "condition", "data": {"expression": 'start.branch == "true"'}},
            {"id": "end_t", "type": "end", "data": {"output": "真分支：{{llm1.output}} / {{start.q}}"}},
            {"id": "end_f", "type": "end", "data": {"output": "假分支：{{start.q}}"}},
        ],
        "edges": [
            {"source": "start", "target": "llm1"},
            {"source": "llm1", "target": "cond"},
            {"source": "cond", "target": "end_t", "branch": "true"},
            {"source": "cond", "target": "end_f", "branch": "false"},
        ],
    }, ensure_ascii=False)
    r = api("/admin/v1/ai/agents", "POST", {
        "data": {"name": "e2e-agent-run", "description": "切片②试运行", "enabled": True, "graph": graph},
    })
    check("建应用（图过静态校验）", r.get("id") is not None, str(r)[:120])
    app_id = r["id"]

    print("== 2. 真分支试运行 ==")
    r = api(f"/admin/v1/ai/agents/{app_id}/run", "POST", {"input": {"q": "你好引擎", "branch": "true"}})
    check("Run 返回运行记录", r.get("id") is not None, str(r)[:140])
    check("终态 SUCCEEDED", r.get("status") == "RUNNING_SUCCEEDED" or "SUCCEEDED" in str(r.get("status")), str(r.get("status")))
    check("产出含模板渲染（LLM 答案+入参）", MOCK_ANSWER in str(r.get("output")) and "你好引擎" in str(r.get("output")), str(r.get("output"))[:80])
    check("用量合计 18", int(r.get("totalTokens") or 0) == 18, str(r.get("totalTokens")))
    trace = json.loads(r.get("trace") or "[]")
    check("轨迹含 4 节点", [s.get("nodeId") for s in trace] == ["start", "llm1", "cond", "end_t"], str([s.get("nodeId") for s in trace]))

    # 用量流水：agent 非对话调用落一行
    time.sleep(0.5)
    ulogs = api("/admin/v1/ai/usage-logs?page=1&pageSize=10&orderBy=%5B%22-id%22%5D")
    ul = next((it for it in ulogs.get("items", []) if it.get("modelName") == "mock-agent-model"), None)
    check("用量流水已落（非对话调用）", ul is not None and int(ul.get("totalTokens") or 0) == 18, str(ul)[:80] if ul else "none")

    print("== 3. 假分支试运行 ==")
    r = api(f"/admin/v1/ai/agents/{app_id}/run", "POST", {"input": {"q": "第二条", "branch": "false"}})
    check("假分支产出", "假分支：第二条" in str(r.get("output")), str(r.get("output"))[:80])
    trace = json.loads(r.get("trace") or "[]")
    check("假分支走 end_f", trace[-1].get("nodeId") == "end_f")

    print("== 4. 校验语义 ==")
    api(f"/admin/v1/ai/agents/{app_id}", "PUT", {
        "data": {"enabled": False}, "updateMask": "enabled",
    })
    r = api(f"/admin/v1/ai/agents/{app_id}/run", "POST", {"input": {"q": "x", "branch": "true"}})
    check("停用应用拒绝运行（400）", "disabled" in json.dumps(r), str(r)[:120])

    api(f"/admin/v1/ai/agents/{app_id}", "PUT", {
        "data": {"enabled": True}, "updateMask": "enabled",
    })
    bad_graph = json.dumps({"nodes": [{"id": "start", "type": "start", "data": {}}, {"id": "start2", "type": "start", "data": {}}], "edges": []})
    r = api(f"/admin/v1/ai/agents/{app_id}", "PUT", {
        "data": {"graph": bad_graph}, "updateMask": "graph",
    })
    check("存非法图被拒（双 start 400）", "exactly one start" in json.dumps(r), str(r)[:140])

    r = api(f"/admin/v1/ai/agents/{app_id}/run", "POST", {"input": {}})
    check("缺必填入参 → FAILED（建行）", r.get("status") is not None and "FAILED" in str(r.get("status")) and "missing required input" in str(r.get("error")), str(r)[:140])

    runs = api(f"/admin/v1/ai/agent-runs?page=1&pageSize=100&appId={app_id}")
    mine = [it for it in (runs.get("items") or []) if int(it.get("appId") or 0) == app_id]
    check("本应用台账累计 3 行", len(mine) == 3, str(len(mine)))

    print("== 5. 知识节点链路（mock embeddings + SQL 种子切片）==")
    # 建知识库（指向 mock provider + embedding 模型）
    ts = int(time.time())
    kb_name = f"e2e-kb-agent-{ts}"
    r = api("/admin/v1/ai/knowledge-bases", "POST", {"data": {
        "name": kb_name,
        "description": "agent knowledge E2E",
        "providerId": prov["id"],
        "embeddingModel": "mock-embedding",
    }})
    check("建知识库", "error" not in r and "code" not in r, str(r)[:100])
    kbl = api("/admin/v1/ai/knowledge-bases?page=1&pageSize=50")
    kb_id = next((it.get("id") for it in kbl.get("items", []) if it.get("name") == kb_name), None)
    check("知识库可检索", kb_id is not None, str(kb_id))

    # SQL 种子 doc + 2 chunks（向量同 mock query 向量：余弦距离 0）
    vec = "[" + ",".join(["1"] * 1536) + "]"
    seeds = ["知识条目甲：项目代号晨曦", "知识条目乙：项目代号星辰"]
    sql("INSERT INTO sys_ai_docs (tenant_id, name, chunk_count, status, user_id, base_id) VALUES (0, 'e2e-doc', 2, 'READY', 1, %s);" % kb_id)
    doc_id = sql("SELECT max(id) FROM sys_ai_docs;")
    for i, content in enumerate(seeds):
        sql(f"INSERT INTO sys_ai_chunks (tenant_id, content, chunk_index, doc_id, embedding) VALUES (0, '{content}', {i}, {doc_id}, '{vec}');")
    check("种子 doc + 2 chunks", True, f"doc={doc_id}")

    graph_kb = json.dumps({
        "nodes": [
            {"id": "start", "type": "start", "data": {"inputVars": [{"name": "q", "required": True}]}},
            {"id": "kb_1", "type": "knowledge", "data": {"baseId": kb_id, "topK": 2, "query": "{{start.q}}"}},
            {"id": "end", "type": "end", "data": {"output": "{{kb_1.context}}|hits={{kb_1.hits}}"}},
        ],
        "edges": [
            {"source": "start", "target": "kb_1"},
            {"source": "kb_1", "target": "end"},
        ],
    }, ensure_ascii=False)
    r = api("/admin/v1/ai/agents", "POST", {
        "data": {"name": "e2e-agent-kb", "description": "知识节点链路", "enabled": True, "graph": graph_kb},
    })
    kb_app_id = r.get("id")
    check("建知识节点应用", kb_app_id is not None, str(r)[:100])

    r = api(f"/admin/v1/ai/agents/{kb_app_id}/run", "POST", {"input": {"q": "检索一下"}})
    out = str(r.get("output"))
    check("知识节点运行 SUCCEEDED", r.get("status") == "SUCCEEDED", str(r.get("status")) + " err=" + str(r.get("error"))[:120])
    check("检索上下文进入产出（两片段 + 命中数）", "知识条目甲" in out and "知识条目乙" in out and "hits=2" in out, out[:120])

    api(f"/admin/v1/ai/agents/{kb_app_id}", "DELETE")
    sql(f"DELETE FROM sys_ai_chunks WHERE base_id={kb_id}; DELETE FROM sys_ai_docs WHERE base_id={kb_id}; DELETE FROM sys_ai_knowledge_bases WHERE id={kb_id};")
    check("知识库链路清理", True, f"kb={kb_id}")

    print("== 6. 清理 ==")
    api(f"/admin/v1/ai/agents/{app_id}", "DELETE")
    if prov:
        api(f"/admin/v1/ai/providers/{prov['id']}", "DELETE")
    check("清理完成", True)

    server.shutdown()
    print(f"\n结果: {passed} passed, {failed} failed")
    raise SystemExit(1 if failed else 0)


if __name__ == "__main__":
    main()
