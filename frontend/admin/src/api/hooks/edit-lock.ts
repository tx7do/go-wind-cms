import { useEffect, useRef, useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { apiClient } from '@/api/client';

// ==============================
// 编辑锁（协同编辑 v1）：编辑感知 + 冲突防护
//
// 进入编辑页即 Acquire（他人持锁时返回归属，前端展示提示横幅）；
// 每 60s 重发 Acquire 兼作心跳续期；卸载时 best-effort Release。
// 锁是 Redis 120s 软锁，只做感知不拦保存（版本化兜底冲突）。
// ==============================

export interface EditLockInfo {
  /** 非空表示当前被该用户持锁（非本人） */
  heldByName?: string;
  heldSince?: string;
}

const EDIT_LOCK_HEARTBEAT_MS = 60_000;

export function useEditLock(targetType: 'post' | 'page', targetId?: number | null) {
  const [lock, setLock] = useState<EditLockInfo>({});

  const acquireMutation = useMutation({
    mutationFn: (req: { targetType: string; targetId: number }) =>
      apiClient.editLockService.Acquire(req as never),
  });
  const releaseMutation = useMutation({
    mutationFn: (req: { targetType: string; targetId: number }) =>
      apiClient.editLockService.Release(req as never),
  });

  const targetRef = useRef<{ targetType: string; targetId: number } | undefined>(undefined);
  targetRef.current = targetId ? { targetType, targetId } : undefined;

  useEffect(() => {
    if (!targetRef.current) return undefined;
    let disposed = false;

    const acquire = async () => {
      const target = targetRef.current;
      if (!target) return;
      try {
        const resp = await acquireMutation.mutateAsync(target);
        if (disposed) return;
        setLock({
          heldByName: resp?.acquired ? undefined : resp?.holderName ?? undefined,
          heldSince: resp?.acquiredAt ? new Date(String(resp.acquiredAt)).toLocaleString() : undefined,
        });
      } catch (error) {
        console.error('Acquire edit lock failed:', error);
      }
    };

    const release = () => {
      const target = targetRef.current;
      if (!target) return;
      releaseMutation.mutateAsync(target).catch((error: Error) =>
        console.error('Release edit lock failed:', error),
      );
    };

    void acquire();
    const timer = setInterval(() => void acquire(), EDIT_LOCK_HEARTBEAT_MS);
    return () => {
      disposed = true;
      clearInterval(timer);
      release();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [targetType, targetId]);

  return lock;
}
