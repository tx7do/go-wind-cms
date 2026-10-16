import { useCallback, useMemo, useState, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import {
  App,
  Avatar,
  Button,
  Drawer,
  Empty,
  Popconfirm,
  Skeleton,
  Space,
  Tag,
  Tooltip,
  Typography,
} from 'antd';
import {
  ArrowLeftOutlined,
  EyeOutlined,
  HistoryOutlined,
  ReloadOutlined,
  UndoOutlined,
  UserOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';
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
import { useListUsers } from '@/api/hooks/user';
import { PaginationQuery } from '@/core';
import { diffLines, type DiffLine } from './lineDiff';

/**
 * 内容版本历史抽屉（帖子/页面通用），交互形态对齐 Notion / GitHub 的版本历史：
 * 列表态：整行可点击进入快照预览，悬停浮现「查看/恢复」，相对时间 + 操作人头像；
 * 预览态：快照元信息 + 变更统计（+N −M）+ 行号 diff（连续未变更行折叠）→ 一键恢复。
 * 恢复走后端 RestoreVersion：恢复本身会生成一个新版本，版本链完整可回溯。
 */

dayjs.extend(relativeTime);

const ABSOLUTE_TIME_FORMAT = 'YYYY-MM-DD HH:mm';

/** 连续未变更行超过该长度时折叠为「展开」行 */
const COLLAPSE_THRESHOLD = 6;

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

/** diff 渲染分段：普通行组 / 可折叠的未变更段 */
interface DiffSegment {
  kind: 'rows' | 'collapsed';
  rows: DiffLine[];
  /** 仅 collapsed 段有值 */
  count?: number;
}

/** 版本号徽标（等宽字体小圆角块；当前版本用主色高亮） */
function VersionBadge({ value, current = false }: { value?: number; current?: boolean }) {
  return (
    <span
      className="inline-flex h-6 min-w-11 shrink-0 items-center justify-center rounded-md border px-1.5 font-mono text-xs font-semibold"
      style={
        current
          ? {
              borderColor: 'var(--ant-color-primary)',
              color: 'var(--ant-color-primary)',
              backgroundColor: 'var(--ant-color-primary-bg)',
            }
          : {
              borderColor: 'var(--ant-color-border)',
              color: 'var(--ant-color-text-secondary)',
              backgroundColor: 'var(--ant-color-fill-quaternary)',
            }
      }
    >
      {`v${value ?? '-'}`}
    </span>
  );
}

export default function VersionHistoryDrawer({
  open,
  onClose,
  kind,
  targetId,
  lang,
  currentContent,
  onRestored,
}: VersionHistoryDrawerProps) {
  const { t, i18n } = useTranslation('content-version');
  const { message } = App.useApp();

  // dayjs 相对时间文案跟随界面语言（zh* → zh-cn）
  dayjs.locale((i18n.language ?? '').startsWith('zh') ? 'zh-cn' : 'en');

  const postVersions = useListPostVersions(kind === 'post' ? targetId : null, {
    enabled: open && kind === 'post',
  });
  const pageVersions = useListPageVersions(kind === 'page' ? targetId : null, {
    enabled: open && kind === 'page',
  });

  const listQuery = kind === 'post' ? postVersions : pageVersions;
  const versions = useMemo<AnyVersion[]>(() => (listQuery.data?.items ?? []) as AnyVersion[], [listQuery.data]);
  // 列表按版本号倒序返回，取最大号作为「当前版本」标记依据
  const currentVersionNumber = useMemo(
    () => Math.max(0, ...versions.map((v) => v.versionNumber ?? 0)) || undefined,
    [versions],
  );
  const totalText = Number(listQuery.data?.total ?? versions.length);

  // 操作人昵称解析：尽力而为，接口失败/无权限时 react-query 保留错误态，UI 回退显示 #ID
  const usersQuery = useListUsers(new PaginationQuery(), {
    enabled: open,
    retry: false,
    staleTime: 5 * 60_000,
  });
  const userMap = useMemo(() => {
    const map = new Map<number, { nickname?: string; username?: string; avatar?: string }>();
    for (const u of usersQuery.data?.items ?? []) {
      if (u.id !== undefined) map.set(u.id, { nickname: u.nickname, username: u.username, avatar: u.avatar });
    }
    return map;
  }, [usersQuery.data]);

  const restorePost = useRestorePostVersion();
  const restorePage = useRestorePageVersion();
  const restoring = restorePost.isPending || restorePage.isPending;

  const [preview, setPreview] = useState<{
    version: AnyVersion;
    snapshot: { title: string; content: string };
    diff: ReturnType<typeof diffLines>;
  } | null>(null);
  const [previewLoading, setPreviewLoading] = useState(false);
  const [expandedSegments, setExpandedSegments] = useState<ReadonlySet<number>>(new Set());

  // 关闭抽屉时清预览态（不走 effect，避免 setState-in-effect）
  const handleClose = useCallback(() => {
    setPreview(null);
    setExpandedSegments(new Set());
    onClose();
  }, [onClose]);

  const handleBackToList = useCallback(() => {
    setPreview(null);
    setExpandedSegments(new Set());
  }, []);

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
        setExpandedSegments(new Set());
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
        handleBackToList();
        await listQuery.refetch();
        onRestored();
      } catch (error) {
        console.error('Failed to restore version:', error);
        // 错误提示交给全局 reason→i18n 拦截器，这里只记日志
      }
    },
    [kind, targetId, restorePost, restorePage, listQuery, onRestored, handleBackToList, t, message],
  );

  /** 操作人展示：优先昵称 → 用户名 → #ID（头像尽力取到） */
  const renderOperator = (createdBy?: number) => {
    if (createdBy === undefined) return null;
    const user = userMap.get(createdBy);
    const name = user?.nickname || user?.username || `#${createdBy}`;
    return (
      <Tooltip title={`${t('operator')}: ${name}`}>
        <span className="inline-flex min-w-0 items-center gap-1.5">
          <Avatar
            size={16}
            src={user?.avatar || undefined}
            icon={<UserOutlined />}
            style={{ backgroundColor: 'var(--ant-color-fill-secondary)', color: 'var(--ant-color-text-tertiary)' }}
          />
          <span className="truncate">{name}</span>
        </span>
      </Tooltip>
    );
  };

  const renderVersionItem = (item: AnyVersion) => {
    const isCurrent = !!item.versionNumber && item.versionNumber === currentVersionNumber;
    const time = item.createdAt ? dayjs(item.createdAt) : null;
    return (
      <div
        key={item.id}
        role="button"
        tabIndex={0}
        className="group flex cursor-pointer items-center gap-3 border-b border-[color:var(--ant-color-split)] px-5 py-3 transition-colors last:border-b-0 hover:bg-[color:var(--ant-color-fill-tertiary)]"
        onClick={() => void handlePreview(item)}
        onKeyDown={(event) => {
          if (event.key === 'Enter') void handlePreview(item);
        }}
      >
        <VersionBadge value={item.versionNumber} current={isCurrent} />
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-[color:var(--ant-color-text)]">
              {time?.isValid() ? time.fromNow() : '-'}
            </span>
            {isCurrent && (
              <Tag color="processing" className="!mr-0 !text-xs">
                {t('currentBadge')}
              </Tag>
            )}
          </div>
          <div className="truncate text-xs text-[color:var(--ant-color-text-secondary)]">
            {item.comment || <span className="italic opacity-70">{t('autosave')}</span>}
          </div>
          <div className="mt-0.5 flex items-center gap-1.5 text-xs text-[color:var(--ant-color-text-tertiary)]">
            {renderOperator(item.createdBy)}
            {time?.isValid() && (
              <>
                <span>·</span>
                <Tooltip title={time.format(ABSOLUTE_TIME_FORMAT)}>
                  <span>{time.format(ABSOLUTE_TIME_FORMAT)}</span>
                </Tooltip>
              </>
            )}
          </div>
        </div>
        <div className="flex shrink-0 items-center gap-0.5 opacity-0 transition-opacity group-focus-within:opacity-100 group-hover:opacity-100">
          <Tooltip title={t('view')}>
            <Button
              type="text"
              size="small"
              icon={<EyeOutlined />}
              onClick={(event) => {
                event.stopPropagation();
                void handlePreview(item);
              }}
            />
          </Tooltip>
          <Popconfirm
            title={t('restoreConfirm')}
            description={t('restoreConfirmDesc')}
            onConfirm={() => void handleRestore(item)}
          >
            <Tooltip title={t('restore')}>
              <Button
                type="text"
                size="small"
                icon={<UndoOutlined />}
                onClick={(event) => event.stopPropagation()}
              />
            </Tooltip>
          </Popconfirm>
        </div>
      </div>
    );
  };

  /** 把 diff 行序列切成「行组 + 可折叠未变更段」 */
  const diffSegments = useMemo<DiffSegment[]>(() => {
    if (!preview) return [];
    const segments: DiffSegment[] = [];
    let rows: DiffLine[] = [];
    let unchangedRun: DiffLine[] = [];
    const flushRun = () => {
      if (!unchangedRun.length) return;
      if (unchangedRun.length > COLLAPSE_THRESHOLD) {
        if (rows.length) {
          segments.push({ kind: 'rows', rows });
          rows = [];
        }
        segments.push({ kind: 'collapsed', rows: unchangedRun, count: unchangedRun.length });
      } else {
        rows.push(...unchangedRun);
      }
      unchangedRun = [];
    };
    for (const line of preview.diff) {
      if (line.type === 'unchanged') {
        unchangedRun.push(line);
        continue;
      }
      flushRun();
      rows.push(line);
    }
    flushRun();
    if (rows.length) segments.push({ kind: 'rows', rows });
    return segments;
  }, [preview]);

  const renderDiffLine = (line: DiffLine, key: string): ReactNode => (
    <div
      key={key}
      className="flex items-start"
      style={{
        backgroundColor:
          line.type === 'added'
            ? 'var(--ant-color-success-bg)'
            : line.type === 'removed'
              ? 'var(--ant-color-error-bg)'
              : undefined,
      }}
    >
      <span className="w-9 shrink-0 select-none border-r border-[color:var(--ant-color-split)] pr-1 text-right leading-5 text-[color:var(--ant-color-text-quaternary)]">
        {line.oldLine ?? ''}
      </span>
      <span className="w-9 shrink-0 select-none border-r border-[color:var(--ant-color-split)] pr-1 text-right leading-5 text-[color:var(--ant-color-text-quaternary)]">
        {line.newLine ?? ''}
      </span>
      <span className="w-5 shrink-0 select-none text-center leading-5 opacity-70">
        {line.type === 'added' ? '+' : line.type === 'removed' ? '−' : ''}
      </span>
      <span className="min-w-0 flex-1 whitespace-pre-wrap break-all pr-2 leading-5">{line.text || ' '}</span>
    </div>
  );

  const renderDiffBody = (): ReactNode => {
    const nodes: ReactNode[] = [];
    let collapsedIndex = -1;
    for (const segment of diffSegments) {
      if (segment.kind === 'rows') {
        segment.rows.forEach((line, index) => nodes.push(renderDiffLine(line, `r${index}-${line.oldLine}-${line.newLine}`)));
        continue;
      }
      collapsedIndex += 1;
      if (expandedSegments.has(collapsedIndex)) {
        segment.rows.forEach((line, index) => nodes.push(renderDiffLine(line, `c${collapsedIndex}-${index}-${line.oldLine}-${line.newLine}`)));
        continue;
      }
      const ordinal = collapsedIndex;
      nodes.push(
        <button
          key={`collapsed-${ordinal}`}
          type="button"
          className="block w-full cursor-pointer border-y border-[color:var(--ant-color-split)] bg-[color:var(--ant-color-fill-quaternary)] py-0.5 text-center text-xs text-[color:var(--ant-color-text-tertiary)] transition-colors hover:text-[color:var(--ant-color-primary)]"
          onClick={() => setExpandedSegments((prev) => new Set(prev).add(ordinal))}
        >
          {t('showUnchanged', { count: segment.count ?? 0 })}
        </button>,
      );
    }
    return nodes;
  };

  const previewVersionNumber = preview?.version.versionNumber;
  const previewIsCurrent =
    !!previewVersionNumber && previewVersionNumber === currentVersionNumber;
  const previewTime = preview?.version.createdAt ? dayjs(preview.version.createdAt) : null;
  const addedCount = preview?.diff.filter((line) => line.type === 'added').length ?? 0;
  const removedCount = preview?.diff.filter((line) => line.type === 'removed').length ?? 0;

  return (
    <Drawer
      open={open}
      onClose={handleClose}
      width={preview ? 760 : 480}
      destroyOnHidden
      styles={{ body: { padding: 0, display: 'flex', flexDirection: 'column' } }}
      title={
        <Space>
          <HistoryOutlined />
          {t('title')}
          {!listQuery.isPending && (
            <Tag className="!mr-0 !font-normal" bordered={false}>
              {t('versionCount', { count: totalText })}
            </Tag>
          )}
        </Space>
      }
      extra={
        <Tooltip title={t('view')}>
          <Button
            type="text"
            size="small"
            icon={<ReloadOutlined />}
            onClick={() => void listQuery.refetch()}
          />
        </Tooltip>
      }
    >
      {preview ? (
        <div className="flex min-h-0 flex-1 flex-col">
          {/* 顶栏：返回 / 当前标识 / 恢复入口 */}
          <div className="flex items-center justify-between gap-2 border-b border-[color:var(--ant-color-split)] px-4 pb-3">
            <Button type="link" className="!px-0" icon={<ArrowLeftOutlined />} onClick={handleBackToList}>
              {t('backToList')}
            </Button>
            {previewIsCurrent && (
              <Tag color="processing" className="!mr-0">
                {t('currentBadge')}
              </Tag>
            )}
            <Popconfirm
              title={t('restoreConfirm')}
              description={t('restoreConfirmDesc')}
              onConfirm={() => void handleRestore(preview.version)}
            >
              <Button type="primary" size="small" icon={<UndoOutlined />} loading={restoring}>
                {t('restore')}
              </Button>
            </Popconfirm>
          </div>

          <div className="min-h-0 flex-1 overflow-y-auto px-5 py-4">
            {/* 元信息：版本徽标 + 相对/绝对时间 + 操作人 */}
            <div className="mb-4 flex flex-wrap items-center gap-x-3 gap-y-1.5">
              <VersionBadge value={previewVersionNumber} current={previewIsCurrent} />
              {previewTime?.isValid() && (
                <Tooltip title={previewTime.format(ABSOLUTE_TIME_FORMAT)}>
                  <span className="text-xs text-[color:var(--ant-color-text-secondary)]">
                    {previewTime.fromNow()}
                  </span>
                </Tooltip>
              )}
              <span className="flex min-w-0 items-center text-xs text-[color:var(--ant-color-text-tertiary)]">
                {renderOperator(preview.version.createdBy)}
              </span>
            </div>

            {/* 快照标题（页面附区块数量） */}
            <div className="mb-4">
              <Typography.Text type="secondary" className="!text-xs">
                {t('snapshotTitle')}
              </Typography.Text>
              <div className="mt-0.5 text-base font-medium text-[color:var(--ant-color-text)]">
                {preview.snapshot.title || '-'}
                {kind === 'page' && (
                  <Tag className="!ml-2 !font-normal" bordered={false}>
                    {t('sectionCount', { count: snapshotSectionCount(preview.version) ?? 0 })}
                  </Tag>
                )}
              </div>
            </div>

            {kind === 'post' ? (
              <>
                {/* 变更统计 + 图例 */}
                <div className="mb-2 flex flex-wrap items-center gap-x-3 gap-y-1">
                  <span className="font-mono text-xs font-semibold" style={{ color: 'var(--ant-color-success)' }}>
                    {`+${addedCount}`}
                  </span>
                  <span className="font-mono text-xs font-semibold" style={{ color: 'var(--ant-color-error)' }}>
                    {`−${removedCount}`}
                  </span>
                  <span className="inline-flex items-center gap-1 text-xs text-[color:var(--ant-color-text-tertiary)]">
                    <span className="h-2 w-2 rounded-sm" style={{ backgroundColor: 'var(--ant-color-error)' }} />
                    {t('diffOnlyVersion')}
                  </span>
                  <span className="inline-flex items-center gap-1 text-xs text-[color:var(--ant-color-text-tertiary)]">
                    <span className="h-2 w-2 rounded-sm" style={{ backgroundColor: 'var(--ant-color-success)' }} />
                    {t('diffOnlyCurrent')}
                  </span>
                </div>
                <div className="overflow-hidden rounded-lg border border-[color:var(--ant-color-split)] py-1 font-mono text-xs">
                  {renderDiffBody()}
                </div>
              </>
            ) : (
              <Typography.Text type="secondary">{t('pageSectionRestoreNote')}</Typography.Text>
            )}
          </div>
        </div>
      ) : (
        <div className="flex min-h-0 flex-1 flex-col overflow-y-auto">
          {listQuery.isPending || (previewLoading && !preview) ? (
            <div className="px-5 py-3">
              <Skeleton active paragraph={{ rows: 6 }} />
            </div>
          ) : versions.length === 0 ? (
            <div className="flex flex-1 items-center justify-center">
              <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={t('noVersions')} />
            </div>
          ) : (
            versions.map(renderVersionItem)
          )}
        </div>
      )}
    </Drawer>
  );
}
