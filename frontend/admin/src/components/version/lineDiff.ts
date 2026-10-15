/**
 * 最小行级 diff（LCS）：版本历史对比展示用。
 * 输入两段文本，输出按行对齐的差异序列，UI 按 added/removed/unchanged 着色。
 */

export type DiffLineType = 'added' | 'removed' | 'unchanged';

export interface DiffLine {
  type: DiffLineType;
  text: string;
  /** 旧文本行号（removed/unchanged 有值，1 起） */
  oldLine?: number;
  /** 新文本行号（added/unchanged 有值，1 起） */
  newLine?: number;
}

/** LCS 长度表（O(n*m)，正文级文本行数一般 <2k，可接受；超长截断保护） */
function lcsTable(a: string[], b: string[]): number[][] {
  const dp: number[][] = Array.from({ length: a.length + 1 }, () =>
    new Array<number>(b.length + 1).fill(0),
  );
  for (let i = 1; i <= a.length; i++) {
    for (let j = 1; j <= b.length; j++) {
      dp[i][j] = a[i - 1] === b[j - 1] ? dp[i - 1][j - 1] + 1 : Math.max(dp[i - 1][j], dp[i][j - 1]);
    }
  }
  return dp;
}

const MAX_DIFF_LINES = 4000;

export function diffLines(oldText: string, newText: string): DiffLine[] {
  const a = oldText.length > 200000 ? oldText.slice(0, 200000).split('\n') : oldText.split('\n');
  const b = newText.length > 200000 ? newText.slice(0, 200000).split('\n') : newText.split('\n');
  if (a.length * b.length > MAX_DIFF_LINES * MAX_DIFF_LINES) {
    // 异常超大输入直接整体替换展示，避免 UI 卡死
    return [
      ...a.map((text, i) => ({ type: 'removed' as const, text, oldLine: i + 1 })),
      ...b.map((text, j) => ({ type: 'added' as const, text, newLine: j + 1 })),
    ];
  }

  const dp = lcsTable(a, b);
  const result: DiffLine[] = [];
  let i = a.length;
  let j = b.length;
  while (i > 0 || j > 0) {
    if (i > 0 && j > 0 && a[i - 1] === b[j - 1]) {
      result.push({ type: 'unchanged', text: a[i - 1], oldLine: i, newLine: j });
      i--;
      j--;
    } else if (j > 0 && (i === 0 || dp[i][j - 1] >= dp[i - 1][j])) {
      result.push({ type: 'added', text: b[j - 1], newLine: j });
      j--;
    } else {
      result.push({ type: 'removed', text: a[i - 1], oldLine: i });
      i--;
    }
  }
  return result.reverse();
}
