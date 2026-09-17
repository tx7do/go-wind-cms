import { useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import { Card, Col, Row, Spin } from 'antd';
import {
  CommentOutlined,
  FileTextOutlined,
  HeartOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { useI18n } from '@/core/i18n';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import {
  useGetContentTrend,
  useGetDashboardOverview,
  useGetInteractionStats,
  useGetLoginActivity,
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

// 无标题帖子与超长标题的兜底展示
function displayTitle(title: string | undefined, fallback: string): string {
  if (!title) return fallback;
  return title.length > 16 ? `${title.slice(0, 16)}…` : title;
}

const chartBaseStyle = { height: 300, width: '100%' } as const;
const spinStyle = { display: 'block', margin: '80px auto' } as const;

/** 概览卡：主数字为总量，footer 为近 7 天新增 */
function OverviewCard({
  icon,
  title,
  value,
  weeklyTitle,
  weeklyValue,
  tone,
}: {
  icon: React.ReactNode;
  title: string;
  value: number;
  weeklyTitle: string;
  weeklyValue: number;
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
            {value.toLocaleString()}
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
          +{weeklyValue.toLocaleString()}
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
        formatter: (params: any) => {
          const idx = Array.isArray(params) ? params[0].dataIndex : params.dataIndex;
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

/** 互动类型分布（点赞/收藏 环形图） */
const InteractionDistributionChart = () => {
  const { t } = useI18n('dashboard');
  const { data, isLoading } = useGetInteractionStats();

  const option = useMemo(
    () => ({
      legend: { bottom: 0 },
      series: [
        {
          avoidLabelOverlap: true,
          data: [
            {
              itemStyle: { color: '#019680' },
              name: t('analytics.totalLikes'),
              value: toNum(data?.totalLikes),
            },
            {
              itemStyle: { color: '#5ab1ef' },
              name: t('analytics.totalWatches'),
              value: toNum(data?.totalWatches),
            },
          ],
          emphasis: {
            label: { fontSize: 16, fontWeight: 'bold', show: true },
          },
          label: { formatter: '{b}: {c}', show: true },
          radius: ['45%', '75%'],
          type: 'pie',
        },
      ],
      tooltip: { trigger: 'item' },
    }),
    [data, t],
  );

  if (isLoading) return <Spin style={spinStyle} />;
  return <ReactECharts option={option} notMerge style={chartBaseStyle} />;
};

/** 登录活跃趋势（近 30 天成功/失败双线） */
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

/**
 * 分析页（对齐 vben analytics）：概览卡（总量+近7天新增）+ 内容增长趋势
 * + 互动 TOP 榜 / 互动类型分布 / 登录活跃。数据源 statsService 4 RPC。
 */
const Dashboard = () => {
  const { t } = useI18n('dashboard');
  const overviewQuery = useGetDashboardOverview();
  const o = overviewQuery.data;

  const overviewItems = [
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
  ];

  return (
    <ContentContainer heightMode="auto" scrollable padding="16px">
      {/* 概览卡：总量 + 近 7 天新增 */}
      <Row gutter={[16, 16]}>
        {overviewItems.map((item) => (
          <Col key={item.title} lg={6} sm={12} xs={24}>
            <OverviewCard
              icon={item.icon}
              title={item.title}
              tone={item.tone}
              value={item.value}
              weeklyTitle={t('analytics.weeklyNew')}
              weeklyValue={item.weeklyValue}
            />
          </Col>
        ))}
      </Row>

      {/* 内容增长趋势 */}
      <Card
        style={{ marginTop: 16 }}
        title={t('analytics.contentTrend')}
        styles={{ body: { paddingTop: 12 } }}
      >
        <ContentTrendChart />
      </Card>

      {/* 互动 TOP 榜 / 互动类型分布 / 登录活跃 */}
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col lg={8} xs={24}>
          <Card
            title={t('analytics.interactionTop')}
            styles={{ body: { paddingTop: 12 } }}
          >
            <InteractionTopChart />
          </Card>
        </Col>
        <Col lg={8} xs={24}>
          <Card
            title={t('analytics.interactionDistribution')}
            styles={{ body: { paddingTop: 12 } }}
          >
            <InteractionDistributionChart />
          </Card>
        </Col>
        <Col lg={8} xs={24}>
          <Card
            title={t('analytics.loginActivity')}
            styles={{ body: { paddingTop: 12 } }}
          >
            <LoginActivityChart />
          </Card>
        </Col>
      </Row>
    </ContentContainer>
  );
};

export default Dashboard;
