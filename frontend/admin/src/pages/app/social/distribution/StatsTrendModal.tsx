import { useMemo } from 'react';
import { Empty, Modal, Spin } from 'antd';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';

import { useSocialDistributionStatsHistory } from '@/api/hooks/social';

interface StatsTrendModalProps {
  open: boolean;
  distributionId?: number;
  onClose: () => void;
}

interface Series {
  key: 'likes' | 'replies' | 'reposts';
  color: string;
  labelKey: string;
}


const SERIES: Series[] = [
  { key: 'likes', color: '#f5222d', labelKey: 'trend.likes' },
  { key: 'replies', color: '#1677ff', labelKey: 'trend.replies' },
  { key: 'reposts', color: '#52c41a', labelKey: 'trend.reposts' },
];

const W = 640;
const H = 280;
const PAD = { top: 16, right: 16, bottom: 28, left: 40 };

/**
 * 互动统计趋势图（轻量 SVG 折线，不引入图表库）：
 * 三条曲线（赞/评/转）随快照时间的变化；数据来自统计回拉的历史快照，
 * 快照密度取决于周期回拉（每 30 分钟）+ 手动刷新。
 */
const StatsTrendModal: React.FC<StatsTrendModalProps> = ({ open, distributionId, onClose }) => {
  const { t } = useTranslation('social');

  const { data, isLoading } = useSocialDistributionStatsHistory(
    { id: distributionId ?? 0, limit: 200 },
    { enabled: open && !!distributionId },
  );

  const points = useMemo(() => data?.items ?? [], [data]);

  const chart = useMemo(() => {
    if (points.length < 2) return null;

    const maxY = Math.max(
      1,
      ...points.map((p) => Math.max(p.likes ?? 0, p.replies ?? 0, p.reposts ?? 0)),
    );
    const t0 = dayjs(points[0].recordedAt).valueOf();
    const t1 = dayjs(points[points.length - 1].recordedAt).valueOf();

    const iw = W - PAD.left - PAD.right;
    const ih = H - PAD.top - PAD.bottom;
    const x = (i: number) => PAD.left + (i / (points.length - 1)) * iw;
    const y = (v: number) => PAD.top + ih - (v / maxY) * ih;

    const paths = SERIES.map((s) => ({
      ...s,
      d: points.map((p, i) => `${i === 0 ? 'M' : 'L'}${x(i).toFixed(1)},${y(p[s.key] ?? 0).toFixed(1)}`).join(' '),
    }));

    const yTicks = [0, maxY / 2, maxY];

    return { x, y, paths, yTicks, t0, t1 };
  }, [points]);

  return (
    <Modal
      open={open}
      title={t('trendTitle')}
      onCancel={onClose}
      footer={null}
      width={720}
      destroyOnHidden
    >
      {isLoading ? (
        <div style={{ textAlign: 'center', padding: 48 }}>
          <Spin />
        </div>
      ) : points.length < 2 ? (
        <Empty description={t('trendNoData')} style={{ padding: 32 }} />
      ) : (
        <>
          <div style={{ display: 'flex', gap: 16, marginBottom: 8 }}>
            {SERIES.map((s) => (
              <span key={s.key} style={{ fontSize: 12 }}>
                <span
                  style={{
                    display: 'inline-block',
                    width: 16,
                    height: 3,
                    background: s.color,
                    marginRight: 6,
                    verticalAlign: 'middle',
                  }}
                />
                {t(s.labelKey)}
              </span>
            ))}
          </div>
          {chart && (
            <svg width="100%" viewBox={`0 0 ${W} ${H}`} role="img">
              {chart.yTicks.map((v) => (
                <g key={v}>
                  <line
                    x1={PAD.left}
                    x2={W - PAD.right}
                    y1={chart.y(v)}
                    y2={chart.y(v)}
                    stroke="var(--ant-color-split)"
                    strokeDasharray="4 4"
                  />
                  <text x={PAD.left - 6} y={chart.y(v) + 4} textAnchor="end" fontSize={11} fill="var(--ant-color-text-tertiary)">
                    {Math.round(v)}
                  </text>
                </g>
              ))}
              {chart.paths.map((p) => (
                <path key={p.key} d={p.d} fill="none" stroke={p.color} strokeWidth={2} />
              ))}
              <text x={PAD.left} y={H - 8} fontSize={11} fill="var(--ant-color-text-tertiary)">
                {dayjs(chart.t0).format('MM-DD HH:mm')}
              </text>
              <text x={W - PAD.right} y={H - 8} textAnchor="end" fontSize={11} fill="var(--ant-color-text-tertiary)">
                {dayjs(chart.t1).format('MM-DD HH:mm')}
              </text>
            </svg>
          )}
          <div style={{ color: 'var(--ant-color-text-tertiary)', fontSize: 12, marginTop: 4 }}>
            {t('trendHint')}
          </div>
        </>
      )}
    </Modal>
  );
};

export default StatsTrendModal;
