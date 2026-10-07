/**
 * 任务模块枚举映射常量
 */

type TFn = (key: string, options?: Record<string, any>) => string;

// ========== 任务类型 ==========

export const TASK_TYPE_COLORS: Record<string, string> = {
  PERIODIC: 'blue',
  DELAY: 'orange',
  WAIT_RESULT: 'green',
};

export function getTaskTypeMap(t: TFn) {
  return {
    PERIODIC: { text: t('taskType.PERIODIC'), color: TASK_TYPE_COLORS.PERIODIC },
    DELAY: { text: t('taskType.DELAY'), color: TASK_TYPE_COLORS.DELAY },
    WAIT_RESULT: { text: t('taskType.WAIT_RESULT'), color: TASK_TYPE_COLORS.WAIT_RESULT },
  };
}

export function getTaskTypeOptions(t: TFn) {
  return [
    { label: t('taskType.PERIODIC'), value: 'PERIODIC' },
    { label: t('taskType.DELAY'), value: 'DELAY' },
    { label: t('taskType.WAIT_RESULT'), value: 'WAIT_RESULT' },
  ];
}

// ========== 启用状态选项 ==========

export function getEnableOptions(t: TFn) {
  return [
    { label: t('enableMap.true'), value: 'true' },
    { label: t('enableMap.false'), value: 'false' },
  ];
}
