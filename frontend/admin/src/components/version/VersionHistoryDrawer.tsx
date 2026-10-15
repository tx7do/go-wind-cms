import { useCallback, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { App, Button, Drawer, Empty, List, Popconfirm, Space, Tag, Typography } from 'antd';
import {
  HistoryOutlined,
  UndoOutlined,
  EyeOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import type {
  contentservicev1_PageVersion,
  contentservicev1_PostVersion,
} from '@/api/generated/admin/service/v1';
import {
  extractSnapshot,
  fetchPageVersion,
  fetchPostVersion,
  snapshotSectionCount,
  useListPageVersions,
  useListPostVersions,
  useRestorePageVersion,
  useRestorePostVersion,
} from '@/api/hooks/content-version';
import { diffLines } from './lineDiff';

/**
 * 内容版本历史抽屉（帖子/页面通用）：
 * 版本列表（元数据）→ 查看快照 + 与当前编辑内容行级 diff → 一键恢复。
 * 恢复走后端 RestoreVersion：恢复本身会生成一个新版本，版本链完整可回溯。
 */

export interface VersionHistoryDrawerProps {
  open: boolean;
  onClose: () => void;
  kind: 'post' | 'page';
  /** 目标 id（帖子/页面），null 时不拉取 */
  targetId: number | null;
  /** 当前编辑语言（快照对比与展示用） */
  lang: string;
  /** 当前编辑器正文（diff 的"新"侧；页面无线性正文传空串） */
  currentContent: string;
  /** 恢复成功回调（父级应重新拉取表单数据） */
  onRestored: () => void;
}

type AnyVersion = contentservicev1_PostVersion | contentservicev1_PageVersion;

const formatTime = (value: unknown): string => {
  if (!value) return '-';
  try {
    return new Date(String(value)).toLocaleString();
  } catch {
    return String(value);
  }
};

export default function VersionHistoryDrawer({
  open,
  onClose,
  kind,
  targetId,
  lang,
  currentContent,
  onRestored,
}: VersionHistoryDrawerProps) {
  const { t } = useTranslation('contentVersion');
  const { message } = App.useApp();

  const postVersions = useListPostVersions(kind === 'post' ? targetId : null, {
    enabled: open && kind === 'post',
  });
  const pageVersions = useListPageVersions(kind === 'page' ? targetId : null, {
    enabled: open && kind === 'page',
  });

  const listQuery = kind === 'post' ? postVersions : pageVersions;
  const versions = useMemo<AnyVersion[]>(() => (listQuery.data?.items ?? []) as AnyVersion[], [listQuery.data]);

  const restorePost = useRestorePostVersion();
  const restorePage = useRestorePageVersion();
  const restoring = restorePost.isPending || restorePage.isPending;

  const [preview, setPreview] = useState<{
    version: AnyVersion;
    snapshot: { title: string; content: string };
    diff: ReturnType<typeof diffLines>;
  } | null>(null);
  const [previewLoading, setPreviewLoading] = useState(false);

  // 关闭时清预览态（不走 effect，避免 setState-in-effect）
  const handleClose = useCallback(() => {
    setPreview(null);
    onClose();
  }, [onClose]);

  const handlePreview = useCallback(
    async (version: AnyVersion) => {
      if (!targetId || !version.id) return;
      setPreviewLoading(true);
      try {
        const detail =
          kind === 'post'
            ? await fetchPostVersion({ id: version.id })
            : await fetchPageVersion({ id: version.id });
        const snapshot = extractSnapshot(detail, lang, kind);
        if (!snapshot) {
          message.warning(t('snapshotUnreadable'));
          return;
        }
        setPreview({
          version: detail,
          snapshot,
          diff: kind === 'post' ? diffLines(snapshot.content, currentContent) : [],
        });
      } catch (error) {
        console.error('Failed to load version detail:', error);
        message.error(t('loadDetailFailed'));
      } finally {
        setPreviewLoading(false);
      }
    },
    [kind, targetId, lang, currentContent, t, message],
  );

  const handleRestore = useCallback(
    async (version: AnyVersion) => {
      if (!targetId || !version.id) return;
      try {
        if (kind === 'post') {
          await restorePost.mutateAsync({ postId: targetId, versionId: version.id });
        } else {
          await restorePage.mutateAsync({ pageId: targetId, versionId: version.id });
        }
        message.success(t('restoreSuccess'));
        setPreview(null);
        await listQuery.refetch();
        onRestored();
      } catch (error) {
        console.error('Failed to restore version:', error);
        // 错误提示交给全局 reason→i18n 拦截器，这里只记日志
      }
    },
    [kind, targetId, restorePost, restorePage, listQuery, onRestored, t, message],
  );

  return (
    <Drawer
      open={open}
      onClose={handleClose}
      width={preview ? 760 : 480}
      title={
        <Space>
          <HistoryOutlined />
          {t('title')}
        </Space>
      }
      destroyOnHidden
    >
      {preview ? (
        <div>
          <Space style={{ marginBottom: 12 }}>
            <Button
              icon={<ReloadOutlined />}
              onClick={() => setPreview(null)}
            >
              {t('backToList')}
            </Button>
            <Tag color="blue">{`v${preview.version.versionNumber ?? '-'}`}</Tag>
            <Typography.Text type="secondary">{formatTime(preview.version.createdAt)}</Typography.Text>
          </Space>

          <Typography.Paragraph>
            <Typography.Text strong>{t('snapshotTitle')}: </Typography.Text>
            {preview.snapshot.title || '-'}
            {kind === 'page' && (
              <Typography.Text type="secondary">
                {' '}
                ({t('sectionCount', { count: snapshotSectionCount(preview.version) ?? 0 })})
              </Typography.Text>
            )}
          </Typography.Paragraph>

          {kind === 'post' ? (
            <>
              <Typography.Paragraph type="secondary">{t('diffLegend')}</Typography.Paragraph>
              <div style={{ border: '1px solid var(--ant-color-split)', borderRadius: 8, overflow: 'hidden' }}>
                {preview.diff.map((line, index) => (
                  <div
                    key={index}
                    style={{
                      padding: '1px 10px',
                      fontFamily: 'monospace',
                      fontSize: 12,
                      whiteSpace: 'pre-wrap',
                      wordBreak: 'break-all',
                      background:
                        line.type === 'added'
                          ? 'var(--ant-color-success-bg)'
                          : line.type === 'removed'
                            ? 'var(--ant-color-error-bg)'
                            : undefined,
                    }}
                  >
                    <span style={{ userSelect: 'none', marginRight: 8, opacity: 0.6 }}>
                      {line.type === 'added' ? '+' : line.type === 'removed' ? '-' : ' '}
                    </span>
                    {line.text || ' '}
                  </div>
                ))}
              </div>
            </>
          ) : (
            <Typography.Text type="secondary">{t('pageSectionRestoreNote')}</Typography.Text>
          )}

          <div style={{ marginTop: 16 }}>
            <Popconfirm
              title={t('restoreConfirm')}
              description={t('restoreConfirmDesc')}
              onConfirm={() => handleRestore(preview.version)}
            >
              <Button type="primary" icon={<UndoOutlined />} loading={restoring}>
                {t('restore')}
              </Button>
            </Popconfirm>
          </div>
        </div>
      ) : (
        <List
          loading={listQuery.isFetching || previewLoading}
          locale={{ emptyText: <Empty description={t('noVersions')} /> }}
          dataSource={versions}
          renderItem={(item) => (
            <List.Item
              actions={[
                <Button
                  key="view"
                  type="link"
                  size="small"
                  icon={<EyeOutlined />}
                  onClick={() => handlePreview(item)}
                >
                  {t('view')}
                </Button>,
                <Popconfirm
                  key="restore"
                  title={t('restoreConfirm')}
                  description={t('restoreConfirmDesc')}
                  onConfirm={() => handleRestore(item)}
                >
                  <Button type="link" size="small" icon={<UndoOutlined />}>
                    {t('restore')}
                  </Button>
                </Popconfirm>,
              ]}
            >
              <List.Item.Meta
                title={
                  <Space>
                    <Tag color="blue">{`v${item.versionNumber ?? '-'}`}</Tag>
                    {item.comment ? <Typography.Text>{item.comment}</Typography.Text> : null}
                  </Space>
                }
                description={
                  <Typography.Text type="secondary" style={{ fontSize: 12 }}>
                    {formatTime(item.createdAt)}
                    {item.createdBy ? ` · ${t('operator')} #${item.createdBy}` : ''}
                  </Typography.Text>
                }
              />
            </List.Item>
          )}
        />
      )}
    </Drawer>
  );
}
