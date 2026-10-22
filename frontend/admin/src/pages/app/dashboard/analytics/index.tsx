import { useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import { useNavigate } from 'react-router-dom';
import { Button, Card, Col, Row, Spin, Tabs } from 'antd';
import {
  AuditOutlined,
  ClockCircleOutlined,
  CommentOutlined,
  DatabaseOutlined,
  FileTextOutlined,
  HeartOutlined,
  PictureOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { useI18n } from '@/core/i18n';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import AiInsightsCard from '../components/AiInsightsCard';
import { useUserStore } from '@/stores/user';
import {
  useGetAiUsageTrend,
  useGetApiCallTrend,
  useGetCommentStatusDistribution,
  useGetContentTrend,
  useGetDashboardOverview,
  useGetInteractionStats,
  useGetLoginActivity,
  useGetMediaAssetTypeDistribution,
  useGetPostStatusDistribution,
} from '@/api/hooks/dashboard';

// proto3 int64 经 JSON 序列化为字符串("2")，生成器类型标注为 number 与运行时
// 不符，统一转 number（对齐 vben 侧 toNum）。
function toNum(v: number | string | undefined): number {
  return Number(v ?? 0);
}

// 日期截短为 MM-DD，减小 x 轴密度
function shortDate(date: string | undefined): string {
  return date?.length === 10 ? date.slice(5) : (date ?? '');
}

// 字节数人性化展示（存储用量卡主数字/增量）
function formatBytes(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes <= 0) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  const v = bytes / 1024 ** i;
  return `${v >= 100 ? v.toFixed(0) : v.toFixed(1)} ${units[i]}`;
}

// 无标题帖子与超长标题的兜底展示
function displayTitle(title: string | undefined, fallback: string): string {
  if (!title) return fallback;
  return title.length > 16 ? `${title.slice(0, 16)}…` : title;
}

const chartBaseStyle = { height: 300, width: '100%' } as const;
const spinStyle = { display: 'block', margin: '80px auto' } as const;

/** 概览卡：主数字为总量，footer 为近 7 天新增（weeklyValueText 优先于 weeklyValue） */
function OverviewCard({
  icon,
  title,
  value,
  weeklyTitle,
  weeklyValue,
  weeklyValueText,
  tone,
}: {
  icon: React.ReactNode;
  title: string;
  value: number | string;
  weeklyTitle: string;
  weeklyValue?: number;
  weeklyValueText?: string;
  tone: string;
}) {
  return (
    <Card styles={{ body: { padding: '20px 24px' } }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
        <div
          style={{
            width: 48,
            height: 48,
            borderRadius: '50%',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: 22,
            color: '#fff',
            background: tone,
            flexShrink: 0,
          }}
        >
          {icon}
        </div>
        <div style={{ flex: 1, minWidth: 0 }}>
          <div style={{ fontSize: 13, color: 'var(--ant-color-text-secondary)' }}>
            {title}
          </div>
          <div style={{ fontSize: 26, fontWeight: 600, lineHeight: 1.3 }}>
            {typeof value === 'number' ? value.toLocaleString() : value}
          </div>
        </div>
      </div>
      <div
        style={{
          marginTop: 12,
          paddingTop: 10,
          borderTop: '1px solid var(--ant-color-split)',
          fontSize: 12,
          color: 'var(--ant-color-text-tertiary)',
        }}
      >
        {weeklyTitle}
        <span
          style={{
            marginLeft: 8,
            color: 'var(--ant-color-success)',
            fontWeight: 500,
          }}
        >
          +{weeklyValueText ?? (weeklyValue ?? 0).toLocaleString()}
        </span>
      </div>
    </Card>
  );
}

/** 内容增长趋势（近 30 天新增用户/帖子/评论三条平滑面积线） */
const ContentTrendChart = () => {
  const { t } = useI18n('dashboard');
  const { data, isLoading } = useGetContentTrend();

  const option = useMemo(() => {
    const seriesDef = [
      { key: 'users' as const, color: '#5ab1ef', name: t('analytics.newUsers') },
      { key: 'posts' as const, color: '#019680', name: t('analytics.newPosts') },
      {
        key: 'comments' as const,
        color: '#d48265',
        name: t('analytics.newComments'),
      },
    ];
    const dates = (data?.users ?? []).map((item) => shortDate(item.date));
    return {
      grid: { bottom: 40, containLabel: true, left: '1%', right: '1%', top: '5%' },
      legend: { bottom: 0, data: seriesDef.map((s) => s.name) },
      series: seriesDef.map((s) => ({
        areaStyle: {},
        data: (data?.[s.key] ?? []).map((item) => toNum(item.value)),
        itemStyle: { color: s.color },
        name: s.name,
        smooth: true,
        type: 'line',
      })),
      tooltip: { axisPointer: { type: 'line' }, trigger: 'axis' },
      xAxis: {
        axisTick: { show: false },
        boundaryGap: false,
        data: dates,
        type: 'category',
      },
      yAxis: { minInterval: 1, type: 'value' },
    };
  }, [data, t]);

  if (isLoading) return <Spin style={spinStyle} />;
  return <ReactECharts option={option} notMerge style={chartBaseStyle} />;
};

/** 互动 TOP 榜（点赞数横向条形图，TOP1 在顶部） */
const InteractionTopChart = () => {
  const { t } = useI18n('dashboard');
  const { data, isLoading } = useGetInteractionStats();

  const option = useMemo(() => {
    const items = data?.topLikedPosts ?? [];
    // echarts 的 y 类轴自下而上，reverse 使 TOP1 显示在顶部
    const titles = items
      .map((item) => displayTitle(item.title, t('analytics.titleFallback')))
      .reverse();
    const counts = items.map((item) => toNum(item.likeCount)).reverse();
    const fullTitles = items
      .map((item) => item.title || t('analytics.titleFallback'))
      .reverse();
    return {
      grid: { bottom: 0, containLabel: true, left: '1%', right: '5%', top: '2%' },
      series: [
        {
          data: counts,
          itemStyle: { color: '#019680' },
          label: { position: 'right', show: true },
          name: t('analytics.totalLikes'),
          type: 'bar',
        },
      ],
      tooltip: {
        formatter: (params: unknown) => {
          const p = params as { dataIndex: number } | { dataIndex: number }[];
          const idx = Array.isArray(p) ? p[0].dataIndex : p.dataIndex;
          return `${fullTitles[idx]}<br/>${t('analytics.totalLikes')}: ${counts[idx]}`;
        },
        trigger: 'axis',
      },
      xAxis: { minInterval: 1, type: 'value' },
      yAxis: {
        axisTick: { show: false },
        data: titles,
        type: 'category',
      },
    };
  }, [data, t]);

  if (isLoading) return <Spin style={spinStyle} />;
  return <ReactECharts option={option} notMerge style={chartBaseStyle} />;
};

/** 浏览 TOP 榜（真实 PV，VIEW 计数横向条形图，TOP1 在顶部） */
const ViewTopChart = () => {
  const { t } = useI18n('dashboard');
  const { data, isLoading } = useGetInteractionStats();

  const option = useMemo(() => {
    const items = data?.topViewedPosts ?? [];
    const titles = items
      .map((item) => displayTitle(item.title, t('analytics.titleFallback')))
      .reverse();
    const counts = items.map((item) => toNum(item.likeCount)).reverse();
    const fullTitles = items
      .map((item) => item.title || t('analytics.titleFallback'))
      .reverse();
    return {
      grid: { bottom: 0, containLabel: true, left: '1%', right: '5%', top: '2%' },
      series: [
        {
          data: counts,
          itemStyle: { color: '#5ab1ef' },
          label: { position: 'right', show: true },
          name: t('analytics.totalViews'),
          type: 'bar',
        },
      ],
      tooltip: {
        formatter: (params: unknown) => {
          const p = params as { dataIndex: number } | { dataIndex: number }[];
          const idx = Array.isArray(p) ? p[0].dataIndex : p.dataIndex;
          return `${fullTitles[idx]}<br/>${t('analytics.totalViews')}: ${counts[idx]}`;
        },
        trigger: 'axis',
      },
      xAxis: { minInterval: 1, type: 'value' },
      yAxis: {
        axisTick: { show: false },
        data: titles,
        type: 'category',
      },
    };
  }, [data, t]);

  if (isLoading) return <Spin style={spinStyle} />;
  return <ReactECharts option={option} notMerge style={chartBaseStyle} />;
};

/** 环形分布图数据项 */
interface DonutDatum {
  color?: string;
  name: string;
  value: number;
}

/**
 * 环形分布图：无扇区标注（长文案被截断、零值标注堆顶、压图例），中心总量走
 * HTML 叠加层（canvas 文字不随 antd 明暗主题换色），底部图例带计数。
 */
const DonutChart = ({ data, unitText }: { data: DonutDatum[]; unitText: string }) => {
  const total = useMemo(() => data.reduce((sum, d) => sum + d.value, 0), [data]);
  const byName = useMemo(
    () => Object.fromEntries(data.map((d) => [d.name, d.value])) as Record<string, number>,
    [data],
  );

  const option = useMemo(
    () => ({
      legend: {
        bottom: 0,
        formatter: (name: string) => `${name} ${byName[name] ?? 0}`,
        icon: 'circle',
        itemHeight: 8,
        itemWidth: 8,
      },
      series: [
        {
          avoidLabelOverlap: true,
          center: ['50%', '42%'],
          data: data.map((d) => ({
            itemStyle: d.color ? { color: d.color } : undefined,
            name: d.name,
            value: d.value,
          })),
          itemStyle: { borderColor: 'transparent', borderRadius: 5, borderWidth: 2 },
          label: { show: false },
          labelLine: { show: false },
          radius: ['50%', '74%'],
          type: 'pie',
        },
      ],
      tooltip: { trigger: 'item' },
    }),
    [data, byName],
  );

  return (
    <div style={{ position: 'relative' }}>
      <ReactECharts option={option} notMerge style={chartBaseStyle} />
      <div
        style={{
          alignItems: 'center',
          display: 'flex',
          flexDirection: 'column',
          inset: 0,
          justifyContent: 'center',
          pointerEvents: 'none',
          position: 'absolute',
          transform: 'translateY(-8%)',
        }}
      >
        <span style={{ fontSize: 24, fontWeight: 600, lineHeight: 1.2 }}>
          {total.toLocaleString()}
        </span>
        <span style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)' }}>
          {unitText}
        </span>
      </div>
    </div>
  );
};

/** 互动类型分布（点赞/收藏/浏览环形，中心为累计总量） */
const InteractionDistributionChart = () => {
  const { t } = useI18n('dashboard');
  const { data, isLoading } = useGetInteractionStats();

  if (isLoading) return <Spin style={spinStyle} />;
  const items: DonutDatum[] = [
    { color: '#019680', name: t('analytics.totalLikes'), value: toNum(data?.totalLikes) },
    { color: '#5ab1ef', name: t('analytics.totalWatches'), value: toNum(data?.totalWatches) },
    { color: '#d48265', name: t('analytics.totalViews'), value: toNum(data?.totalViews) },
  ];
  return <DonutChart data={items} unitText={t('analytics.unitInteractions')} />;
};

/** 登录活跃趋势（近 30 天成功/失败/活跃用户三线） */
const LoginActivityChart = () => {
  const { t } = useI18n('dashboard');
  const { data, isLoading } = useGetLoginActivity();

  const option = useMemo(() => {
    const success = data?.success ?? [];
    const seriesDef = [
      { data: success, color: '#019680', name: t('analytics.loginSuccess') },
      {
        data: data?.failed ?? [],
        color: '#d48265',
        name: t('analytics.loginFailed'),
      },
      {
        data: data?.activeUsers ?? [],
        color: '#9254de',
        name: t('analytics.activeUsers'),
      },
    ];
    return {
      grid: { bottom: 40, containLabel: true, left: '1%', right: '1%', top: '5%' },
      legend: { bottom: 0, data: seriesDef.map((s) => s.name) },
      series: seriesDef.map((s) => ({
        areaStyle: {},
        data: s.data.map((item) => toNum(item.value)),
        itemStyle: { color: s.color },
        name: s.name,
        smooth: true,
        type: 'line',
      })),
      tooltip: { axisPointer: { type: 'line' }, trigger: 'axis' },
      xAxis: {
        axisTick: { show: false },
        boundaryGap: false,
        data: success.map((item) => shortDate(item.date)),
        type: 'category',
      },
      yAxis: { minInterval: 1, type: 'value' },
    };
  }, [data, t]);

  if (isLoading) return <Spin style={spinStyle} />;
  return <ReactECharts option={option} notMerge style={chartBaseStyle} />;
};

/** AI 用量趋势（近 30 天调用次数柱 + Token 消耗线，双 Y 轴） */
const AiUsageChart = () => {
  const { t } = useI18n('dashboard');
  const { data, isLoading } = useGetAiUsageTrend();

  const option = useMemo(() => {
    const calls = data?.calls ?? [];
    const tokens = data?.tokens ?? [];
    return {
      grid: { bottom: 40, containLabel: true, left: '1%', right: '1%', top: '5%' },
      legend: {
        bottom: 0,
        data: [t('analytics.aiCalls'), t('analytics.aiTokens')],
      },
      series: [
        {
          barMaxWidth: 24,
          data: calls.map((item) => toNum(item.value)),
          itemStyle: { color: '#5ab1ef' },
          name: t('analytics.aiCalls'),
          type: 'bar',
          yAxisIndex: 0,
        },
        {
          areaStyle: {},
          data: tokens.map((item) => toNum(item.value)),
          itemStyle: { color: '#d48265' },
          name: t('analytics.aiTokens'),
          smooth: true,
          type: 'line',
          yAxisIndex: 1,
        },
      ],
      tooltip: { axisPointer: { type: 'shadow' }, trigger: 'axis' },
      xAxis: {
        axisTick: { show: false },
        boundaryGap: true,
        data: calls.map((item) => shortDate(item.date)),
        type: 'category',
      },
      yAxis: [
        { minInterval: 1, type: 'value' },
        { splitLine: { show: false }, type: 'value' },
      ],
    };
  }, [data, t]);

  if (isLoading) return <Spin style={spinStyle} />;
  return <ReactECharts option={option} notMerge style={chartBaseStyle} />;
};

/** API 调用量趋势（近 30 天网关请求单线面积图） */
const ApiCallChart = () => {
  const { t } = useI18n('dashboard');
  const { data, isLoading } = useGetApiCallTrend();

  const option = useMemo(() => {
    const calls = data?.calls ?? [];
    return {
      grid: { bottom: 40, containLabel: true, left: '1%', right: '1%', top: '5%' },
      legend: { bottom: 0, data: [t('analytics.apiCalls')] },
      series: [
        {
          areaStyle: {},
          data: calls.map((item) => toNum(item.value)),
          itemStyle: { color: '#019680' },
          name: t('analytics.apiCalls'),
          smooth: true,
          type: 'line',
        },
      ],
      tooltip: { axisPointer: { type: 'line' }, trigger: 'axis' },
      xAxis: {
        axisTick: { show: false },
        boundaryGap: false,
        data: calls.map((item) => shortDate(item.date)),
        type: 'category',
      },
      yAxis: { minInterval: 1, type: 'value' },
    };
  }, [data, t]);

  if (isLoading) return <Spin style={spinStyle} />;
  return <ReactECharts option={option} notMerge style={chartBaseStyle} />;
};

// 帖子状态 → 环形图颜色/i18n 键；未识别状态回退原始值展示
const POST_STATUS_META: Record<string, { color: string; labelKey: string }> = {
  POST_STATUS_DRAFT: { color: '#5ab1ef', labelKey: 'analytics.statusDraft' },
  POST_STATUS_PUBLISHED: { color: '#019680', labelKey: 'analytics.statusPublished' },
  POST_STATUS_SCHEDULED: { color: '#faad14', labelKey: 'analytics.statusScheduled' },
};

/** 帖子状态分布（环形图，排除回收站） */
const PostStatusChart = () => {
  const { t } = useI18n('dashboard');
  const { data, isLoading } = useGetPostStatusDistribution();

  if (isLoading) return <Spin style={spinStyle} />;
  const items: DonutDatum[] = (data?.items ?? []).map((item) => {
    const meta = item.status ? POST_STATUS_META[item.status] : undefined;
    return {
      color: meta?.color,
      name: meta ? t(meta.labelKey) : (item.status ?? ''),
      value: toNum(item.count),
    };
  });
  return <DonutChart data={items} unitText={t('analytics.unitPosts')} />;
};

// 评论状态 → 环形图颜色/i18n 键（后端把 NULL 状态归并为已通过）
const COMMENT_STATUS_META: Record<string, { color: string; labelKey: string }> = {
  STATUS_PENDING: { color: '#faad14', labelKey: 'analytics.statusPending' },
  STATUS_APPROVED: { color: '#019680', labelKey: 'analytics.statusApproved' },
  STATUS_REJECTED: { color: '#ff7875', labelKey: 'analytics.statusRejected' },
  STATUS_SPAM: { color: '#8c8c8c', labelKey: 'analytics.statusSpam' },
};

/** 评论状态分布（环形图，NULL 状态按已通过计） */
const CommentStatusChart = () => {
  const { t } = useI18n('dashboard');
  const { data, isLoading } = useGetCommentStatusDistribution();

  if (isLoading) return <Spin style={spinStyle} />;
  const items: DonutDatum[] = (data?.items ?? []).map((item) => {
    const meta = item.status ? COMMENT_STATUS_META[item.status] : undefined;
    return {
      color: meta?.color,
      name: meta ? t(meta.labelKey) : (item.status ?? ''),
      value: toNum(item.count),
    };
  });
  return <DonutChart data={items} unitText={t('analytics.unitComments')} />;
};

// 媒体类型 → 环形图颜色/i18n 键（后端把 NULL 类型归并为其他）
const MEDIA_TYPE_META: Record<string, { color: string; labelKey: string }> = {
  ASSET_TYPE_IMAGE: { color: '#5ab1ef', labelKey: 'analytics.assetImage' },
  ASSET_TYPE_VIDEO: { color: '#722ed1', labelKey: 'analytics.assetVideo' },
  ASSET_TYPE_DOCUMENT: { color: '#faad14', labelKey: 'analytics.assetDocument' },
  ASSET_TYPE_AUDIO: { color: '#d48265', labelKey: 'analytics.assetAudio' },
  ASSET_TYPE_ARCHIVE: { color: '#8c8c8c', labelKey: 'analytics.assetArchive' },
  ASSET_TYPE_OTHER: { color: '#13c2c2', labelKey: 'analytics.assetOther' },
};

/** 媒体类型分布（环形图，NULL 类型按其他计） */
const MediaTypeChart = () => {
  const { t } = useI18n('dashboard');
  const { data, isLoading } = useGetMediaAssetTypeDistribution();

  if (isLoading) return <Spin style={spinStyle} />;
  const items: DonutDatum[] = (data?.items ?? []).map((item) => {
    const meta = item.status ? MEDIA_TYPE_META[item.status] : undefined;
    return {
      color: meta?.color,
      name: meta ? t(meta.labelKey) : (item.status ?? ''),
      value: toNum(item.count),
    };
  });
  return <DonutChart data={items} unitText={t('analytics.unitAssets')} />;
};

/** 待办提醒条：待审评论/草稿/定时发布，点击跳转对应列表页 */
const TodoStrip = () => {
  const { t } = useI18n('dashboard');
  const navigate = useNavigate();
  const { data, isLoading } = useGetDashboardOverview();

  if (isLoading) return null;
  const items = [
    {
      count: toNum(data?.pendingCommentCount),
      icon: <AuditOutlined />,
      key: 'pendingComments',
      path: '/engagement/comments',
    },
    {
      count: toNum(data?.draftPostCount),
      icon: <FileTextOutlined />,
      key: 'draftPosts',
      path: '/content/posts',
    },
    {
      count: toNum(data?.scheduledPostCount),
      icon: <ClockCircleOutlined />,
      key: 'scheduledPosts',
      path: '/content/posts',
    },
  ];
  return (
    <Card styles={{ body: { padding: '12px 24px' } }}>
      <div style={{ alignItems: 'center', display: 'flex', flexWrap: 'wrap', gap: 12 }}>
        <span style={{ color: 'var(--ant-color-text-secondary)', fontSize: 13 }}>
          {t('analytics.todoTitle')}
        </span>
        {items.map((item) => (
          <Button
            key={item.key}
            onClick={() => navigate(item.path)}
            style={{ paddingInline: 12 }}
            type="text"
          >
            {item.icon}
            <span style={{ marginInline: 6 }}>{t(`analytics.${item.key}`)}</span>
            <span
              style={{
                color: item.count > 0 ? 'var(--ant-color-warning)' : 'var(--ant-color-text-tertiary)',
                fontSize: 16,
                fontWeight: 600,
              }}
            >
              {item.count.toLocaleString()}
            </span>
          </Button>
        ))}
      </div>
    </Card>
  );
};

/**
 * 分析页（对齐 vben analytics）：概览卡（6 张：内容/互动/媒体/存储）+ 待办提醒
 * + AI 洞察（平台专属）+ 趋势 Tab 卡（内容/登录/AI 用量/API 调用）+ 互动 TOP 榜
 * / 互动·帖子·评论状态分布。数据源 statsService 各 RPC。
 */
const Dashboard = () => {
  const { t } = useI18n('dashboard');
  const overviewQuery = useGetDashboardOverview();
  const o = overviewQuery.data;

  const overviewItems: Array<{
    icon: React.ReactNode;
    title: string;
    tone: string;
    value: number | string;
    weeklyValue?: number;
    weeklyValueText?: string;
  }> = [
    {
      icon: <UserOutlined />,
      title: t('analytics.userCountTotal'),
      tone: '#409eff',
      value: toNum(o?.userCount),
      weeklyValue: toNum(o?.newUserCountWeek),
    },
    {
      icon: <FileTextOutlined />,
      title: t('analytics.postCountTotal'),
      tone: '#13c2c2',
      value: toNum(o?.postCount),
      weeklyValue: toNum(o?.newPostCountWeek),
    },
    {
      icon: <CommentOutlined />,
      title: t('analytics.commentCountTotal'),
      tone: '#722ed1',
      value: toNum(o?.commentCount),
      weeklyValue: toNum(o?.newCommentCountWeek),
    },
    {
      icon: <HeartOutlined />,
      title: t('analytics.interactionCountTotal'),
      tone: '#52c41a',
      value: toNum(o?.interactionCount),
      weeklyValue: toNum(o?.newLikeCountWeek),
    },
    {
      icon: <PictureOutlined />,
      title: t('analytics.mediaCountTotal'),
      tone: '#fa8c16',
      value: toNum(o?.mediaAssetCount),
      weeklyValue: toNum(o?.newMediaCountWeek),
    },
    {
      icon: <DatabaseOutlined />,
      title: t('analytics.storageUsed'),
      tone: '#eb2f96',
      value: formatBytes(toNum(o?.mediaStorageBytes)),
      weeklyValueText: formatBytes(toNum(o?.newMediaStorageBytesWeek)),
    },
  ];

  return (
    <ContentContainer heightMode="auto" scrollable padding="16px">
      {/* 概览卡：总量 + 近 7 天新增 */}
      <Row gutter={[16, 16]}>
        {overviewItems.map((item) => (
          <Col key={item.title} lg={4} sm={12} xs={24}>
            <OverviewCard
              icon={item.icon}
              title={item.title}
              tone={item.tone}
              value={item.value}
              weeklyTitle={t('analytics.weeklyNew')}
              weeklyValue={item.weeklyValue}
              weeklyValueText={item.weeklyValueText}
            />
          </Col>
        ))}
      </Row>

      {/* 待办提醒：待审评论 / 草稿 / 定时发布 */}
      <div style={{ marginTop: 16 }}>
        <TodoStrip />
      </div>

      {/* AI 安全与异常洞察（平台用户专属：明细为全平台数据且外发模型端点） */}
      {!useUserStore((st) => st.isTenantUser) && (
        <div style={{ marginTop: 16 }}>
          <AiInsightsCard />
        </div>
      )}

      {/* 趋势卡：内容增长 / 登录活跃 / AI 用量 / API 调用（pane 默认懒挂载，ECharts 不会在隐藏容器里以 0 宽初始化） */}
      <Card style={{ marginTop: 16 }} styles={{ body: { paddingTop: 12 } }}>
        <Tabs
          items={[
            {
              key: 'contentTrend',
              label: t('analytics.contentTrend'),
              children: <ContentTrendChart />,
            },
            {
              key: 'loginActivity',
              label: t('analytics.loginActivity'),
              children: <LoginActivityChart />,
            },
            {
              key: 'aiUsage',
              label: t('analytics.aiUsageTrend'),
              children: <AiUsageChart />,
            },
            {
              key: 'apiCalls',
              label: t('analytics.apiCallTrend'),
              children: <ApiCallChart />,
            },
          ]}
        />
      </Card>

      {/* 浏览 TOP 榜（真实 PV）/ 互动 TOP 榜 / 互动·帖子·评论·媒体四分布环形图（满排） */}
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col lg={12} xs={24}>
          <Card
            title={t('analytics.viewTop')}
            styles={{ body: { paddingTop: 12 } }}
          >
            <ViewTopChart />
          </Card>
        </Col>
        <Col lg={12} xs={24}>
          <Card
            title={t('analytics.interactionTop')}
            styles={{ body: { paddingTop: 12 } }}
          >
            <InteractionTopChart />
          </Card>
        </Col>
        <Col lg={6} sm={12} xs={24}>
          <Card
            title={t('analytics.interactionDistribution')}
            styles={{ body: { paddingTop: 12 } }}
          >
            <InteractionDistributionChart />
          </Card>
        </Col>
        <Col lg={6} sm={12} xs={24}>
          <Card
            title={t('analytics.postStatusDistribution')}
            styles={{ body: { paddingTop: 12 } }}
          >
            <PostStatusChart />
          </Card>
        </Col>
        <Col lg={6} sm={12} xs={24}>
          <Card
            title={t('analytics.commentStatusDistribution')}
            styles={{ body: { paddingTop: 12 } }}
          >
            <CommentStatusChart />
          </Card>
        </Col>
        <Col lg={6} sm={12} xs={24}>
          <Card
            title={t('analytics.mediaTypeDistribution')}
            styles={{ body: { paddingTop: 12 } }}
          >
            <MediaTypeChart />
          </Card>
        </Col>
      </Row>
    </ContentContainer>
  );
};

export default Dashboard;
