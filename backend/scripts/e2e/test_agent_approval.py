"""AI Agent 人工审批节点 E2E：approval 挂起 → ApproveRun 批准/驳回。

前置（同 test_ai_agent_run.py）：core/admin 已启动（本特性代码版本）。

断言（12 条）：
  1. 含 approval 节点的图过静态校验（建应用成功）
  2. Run 到 approval 挂起：status=WAITING_APPROVAL
  3. 挂起运行带 pending_node_id 与 resume_state
  4. trace 记录到挂起节点为止
  5. LLM 产物在挂起前已产出（output 暂空、tokens 已计）
  6. 批准后从下游继续：SUCCEEDED + output 含审批结论 note
  7. 批准后 tokens 累加（挂起前+恢复后）
  8. 第二次运行走驳回：FAILED + error 含审批意见
  9. 对非挂起运行 ApproveRun 被拒（400）
 10. 双审批节点串联：第一次批准→第二个 approval 再挂起→第二次批准→SUCCEEDED
 11. 清理完整
"""
import json
import sys
import threading
import time
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

sys.path.insert(0, __file__.rsplit("\\", 1)[0])
from common import AdminApi, admin_login  # noqa: E402

MOCK_PORT = 6903
MOCK_ANSWER = "mock-approval-answer"
PASS = 0


def check(name, cond, extra=""):
    global PASS
    assert cond, f"FAIL: {name} {extra}"
    PASS += 1
    print(f"  ok {PASS:02d} - {name} {extra}")


def sse_chunk(content, finish_reason, usage=None):
    """OpenAI 流式分块（工具循环走 CreateChatCompletionStream，非流式 JSON 会静默空答案）。"""
    chunk = {
        "id": "chatcmpl-mock",
        "object": "chat.completion.chunk",
        "created": int(time.time()),
        "model": "mock-approval-model",
        "choices": [{"index": 0, "delta": {}, "finish_reason": finish_reason}],
    }
    if content is not None:
        chunk["choices"][0]["delta"] = {"role": "assistant", "content": content}
    if usage is not None:
        chunk["usage"] = usage
    return "data: " + json.dumps(chunk) + "\n"


class MockOpenAIHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers.get("Content-Length", 0))
        self.rfile.read(length)

        self.send_response(200)
        self.send_header("Content-Type", "text/event-stream; charset=utf-8")
        self.send_header("Cache-Control", "no-cache")
        self.end_headers()
        self.wfile.write(sse_chunk(MOCK_ANSWER[:6], None).encode())
        self.wfile.write(sse_chunk(MOCK_ANSWER[6:], "stop", {
            "prompt_tokens": 10, "completion_tokens": 8, "total_tokens": 18,
        }).encode())
        self.wfile.write(b"data: [DONE]\n")

    def log_message(self, *args):
        pass


def main():
    server = ThreadingHTTPServer(("127.0.0.1", MOCK_PORT), MockOpenAIHandler)
    threading.Thread(target=server.serve_forever, daemon=True).start()
    print(f"mock openai at :{MOCK_PORT}")

    hdr = admin_login("admin", "admin")
    api = AdminApi(hdr)
    TS = int(time.time())

    # ---- mock provider（唯一名强制新建，避免复用脏 baseUrl）----
    prov_name = f"e2e-approval-provider-{TS}"
    r = api("/admin/v1/ai/providers", "POST", {
        "data": {
            "name": prov_name,
            "baseUrl": f"http://127.0.0.1:{MOCK_PORT}/v1",
            "apiKey": "sk-mock",
            "modelName": "mock-approval-model",
            "modelType": "CLOUD",
            "isEnabled": True,
            "isDefault": False,
        },
    })
    assert "error" not in r, str(r)[:120]
    lst = api("/admin/v1/ai/providers?page=1&pageSize=50")
    prov = next((it for it in lst.get("items", []) if it.get("name") == prov_name), None)
    provider_id = prov["id"]

    # ---- 1. 建图：start → llm → approval → condition(approved) → 双 end ----
    graph = json.dumps({
        "nodes": [
            {"id": "start", "type": "start", "data": {"inputVars": [{"name": "q", "required": True}]}},
            {"id": "llm1", "type": "llm", "data": {"providerId": provider_id, "systemPrompt": "测试", "userPrompt": "{{start.q}}"}},
            {"id": "appr1", "type": "approval", "data": {"label": "发布前审批"}},
            {"id": "cond", "type": "condition", "data": {"expression": 'appr1.approved == "true"'}},
            {"id": "end_ok", "type": "end", "data": {"output": "已发布：{{llm1.output}}（{{appr1.note}}）"}},
            {"id": "end_no", "type": "end", "data": {"output": "已驳回"}},
        ],
        "edges": [
            {"source": "start", "target": "llm1"},
            {"source": "llm1", "target": "appr1"},
            {"source": "appr1", "target": "cond"},
            {"source": "cond", "target": "end_ok", "branch": "true"},
            {"source": "cond", "target": "end_no", "branch": "false"},
        ],
    }, ensure_ascii=False)
    r = api("/admin/v1/ai/agents", "POST", {
        "data": {"name": f"e2e-agent-approval-{TS}", "description": "审批节点 E2E", "enabled": True, "graph": graph},
    })
    app_id = r.get("id")
    check(bool(app_id), "含 approval 节点的图过静态校验", f"id={app_id}")

    # ---- 2. 首跑：挂起 ----
    r = api(f"/admin/v1/ai/agents/{app_id}/run", "POST", {"input": {"q": "审批流测试"}})
    run_id = r.get("id")
    status = str(r.get("status"))
    check("WAITING_APPROVAL" in status, "Run 挂起为 WAITING_APPROVAL", status)
    check(bool(r.get("pendingNodeId")) or bool(r.get("pending_node_id")),
          "挂起带 pendingNodeId", str(r.get("pendingNodeId") or r.get("pending_node_id")))
    check(bool(r.get("resumeState") or r.get("resume_state")), "挂起带 resumeState")
    trace = json.loads(r.get("trace") or "[]")
    check([s.get("nodeId") for s in trace] == ["start", "llm1"],
          "trace 记录到挂起前（start/llm1）", str([s.get("nodeId") for s in trace]))
    check(int(r.get("totalTokens") or 0) == 18, "挂起前 LLM tokens 已计", str(r.get("totalTokens")))

    # ---- 3. 批准 → 继续 ----
    r = api(f"/admin/v1/ai/agent-runs/{run_id}/approve", "POST",
            {"runId": run_id, "approved": True, "note": "同意发布"})
    check("SUCCEEDED" in str(r.get("status")), "批准后 SUCCEEDED", str(r.get("status")))
    check("已发布：mock-approval-answer（同意发布）" in str(r.get("output")),
          "产出含审批结论", str(r.get("output"))[:80])
    # 恢复段无 LLM 节点：tokens 应保持挂起时的 18（挂起+终态不重复计）
    check(int(r.get("totalTokens") or 0) == 18, "tokens 保持 18（恢复段无新增）", str(r.get("totalTokens")))
    trace = json.loads(r.get("trace") or "[]")
    check([s.get("nodeId") for s in trace][-2:] == ["cond", "end_ok"],
          "trace 追加恢复后节点", str([s.get("nodeId") for s in trace]))

    # ---- 4. 驳回流 ----
    r = api(f"/admin/v1/ai/agents/{app_id}/run", "POST", {"input": {"q": "第二条"}})
    run2 = r.get("id")
    check("WAITING_APPROVAL" in str(r.get("status")), "第二条挂起", str(r.get("status")))
    r = api(f"/admin/v1/ai/agent-runs/{run2}/approve", "POST",
            {"runId": run2, "approved": False, "note": "质量不达标"})
    check("FAILED" in str(r.get("status")), "驳回置 FAILED", str(r.get("status")))
    check("质量不达标" in str(r.get("error") or ""), "驳回理由进 error", str(r.get("error"))[:60])

    # ---- 5. 非 WAITING_APPROVAL 的 ApproveRun 被拒 ----
    r = api(f"/admin/v1/ai/agent-runs/{run_id}/approve", "POST",
            {"runId": run_id, "approved": True, "note": ""})
    check("_raw" in r or r.get("code") == 400, "非挂起运行 ApproveRun 被拒", str(r)[:80])

    # ---- 6. 双审批节点串联 ----
    graph2 = json.dumps({
        "nodes": [
            {"id": "start", "type": "start", "data": {"inputVars": [{"name": "q", "required": True}]}},
            {"id": "a1", "type": "approval", "data": {}},
            {"id": "a2", "type": "approval", "data": {}},
            {"id": "end", "type": "end", "data": {"output": "双审通过"}},
        ],
        "edges": [
            {"source": "start", "target": "a1"},
            {"source": "a1", "target": "a2"},
            {"source": "a2", "target": "end"},
        ],
    })
    r = api("/admin/v1/ai/agents", "POST", {
        "data": {"name": f"e2e-agent-approval2-{TS}", "enabled": True, "graph": graph2},
    })
    app2 = r.get("id")
    r = api(f"/admin/v1/ai/agents/{app2}/run", "POST", {"input": {"q": "双审批"}})
    run3 = r.get("id")
    check("WAITING_APPROVAL" in str(r.get("status")), "双审批：第一次挂起", str(r.get("status")))
    r = api(f"/admin/v1/ai/agent-runs/{run3}/approve", "POST", {"runId": run3, "approved": True, "note": ""})
    check("WAITING_APPROVAL" in str(r.get("status")), "第一审批准→第二审再挂起", str(r.get("status")))
    check(str(r.get("pendingNodeId") or r.get("pending_node_id")) == "a2",
          "第二挂起节点为 a2", str(r.get("pendingNodeId") or r.get("pending_node_id")))
    r = api(f"/admin/v1/ai/agent-runs/{run3}/approve", "POST", {"runId": run3, "approved": True, "note": ""})
    check("SUCCEEDED" in str(r.get("status")) and "双审通过" in str(r.get("output")),
          "第二审批准→SUCCEEDED", str(r.get("output"))[:60])

    # ---- 清理 ----
    api(f"/admin/v1/ai/agents/{app_id}", "DELETE")
    api(f"/admin/v1/ai/agents/{app2}", "DELETE")

    print(f"\nALL PASS ({PASS} assertions)")


if __name__ == "__main__":
    main()
