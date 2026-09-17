import type { BreadcrumbItem } from '../types';

export const renderBreadcrumbItems = (
  items: BreadcrumbItem[],
  styleType: 'normal' | 'background',
) => {
  return {
    items: items.map((item) => ({
      title: (
        <span
          className={`breadcrumb-item ${styleType === 'background' ? 'breadcrumb-item-bg' : ''}`}
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: 4,
          }}
        >
          {item.icon}
          <span>{item.breadcrumbName}</span>
        </span>
      ),
      onClick: item.onClick,
    })),
  };
};
